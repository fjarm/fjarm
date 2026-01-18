package passwords

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// delimiter describes the character used to delimit/indicate different parts of an encoded credential.
const delimiter = "$"

// decodedHash describes a hashed and salted credential as well as the parameters used to encrypt the credential.
type decodedHash struct {
	params *hashParams
	salt   []byte
	hash   []byte
}

// HashPassword creates a new hash of a plain-text password using Argon2id.
func HashPassword(password string) (string, error) {
	params := defaultParams()

	salt, err := generateSalt(params.saltLength)
	if err != nil {
		return "", err
	}

	// Hash the password using Argon2id
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength,
	)

	// Convert the hash and salt to base64 encoded strings.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=memory,t=iterations,p=parallelism$salt$hash.
	// The format represents the algorithm identifier, algorithm version, hashing parameters (memory cost,
	// iterations/time cost, parallelism/number of threads), salt, and hash.
	encodedHash := strings.Join([]string{
		"",         // Leading delimiter
		"argon2id", // Hashing algorithm
		"v=19",     // Hashing algorithm version
		fmt.Sprintf("m=%d,t=%d,p=%d", params.memory, params.iterations, params.parallelism),
		b64Salt,
		b64Hash,
	}, delimiter)

	return encodedHash, nil
}

// VerifyPassword checks if a supplied password matches a generated hash. Return `1` if the password matches the hash,
// and `0` if it does not.
func VerifyPassword(password string, encodedCredential string) (bool, error) {
	// Parse the parameters, salt, and hash from the encoded string
	creds, err := decodeHash(encodedCredential)
	if err != nil {
		return false, err
	}

	// Hash the password with the same parameters and salt
	otherHash := argon2.IDKey(
		[]byte(password),
		creds.salt,
		creds.params.iterations,
		creds.params.memory,
		creds.params.parallelism,
		creds.params.keyLength,
	)

	// Compare the hashes in constant time to prevent timing attacks
	return subtle.ConstantTimeCompare(creds.hash, otherHash) == 1, nil
}

// decodeHash parses an encoded hash string into its components - parameters, salt, and hash.
func decodeHash(encodedHash string) (*decodedHash, error) {
	parts := strings.Split(encodedHash, delimiter)
	if len(parts) != 6 {
		return nil, ErrInvalidHashFormat
	}

	// Verify the algorithm used to hash the original credential is argon2id.
	if parts[1] != "argon2id" {
		return nil, ErrUnsupportedHashAlgorithm
	}

	// Verify the version of argon2id that was used to hash the original credential is the one imported in argon2.
	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, ErrInvalidHashAlgorithmVersionFormat
	}
	if version != argon2.Version {
		return nil, fmt.Errorf("%w: got %d, expected %d", ErrIncompatibleHashAlgorithmVersion, version, argon2.Version)
	}

	// Parse the hashing parameters.
	params := hashParams{}
	_, err = fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&params.memory,
		&params.iterations,
		&params.parallelism,
	)
	if err != nil {
		return nil, err
	}

	salt, err := readSalt(parts[4])
	if err != nil {
		return nil, err
	}
	params.saltLength = len(salt)

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, err
	}
	params.keyLength = uint32(len(hash))

	return &decodedHash{
		params: &params,
		salt:   salt,
		hash:   hash,
	}, nil
}
