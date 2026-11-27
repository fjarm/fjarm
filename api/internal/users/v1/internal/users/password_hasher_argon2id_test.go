package users

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/fjarm/fjarm/api/internal/authentication/v1/pkg/passwords"
)

func TestArgon2idPasswordHasher_hash(t *testing.T) {
	sharedHasher := passwords.NewHasher()

	tests := map[string]struct {
		givenHasher      *passwords.Hasher
		givenPassword    string
		givenCanceledCtx bool
		wantErr          bool
		wantErrKind      error
	}{
		"valid_password_with_injected_hasher": {
			givenHasher:      sharedHasher,
			givenPassword:    "password",
			givenCanceledCtx: false,
			wantErr:          false,
		},
		"valid_password_with_nil_default_hasher": {
			givenHasher:      nil,
			givenPassword:    "password",
			givenCanceledCtx: false,
			wantErr:          false,
		},
		"canceled_context_returns_context_canceled": {
			givenHasher:      sharedHasher,
			givenPassword:    "password",
			givenCanceledCtx: true,
			wantErr:          true,
			wantErrKind:      context.Canceled,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Given
			adapter := NewArgon2idPasswordHasher(tc.givenHasher)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tc.givenCanceledCtx {
				cancel()
			}

			// When
			hashed, err := adapter.hash(ctx, tc.givenPassword)

			// Then
			if tc.wantErr {
				if err == nil {
					t.Fatal("hash expected an error, got nil")
				}
				if !errors.Is(err, tc.wantErrKind) {
					t.Errorf("hash got error %v, want %v", err, tc.wantErrKind)
				}
				if hashed != "" {
					t.Errorf("hash expected empty string on error, got %q", hashed)
				}
				return
			}

			if err != nil {
				t.Fatalf("hash got an unexpected error: %v", err)
			}
			if hashed == "" || hashed == tc.givenPassword {
				t.Fatalf("hash returned invalid digest %q", hashed)
			}
			if !strings.HasPrefix(hashed, "$argon2id$v=19$") {
				t.Errorf("hash expected $argon2id$v=19$ prefix, got: %q", hashed)
			}

			matched, verifyErr := sharedHasher.VerifyPassword(context.Background(), tc.givenPassword, hashed)
			if verifyErr != nil {
				t.Fatalf("VerifyPassword got an unexpected error: %v", verifyErr)
			}
			if !matched {
				t.Error("VerifyPassword expected hashed output to match plaintext password, got false")
			}
		})
	}
}
