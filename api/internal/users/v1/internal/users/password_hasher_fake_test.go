package users

import (
	"context"
	"fmt"
)

// Compile-time check ensuring fakePasswordHasher satisfies PasswordHasher.
var _ PasswordHasher = (*fakePasswordHasher)(nil)

// fakePasswordHasher is a fast, deterministic test double for PasswordHasher.
// It avoids Argon2id CPU and 64 MiB memory overhead in unit tests and supports configuring an error at construction.
type fakePasswordHasher struct {
	err error
}

func newFakePasswordHasher() *fakePasswordHasher {
	return &fakePasswordHasher{}
}

func (f *fakePasswordHasher) hash(ctx context.Context, password string) (string, error) {
	ctxErr := ctx.Err()
	if ctxErr != nil {
		return "", ctxErr
	}
	if f.err != nil {
		return "", f.err
	}
	return fmt.Sprintf("fake_hashed:%s", password), nil
}
