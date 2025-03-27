package remote

import (
	"context"
	"errors"
	"github.com/fjarm/fjarm/api/internal/cache"
	"io"
	"log/slog"
	"testing"
	"time"
)

func TestRedisLock_SafeReleaseLock(t *testing.T) {
	// This test verifies that when the supplied lock key and value match what is contained in Redis, the safe release
	// script works correctly - i.e. the key is deleted and the returned error is nil.
	defer func() {
		err := rdb.Do(context.Background(), rdb.B().Flushall().Build()).Error()
		if err != nil {
			t.Errorf("failed to flush Redis cache: %v", err)
		}
	}()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	redisLock := NewRedisCache(rdb, logger)

	// Acquiring a free lock should succeed.
	lockKey := "testredislock:safereleaselock:lock"
	acquiredLock, err := redisLock.AcquireLock(context.Background(), lockKey, 10*time.Second)
	if err != nil || acquiredLock == "" {
		t.Errorf("failed to acquire Redis lock: %v", err)
	}

	// Attempting to acquire an existing lock will fail.
	failedLock, err := redisLock.AcquireLock(context.Background(), lockKey, 10*time.Second)
	if err == nil || !errors.Is(err, cache.ErrKeyExists) || failedLock != "" {
		t.Errorf("acquired lock that should not be accessible: %v", err)
	}

	// Releasing a lock we do not own should fail.
	err = redisLock.SafeReleaseLock(context.Background(), lockKey, "a lock we do not own")
	if err == nil {
		t.Errorf("incorrectly released an unowned lock: %v", err)
	}

	// Releasing a lock that we own should succeed.
	err = redisLock.SafeReleaseLock(context.Background(), lockKey, acquiredLock)
	if err != nil {
		t.Errorf("failed to release lock: %v", err)
	}

	// Verify that the released lock does not exist.
	existsCmd := redisLock.rdb.B().Exists().Key(lockKey).Build()
	exists, err := redisLock.rdb.Do(context.Background(), existsCmd).AsBool()
	if err != nil || exists {
		t.Errorf("failed to verify that lock was released: %v", err)
	}
}

func TestRedisLock_VerifyLock(t *testing.T) {
	defer func() {
		err := rdb.Do(context.Background(), rdb.B().Flushall().Build()).Error()
		if err != nil {
			t.Errorf("failed to flush Redis cache: %v", err)
		}
	}()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	redisLock := NewRedisCache(rdb, logger)

	// Acquiring a free lock should succeed.
	lockKey := "testredislock:verifylock:lock"
	acquiredLock, err := redisLock.AcquireLock(context.Background(), lockKey, 10*time.Second)
	if err != nil || acquiredLock == "" {
		t.Errorf("failed to acquire Redis lock: %v", err)
	}

	own, err := redisLock.VerifyLock(context.Background(), lockKey, acquiredLock)
	if err != nil || !own {
		t.Errorf("failed to verify that lock was acquired: %v", err)
	}

	own, err = redisLock.VerifyLock(context.Background(), lockKey, "a lock we do not own")
	if err == nil || !errors.Is(err, cache.ErrLockVerifyFailed) || own {
		t.Errorf("acquired lock that is not owned by current goroutine: %v", err)
	}

	err = redisLock.SafeReleaseLock(context.Background(), lockKey, acquiredLock)
	if err != nil {
		t.Errorf("failed to safely release lock: %v", err)
	}

	own, err = redisLock.VerifyLock(context.Background(), lockKey, acquiredLock)
	if err == nil || !errors.Is(err, cache.ErrCacheMiss) || own {
		t.Errorf("acquired lock that no longer exists: %v", err)
	}
}
