package passwords

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"runtime"
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

// Hasher bounds concurrent Argon2id password hashing and verification operations using an in-memory semaphore
// so traffic spikes cannot exhaust host memory or starve CPU cores.
type Hasher struct {
	semaphore chan struct{}
	params    *hashParams
}

func defaultMaxConcurrency(params *hashParams) int {
	parallelism := max(1, int(params.parallelism))
	return max(1, runtime.GOMAXPROCS(0)/parallelism)
}

// NewHasher creates a Hasher with default Argon2id parameters whose concurrency limit defaults to
// max(1, GOMAXPROCS / parallelism).
func NewHasher() *Hasher {
	params := defaultParams()
	maxConcurrency := defaultMaxConcurrency(params)
	return &Hasher{
		semaphore: make(chan struct{}, maxConcurrency),
		params:    params,
	}
}

// NewHasherWithMaxConcurrency creates a Hasher that permits at most maxConcurrency concurrent Argon2id
// operations across HashPassword and VerifyPassword. Returns ErrInvalidMaxConcurrency if maxConcurrency < 1.
func NewHasherWithMaxConcurrency(maxConcurrency int) (*Hasher, error) {
	if maxConcurrency < 1 {
		return nil, fmt.Errorf("%w: got %d, must be >= 1", ErrInvalidMaxConcurrency, maxConcurrency)
	}
	return &Hasher{
		semaphore: make(chan struct{}, maxConcurrency),
		params:    defaultParams(),
	}, nil
}

// MaxConcurrency returns the maximum number of Argon2id operations allowed to execute concurrently.
func (h *Hasher) MaxConcurrency() int {
	return cap(h.semaphore)
}

func (h *Hasher) derive(ctx context.Context, password []byte, salt []byte, params *hashParams) ([]byte, error) {
	err := ctx.Err()
	if err != nil {
		return nil, err
	}

	select {
	case h.semaphore <- struct{}{}:
		defer func() { <-h.semaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	hash := argon2.IDKey(
		password,
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength,
	)
	return hash, nil
}

// HashPassword creates a new Argon2id hash of a plain-text password, blocking until a concurrency slot is
// available or ctx is canceled.
func (h *Hasher) HashPassword(ctx context.Context, password string) (string, error) {
	err := ctx.Err()
	if err != nil {
		return "", err
	}

	salt, err := generateSalt(h.params.saltLength)
	if err != nil {
		return "", err
	}

	hash, err := h.derive(ctx, []byte(password), salt, h.params)
	if err != nil {
		return "", err
	}

	// Convert the hash and salt to base64 encoded strings.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=memory,t=iterations,p=parallelism$salt$hash.
	// The format represents the algorithm identifier, algorithm version, hashing parameters (memory cost,
	// iterations/time cost, parallelism/number of threads), salt, and hash.
	encodedHash := strings.Join([]string{
		"",                                  // Leading delimiter
		"argon2id",                          // Hashing algorithm
		fmt.Sprintf("v=%d", argon2.Version), // Hashing algorithm version (should be 19)
		fmt.Sprintf("m=%d,t=%d,p=%d", h.params.memory, h.params.iterations, h.params.parallelism),
		b64Salt,
		b64Hash,
	}, delimiter)

	return encodedHash, nil
}

// VerifyPassword checks in constant time whether password matches encodedCredential, blocking until a
// concurrency slot is available or ctx is canceled.
func (h *Hasher) VerifyPassword(ctx context.Context, password string, encodedCredential string) (bool, error) {
	err := ctx.Err()
	if err != nil {
		return false, err
	}

	// Parse and validate the encoded hash format before waiting on an Argon2id concurrency slot.
	creds, err := decodeHash(encodedCredential)
	if err != nil {
		return false, err
	}

	otherHash, err := h.derive(ctx, []byte(password), creds.salt, creds.params)
	if err != nil {
		return false, err
	}

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
