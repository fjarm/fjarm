package users

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/fjarm/fjarm/api/internal/cache/v1/pkg/remote"
)

func TestUserUseCase_createUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cache := remote.NewFakeRedisCache()
	repo := newUserStore()
	idemStore := newFakeRedisIdempotencyKeyStore(cache)
	distLocker := newFakeRedisDistributedLocker(cache)
	pwdHasher := newFakePasswordHasher()
	uc := NewUserUseCase(logger, idemStore, distLocker, repo, pwdHasher)

	tests := map[string]struct {
		idemKeys  []string
		emails    []string
		passwords []string
		hashErr   error
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
		"hashing_failure_returns_operation_failed": {
			idemKeys:  []string{"123e4567-e89b-12d3-a456-426614174999"},
			emails:    []string{"foo1@bar.com"},
			passwords: []string{"password"},
			hashErr:   ErrAuthenticationIssue,
			errs:      []bool{true},
			kind:      []error{ErrOperationFailed},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Cleanup(func() {
				cache.Reset()
				repo.reset()
				pwdHasher.reset()
			})
			pwdHasher.setErr(tc.hashErr)

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

					repo.mu.RLock()
					userCount := len(repo.users)
					for _, storedUser := range repo.users {
						if storedUser.Password == tc.passwords[i] {
							t.Errorf("createUser stored plaintext password %q; expected hashed output", storedUser.Password)
						}
						wantHash := fmt.Sprintf("fake_hashed:%s", tc.passwords[i])
						if storedUser.Password != wantHash {
							t.Errorf("createUser stored password = %q, want %q", storedUser.Password, wantHash)
						}
					}
					repo.mu.RUnlock()
					if userCount == 0 {
						t.Error("createUser expected user to be persisted in repository, got 0 users")
					}
				}
			}
		})
	}
}

func TestUserUseCase_createUser_HashingFailureSideEffectsAndLogs(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuf, nil))
	cache := remote.NewFakeRedisCache()
	repo := newUserStore()
	idemStore := newFakeRedisIdempotencyKeyStore(cache)
	distLocker := newFakeRedisDistributedLocker(cache)
	pwdHasher := newFakePasswordHasher()
	pwdHasher.setErr(errors.New("simulated hash failure"))
	uc := NewUserUseCase(logger, idemStore, distLocker, repo, pwdHasher)

	ctx := context.Background()
	idemKey := "123e4567-e89b-12d3-a456-426614174999"
	plaintextPwd := "super_secret_pwd_123"

	err := uc.createUser(ctx, idemKey, "alice@example.com", plaintextPwd)
	if err == nil {
		t.Fatal("createUser expected an error when password hashing fails, got nil")
	}
	if !errors.Is(err, ErrOperationFailed) {
		t.Errorf("createUser got error %v, want error wrapping %v", err, ErrOperationFailed)
	}
	if strings.Contains(err.Error(), "simulated hash failure") {
		t.Errorf("createUser leaked internal hashing error %q", err.Error())
	}

	repo.mu.RLock()
	userCount := len(repo.users)
	repo.mu.RUnlock()
	if userCount != 0 {
		t.Errorf("createUser persisted %d users on hashing failure, want 0", userCount)
	}

	idempotencyKey := fmt.Sprintf("%s:%s", createUserCacheKey, idemKey)
	dup, dupErr := idemStore.isDuplicate(ctx, idempotencyKey)
	if dupErr != nil {
		t.Fatalf("isDuplicate got an unexpected error: %v", dupErr)
	}
	if dup {
		t.Errorf("createUser cached idempotency key %q despite hashing failure", idempotencyKey)
	}

	lockKey := fmt.Sprintf("%s:%s", createUserLockKey, idemKey)
	lockVal, lockErr := distLocker.lock(ctx, lockKey)
	if lockErr != nil {
		t.Errorf("distLocker.lock failed after hashing failure (lock was not released): %v", lockErr)
	}
	if lockVal == "" {
		t.Error("distLocker.lock returned empty lockVal after hashing failure")
	}

	logOutput := logBuf.String()
	for _, expected := range []string{
		"ERROR",
		"msg=\"failed to hash credentials supplied in user message\"",
		"simulated hash failure",
	} {
		if !strings.Contains(logOutput, expected) {
			t.Errorf("createUser log output missing %q; got: %s", expected, logOutput)
		}
	}
	if strings.Contains(logOutput, plaintextPwd) {
		t.Errorf("createUser log output leaked plaintext password %q; got: %s", plaintextPwd, logOutput)
	}
}
