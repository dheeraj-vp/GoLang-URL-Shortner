package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	dynamoTableName string // DynamoDB table name
	redisAddress    string // Redis address
	redisPassword   string // Redis password
	redisDB         int    // Redis DB
	slackToken      string // Slack token
	slackChannelID  string // Slack channel ID
}

func NewConfig() *AppConfig {
	return &AppConfig{
		dynamoTableName: "UrlShortenerTable", // default value
		redisAddress:    "localhost:6379",    // default value
		redisPassword:   "",                  // default value
		redisDB:         0,                   // default value
		slackToken:      "",                  // default value
		slackChannelID:  "",                  // default value
	}
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Error loading .env file: ", err)
	}
}

func (c *AppConfig) GetSlackParams() (string, string) {
	slackToken, tokenOK := os.LookupEnv("SlackToken")
	slackChannelID, channelOK := os.LookupEnv("SlackChannelID")
	if !tokenOK || !channelOK {
		return os.Getenv("SlackToken"), os.Getenv("SlackChannelID")
	}
	return slackToken, slackChannelID
}

func (c *AppConfig) GetLinkTableName() string {
	tableName, ok := os.LookupEnv("LinkTableName")
	if !ok {
		log.Printf("Warning: LinkTableName environment variable not set, using default: %s", c.dynamoTableName)
		return c.dynamoTableName
	}
	if tableName == "" {
		log.Printf("Warning: LinkTableName is empty, using default: %s", c.dynamoTableName)
		return c.dynamoTableName
	}
	return tableName
}

func (c *AppConfig) GetStatsTableName() string {
	tableName, ok := os.LookupEnv("StatsTableName")
	if !ok {
		log.Printf("Warning: StatsTableName environment variable not set, using default")
		return "" // Return empty string - caller should handle this
	}
	if tableName == "" {
		log.Printf("Warning: StatsTableName is empty")
		return ""
	}
	return tableName
}

func (c *AppConfig) GetRedisParams() (string, string, int) {
	address, ok := os.LookupEnv("RedisAddress")
	if !ok {
		log.Printf("Warning: RedisAddress environment variable not set, using default: %s", c.redisAddress)
		return c.redisAddress, c.redisPassword, c.redisDB
	}

	password, ok := os.LookupEnv("RedisPassword")
	if !ok {
		log.Printf("Warning: RedisPassword environment variable not set, using default (empty)")
		return address, c.redisPassword, c.redisDB
	}

	dbStr, ok := os.LookupEnv("RedisDB")
	if !ok {
		log.Printf("Warning: RedisDB environment variable not set, using default: %d", c.redisDB)
		return address, password, c.redisDB
	}

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		log.Printf("Warning: RedisDB environment variable is not a valid integer (%v), using default: %d", err, c.redisDB)
		return address, password, c.redisDB
	}

	return address, password, db
}
