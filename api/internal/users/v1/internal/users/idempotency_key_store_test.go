package users

import (
	"context"
	"testing"

	"github.com/fjarm/fjarm/api/internal/cache/v1/pkg/remote"
)

func TestFakeRedisIdempotencyKeyStore_isDuplicate(t *testing.T) {
	tests := map[string]struct {
		setupKey    string
		queryKey    string
		shouldExist bool
		dup         bool
		err         bool
	}{
		"non_existent_key": {
			queryKey:    "userservice:create:idempotency:user:123",
			shouldExist: false,
			dup:         false,
			err:         false,
		},
		"existing_key": {
			setupKey:    "userservice:create:idempotency:user:456",
			queryKey:    "userservice:create:idempotency:user:456",
			shouldExist: true,
			dup:         true,
			err:         false,
		},
		"invalid_empty_key": {
			queryKey:    "",
			shouldExist: false,
			dup:         false,
			err:         true,
		},
		"invalid_whitespace_key": {
			queryKey:    "   ",
			shouldExist: false,
			dup:         false,
			err:         true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cache := remote.NewFakeRedisCache()
			store := newFakeRedisIdempotencyKeyStore(cache)
			ctx := context.Background()

			if tc.shouldExist {
				if err := store.saveResult(ctx, tc.setupKey); err != nil {
					t.Fatalf("failed to pre-populate key: %v", err)
				}
			}

			gotDup, err := store.isDuplicate(ctx, tc.queryKey)
			if err != nil && !tc.err {
				t.Errorf("isDuplicate got an unexpected error: %v", err)
			}
			if err == nil && tc.err {
				t.Errorf("isDuplicate expected an error but got nil")
			}
			if gotDup != tc.dup {
				t.Errorf("isDuplicate got %v, want %v", gotDup, tc.dup)
			}
		})
	}
}

func TestFakeRedisIdempotencyKeyStore_saveResult(t *testing.T) {
	tests := map[string]struct {
		key string
		err bool
	}{
		"valid_key": {
			key: "userservice:create:idempotency:user:123",
			err: false,
		},
		"invalid_empty_key": {
			key: "",
			err: true,
		},
		"invalid_whitespace_key": {
			key: "   ",
			err: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cache := remote.NewFakeRedisCache()
			store := newFakeRedisIdempotencyKeyStore(cache)
			ctx := context.Background()

			err := store.saveResult(ctx, tc.key)
			if err != nil && !tc.err {
				t.Errorf("saveResult got an unexpected error: %v", err)
			}
			if err == nil && tc.err {
				t.Errorf("saveResult expected an error but got nil")
			}
			if !tc.err {
				gotDup, dupErr := store.isDuplicate(ctx, tc.key)
				if dupErr != nil {
					t.Errorf("isDuplicate got an unexpected error: %v", dupErr)
				}
				if !gotDup {
					t.Errorf("isDuplicate got %v, want %v", gotDup, true)
				}
			}
		})
	}
}
