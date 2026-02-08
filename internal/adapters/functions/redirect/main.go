package main

import (
	"context"
	"log"

	"github.com/dheeraj-vp/golang-url-shortener/internal/adapters/cache"
	"github.com/dheeraj-vp/golang-url-shortener/internal/adapters/handlers"
	"github.com/dheeraj-vp/golang-url-shortener/internal/adapters/repository"
	"github.com/dheeraj-vp/golang-url-shortener/internal/config"
	"github.com/dheeraj-vp/golang-url-shortener/internal/core/services"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	ctx := context.Background()
	appConfig := config.NewConfig()
	redisAddress, redisPassword, redisDB := appConfig.GetRedisParams()
	cache := cache.NewRedisCache(redisAddress, redisPassword, redisDB)
	linkTableName := appConfig.GetLinkTableName()
	statsTableName := appConfig.GetStatsTableName()

	linkRepo, err := repository.NewLinkRepository(ctx, linkTableName)
	if err != nil {
		log.Fatalf("failed to create link repository: %v", err)
	}
	linkService := services.NewLinkService(linkRepo, cache)

	statsRepo, err := repository.NewStatsRepository(ctx, statsTableName)
	if err != nil {
		log.Fatalf("failed to create stats repository: %v", err)
	}
	statsService := services.NewStatsService(statsRepo, cache)

	handler := handlers.NewRedirectFunctionHandler(linkService, statsService)

	lambda.Start(handler.Redirect)
}
