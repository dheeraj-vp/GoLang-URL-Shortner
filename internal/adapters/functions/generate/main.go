package main

import (
	"context"

	"github.com/dheeraj-vp/golang-url-shortener/internal/adapters/cache"
	"github.com/dheeraj-vp/golang-url-shortener/internal/adapters/handlers"
	"github.com/dheeraj-vp/golang-url-shortener/internal/adapters/repository"
	"github.com/dheeraj-vp/golang-url-shortener/internal/config"
	"github.com/dheeraj-vp/golang-url-shortener/internal/core/services"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	appConfig := config.NewConfig()
	redisAddress, redisPassword, redisDB := appConfig.GetRedisParams()
	cache := cache.NewRedisCache(redisAddress, redisPassword, redisDB)
	linkTableName := appConfig.GetLinkTableName()
	statsTableName := appConfig.GetStatsTableName()

	linkRepo := repository.NewLinkRepository(context.TODO(), linkTableName)
	linkService := services.NewLinkService(linkRepo, cache)

	statsRepo := repository.NewStatsRepository(context.TODO(), statsTableName)
	statsService := services.NewStatsService(statsRepo, cache)

	handler := handlers.NewGenerateLinkFunctionHandler(linkService, statsService)
	lambda.Start(handler.CreateShortLink)
}
