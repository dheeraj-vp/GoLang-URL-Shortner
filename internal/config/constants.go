package config

import "time"

// URL validation constants
const (
	MinURLLength  = 15
	ShortIDLength = 8
	MaxRetries    = 3
)

// Cache constants
const (
	DefaultCacheTTL = 24 * time.Hour
	CacheKeyPrefix  = "url:"
)

// HTTP status codes
const (
	StatusCreated     = 201
	StatusNoContent   = 204
	StatusBadRequest  = 400
	StatusNotFound    = 404
	StatusServerError = 500
)

// DynamoDB constants
const (
	DefaultScanLimit   = 20
	MaxBatchGetItems   = 100
	DefaultQueryLimit  = 50
)

// Lambda constants
const (
	DefaultTimeout = 4 * time.Second
	MaxTimeout     = 29 * time.Second // API Gateway timeout is 30s
)

// Platform detection patterns
const (
	InstagramUserAgent = "Instagram"
	TwitterUserAgent   = "Twitter"
	YouTubeUserAgent   = "YouTube"
	InstagramReferer   = "instagram.com"
	TwitterReferer     = "twitter.com"
	YouTubeReferer     = "youtube.com"
)
