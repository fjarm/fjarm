package remote

import (
	"context"
	"fmt"
	cachev1 "github.com/fjarm/fjarm/api/internal/cache"
	"github.com/google/uuid"
	"strings"
	"time"
)

func (c *FakeRedisCache) AcquireLock(ctx context.Context, key string, ttl time.Duration) (string, error) {
	c.removeExpiredValues()
	lockVal := uuid.NewString()
	err := c.Set(ctx, key, []byte(lockVal), ttl)
	if err != nil {
		return "", err
	}
	return lockVal, nil
}

func (c *FakeRedisCache) SafeReleaseLock(ctx context.Context, key string, value string) error {
	if strings.TrimSpace(key) == "" || strings.Contains(key, " ") {
		return fmt.Errorf("%w: %s", cachev1.ErrInvalidKey, "key cannot be empty or whitespace")
	}
	val, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	if string(val) != value {
		return fmt.Errorf("%w: %s", cachev1.ErrLockReleaseFailed, "failed to release lock. value does not match")
	}
	c.removeExpiredValues()
	delete(c.rdb, key)
	return nil
}

func (c *FakeRedisCache) VerifyLock(ctx context.Context, lockKey string, lockVal string) (bool, error) {
	c.removeExpiredValues()
	val, err := c.Get(ctx, lockKey)
	if err != nil {
		return false, err
	}
	if string(val) == lockVal {
		return true, nil
	}
	return false, nil
}
