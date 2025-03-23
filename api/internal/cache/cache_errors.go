package cache

import (
	"fmt"
)

// ErrCacheMiss is returned when a key is not found in the cache.
var ErrCacheMiss = fmt.Errorf("cache miss")

// ErrInvalidExpiration is returned when the expiration time is invalid.
var ErrInvalidExpiration = fmt.Errorf("invalid expiration")

// ErrInvalidKey is returned when a key is invalid - i.e. empty strings or whitespace only strings.
var ErrInvalidKey = fmt.Errorf("invalid key")

// ErrLockReleaseFailed is returned when a process attempts to delete a key that does not belong to it. This indicates
// that either another process owns the key/lock or that the key/lock has already been deleted.
var ErrLockReleaseFailed = fmt.Errorf("failed to release lock")

// ErrKeyExists is returned when a key already exists in the cache.
var ErrKeyExists = fmt.Errorf("key exists")
