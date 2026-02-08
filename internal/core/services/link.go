package services

import (
	"context"
	"fmt"
	"log"

	"github.com/dheeraj-vp/golang-url-shortener/internal/core/domain"
	"github.com/dheeraj-vp/golang-url-shortener/internal/core/ports"
)

type LinkService struct {
	port  ports.LinkPort
	cache ports.Cache
}

func NewLinkService(p ports.LinkPort, c ports.Cache) *LinkService {
	return &LinkService{port: p, cache: c}
}

func (service *LinkService) GetAll(ctx context.Context) ([]domain.Link, error) {
	links, err := service.port.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all links: %w", err)
	}
	return links, nil
}

func (service *LinkService) GetOriginalURL(ctx context.Context, shortLinkKey string) (*string, error) {
	// Try cache first (cache-aside pattern)
	cachedURL, err := service.cache.Get(ctx, shortLinkKey)
	if err == nil && cachedURL != "" {
		// Cache hit
		log.Printf("Cache hit for key: %s", shortLinkKey)
		return &cachedURL, nil
	}

	// Cache miss - fetch from database
	log.Printf("Cache miss for key: %s, fetching from database", shortLinkKey)
	data, err := service.port.Get(ctx, shortLinkKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get short URL for identifier '%s': %w", shortLinkKey, err)
	}

	// Validate link exists and has URL
	if data.OriginalURL == "" {
		return nil, fmt.Errorf("link '%s' not found or has no URL", shortLinkKey)
	}

	// Populate cache asynchronously to avoid blocking the response
	go func() {
		if err := service.cache.Set(context.Background(), shortLinkKey, data.OriginalURL); err != nil {
			log.Printf("Failed to populate cache for key '%s': %v", shortLinkKey, err)
		}
	}()

	return &data.OriginalURL, nil
}

func (service *LinkService) Create(ctx context.Context, link domain.Link) error {
	// Create in database first
	if err := service.port.Create(ctx, link); err != nil {
		return fmt.Errorf("failed to create short URL: %w", err)
	}

	// Populate cache asynchronously
	go func() {
		if err := service.cache.Set(context.Background(), link.Id, link.OriginalURL); err != nil {
			log.Printf("Failed to populate cache for new link '%s': %v", link.Id, err)
		}
	}()

	return nil
}

func (service *LinkService) Delete(ctx context.Context, short string) error {
	// Delete from database
	if err := service.port.Delete(ctx, short); err != nil {
		return fmt.Errorf("failed to delete short URL for identifier '%s': %w", short, err)
	}

	// Delete from cache asynchronously
	go func() {
		if err := service.cache.Delete(context.Background(), short); err != nil {
			log.Printf("Failed to delete from cache for key '%s': %v", short, err)
		}
	}()

	return nil
}
