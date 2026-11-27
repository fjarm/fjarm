package users

import (
	"context"

	"github.com/fjarm/fjarm/api/internal/authentication/v1/pkg/passwords"
)

// Compile-time check ensuring Argon2idPasswordHasher satisfies PasswordHasher.
var _ PasswordHasher = (*Argon2idPasswordHasher)(nil)

// Argon2idPasswordHasher implements PasswordHasher by delegating to a concurrency-limited passwords.Hasher.
type Argon2idPasswordHasher struct {
	hasher *passwords.Hasher
}

// NewArgon2idPasswordHasher creates a new Argon2idPasswordHasher backed by the supplied passwords.Hasher.
// If hasher is nil, a default passwords.NewHasher() instance is constructed.
func NewArgon2idPasswordHasher(hasher *passwords.Hasher) *Argon2idPasswordHasher {
	if hasher == nil {
		hasher = passwords.NewHasher()
	}
	return &Argon2idPasswordHasher{
		hasher: hasher,
	}
}

func (h *Argon2idPasswordHasher) hash(ctx context.Context, password string) (string, error) {
	return h.hasher.HashPassword(ctx, password)
}
