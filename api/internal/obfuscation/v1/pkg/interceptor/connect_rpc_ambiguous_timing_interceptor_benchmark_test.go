package interceptor

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"connectrpc.com/connect"
)

func Benchmark_NewConnectRPCAmbiguousTimingInterceptor(b *testing.B) {
	dl := slog.Default()
	defer slog.SetDefault(dl)

	var buf bytes.Buffer
	l := slog.New(slog.NewTextHandler(&buf, nil))
	slog.SetDefault(l)

	next := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		res := "a cool response"
		return connect.NewResponse(&res), nil
	}

	interceptor := NewConnectRPCAmbiguousTimingInterceptor(l, DelayDuration(1))(next)

	for i := 0; i < b.N; i++ {
		req := connect.NewRequest(
			&[]string{"a", "cool", "request"},
		)
		_, _ = interceptor(context.Background(), req)
	}
}
