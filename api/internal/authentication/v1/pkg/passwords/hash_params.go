package passwords

// hashParams defines the parameters used by the Argon2id algorithm.
type hashParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  int
	keyLength   uint32
}

// DefaultParams provides recommended parameters for Argon2id. These values should be adjusted based on the specific
// security requirements.
func defaultParams() *hashParams {
	return &hashParams{
		memory:      64 * 1024, // 64MB
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}
