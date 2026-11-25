package users

import (
	"context"
	"strings"
	"testing"

	"github.com/fjarm/fjarm/api/internal/authentication/v1/pkg/passwords"
)

func TestArgon2idPasswordHasher_hash(t *testing.T) {
	tests := map[string]struct {
		givenPassword string
		wantPrefix    string
		wantErr       bool
	}{
		"valid_alphabetic_password": {
			givenPassword: "password",
			wantPrefix:    "$argon2id$v=19$",
			wantErr:       false,
		},
		"valid_mixed_character_password": {
			givenPassword: "$$%%$abcd$efgh$%%%",
			wantPrefix:    "$argon2id$v=19$",
			wantErr:       false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Given
			hasher := NewArgon2idPasswordHasher()

			// When
			hashed, err := hasher.hash(context.Background(), tc.givenPassword)

			// Then
			if tc.wantErr && err == nil {
				t.Fatal("hash expected an error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("hash got an unexpected error: %v", err)
			}
			if hashed == "" || hashed == tc.givenPassword {
				t.Fatalf("hash returned invalid digest %q", hashed)
			}
			if !strings.HasPrefix(hashed, tc.wantPrefix) {
				t.Errorf("hash expected %s prefix, got: %q", tc.wantPrefix, hashed)
			}

			matched, verifyErr := passwords.VerifyPassword(tc.givenPassword, hashed)
			if verifyErr != nil {
				t.Fatalf("VerifyPassword got an unexpected error: %v", verifyErr)
			}
			if !matched {
				t.Error("VerifyPassword expected hashed output to match plaintext password, got false")
			}
		})
	}
}
