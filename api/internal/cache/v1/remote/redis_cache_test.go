package remote

import (
	"context"
	"errors"
	cachev1 "github.com/fjarm/fjarm/api/internal/cache/v1"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/redis/rueidis"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"
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
	}()

	// Get the connection URI directly from the Redis module
	connectionURI, err := container.ConnectionString(ctx)
	if err != nil {
		slog.Error("failed to get Redis connection URI", slog.Any(logkeys.Err, err))
	}

	addrs := rueidis.MustParseURL(connectionURI).InitAddress
	rdb, err = newRedisClient(addrs)
	if err != nil {
		slog.Error("failed to create Redis client", slog.Any(logkeys.Err, err))
	}

	// Run the tests
	code := m.Run()

	// Exit with the appropriate code
	os.Exit(code)
}

func TestRedisCache_GetAndSet(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	type testcase struct {
		val  string
		err  bool
		kind error
	}
	tests := map[string]struct {
		set map[string]testcase
		get map[string]testcase
	}{
		"valid_set_and_get": {
			set: map[string]testcase{
				"key1": {val: "value1", err: false, kind: nil},
				"key2": {val: "value2", err: false, kind: nil},
				"key3": {val: "value3", err: false, kind: nil},
			},
			get: map[string]testcase{
				"key1": {val: "value1", err: false, kind: nil},
				"key2": {val: "value2", err: false, kind: nil},
				"key3": {val: "value3", err: false, kind: nil},
				"key4": {val: "", err: true, kind: cachev1.ErrCacheMiss},
				"key5": {val: "", err: true, kind: cachev1.ErrCacheMiss},
				"":     {val: "", err: true, kind: cachev1.ErrCacheMiss},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rc := NewRedisCache(rdb, logger)
			for key, val := range tc.set {
				se := rc.Set(ctx, key, []byte(val.val), 1*time.Minute)
				if se != nil && !val.err {
					t.Errorf("Set got an unexpected error: %v", se)
				}
				if se == nil && val.err {
					t.Errorf("Set expected an error but got none")
				}
				if !errors.Is(se, val.kind) {
					t.Errorf("Set got an unexpected error type: %v, wanted: %v", se, val.kind)
				}
			}
			for key, val := range tc.get {
				actual, ge := rc.Get(ctx, key)
				if ge != nil && !val.err {
					t.Errorf("Set got an unexpected error: %v", ge)
				}
				if ge == nil && val.err {
					t.Errorf("Set expected an error but got none")
				}
				if !errors.Is(ge, val.kind) {
					t.Errorf("Set got an unexpected error type: %v, wanted: %v", ge, val.kind)
				}
				if string(actual) != val.val {
					t.Errorf("Get got an unexpected value: %v, wanted: %v", string(actual), val.val)
				}
			}
		})
	}
}
