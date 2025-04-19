package startup

import (
	"context"
	"fmt"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"io"
	"log/slog"
	"os"
	"text/template"
)

// TODO(2025-03-15): Create Redis admin client.
// TODO(2025-03-15): Use Redis admin client to setup ACLs and TLS.

type RedisServerStarter struct {
	logger *slog.Logger
}

// CreateConfigFile creates an empty redis.conf file with the name specified by fname. The caller is responsible for
// closing the file if non-nil, non-error values are returned.
func (rs *RedisServerStarter) CreateConfigFile(fname string) (*os.File, error) {
	f, err := os.Create(fname)
	if err != nil && f != nil {
		return nil, fmt.Errorf(
			"failed to create file name with path %s. %w. closing file: %w",
			fname,
			err,
			f.Close(),
		)
	} else if err != nil {
		return nil, err
	}
	return f, nil
}

func (rs *RedisServerStarter) fetchCertificates() error {
	return ErrUnimplemented
}

func (rs *RedisServerStarter) saveCertificates() error {
	return ErrUnimplemented
}

// WriteNewRedisPrimaryConfigFile creates a redisServerConfig struct with parameters that fit a Redis primary node and
// translates that struct to a redis.conf file using text templating.
func (rs *RedisServerStarter) WriteNewRedisPrimaryConfigFile(
	ctx context.Context,
	masteruser string,
	masterauth string,
	destination io.WriteCloser,
) error {
	defer func() {
		err := destination.Close()
		if err != nil {
			rs.logger.WarnContext(ctx, "failed to close Redis config file", slog.Any(logkeys.Err, err))
		}
	}()

	replicaUser, err := newReplicaUser(masteruser, masterauth)
	if err != nil {
		return err
	}

	configStruct, err := newDefaultRedisPrimaryServerConfig(replicaUser.Username, replicaUser.Password, replicaUser)
	if err != nil {
		return err
	}

	configTemplate := template.Must(template.New("config").Parse(config))
	err = configTemplate.Execute(destination, configStruct)
	if err != nil {
		return err
	}
	return nil
}

func NewRedisServerStarter() *RedisServerStarter {
	return &RedisServerStarter{}
}
