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
	delay := DelayDuration(1)
	interceptor := NewConnectRPCAmbiguousTimingInterceptor(l, delay)(next)

	tests := map[string]struct {
		delay  DelayDuration
		output []string
		err    bool
	}{
		"valid_delay": {
			delay:  delay,
			output: []string{"level=\"INFO\"", "msg=\"introduced ambiguous delay\"", "delay=1000ms"},
		},
	}
	t.Parallel()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := connect.NewRequest(
				&[]string{"a", "cool", "request"},
			)
			_, err := interceptor(context.Background(), req)
			if err != nil && !tc.err {
				t.Errorf("NewConnectRPCAmbiguousTimingInterceptor got an unexpected error: %v", err)
			}
			if err == nil && tc.err {
				t.Errorf("NewConnectRPCAmbiguousTimingInterceptor expected an error but got none")
			}
		})
	}
}
