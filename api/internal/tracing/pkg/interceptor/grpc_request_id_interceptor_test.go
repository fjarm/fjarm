package interceptor

import (
	"bytes"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"strings"
	"testing"
)

func TestRequestIDLoggingInterceptor_LogOutput(t *testing.T) {
	dl := slog.Default()
	defer slog.SetDefault(dl)

	var buf bytes.Buffer
	l := slog.New(slog.NewTextHandler(&buf, nil))

	slog.SetDefault(l)

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test/method",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return nil, nil
	}
	si := RequestIDLoggingGRPCInterceptor(l)

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
			ctx := metadata.NewIncomingContext(context.Background(), metadata.New(tc.headers))
			_, err := si(ctx, nil, info, handler)
			if err != nil && !tc.err {
				t.Errorf("RequestIDLoggingGRPCInterceptor got an unexpected error: %v", err)
			}

			actual := buf.String()
			if !strings.Contains(actual, tc.expected) {
				t.Errorf("RequestIDLoggingGRPCInterceptor got: %v, want: %v", actual, tc.expected)
			}
		})
	}
}
