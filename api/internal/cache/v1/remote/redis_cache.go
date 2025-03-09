package remote

import (
	"context"
	"errors"
	"fmt"
	cachev1 "github.com/fjarm/fjarm/api/internal/cache/v1"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

const redisCacheTag = "redis_cache"

// TODO(2025-03-09): Test RedisCache implementation.
// TODO(2025-03-09): Implement TLS client and server support.
// TODO(2025-03-09): Use Redis Sentinel.
// TODO(2025-03-09): Manually configure Redis connection pooling.
// TODO(2025-03-09): Add OpenTelemetry-based monitoring to Redis.

// RedisCache is a distributed cache that uses Redis Sentinel for high availability.
type RedisCache struct {
	client *redis.Client
	logger *slog.Logger
}

// Get retrieves the value associated with the supplied key from the remote Redis cache. If no such key/value pair
// exists, a v1.ErrCacheMiss error is returned. Other errors indicate something more serious.
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	logger := c.logger.With(slog.String(logkeys.Tag, redisCacheTag), slog.String("key", key))
	logger.DebugContext(ctx, "attempted to get a key from Redis cache")

	res, err := c.client.Get(ctx, key).Bytes()
	if err != nil && errors.Is(err, redis.Nil) {
		// The key doesn't exist in the cache. This is an innocuous error.
		logger.DebugContext(ctx, "failed to find key in Redis cache")
		return nil, fmt.Errorf("%w: %w", cachev1.ErrCacheMiss, err)
	} else if err != nil {
		return nil, err
	}
	return res, nil
}

// Set adds the supplied key/value pair to the Redis cache. If the key already exists, a v1.ErrKeyExists error is
// returned. Other errors indicate something more serious.
func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	logger := c.logger.With(slog.String(logkeys.Tag, redisCacheTag), slog.String("key", key))
	logger.DebugContext(ctx, "attempted to set a key in Redis cache")

	success, err := c.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		logger.WarnContext(ctx, "failed to set key in Redis cache")
		return err
	}
	if !success {
		logger.DebugContext(ctx, "key already exists in Redis cache")
		return fmt.Errorf("%w: key %s already exists", cachev1.ErrKeyExists, key)
	}
	return nil
}

// Update adds the supplied key/value pair to the Redis cache. If the key already exists, the associated value is
// overwritten.
func (c *RedisCache) Update(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	logger := c.logger.With(slog.String(logkeys.Tag, redisCacheTag), slog.String("key", key))
	logger.DebugContext(ctx, "attempted to update a key in Redis cache")

	_, err := c.client.Set(ctx, key, value, ttl).Result()
	if err != nil {
		logger.WarnContext(ctx, "failed to set key in Redis cache")
		return err
	}
	return nil
}

// NewRedisCache creates a new RedisCache instance.
func NewRedisCache(client *redis.Client, logger *slog.Logger) *RedisCache {
	return &RedisCache{
		client: client,
		logger: logger,
	}
}
