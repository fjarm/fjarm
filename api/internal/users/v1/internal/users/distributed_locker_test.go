package users

import (
	"context"
	"errors"
	"testing"

	"github.com/fjarm/fjarm/api/internal/cache"
	"github.com/fjarm/fjarm/api/internal/cache/v1/pkg/remote"
)

func TestFakeRedisDistributedLocker_lock(t *testing.T) {
	tests := map[string]struct {
		setupKey string
		lockKey  string
		err      bool
		kind     error
	}{
		"valid_free_lock": {
			lockKey: "userservice:create:lock:user:1",
			err:     false,
			kind:    nil,
		},
		"already_locked_key": {
			setupKey: "userservice:create:lock:user:2",
			lockKey:  "userservice:create:lock:user:2",
			err:      true,
			kind:     cache.ErrKeyExists,
		},
		"invalid_empty_key": {
			lockKey: "",
			err:     true,
			kind:    cache.ErrInvalidKey,
		},
		"invalid_whitespace_key": {
			lockKey: "   ",
			err:     true,
			kind:    cache.ErrInvalidKey,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cacheInst := remote.NewFakeRedisCache()
			locker := newFakeRedisDistributedLocker(cacheInst)
			ctx := context.Background()

			if tc.setupKey != "" {
				if _, err := locker.lock(ctx, tc.setupKey); err != nil {
					t.Fatalf("failed to pre-acquire lock: %v", err)
				}
			}

			lockVal, err := locker.lock(ctx, tc.lockKey)
			if err != nil && !tc.err {
				t.Errorf("lock got an unexpected error: %v", err)
			}
			if err == nil && tc.err {
				t.Errorf("lock expected an error but got nil")
			}
			if tc.kind != nil && !errors.Is(err, tc.kind) {
				t.Errorf("lock got an unexpected error kind: %v, wanted: %v", err, tc.kind)
			}
			if !tc.err && lockVal == "" {
				t.Errorf("lock got empty lockVal, want non-empty")
			}
		})
	}
}

func TestFakeRedisDistributedLocker_verifyLockAcquired(t *testing.T) {
	tests := map[string]struct {
		acquireLock bool
		lockKey     string
		useWrongVal bool
		verified    bool
		err         bool
		kind        error
	}{
		"valid_matching_lock": {
			acquireLock: true,
			lockKey:     "userservice:create:lock:user:1",
			useWrongVal: false,
			verified:    true,
			err:         false,
			kind:        nil,
		},
		"unowned_mismatched_lock": {
			acquireLock: true,
			lockKey:     "userservice:create:lock:user:2",
			useWrongVal: true,
			verified:    false,
			err:         false,
			kind:        nil,
		},
		"non_existent_lock": {
			acquireLock: false,
			lockKey:     "userservice:create:lock:user:3",
			useWrongVal: false,
			verified:    false,
			err:         true,
			kind:        cache.ErrCacheMiss,
		},
		"invalid_empty_key": {
			acquireLock: false,
			lockKey:     "",
			useWrongVal: false,
			verified:    false,
			err:         true,
			kind:        cache.ErrInvalidKey,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cacheInst := remote.NewFakeRedisCache()
			locker := newFakeRedisDistributedLocker(cacheInst)
			ctx := context.Background()

			var lockVal string
			if tc.acquireLock {
				var err error
				lockVal, err = locker.lock(ctx, tc.lockKey)
				if err != nil {
					t.Fatalf("failed to acquire lock: %v", err)
				}
			}

			verifyVal := lockVal
			if tc.useWrongVal {
				verifyVal = "wrong-lock-value"
			}

			verified, err := locker.verifyLockAcquired(ctx, tc.lockKey, verifyVal)
			if err != nil && !tc.err {
				t.Errorf("verifyLockAcquired got an unexpected error: %v", err)
			}
			if err == nil && tc.err {
				t.Errorf("verifyLockAcquired expected an error but got nil")
			}
			if tc.kind != nil && !errors.Is(err, tc.kind) {
				t.Errorf("verifyLockAcquired got an unexpected error kind: %v, wanted: %v", err, tc.kind)
			}
			if verified != tc.verified {
				t.Errorf("verifyLockAcquired got %v, want %v", verified, tc.verified)
			}
		})
	}
}

func TestFakeRedisDistributedLocker_unlock(t *testing.T) {
	tests := map[string]struct {
		acquireLock bool
		lockKey     string
		useWrongVal bool
		err         bool
		kind        error
	}{
		"valid_unlock": {
			acquireLock: true,
			lockKey:     "userservice:create:lock:user:1",
			useWrongVal: false,
			err:         false,
			kind:        nil,
		},
		"unlock_with_mismatched_value": {
			acquireLock: true,
			lockKey:     "userservice:create:lock:user:2",
			useWrongVal: true,
			err:         true,
			kind:        cache.ErrLockReleaseFailed,
		},
		"unlock_non_existent_lock": {
			acquireLock: false,
			lockKey:     "userservice:create:lock:user:3",
			useWrongVal: false,
			err:         true,
			kind:        cache.ErrCacheMiss,
		},
		"invalid_empty_key": {
			acquireLock: false,
			lockKey:     "",
			useWrongVal: false,
			err:         true,
			kind:        cache.ErrInvalidKey,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cacheInst := remote.NewFakeRedisCache()
			locker := newFakeRedisDistributedLocker(cacheInst)
			ctx := context.Background()

			var lockVal string
			if tc.acquireLock {
				var err error
				lockVal, err = locker.lock(ctx, tc.lockKey)
				if err != nil {
					t.Fatalf("failed to acquire lock: %v", err)
				}
			}

			unlockVal := lockVal
			if tc.useWrongVal {
				unlockVal = "wrong-lock-value"
			}

			err := locker.unlock(ctx, tc.lockKey, unlockVal)
			if err != nil && !tc.err {
				t.Errorf("unlock got an unexpected error: %v", err)
			}
			if err == nil && tc.err {
				t.Errorf("unlock expected an error but got nil")
			}
			if tc.kind != nil && !errors.Is(err, tc.kind) {
				t.Errorf("unlock got an unexpected error kind: %v, wanted: %v", err, tc.kind)
			}
			if !tc.err && tc.acquireLock {
				// Verify lock can be re-acquired after successful unlock
				newLockVal, relockErr := locker.lock(ctx, tc.lockKey)
				if relockErr != nil {
					t.Errorf("lock got an unexpected error: %v", relockErr)
				}
				if newLockVal == "" {
					t.Errorf("lock got empty lockVal, want non-empty")
				}
			}
		})
	}
}
