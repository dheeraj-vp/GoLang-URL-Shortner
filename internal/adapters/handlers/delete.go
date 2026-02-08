package handlers

import (
	"context"
	"net/http"

	"github.com/dheeraj-vp/golang-url-shortener/internal/config"
	"github.com/dheeraj-vp/golang-url-shortener/internal/core/services"
	"github.com/aws/aws-lambda-go/events"
)

type DeleteFunctionHandler struct {
	statsService *services.StatsService
	linkService  *services.LinkService
}

func NewDeleteFunctionHandler(l *services.LinkService, s *services.StatsService) *DeleteFunctionHandler {
	return &DeleteFunctionHandler{linkService: l, statsService: s}
}

func (h *DeleteFunctionHandler) Delete(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	// Add context timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()

	id := req.PathParameters["id"]
	if id == "" {
		return ClientError(http.StatusBadRequest, "ID parameter is required")
	}

	// Delete link first
	err := h.linkService.Delete(timeoutCtx, id)
	if err != nil {
		return ServerError(err)
	}

	// Delete associated stats
	err = h.statsService.Delete(timeoutCtx, id)
	if err != nil {
		// Log the error but don't fail the request
		// Stats deletion is not critical
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNoContent,
			Body:       "Link deleted but stats deletion failed",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}
