package pkg

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"log/slog"
	"math/rand"
	"time"
)

const connectRPCAmbiguousTimingInterceptorTag = "connect_rpc_ambiguous_timing_interceptor"

// NewConnectRPCAmbiguousTimingInterceptor creates an interceptor that introduces random delays to requests.
func NewConnectRPCAmbiguousTimingInterceptor(l *slog.Logger, dd DelayDuration) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			logger := l.With(
				slog.String(logkeys.Tag, connectRPCAmbiguousTimingInterceptorTag),
			)

			// Introduce a random delay between 0 and dd milliseconds
			delay := time.Duration(rand.Intn(int(dd))) * time.Millisecond

			logger.InfoContext(ctx, "introduced ambiguous delay", slog.Duration("delay", delay))

			select {
			case <-time.After(delay):
				// Proceed with the next handler after the delay
				return next(ctx, req)
			case <-ctx.Done():
				// Handle context cancellation
				logger.ErrorContext(ctx, "terminated request", slog.Any("err", ctx.Err()))
				return nil, connect.NewError(connect.CodeAborted, fmt.Errorf("terminated request"))
			}
		}
	}
}
