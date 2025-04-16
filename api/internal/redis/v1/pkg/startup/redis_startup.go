package startup

// TODO(2025-03-15): Create Redis admin client.
// TODO(2025-03-15): Use Redis admin client to setup ACLs and TLS.

type redisServerStarter struct{}

func (rs *redisServerStarter) createConfigFile() error {
	return ErrUnimplemented
}

func (rs *redisServerStarter) fetchCertificates() error {
	return ErrUnimplemented
}

func (rs *redisServerStarter) saveCertificates() error {
	return ErrUnimplemented
}

func newRedisServerStarter() *redisServerStarter {
	return &redisServerStarter{}
}
