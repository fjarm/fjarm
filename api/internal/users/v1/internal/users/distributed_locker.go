package users

import (
	"context"

	"github.com/fjarm/fjarm/api/internal/cache/v1/pkg/remote"
)

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
	return l.locker.AcquireLock(ctx, key, idempotencyLockTTL)
}

func (l *fakeRedisDistributedLocker) unlock(ctx context.Context, key string, value string) error {
	return l.locker.SafeReleaseLock(ctx, key, value)
}

func (l *fakeRedisDistributedLocker) verifyLockAcquired(ctx context.Context, lockKey string, lockVal string) (bool, error) {
	return l.locker.VerifyLock(ctx, lockKey, lockVal)
}
