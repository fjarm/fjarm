package remote

import (
	"context"
	"fmt"
	"strings"
	"time"

	cachev1 "github.com/fjarm/fjarm/api/internal/cache"
)

// FakeRedisCache is a duplicate of RedisCache that doesn't actually connect to Redis. Instead, it uses an in-memory map
// to replicate the behavior of Redis for testing purposes.
type FakeRedisCache struct {
	rdb map[string]*fakeRedisCacheValue
}

type fakeRedisCacheValue struct {
	val       []byte
	expiredAt time.Time
}

func (c *FakeRedisCache) removeExpiredValues() {
	keysToRemove := []string{}
	for key, value := range c.rdb {
		if time.Now().After(value.expiredAt) {
			keysToRemove = append(keysToRemove, key)
		}
	}
	for _, key := range keysToRemove {
		delete(c.rdb, key)
	}
}

// Get retrieves the value associated with the supplied key from the in-memory map. If no such key/value pair exists, a
// pkg.ErrCacheMiss error is returned. Other errors indicate something more serious.
func (c *FakeRedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	if strings.TrimSpace(key) == "" || strings.Contains(key, " ") {
		return nil, fmt.Errorf("%w: %s", cachev1.ErrInvalidKey, "key cannot be empty or contain whitespace")
	}
	c.removeExpiredValues()

	value, exists := c.rdb[key]
	if !exists {
		return nil, fmt.Errorf("%w: key not found", cachev1.ErrCacheMiss)
	}

	return value.val, nil
}

// Set adds the supplied key/value pair to the in-memory map. If the key already exists, an error is returned.
func (c *FakeRedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if strings.TrimSpace(key) == "" || strings.Contains(key, " ") {
		return fmt.Errorf("%w: %s", cachev1.ErrInvalidKey, "key cannot be empty or whitespace")
	}
	if ttl <= 0 {
		return fmt.Errorf("%w: %s", cachev1.ErrInvalidExpiration, "ttl must be greater than 0")
	}
	c.removeExpiredValues()

	_, exists := c.rdb[key]
	if exists {
		return fmt.Errorf("%w: key already exists", cachev1.ErrKeyExists)
	}

	cacheValue := &fakeRedisCacheValue{
		val:       value,
		expiredAt: time.Now().Add(ttl),
	}
	c.rdb[key] = cacheValue
	return nil
}

// Update adds the supplied key/value pair to the in-memory map. If the key already exists, the value is overwritten.
func (c *FakeRedisCache) Update(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if strings.TrimSpace(key) == "" || strings.Contains(key, " ") {
		return fmt.Errorf("%w: %s", cachev1.ErrInvalidKey, "key cannot be empty or whitespace")
	}
	if ttl <= 0 {
		return fmt.Errorf("%w: %s", cachev1.ErrInvalidExpiration, "ttl must be greater than 0")
	}
	c.removeExpiredValues()

	cacheValue := &fakeRedisCacheValue{
		val:       value,
		expiredAt: time.Now().Add(ttl),
	}
	c.rdb[key] = cacheValue
	return nil
}

// NewFakeRedisCache creates a new instance of FakeRedisCache.
func NewFakeRedisCache() *FakeRedisCache {
	return &FakeRedisCache{
		rdb: make(map[string]*fakeRedisCacheValue),
	}
}
