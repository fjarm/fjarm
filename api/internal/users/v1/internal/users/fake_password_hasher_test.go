package users

import (
	"context"
	"fmt"
	"sync"
)

// Compile-time check ensuring fakePasswordHasher satisfies PasswordHasher.
var _ PasswordHasher = (*fakePasswordHasher)(nil)

// fakePasswordHasher is a fast, deterministic, thread-safe test double for PasswordHasher.
// It avoids Argon2id CPU and 64 MiB memory overhead in unit tests and supports injecting errors.
type fakePasswordHasher struct {
	mu  sync.RWMutex
	err error
}

func newFakePasswordHasher() *fakePasswordHasher {
	return &fakePasswordHasher{}
}

func (f *fakePasswordHasher) setErr(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.err = err
}

func (f *fakePasswordHasher) reset() {
	f.setErr(nil)
}

func (f *fakePasswordHasher) hash(ctx context.Context, password string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	f.mu.RLock()
	err := f.err
	f.mu.RUnlock()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("fake_hashed:%s", password), nil
}
