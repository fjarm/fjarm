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
			// Given
			cache := remote.NewFakeRedisCache()
			repo := newUserStore()
			idemStore := newFakeRedisIdempotencyKeyStore(cache)
			distLocker := newFakeRedisDistributedLocker(cache)
			pwdHasher := &fakePasswordHasher{err: tc.hashErr}
			uc := NewUserUseCase(logger, idemStore, distLocker, repo, pwdHasher)

			for i := range tc.idemKeys {
				// When
				err := uc.createUser(context.Background(), tc.idemKeys[i], tc.emails[i], tc.passwords[i])

				// Then
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

func TestUserUseCase_createUser_SideEffectsAndLogs(t *testing.T) {
	tests := map[string]struct {
		givenIdemKey       string
		givenEmail         string
		givenPassword      string
		givenHashErr       error
		wantErr            bool
		wantErrKind        error
		wantErrExcludes    string
		wantPersistedUsers int
		wantIdemKeyCached  bool
		wantLockReleased   bool
		wantLogContains    []string
		wantLogExcludes    []string
	}{
		"valid_user_persists_hash_caches_idempotency_key_and_releases_lock": {
			givenIdemKey:       "123e4567-e89b-12d3-a456-426614174999",
			givenEmail:         "alice@example.com",
			givenPassword:      "super_secret_pwd_123",
			givenHashErr:       nil,
			wantErr:            false,
			wantErrKind:        nil,
			wantPersistedUsers: 1,
			wantIdemKeyCached:  true,
			wantLockReleased:   true,
			wantLogContains:    []string{},
			wantLogExcludes:    []string{"super_secret_pwd_123"},
		},
		"hashing_failure_obscures_error_skips_persistence_and_cache_releases_lock_and_redacts_logs": {
			givenIdemKey:       "123e4567-e89b-12d3-a456-426614174999",
			givenEmail:         "alice@example.com",
			givenPassword:      "super_secret_pwd_123",
			givenHashErr:       errors.New("simulated hash failure"),
			wantErr:            true,
			wantErrKind:        ErrOperationFailed,
			wantErrExcludes:    "simulated hash failure",
			wantPersistedUsers: 0,
			wantIdemKeyCached:  false,
			wantLockReleased:   true,
			wantLogContains: []string{
				"ERROR",
				"msg=\"failed to hash credentials supplied in user message\"",
				"simulated hash failure",
			},
			wantLogExcludes: []string{"super_secret_pwd_123"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Given
			var logBuf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&logBuf, nil))
			cache := remote.NewFakeRedisCache()
			repo := newUserStore()
			idemStore := newFakeRedisIdempotencyKeyStore(cache)
			distLocker := newFakeRedisDistributedLocker(cache)
			pwdHasher := &fakePasswordHasher{err: tc.givenHashErr}
			uc := NewUserUseCase(logger, idemStore, distLocker, repo, pwdHasher)
			ctx := context.Background()

			// When
			err := uc.createUser(ctx, tc.givenIdemKey, tc.givenEmail, tc.givenPassword)

			// Then
			if tc.wantErr && err == nil {
				t.Fatal("createUser expected an error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("createUser got an unexpected error: %v", err)
			}
			if !errors.Is(err, tc.wantErrKind) {
				t.Errorf("createUser got error %v, want error wrapping %v", err, tc.wantErrKind)
			}
			if tc.wantErrExcludes != "" && err != nil && strings.Contains(err.Error(), tc.wantErrExcludes) {
				t.Errorf("createUser leaked internal error %q in %q", tc.wantErrExcludes, err.Error())
			}

			repo.mu.RLock()
			userCount := len(repo.users)
			for _, storedUser := range repo.users {
				wantHash := fmt.Sprintf("fake_hashed:%s", tc.givenPassword)
				if storedUser.Password != wantHash {
					t.Errorf("createUser stored password = %q, want %q", storedUser.Password, wantHash)
				}
			}
			repo.mu.RUnlock()
			if userCount != tc.wantPersistedUsers {
				t.Errorf("createUser persisted %d users, want %d", userCount, tc.wantPersistedUsers)
			}

			idempotencyKey := fmt.Sprintf("%s:%s", createUserCacheKey, tc.givenIdemKey)
			dup, dupErr := idemStore.isDuplicate(ctx, idempotencyKey)
			if dupErr != nil {
				t.Fatalf("isDuplicate got an unexpected error: %v", dupErr)
			}
			if dup != tc.wantIdemKeyCached {
				t.Errorf("isDuplicate = %v, want %v", dup, tc.wantIdemKeyCached)
			}

			if tc.wantLockReleased {
				lockKey := fmt.Sprintf("%s:%s", createUserLockKey, tc.givenIdemKey)
				lockVal, lockErr := distLocker.lock(ctx, lockKey)
				if lockErr != nil {
					t.Errorf("distLocker.lock failed after createUser (lock was not released): %v", lockErr)
				}
				if lockVal == "" {
					t.Error("distLocker.lock returned empty lockVal after createUser")
				}
			}

			logOutput := logBuf.String()
			for _, expected := range tc.wantLogContains {
				if !strings.Contains(logOutput, expected) {
					t.Errorf("createUser log output missing %q; got: %s", expected, logOutput)
				}
			}
			for _, excluded := range tc.wantLogExcludes {
				if strings.Contains(logOutput, excluded) {
					t.Errorf("createUser log output leaked %q; got: %s", excluded, logOutput)
				}
			}
		})
	}
}
