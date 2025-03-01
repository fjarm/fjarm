package pkg

import (
	"bytes"
	"connectrpc.com/connect"
	"context"
	"log/slog"
	"testing"
)

func TestNewConnectRPCAmbiguousTimingInterceptor_LogOutput(t *testing.T) {
	dl := slog.Default()
	defer slog.SetDefault(dl)

	var buf bytes.Buffer
	l := slog.New(slog.NewJSONHandler(&buf, nil))
	slog.SetDefault(l)

	next := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return nil, nil
	}
	interceptor := NewConnectRPCAmbiguousTimingInterceptor(l, DelayDuration_15000ms)(next)

	tests := map[string]struct {
	}{}
	for name, _ := range tests {
		t.Run(name, func(t *testing.T) {
			req := connect.NewRequest(
				&[]string{"a", "cool", "request"},
			)
			_, err := interceptor(context.Background(), req)
			if err != nil {

			}
		})
	}
}
