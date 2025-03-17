package remote

import (
	"context"
	"fmt"
	cachev1 "github.com/fjarm/fjarm/api/internal/cache/v1"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/redis/rueidis"
	"log/slog"
	"strings"
	"time"
)

const redisCacheTag = "redis_cache"

// TODO(2025-03-09): Setup Redis Sentinel in docker-compose.yaml file.
// TODO(2025-03-09): Manually configure Redis connection pooling.
// TODO(2025-03-09): Add OpenTelemetry-based monitoring to Redis.

// RedisCache is a distributed cache that uses Redis Sentinel for high availability.
type RedisCache struct {
	logger *slog.Logger
	rdb    rueidis.Client
}

// Get retrieves the value associated with the supplied key from the remote Redis cache. If no such key/value pair
// exists, a v1.ErrCacheMiss error is returned. Other errors indicate something more serious.
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	if strings.TrimSpace(key) == "" || strings.Contains(key, " ") {
		return nil, fmt.Errorf("%w: %s", cachev1.ErrInvalidKey, "key cannot be empty or contain whitespace")
	}

	logger := c.logger.With(slog.String(logkeys.Tag, redisCacheTag), slog.String("key", key))
	logger.DebugContext(ctx, "attempted to get a key from Redis cache")

	cmd := c.rdb.B().Get().Key(key).Build()
	res, err := c.rdb.Do(ctx, cmd).AsBytes()
	if err != nil && rueidis.IsRedisNil(err) {
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
	if strings.TrimSpace(key) == "" || strings.Contains(key, " ") {
		return fmt.Errorf("%w: %s", cachev1.ErrInvalidKey, "key cannot be empty or whitespace")
	}
	if ttl <= 0 {
		return fmt.Errorf("%w: %s", cachev1.ErrInvalidExpiration, "ttl must be greater than 0")
	}

	logger := c.logger.With(slog.String(logkeys.Tag, redisCacheTag), slog.String("key", key))
	logger.DebugContext(ctx, "attempted to set a key in Redis cache")

	cmd := c.rdb.B().Set().Key(key).Value(rueidis.BinaryString(value)).Nx().Ex(ttl).Build()
	err := c.rdb.Do(ctx, cmd).Error()
	if err != nil && rueidis.IsRedisNil(err) {
		// The key already exists in the cache. This is an innocuous error.
		logger.DebugContext(ctx, "failed to set existing key in Redis cache", slog.String("key", key))
		return fmt.Errorf("%w: %w", cachev1.ErrKeyExists, err)
	} else if err != nil {
		logger.WarnContext(ctx, "failed to set key in Redis cache", slog.Any(logkeys.Err, err))
		return err
	}
	return nil
}

// Update adds the supplied key/value pair to the Redis cache. If the key already exists, the associated value is
// overwritten.
func (c *RedisCache) Update(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if strings.TrimSpace(key) == "" || strings.Contains(key, " ") {
		return fmt.Errorf("%w: %s", cachev1.ErrInvalidKey, "key cannot be empty or whitespace")
	}
	if ttl <= 0 {
		return fmt.Errorf("%w: %s", cachev1.ErrInvalidExpiration, "ttl must be greater than 0")
	}

	logger := c.logger.With(slog.String(logkeys.Tag, redisCacheTag), slog.String("key", key))
	logger.DebugContext(ctx, "attempted to update a key in Redis cache")

	cmd := c.rdb.B().Set().Key(key).Value(rueidis.BinaryString(value)).Ex(ttl).Build()
	err := c.rdb.Do(ctx, cmd).Error()
	if err != nil {
		logger.WarnContext(ctx, "failed to update key in Redis cache")
		return err
	}
	return nil
}

// NewRedisCache creates a new RedisCache instance.
func NewRedisCache(rdb rueidis.Client, logger *slog.Logger) *RedisCache {
	return &RedisCache{
		logger: logger,
		rdb:    rdb,
	}
}
