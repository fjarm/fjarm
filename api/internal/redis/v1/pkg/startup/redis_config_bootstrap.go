package startup

// redisServerConfig represents the tunable parameters in a redis.conf file.
type redisServerConfig struct {
	Port              int
	TLSPort           int
	TLSCertFile       string // The path to the .crt file used as a redis TLS certificate
	TLSKeyFile        string // The path to the .key file used as a redis private TLS key
	TLSCACertFile     string // The path to the .crt file used to authenticate TLS clients/peers
	Replica           bool   // Indicates that the Redis node is a replica and therefore should specify replicaof in its config
	MasterIP          string
	MasterPort        int
	MasterAuth        string // The password that should be used to authenticate with the master node
	MasterUser        string // The user that is capable of running replication commands
	EnableDefaultUser bool
	Users             []*redisUser
}

// newDefaultRedisPrimaryServerConfig generates a redisServerConfig with parameters that are appropriate for a primary
// Redis node.
func newDefaultRedisPrimaryServerConfig(
	masterAuth string,
	masterUser string,
	users ...*redisUser,
) (*redisServerConfig, error) {
	replicaUser, err := newReplicaUser(masterUser, masterAuth)
	if err != nil {
		return nil, err
	}

	return &redisServerConfig{
		Port:              0,
		TLSPort:           6379,
		TLSCertFile:       "redis.crt",
		TLSKeyFile:        "redis.key",
		TLSCACertFile:     "ca.crt",
		Replica:           false,
		MasterAuth:        masterAuth,
		MasterUser:        masterUser,
		EnableDefaultUser: false,
		Users:             append(users, replicaUser),
	}, nil
}
