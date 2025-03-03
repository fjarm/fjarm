package pkg

import (
	"bytes"
	"connectrpc.com/connect"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestNewConnectRPCAmbiguousTimingInterceptor_LogOutput(t *testing.T) {
	dl := slog.Default()
	defer slog.SetDefault(dl)

	var buf bytes.Buffer
	l := slog.New(slog.NewTextHandler(&buf, nil))
	slog.SetDefault(l)

	next := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return nil, nil
	}

	tests := map[string]struct {
		delay  DelayDuration
		output []string
		err    bool
	}{
		"valid_delay": {
			delay:  DelayDuration(1000),
			output: []string{"level=INFO", "msg=\"introduced ambiguous delay\"", "delay"},
		},
		"invalid_negative_delay": {
			delay:  DelayDuration(-1),
			output: []string{"level=INFO", "level=WARN", "msg=\"introduced ambiguous delay\"", "msg=\"invalid delay duration\"", "delay"},
		},
	}
	t.Parallel()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			interceptor := NewConnectRPCAmbiguousTimingInterceptor(l, tc.delay)(next)
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
			actual := buf.String()
			for _, exp := range tc.output {
				if !strings.Contains(actual, exp) {
					t.Errorf("NewConnectRPCAmbiguousTimingInterceptor got: %s, want: %s", actual, tc.output)
				}
			}
		})
	}
}
