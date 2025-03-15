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
		ttl  time.Duration
		kind error
	}
	tests := map[string]struct {
		set map[string]testcase
		get map[string]testcase
	}{
		"valid_set_and_get": {
			set: map[string]testcase{
				"key1": {val: "value1", err: false, kind: nil, ttl: 10 * time.Second},
				"key2": {val: "value2", err: false, kind: nil, ttl: 10 * time.Second},
				"key3": {val: "value3", err: false, kind: nil, ttl: 10 * time.Second},
				"key4": {val: "value4", err: true, kind: cachev1.ErrInvalidExpiration},
				"key5": {val: "value5", err: true, kind: cachev1.ErrInvalidExpiration, ttl: 0},
				"key6": {val: "value6", err: true, kind: cachev1.ErrInvalidExpiration, ttl: -1 * time.Second},
				"":     {val: "value7", err: true, kind: cachev1.ErrInvalidKey, ttl: 10 * time.Second},
				" ":    {val: "value7", err: true, kind: cachev1.ErrInvalidKey, ttl: 10 * time.Second},
			},
			get: map[string]testcase{
				"key1": {val: "value1", err: false, kind: nil},
				"key2": {val: "value2", err: false, kind: nil},
				"key3": {val: "value3", err: false, kind: nil},
				"key4": {val: "", err: true, kind: cachev1.ErrCacheMiss},
				"key5": {val: "", err: true, kind: cachev1.ErrCacheMiss},
				"":     {val: "", err: true, kind: cachev1.ErrInvalidKey},
				" ":    {val: "", err: true, kind: cachev1.ErrInvalidKey},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rc := NewRedisCache(rdb, logger)
			for key, val := range tc.set {
				se := rc.Set(ctx, key, []byte(val.val), val.ttl)
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
					t.Errorf("Get got an unexpected error: %v", ge)
				}
				if ge == nil && val.err {
					t.Errorf("Get expected an error but got none")
				}
				if !errors.Is(ge, val.kind) {
					t.Errorf("Get got an unexpected error type: %v, wanted: %v", ge, val.kind)
				}
				if string(actual) != val.val {
					t.Errorf("Get got an unexpected value: %v, wanted: %v", string(actual), val.val)
				}
			}
		})
	}
}

func TestRedisCache_Update(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	type testcase struct {
		key  string
		val  string
		err  bool
		ttl  time.Duration
		kind error
	}
	tests := map[string]struct {
		upd map[string]testcase
		get map[string]testcase
	}{
		"valid_update_and_get": {
			upd: map[string]testcase{
				"upd1": {key: "key1", val: "value1", err: false, kind: nil, ttl: 10 * time.Second},
				"upd2": {key: "key2", val: "value2", err: false, kind: nil, ttl: 10 * time.Second},
				"upd3": {key: "key3", val: "value3", err: false, kind: nil, ttl: 10 * time.Second},
				"upd4": {key: "key4", val: "value4", err: true, kind: cachev1.ErrInvalidExpiration},
				"upd5": {key: "key5", val: "value5", err: true, kind: cachev1.ErrInvalidExpiration, ttl: 0},
				"upd6": {key: "key1", val: "value6", err: false, kind: nil, ttl: 10 * time.Second},
				"upd7": {key: "key2", val: "value6", err: false, kind: nil, ttl: 10 * time.Second},
				"upd8": {key: "", val: "value6", err: true, kind: cachev1.ErrInvalidKey, ttl: 10 * time.Second},
				"upd9": {key: " ", val: "value6", err: true, kind: cachev1.ErrInvalidKey, ttl: 10 * time.Second},
			},
			get: map[string]testcase{
				"get1": {key: "key1", val: "value6", err: false, kind: nil},
				"get2": {key: "key2", val: "value6", err: false, kind: nil},
				"get3": {key: "key3", val: "value3", err: false, kind: nil},
				"get4": {key: "key4", err: true, kind: cachev1.ErrCacheMiss},
				"get5": {key: "key5", err: true, kind: cachev1.ErrCacheMiss},
				"get6": {key: "key1", val: "value6", err: false, kind: nil},
				"get7": {key: "key2", val: "value6", err: false, kind: nil},
				"get8": {key: "", err: true, kind: cachev1.ErrInvalidKey},
				"get9": {key: " ", err: true, kind: cachev1.ErrInvalidKey},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rc := NewRedisCache(rdb, logger)
			for _, val := range tc.upd {
				se := rc.Update(ctx, val.key, []byte(val.val), val.ttl)
				if se != nil && !val.err {
					t.Errorf("Update got an unexpected error: %v", se)
				}
				if se == nil && val.err {
					t.Errorf("Update expected an error but got none")
				}
				if !errors.Is(se, val.kind) {
					t.Errorf("Update got an unexpected error type: %v, wanted: %v", se, val.kind)
				}
			}
			for _, val := range tc.get {
				actual, ge := rc.Get(ctx, val.key)
				if ge != nil && !val.err {
					t.Errorf("Get got an unexpected error: %v", ge)
				}
				if ge == nil && val.err {
					t.Errorf("Get expected an error but got none")
				}
				if !errors.Is(ge, val.kind) {
					t.Errorf(
						"Get got an unexpected error type: %v, wanted: %v, key: %v", ge, val.kind, val.key,
					)
				}
				if string(actual) != val.val {
					t.Errorf(
						"Get got an unexpected value: %v, wanted: %v, key: %v", string(actual), val.val, val.key,
					)
				}
			}
		})
	}
}
