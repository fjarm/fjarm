package passwords

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// generateSalt creates a cryptographically secure random salt. The size parameter determines the number of bytes of
// random data. The size would typically be 16 bytes (128 bits) for a strong salt.
func generateSalt(size int) (string, error) {
	// Allocate buffer for the random bytes
	bytes := make([]byte, size)

	// Read random data from crypto/rand
	_, err := rand.Read(bytes)
	if err != nil {
		// The error is not recoverable and therefore isn't defined in common_errors.go. This should result in a panic
		// and connect.CodeUnknown being returned to the client.
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Encode as base64 for safe storage
	return base64.RawStdEncoding.EncodeToString(bytes), nil
}

func readSalt(salt string) ([]byte, error) {
	// Decode the base64 encoded salt
	return base64.RawStdEncoding.DecodeString(salt)
}
