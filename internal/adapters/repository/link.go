package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/dheeraj-vp/golang-url-shortener/internal/core/domain"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type LinkRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewLinkRepository(ctx context.Context, tableName string) (*LinkRepository, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	return &LinkRepository{
		client:    client,
		tableName: tableName,
	}, nil
}

func (d *LinkRepository) All(ctx context.Context) ([]domain.Link, error) {
	var links []domain.Link
	var lastEvaluatedKey map[string]ddbtypes.AttributeValue

	// Paginate through all results
	for {
		input := &dynamodb.ScanInput{
			TableName:         &d.tableName,
			Limit:             aws.Int32(20),
			ExclusiveStartKey: lastEvaluatedKey,
		}

		result, err := d.client.Scan(ctx, input)
		if err != nil {
			return links, fmt.Errorf("failed to get items from DynamoDB: %w", err)
		}

		var pageLinks []domain.Link
		err = attributevalue.UnmarshalListOfMaps(result.Items, &pageLinks)
		if err != nil {
			return links, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
		}

		links = append(links, pageLinks...)

		// Check if there are more pages
		lastEvaluatedKey = result.LastEvaluatedKey
		if lastEvaluatedKey == nil {
			break
		}
	}

	return links, nil
}

// AllWithPagination returns links with pagination support
func (d *LinkRepository) AllWithPagination(ctx context.Context, limit int32, lastKey map[string]ddbtypes.AttributeValue) ([]domain.Link, map[string]ddbtypes.AttributeValue, error) {
	var links []domain.Link

	input := &dynamodb.ScanInput{
		TableName:         &d.tableName,
		Limit:             aws.Int32(limit),
		ExclusiveStartKey: lastKey,
	}

	result, err := d.client.Scan(ctx, input)
	if err != nil {
		return links, nil, fmt.Errorf("failed to get items from DynamoDB: %w", err)
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &links)
	if err != nil {
		return links, nil, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return links, result.LastEvaluatedKey, nil
}

func (d *LinkRepository) Get(ctx context.Context, id string) (domain.Link, error) {
	link := domain.Link{}

	input := &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: id},
		},
	}

	result, err := d.client.GetItem(ctx, input)
	if err != nil {
		return link, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	err = attributevalue.UnmarshalMap(result.Item, &link)
	if err != nil {
		return link, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	return link, nil
}

func (d *LinkRepository) Create(ctx context.Context, link domain.Link) error {
	item, err := attributevalue.MarshalMap(link)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      item,
		ConditionExpression: aws.String("attribute_not_exists(id)"), // Prevent overwrites (collision detection)
	}

	_, err = d.client.PutItem(ctx, input)
	if err != nil {
		// Check if it's a conditional check failure (collision)
		var condCheckErr *ddbtypes.ConditionalCheckFailedException
		if errors.As(err, &condCheckErr) {
			return fmt.Errorf("link with id '%s' already exists: %w", link.Id, err)
		}
		return fmt.Errorf("failed to put item to DynamoDB: %w", err)
	}

	return nil
}

func (d *LinkRepository) Delete(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: id},
		},
	}

	_, err := d.client.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item from DynamoDB: %w", err)
	}
	return nil
}
