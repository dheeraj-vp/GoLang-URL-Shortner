package unit

import (
	"context"
	"testing"

	"github.com/dheeraj-vp/golang-url-shortener/internal/adapters/cache"
	"github.com/dheeraj-vp/golang-url-shortener/internal/adapters/handlers"
	"github.com/dheeraj-vp/golang-url-shortener/internal/core/services"
	"github.com/dheeraj-vp/golang-url-shortener/internal/tests/mock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestGenerateLinkUnit(t *testing.T) {
	mockLinkRepo := mock.NewMockLinkRepo()
	mockStats := mock.NewMockStatsRepo()
	mockCache := mock.NewImprovedMockCache() // Use mock cache instead of real Redis
	linkService := services.NewLinkService(mockLinkRepo, mockCache)
	statsService := services.NewStatsService(mockStats, mockCache)
	apiHandler := handlers.NewGenerateLinkFunctionHandler(linkService, statsService)

	tests := []struct {
		longURL            string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			longURL:            "https://example.com/link1",
			expectedStatusCode: 200,
			expectedBody:       "",
		},
		{
			longURL:            "",
			expectedStatusCode: 400,
			expectedBody:       "URL cannot be empty",
		},
		{
			longURL:            "invalid",
			expectedStatusCode: 400,
			expectedBody:       "URL must be at least 15 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.longURL, func(t *testing.T) {
			body := `{"long": "` + tt.longURL + `"}`
			request := events.APIGatewayV2HTTPRequest{Body: body}
			response, err := apiHandler.CreateShortLink(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, response.StatusCode)

			if tt.expectedStatusCode != 200 {
				assert.Equal(t, tt.expectedBody, response.Body)
			}
		})
	}

}
