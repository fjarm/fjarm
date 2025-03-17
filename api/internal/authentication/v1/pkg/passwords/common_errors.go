package passwords

import "fmt"

// ErrInvalidHashFormat is returned when a hash is not in the expected 6-part, Argon2ID format. This is an unrecoverable
// error and should result in connect.CodeInternal being returned to the client.
var ErrInvalidHashFormat = fmt.Errorf("invalid hash format")

// ErrUnsupportedHashAlgorithm is returned when a hash is not using the Argon2ID algorithm. This is an unrecoverable
// error and should result in connect.CodeInternal being returned to the client.
var ErrUnsupportedHashAlgorithm = fmt.Errorf("unsupported hash algorithm")
