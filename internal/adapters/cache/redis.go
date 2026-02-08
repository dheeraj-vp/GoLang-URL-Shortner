package cache

import (
	"context"
	"time"

	"github.com/dheeraj-vp/golang-url-shortener/internal/config"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(address string, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: client,
		ttl:    config.DefaultCacheTTL,
	}
}

func NewRedisCacheWithTTL(address string, password string, db int, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

func (r *RedisCache) Set(ctx context.Context, key string, val string) error {
	// Add key prefix for better organization
	fullKey := config.CacheKeyPrefix + key
	return r.client.Set(ctx, fullKey, val, r.ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	fullKey := config.CacheKeyPrefix + key
	val, err := r.client.Get(ctx, fullKey).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	fullKey := config.CacheKeyPrefix + key
	return r.client.Del(ctx, fullKey).Err()
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
