package interceptor

import (
	"bytes"
	"connectrpc.com/connect"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestNewConnectRPCRequestIDLoggingInterceptor_LogOutput(t *testing.T) {
	dl := slog.Default()
	defer slog.SetDefault(dl)

	var buf bytes.Buffer
	l := slog.New(slog.NewTextHandler(&buf, nil))
	slog.SetDefault(l)

	next := func(context.Context, connect.AnyRequest) (connect.AnyResponse, error) {
		return nil, nil
	}

	si := NewConnectRPCRequestIDLoggingInterceptor(l)(next)

	tests := map[string]struct {
		headers  map[string]string
		expected string
		err      bool
	}{
		"valid_non_empty_request_id": {
			headers:  map[string]string{"request-id": "abc123"},
			expected: "INFO msg=\"intercepted request\" request-id=abc123",
			err:      false,
		},
		"invalid_empty_value_request_id": {
			headers:  map[string]string{"request-id": ""},
			expected: "WARN msg=\"intercepted request\" request-id=\"\"",
			err:      true,
		},
		"invalid_empty_request_id": {
			headers:  map[string]string{},
			expected: "WARN msg=\"intercepted request\" request-id=\"\"",
			err:      true,
		},
		"invalid_incorrect_key_request_id": {
			headers:  map[string]string{"Request-id": "abc123"},
			expected: "WARN msg=\"intercepted request\" request-id=\"\"",
			err:      true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := connect.NewRequest(
				&[]string{"a", "cool", "request"},
			)
			for key, val := range tc.headers {
				req.Header().Add(key, val)
			}
			_, err := si(context.Background(), req)
			if err != nil && !tc.err {
				t.Errorf("NewConnectRPCRequestIDLoggingInterceptor got an unexpected error: %v", err)
			}
			actual := buf.String()
			if !strings.Contains(actual, tc.expected) {
				t.Errorf("NewConnectRPCRequestIDLoggingInterceptor got: %v, want: %v", actual, tc.expected)
			}
		})
	}
}
