package handlers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/dheeraj-vp/golang-url-shortener/internal/config"
	"github.com/dheeraj-vp/golang-url-shortener/internal/core/domain"
	"github.com/dheeraj-vp/golang-url-shortener/internal/core/services"
	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type RequestBody struct {
	Long string `json:"long"`
}

type GenerateLinkFunctionHandler struct {
	linkService  *services.LinkService
	statsService *services.StatsService
}

func NewGenerateLinkFunctionHandler(l *services.LinkService, s *services.StatsService) *GenerateLinkFunctionHandler {
	return &GenerateLinkFunctionHandler{linkService: l, statsService: s}
}

func (h *GenerateLinkFunctionHandler) CreateShortLink(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	// Add context timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()

	var requestBody RequestBody
	err := json.Unmarshal([]byte(req.Body), &requestBody)
	if err != nil {
		return ClientError(http.StatusBadRequest, "Invalid JSON")
	}

	// Validation
	if requestBody.Long == "" {
		return ClientError(http.StatusBadRequest, "URL cannot be empty")
	}
	if len(requestBody.Long) < config.MinURLLength {
		return ClientError(http.StatusBadRequest, fmt.Sprintf("URL must be at least %d characters long", config.MinURLLength))
	}
	if !IsValidLink(requestBody.Long) {
		return ClientError(http.StatusBadRequest, "Invalid URL format")
	}
	if IsMaliciousURL(requestBody.Long) {
		return ClientError(http.StatusBadRequest, "URL contains malicious patterns")
	}

	// Generate short URL with collision detection
	var link domain.Link
	var createErr error
	for i := 0; i < config.MaxRetries; i++ {
		link = domain.Link{
			Id:          GenerateShortURLID(config.ShortIDLength),
			OriginalURL: requestBody.Long,
			CreatedAt:   time.Now(),
		}

		createErr = h.linkService.Create(timeoutCtx, link)
		if createErr == nil {
			break // Success
		}

		// Check if it's a collision (DynamoDB conditional check failed)
		// In production, you'd want to check the specific error type
		log.Printf("Failed to create link (attempt %d/%d): %v", i+1, config.MaxRetries, createErr)
	}

	if createErr != nil {
		return ServerError(createErr)
	}

	js, err := json.Marshal(link)
	if err != nil {
		return ServerError(err)
	}

	// Send notification asynchronously
	go sendMessageToQueue(context.Background(), link)

	// Return 201 Created (proper REST status code)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(js),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func sendMessageToQueue(ctx context.Context, link domain.Link) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("unable to load SDK config: %v", err)
		return
	}

	sqsClient := sqs.NewFromConfig(cfg)
	queueUrl := os.Getenv("QueueUrl")

	if queueUrl == "" {
		log.Println("QueueUrl is not set, skipping notification")
		return
	}

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: aws.String("The system generated a short URL with the ID " + link.Id),
	})

	if err != nil {
		log.Printf("Failed to send message to SQS: %v", err)
	}
}

func GenerateShortURLID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to timestamp-based generation if crypto/rand fails
			log.Printf("Failed to generate random number: %v", err)
			charIndex = big.NewInt(int64(time.Now().UnixNano() % int64(len(charset))))
		}
		result[i] = charset[charIndex.Int64()]
	}
	return string(result)
}
