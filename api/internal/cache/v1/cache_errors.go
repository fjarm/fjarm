package v1

import (
	"fmt"
)

// ErrCacheMiss is returned when a key is not found in the cache.
var ErrCacheMiss = fmt.Errorf("cache miss")

// ErrKeyExists is returned when a key already exists in the cache.
var ErrKeyExists = fmt.Errorf("key exists")
