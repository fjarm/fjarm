package users

import (
	"context"
	"errors"

	"github.com/fjarm/fjarm/api/internal/cache"
	"github.com/fjarm/fjarm/api/internal/cache/v1/pkg/remote"
)

// fakeRedisIdempotencyKeyStore implements IdempotencyKeyStore backed by a FakeRedisCache.
type fakeRedisIdempotencyKeyStore struct {
	cache *remote.FakeRedisCache
}

func newFakeRedisIdempotencyKeyStore(cache *remote.FakeRedisCache) *fakeRedisIdempotencyKeyStore {
	return &fakeRedisIdempotencyKeyStore{
		cache: cache,
	}
}

func (s *fakeRedisIdempotencyKeyStore) isDuplicate(ctx context.Context, key string) (bool, error) {
	_, err := s.cache.Get(ctx, key)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, cache.ErrCacheMiss) {
		return false, nil
	}
	return false, err
}

func (s *fakeRedisIdempotencyKeyStore) saveResult(ctx context.Context, key string) error {
	return s.cache.Set(ctx, key, []byte{}, idempotencyKeyTTL)
}
