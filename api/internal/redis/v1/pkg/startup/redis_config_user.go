package startup

type redisUser struct {
	Username        string
	Password        string
	EnabledCommands []string
}

// newRedisUser creates a Redis user with the given username, password, and enabled commands. By default, command groups
// like @admin, @blocking, and @dangerous will be removed.
func newRedisUser(username string, password string, enabledCommands ...string) (*redisUser, error) {
	allCommands := []string{}

	added, err := addCommands(enabledCommands...)
	if err != nil {
		return nil, err
	}
	allCommands = append(allCommands, added...)

	commandsToDisable := []string{admin, blocking, dangerous}
	removed, err := removeCommands(commandsToDisable...)
	if err != nil {
		return nil, err
	}
	allCommands = append(allCommands, removed...)

	return &redisUser{
		Username:        username,
		Password:        password,
		EnabledCommands: allCommands,
	}, nil
}

// newReplicaUser creates a Redis user with the PING, PSYNC, REPLCONF, and other commands required for replication
// enabled.
func newReplicaUser(username string, password string) (*redisUser, error) {
	allCommands := []string{}

	commandsToDisable := []string{admin, blocking, dangerous}
	removed, err := removeCommands(commandsToDisable...)
	if err != nil {
		return nil, err
	}
	allCommands = append(allCommands, removed...)

	enabledCommands := []string{acl, auth, info, ping, psync, replicaof, replconf, role}
	added, err := addCommands(enabledCommands...)
	if err != nil {
		return nil, err
	}
	allCommands = append(allCommands, added...)

	replicaUser := redisUser{
		Username:        username,
		Password:        password,
		EnabledCommands: allCommands,
	}
	return &replicaUser, nil
}
