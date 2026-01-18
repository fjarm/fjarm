package remote

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	cachev1 "github.com/fjarm/fjarm/api/internal/cache"
	"github.com/fjarm/fjarm/api/internal/logkeys"
)

// AcquireLock attempts to "acquire a lock" - that is, write a key in Redis using SET NX (if the write is successful, we
// have the lock, otherwise another process has the lock). If another process already owns the lock, an empty string and
// an error are returned. Otherwise, the lock (the value associated with the key) is returned with no error.
func (c *RedisCache) AcquireLock(ctx context.Context, key string, ttl time.Duration) (string, error) {
	if strings.TrimSpace(key) == "" || strings.Contains(key, " ") {
		return "", fmt.Errorf("%w: %s", cachev1.ErrInvalidKey, "key cannot be empty or whitespace")
	}
	if ttl <= 0 {
		return "", fmt.Errorf("%w: %s", cachev1.ErrInvalidExpiration, "ttl must be greater than 0")
	}

	lockVal := uuid.NewString()
	err := c.Set(ctx, key, []byte(lockVal), ttl)
	if err != nil {
		c.logger.WarnContext(ctx, "failed to acquire lock in Redis cache", slog.Any(logkeys.Err, err))
		return "", err
	}
	return lockVal, nil
}

// SafeReleaseLock attempts to release a lock in Reds - that is, delete a key only if it belongs to the calling process.
// The owning process is identified by the value passed in. If the key doesn't exist or if the value doesn't match the
// value of the key in Redis, an error is returned.
//
// Using GET followed by DEL is not atomic (a process's lock could've expired between each call), so a Lua script is
// used to ensure atomicity.
func (c *RedisCache) SafeReleaseLock(ctx context.Context, key string, value string) error {
	if strings.TrimSpace(key) == "" || strings.Contains(key, " ") {
		return fmt.Errorf("%w: %s", cachev1.ErrInvalidKey, "key cannot be empty or whitespace")
	}

	cmd := c.rdb.B().Eval().Script(safeDeleteScript).
		Numkeys(1).
		Key(key).
		Arg(value).
		Build()
	result, err := c.rdb.Do(ctx, cmd).AsInt64()
	if err != nil {
		// This is a real Redis error. Do not wrap it and let the cache client handle it accordingly.
		c.logger.WarnContext(ctx, "failed to safe release lock in Redis cache", slog.Any(logkeys.Err, err))
		return err
	}

	// Check the script result (1 = success, 0 = lock not owned by us)
	if result != 1 {
		c.logger.WarnContext(
			ctx,
			"encountered already released lock or lock not owned by us",
			slog.String(logkeys.Key, key),
		)
		return fmt.Errorf("%w: %s", cachev1.ErrLockReleaseFailed, "failed to release lock")
	}
	return nil
}

// VerifyLock checks if the calling process owns the lock in Redis. If the key doesn't exist or if the value doesn't
// match the value of the key in Redis, an error is returned.
func (c *RedisCache) VerifyLock(ctx context.Context, lockKey string, lockVal string) (bool, error) {
	if strings.TrimSpace(lockKey) == "" || strings.Contains(lockKey, " ") {
		return false, fmt.Errorf("%w: %s", cachev1.ErrInvalidKey, "key cannot be empty or whitespace")
	}

	val, err := c.Get(ctx, lockKey)
	if err != nil {
		// This could be a real Redis error or a cache miss because the key doesn't exist. If it's a cache miss, the
		// lock key could have been deleted or expired. In either case, the client may choose to abort the operation.
		return false, err
	}
	if string(val) == lockVal {
		return true, nil
	}
	return false, cachev1.ErrLockVerifyFailed
}
