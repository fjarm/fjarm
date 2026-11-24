package users

import (
	"context"
	"strings"
	"testing"

	"github.com/fjarm/fjarm/api/internal/authentication/v1/pkg/passwords"
)

func TestArgon2idPasswordHasher_hash(t *testing.T) {
	hasher := NewArgon2idPasswordHasher()
	plaintext := "password"

	hashed, err := hasher.hash(context.Background(), plaintext)
	if err != nil {
		t.Fatalf("hash got an unexpected error: %v", err)
	}
	if hashed == "" || hashed == plaintext {
		t.Fatalf("hash returned invalid digest %q", hashed)
	}
	if !strings.HasPrefix(hashed, "$argon2id$v=19$") {
		t.Errorf("hash expected $argon2id$v=19$ prefix, got: %q", hashed)
	}

	matched, verifyErr := passwords.VerifyPassword(plaintext, hashed)
	if verifyErr != nil {
		t.Fatalf("VerifyPassword got an unexpected error: %v", verifyErr)
	}
	if !matched {
		t.Error("VerifyPassword expected hashed output to match plaintext password, got false")
	}
}
