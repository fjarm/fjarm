package remote

import (
	"context"
	"errors"
	cachev1 "github.com/fjarm/fjarm/api/internal/cache/v1"
	"io"
	"log/slog"
	"testing"
)

func TestRedisCache_Get(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	addrs := []string{""}
	rdb, err := newRedisClient(addrs)
	if err != nil {
		t.Fatalf("failed to create Redis client: %v", err)
	}
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
				se := rc.Set(context.Background(), key, []byte(val.val), 0)
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
				actual, ge := rc.Get(context.Background(), key)
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
