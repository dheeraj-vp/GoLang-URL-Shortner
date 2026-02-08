package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/dheeraj-vp/golang-url-shortener/internal/config"
	"github.com/dheeraj-vp/golang-url-shortener/internal/core/services"
	"github.com/aws/aws-lambda-go/events"
)

type StatsFunctionHandler struct {
	statsService *services.StatsService
	linkService  *services.LinkService
}

func NewStatsFunctionHandler(l *services.LinkService, s *services.StatsService) *StatsFunctionHandler {
	return &StatsFunctionHandler{linkService: l, statsService: s}
}

func (h *StatsFunctionHandler) Stats(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	// Add context timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()

	links, err := h.linkService.GetAll(timeoutCtx)
	if err != nil {
		return ServerError(err)
	}

	// Use goroutines to fetch stats concurrently (mitigate N+1 problem)
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	for i := range links {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			
			stats, err := h.statsService.GetStatsByLinkID(timeoutCtx, links[index].Id)
			if err != nil {
				log.Printf("Error getting stats for link '%s': %v", links[index].Id, err)
				return
			}
			
			mu.Lock()
			links[index].Stats = stats
			mu.Unlock()
		}(i)
	}
	
	wg.Wait()

	jsonResponse, err := json.Marshal(links)
	if err != nil {
		return ServerError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonResponse),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// GetLinkStats returns stats for a specific link ID
func (h *StatsFunctionHandler) GetLinkStats(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()

	linkID := req.PathParameters["id"]
	if linkID == "" {
		return ClientError(http.StatusBadRequest, "Link ID is required")
	}

	stats, err := h.statsService.GetStatsByLinkID(timeoutCtx, linkID)
	if err != nil {
		return ServerError(err)
	}

	// Count by platform
	platformCounts := make(map[string]int)
	for _, stat := range stats {
		platformCounts[stat.Platform.String()]++
	}

	response := map[string]interface{}{
		"link_id":         linkID,
		"total_clicks":    len(stats),
		"platform_counts": platformCounts,
		"details":         stats,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return ServerError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonResponse),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
