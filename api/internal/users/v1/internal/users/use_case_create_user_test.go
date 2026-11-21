package users

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/fjarm/fjarm/api/internal/cache/v1/pkg/remote"
)

func TestUserUseCase_createUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cache := remote.NewFakeRedisCache()
	repo := newUserStore()
	idemStore := newFakeRedisIdempotencyKeyStore(cache)
	distLocker := newFakeRedisDistributedLocker(cache)
	uc := NewUserUseCase(logger, idemStore, distLocker, repo)

	tests := map[string]struct {
		idemKeys  []string
		emails    []string
		passwords []string
		errs      []bool
		kind      []error
	}{
		"validation_one_valid_user": {
			idemKeys:  []string{"123e4567-e89b-12d3-a456-426614174999"},
			emails:    []string{"foo1@bar.com"},
			passwords: []string{"password"},
			errs:      []bool{false},
			kind:      []error{nil},
		},
		"validation_one_nil_user": {
			idemKeys:  []string{"123e4567-e89b-12d3-a456-426614174999"},
			emails:    []string{""},
			passwords: []string{""},
			errs:      []bool{true},
			kind:      []error{ErrInvalidArgument},
		},
		"validation_one_no_password_user": {
			idemKeys:  []string{"123e4567-e89b-12d3-a456-426614174999"},
			emails:    []string{"foo1@bar.com"},
			passwords: []string{""},
			errs:      []bool{true},
			kind:      []error{ErrInvalidArgument},
		},
		"validation_one_no_idempotency_key_request": {
			idemKeys:  []string{""},
			emails:    []string{"foo1@bar.com"},
			passwords: []string{"password"},
			errs:      []bool{true},
			kind:      []error{ErrInvalidArgument},
		},
		"validation_one_invalid_idempotency_key_request": {
			idemKeys:  []string{"not-a-uuid"},
			emails:    []string{"foo1@bar.com"},
			passwords: []string{"password"},
			errs:      []bool{true},
			kind:      []error{ErrInvalidArgument},
		},
		"validation_one_no_email_user": {
			idemKeys:  []string{"123e4567-e89b-12d3-a456-426614174999"},
			emails:    []string{""},
			passwords: []string{"password"},
			errs:      []bool{true},
			kind:      []error{ErrInvalidArgument},
		},
		"validation_one_invalid_email_user": {
			idemKeys:  []string{"123e4567-e89b-12d3-a456-426614174999"},
			emails:    []string{"not-an-email"},
			passwords: []string{"password"},
			errs:      []bool{true},
			kind:      []error{ErrInvalidArgument},
		},
		"idempotency_two_distinct_valid_users": {
			idemKeys:  []string{"123e4567-e89b-12d3-a456-426614174999", "123e4567-e89b-12d3-a456-426614174888"},
			emails:    []string{"foo1@bar.com", "foo@bar.com"},
			passwords: []string{"password", "password"},
			errs:      []bool{false, false},
			kind:      []error{nil, nil},
		},
		"idempotency_two_identical_email_users": {
			idemKeys:  []string{"123e4567-e89b-12d3-a456-426614174999", "123e4567-e89b-12d3-a456-426614174888"},
			emails:    []string{"foo@bar.com", "foo@bar.com"},
			passwords: []string{"password", "password"},
			errs:      []bool{false, false},
			kind:      []error{nil, nil},
		},
		"idempotency_retry_identical_email_user_with_same_idempotency_key": {
			idemKeys: []string{
				"123e4567-e89b-12d3-a456-426614174999",
				"123e4567-e89b-12d3-a456-426614174888",
				"123e4567-e89b-12d3-a456-426614174888",
			},
			emails:    []string{"foo@bar.com", "foo@bar.com", "foo@bar.com"},
			passwords: []string{"password", "password", "password"},
			errs:      []bool{false, false, false},
			kind:      []error{nil, nil, nil},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for i := range tc.idemKeys {
				err := uc.createUser(context.Background(), tc.idemKeys[i], tc.emails[i], tc.passwords[i])
				if err != nil && !tc.errs[i] {
					t.Errorf("createUser got an unexpected error: %v", err)
				}
				if err == nil && tc.errs[i] {
					t.Errorf("createUser expected an error but got nil")
				}
				if !errors.Is(err, tc.kind[i]) {
					t.Errorf("createUser got an unexpected error type: %v", err)
				}
				if !tc.errs[i] {
					idempotencyKey := fmt.Sprintf("%s:%s", createUserCacheKey, tc.idemKeys[i])
					dup, dupErr := idemStore.isDuplicate(context.Background(), idempotencyKey)
					if dupErr != nil {
						t.Errorf("isDuplicate got an unexpected error: %v", dupErr)
					}
					if !dup {
						t.Errorf("createUser expected idempotency key %q to be cached, but it was not", idempotencyKey)
					}
				}
			}
		})
		// Reset the database for each test run.
		cache.Reset()
		repo.reset()
	}
}
