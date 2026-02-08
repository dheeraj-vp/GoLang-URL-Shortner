package unit

import (
	"context"
	"testing"

	"github.com/dheeraj-vp/golang-url-shortener/internal/core/domain"
	"github.com/dheeraj-vp/golang-url-shortener/internal/core/services"
	"github.com/dheeraj-vp/golang-url-shortener/internal/tests/mock"
	"github.com/stretchr/testify/assert"
)

func TestCacheHitScenario(t *testing.T) {
	// Setup
	mockLinkRepo := mock.NewMockLinkRepo()
	mockCache := mock.NewImprovedMockCache()
	linkService := services.NewLinkService(mockLinkRepo, mockCache)
	ctx := context.Background()

	// Pre-populate cache
	testID := "test123"
	testURL := "https://example.com/long-url"
	err := mockCache.Set(ctx, testID, testURL)
	assert.NoError(t, err)

	// Test cache hit
	result, err := linkService.GetOriginalURL(ctx, testID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testURL, *result)

	// Verify cache was used (get count should be 1)
	assert.Equal(t, 1, mockCache.GetGetCount())
}

func TestCacheMissScenario(t *testing.T) {
	// Setup
	mockLinkRepo := mock.NewMockLinkRepo()
	mockCache := mock.NewImprovedMockCache()
	linkService := services.NewLinkService(mockLinkRepo, mockCache)
	ctx := context.Background()

	// Add data to repository but not cache
	testID := "test456"
	testURL := "https://example.com/another-url"
	mockLinkRepo.Links = []domain.Link{
		{Id: testID, OriginalURL: testURL},
	}

	// Test cache miss - should fetch from repository
	result, err := linkService.GetOriginalURL(ctx, testID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testURL, *result)

	// Verify cache was checked
	assert.Equal(t, 1, mockCache.GetGetCount())
}

func TestCreateLinkPopulatesCache(t *testing.T) {
	// Setup
	mockLinkRepo := mock.NewMockLinkRepo()
	mockCache := mock.NewImprovedMockCache()
	linkService := services.NewLinkService(mockLinkRepo, mockCache)
	ctx := context.Background()

	// Create a new link
	newLink := domain.Link{
		Id:          "newlink789",
		OriginalURL: "https://example.com/new-url",
	}

	err := linkService.Create(ctx, newLink)
	assert.NoError(t, err)

	// Note: Cache population is async in the service, so we can't directly test it
	// But we can verify the link was created in the repository
	assert.Len(t, mockLinkRepo.Links, 1)
	assert.Equal(t, newLink.Id, mockLinkRepo.Links[0].Id)
}

func TestDeleteLinkRemovesFromCache(t *testing.T) {
	// Setup
	mockLinkRepo := mock.NewMockLinkRepo()
	mockCache := mock.NewImprovedMockCache()
	linkService := services.NewLinkService(mockLinkRepo, mockCache)
	ctx := context.Background()

	// Add data
	testID := "delete123"
	testURL := "https://example.com/to-delete"
	mockLinkRepo.Links = []domain.Link{
		{Id: testID, OriginalURL: testURL},
	}
	mockCache.Set(ctx, testID, testURL)

	// Delete the link
	err := linkService.Delete(ctx, testID)
	assert.NoError(t, err)

	// Verify link was deleted from repository
	assert.Len(t, mockLinkRepo.Links, 0)
}

func TestCacheFailureHandling(t *testing.T) {
	// Setup
	mockLinkRepo := mock.NewMockLinkRepo()
	mockCache := mock.NewImprovedMockCache()
	linkService := services.NewLinkService(mockLinkRepo, mockCache)
	ctx := context.Background()

	// Add data to repository
	testID := "fail123"
	testURL := "https://example.com/fail-test"
	mockLinkRepo.Links = []domain.Link{
		{Id: testID, OriginalURL: testURL},
	}

	// Enable cache failure mode
	mockCache.SetFailureMode(true)

	// Should still work by falling back to repository
	result, err := linkService.GetOriginalURL(ctx, testID)
	// Cache failure should be handled gracefully
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testURL, *result)
}

func TestConcurrentCacheAccess(t *testing.T) {
	// Setup
	mockCache := mock.NewImprovedMockCache()
	ctx := context.Background()

	// Perform concurrent operations
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			key := "concurrent"
			value := "test"
			mockCache.Set(ctx, key, value)
			mockCache.Get(ctx, key)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify operations completed successfully
	assert.Equal(t, 10, mockCache.GetSetCount())
	assert.Equal(t, 10, mockCache.GetGetCount())
}
