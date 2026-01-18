package passwords

import "fmt"

// ErrIncompatibleHashAlgorithmVersion is returned when the version of the hash algorithm (argon2id) used to hash the
// credential does not match the version specified in argon2.Version.
var ErrIncompatibleHashAlgorithmVersion = fmt.Errorf("incompatible hash algorithm")

// ErrInvalidHashAlgorithmVersionFormat is returned when the encoded hash algorithm version is not able to be parsed.
var ErrInvalidHashAlgorithmVersionFormat = fmt.Errorf("invalid hash algorithm version format")

// ErrInvalidHashFormat is returned when a hash is not in the expected 6-part, Argon2ID format. This is an unrecoverable
// error and should result in connect.CodeInternal being returned to the client.
var ErrInvalidHashFormat = fmt.Errorf("invalid hash format")

// ErrUnsupportedHashAlgorithm is returned when a hash is not using the Argon2ID algorithm. This is an unrecoverable
// error and should result in connect.CodeInternal being returned to the client.
var ErrUnsupportedHashAlgorithm = fmt.Errorf("unsupported hash algorithm")
