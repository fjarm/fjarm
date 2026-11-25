package users

import (
	"context"

	"github.com/fjarm/fjarm/api/internal/authentication/v1/pkg/passwords"
)

// Compile-time check ensuring Argon2idPasswordHasher satisfies PasswordHasher.
var _ PasswordHasher = (*Argon2idPasswordHasher)(nil)

// Argon2idPasswordHasher implements PasswordHasher using the production Argon2id algorithm and parameters.
type Argon2idPasswordHasher struct{}

// NewArgon2idPasswordHasher creates a new Argon2idPasswordHasher instance for production wiring.
func NewArgon2idPasswordHasher() *Argon2idPasswordHasher {
	return &Argon2idPasswordHasher{}
}

func (h *Argon2idPasswordHasher) hash(_ context.Context, password string) (string, error) {
	return passwords.HashPassword(password)
}
