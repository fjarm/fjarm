package remote

import (
	"context"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/redis/rueidis"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"log/slog"
	"os"
	"testing"
)

var rdb rueidis.Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Use redis.Run instead of the deprecated redis.RunContainer.
	// Specify the image with digest using WithImage option.
	container, err := redis.Run(
		ctx,
		"redis@sha256:9cabfa9c15e13f9e4faee0f80d4373cd76e7b8d5a678b9036402b1b0ed9c661b",
	)
	if err != nil {
		slog.Error("failed to start Redis container", slog.Any(logkeys.Err, err))
	}
	// Clean up the container after the test
	defer func() {
		if te := container.Terminate(ctx); err != nil {
			slog.Error("failed to terminate Redis container", slog.Any(logkeys.Err, te))
		}
		rdb.Close()
	}()

	// Get the connection URI directly from the Redis module
	connectionURI, err := container.ConnectionString(ctx)
	if err != nil {
		slog.Error("failed to get Redis connection URI", slog.Any(logkeys.Err, err))
	}

	addrs := rueidis.MustParseURL(connectionURI).InitAddress
	rdb, err = NewRedisClient(nil, rueidis.AuthCredentials{}, addrs)
	if err != nil {
		slog.Error("failed to create Redis client", slog.Any(logkeys.Err, err))
	}

	// Run the tests
	code := m.Run()

	// Exit with the appropriate code
	os.Exit(code)
}
