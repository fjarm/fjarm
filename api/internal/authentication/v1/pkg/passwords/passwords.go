package passwords

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// HashPassword creates a new hash of a plain-text password using Argon2id.
func HashPassword(password string) (string, error) {
	params := DefaultParams()

	saltString, err := generateSalt(params.SaltLength)
	if err != nil {
		return "", err
	}

	saltBytes, err := readSalt(saltString)
	if err != nil {
		return "", err
	}

	// Hash the password using Argon2id
	hash := argon2.IDKey(
		[]byte(password),
		saltBytes,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)
	hashString := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=memory,t=iterations,p=parallelism$salt$hash
	encodedHash := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		params.Memory,
		params.Iterations,
		params.Parallelism,
		saltString,
		hashString,
	)

	return encodedHash, nil
}

// VerifyPassword checks if a supplied password matches a generated hash. Return `1` if the password matches the hash,
// and `0` if it does not.
func VerifyPassword(password string, encodedHash string) (bool, error) {
	// Parse the parameters, salt, and hash from the encoded string
	params, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Hash the password with the same parameters and salt
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Compare the hashes in constant time to prevent timing attacks
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

// decodeHash parses an encoded hash string into its components - parameters, salt, and hash.
func decodeHash(encodedHash string) (*HashParams, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHashFormat
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, ErrUnsupportedHashAlgorithm
	}

	params := HashParams{}
	_, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&params.Memory,
		&params.Iterations,
		&params.Parallelism,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := readSalt(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = len(salt)

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(hash))

	return &params, salt, hash, nil
}
