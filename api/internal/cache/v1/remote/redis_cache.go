package remote

import (
	"context"
	"errors"
	"fmt"
	cachev1 "github.com/fjarm/fjarm/api/internal/cache/v1"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/redis/go-redis/v9"
	"github.com/redis/rueidis"
	"log/slog"
	"time"
)

const redisCacheTag = "redis_cache"

// TODO(2025-03-09): Replace go-redis with rueidis.
// TODO(2025-03-09): Test RedisCache implementation.
// TODO(2025-03-09): Use Redis Sentinel.
// TODO(2025-03-09): Manually configure Redis connection pooling.
// TODO(2025-03-09): Add OpenTelemetry-based monitoring to Redis.

// RedisCache is a distributed cache that uses Redis Sentinel for high availability.
type RedisCache struct {
	client *redis.Client
	logger *slog.Logger
	rdb    rueidis.Client
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

// newRedisClient creates a new Redis client using rueidis.
func newRedisClient() (rueidis.Client, error) {
	client, err := rueidis.NewClient(
		rueidis.ClientOption{
			// TODO(2025-03-09): Implement TLS client and server support. Load TLS cert/key and CA cert using infisical.
			TLSConfig: nil,
			// TODO(2025-03-09): Supply AuthCredentialsFn to provide username and password for ACL support.
			AuthCredentialsFn: nil,
			InitAddress: []string{
				// When running Sentinel mode, all node addresses need to be supplied. In Cluster mode, only the one
				// address needs to be supplied.
				"redis-cluster.railway.internal:6379",
			},
			ClientTrackingOptions: []string{
				// This is the default value. Keys mentioned in read operations aren't cached. Caching must be
				// proactively turned on immediately before the actual command to enable client-side caching.
				"OPTIN",
			},
			// TODO(2025-03-09): Allow specifying CacheSizeEachConn when client-side caching is enabled.
			BlockingPoolCleanup: 30 * time.Second,
			MaxFlushDelay:       0,
			// TODO(2025-03-09): Set ShardsRefreshInterval to non-zero value after enabling Redis Cluster.
			//ClusterOption:         rueidis.ClusterOption{
			//	ShardsRefreshInterval: 0,
			//},
			DisableCache:          true, // Disable client-side caching.
			DisableAutoPipelining: true, // Manual pipelining can be enabled using client.DoMulti().
			// Toggled to true for read-only clients. But this should be accomplished using ACLs.
			ReplicaOnly: false,
		},
	)
	return client, err
}

// NewRedisCache creates a new RedisCache instance.
func NewRedisCache(logger *slog.Logger) *RedisCache {
	client, err := newRedisClient()
	if err != nil {
		logger.Error("failed to create Redis client", slog.Any(logkeys.Err, err))
		return nil
	}
	return &RedisCache{
		logger: logger,
		rdb:    client,
	}
}
