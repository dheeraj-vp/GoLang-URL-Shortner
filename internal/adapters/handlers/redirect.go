package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dheeraj-vp/golang-url-shortener/internal/config"
	"github.com/dheeraj-vp/golang-url-shortener/internal/core/domain"
	"github.com/dheeraj-vp/golang-url-shortener/internal/core/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type RedirectFunctionHandler struct {
	linkService  *services.LinkService
	statsService *services.StatsService
}

func NewRedirectFunctionHandler(l *services.LinkService, s *services.StatsService) *RedirectFunctionHandler {
	return &RedirectFunctionHandler{linkService: l, statsService: s}
}

func (h *RedirectFunctionHandler) Redirect(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	// Add context timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()

	pathSegments := strings.Split(req.RawPath, "/")
	if len(pathSegments) < 2 {
		return ClientError(http.StatusBadRequest, "Invalid URL path")
	}

	shortLinkKey := pathSegments[len(pathSegments)-1]
	if shortLinkKey == "" {
		return ClientError(http.StatusBadRequest, "Short link key cannot be empty")
	}

	longLink, err := h.linkService.GetOriginalURL(timeoutCtx, shortLinkKey)
	if err != nil || longLink == nil || *longLink == "" {
		return ClientError(http.StatusNotFound, "Link not found")
	}

	// Extract platform from request headers
	platform := ExtractPlatformFromRequest(req)

	// Create stats asynchronously to not block the redirect
	go func() {
		statsCtx := context.Background()
		if err := h.statsService.Create(statsCtx, domain.Stats{
			Id:        uuid.NewString(),
			LinkID:    shortLinkKey,
			CreatedAt: time.Now(),
			Platform:  platform,
		}); err != nil {
			log.Printf("Failed to create stats for link '%s': %v", shortLinkKey, err)
		}
	}()

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMovedPermanently,
		Headers: map[string]string{
			"Location":      *longLink,
			"Cache-Control": "public, max-age=300", // Cache for 5 minutes
		},
	}, nil
}
