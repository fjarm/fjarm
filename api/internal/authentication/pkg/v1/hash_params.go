package v1

// HashParams defines the parameters used by the Argon2id algorithm.
type HashParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  int
	KeyLength   uint32
}

// DefaultParams provides recommended parameters for Argon2id. These values should be adjusted based on the specific
// security requirements.
func DefaultParams() *HashParams {
	return &HashParams{
		Memory:      64 * 1024, // 64MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}
