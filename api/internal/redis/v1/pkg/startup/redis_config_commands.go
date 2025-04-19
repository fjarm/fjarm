package startup

import (
	"fmt"
	"strings"
)

const admin = "@admin"
const all = "@all"
const blocking = "@blocking"
const dangerous = "@dangerous"

const acl = "ACL"
const auth = "AUTH"
const info = "INFO"
const ping = "PING"
const psync = "PSYNC"
const replicaof = "REPLICAOF"
const replconf = "REPLCONF"
const role = "ROLE"

func addCommands(commands ...string) ([]string, error) {
	cmds := []string{}
	for _, command := range commands {
		err := validateCommand(command)
		if err != nil {
			return nil, err
		}
		added := fmt.Sprintf("+%s", command)
		cmds = append(cmds, added)
	}
	return cmds, nil
}

func removeCommands(commands ...string) ([]string, error) {
	cmds := []string{}
	for _, command := range commands {
		err := validateCommand(command)
		if err != nil {
			return nil, err
		}
		removed := fmt.Sprintf("-%s", command)
		cmds = append(cmds, removed)
	}
	return cmds, nil
}

func validateCommand(command string) error {
	if len(command) == 0 {
		return ErrInvalidCommand
	}
	if strings.Contains(command, "+") {
		return ErrInvalidCommand
	}
	if strings.Contains(command, "-") {
		return ErrInvalidCommand
	}
	return nil
}
