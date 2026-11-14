package users

import (
	"context"
	"errors"
	"time"

	"github.com/fjarm/fjarm/api/internal/cache"
	"github.com/fjarm/fjarm/api/internal/cache/v1/pkg/remote"
)

const idempotencyLockTTL = 30 * time.Second // Short-lived to coordinate concurrent requests

// fakeRedisDistributedLocker implements DistributedLocker backed by a FakeRedisCache.
type fakeRedisDistributedLocker struct {
	locker *remote.FakeRedisCache
}

func newFakeRedisDistributedLocker(locker *remote.FakeRedisCache) *fakeRedisDistributedLocker {
	return &fakeRedisDistributedLocker{
		locker: locker,
	}
}

func (l *fakeRedisDistributedLocker) lock(ctx context.Context, key string) (string, error) {
	val, err := l.locker.AcquireLock(ctx, key, idempotencyLockTTL)
	if errors.Is(err, cache.ErrKeyExists) {
		return "", ErrLockAlreadyHeld
	}
	return val, err
}

func (l *fakeRedisDistributedLocker) unlock(ctx context.Context, key string, value string) error {
	return l.locker.SafeReleaseLock(ctx, key, value)
}

func (l *fakeRedisDistributedLocker) verifyLockAcquired(ctx context.Context, lockKey string, lockVal string) (bool, error) {
	return l.locker.VerifyLock(ctx, lockKey, lockVal)
}
