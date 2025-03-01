package pkg

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
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
				slog.Any(tracing.RequestIDKey, ctx.Value(tracing.RequestIDKey)),
				slog.String(logkeys.Rpc, req.Spec().Procedure),
			)

			start := time.Now()
			// Introduce a random delay between 0 and dd milliseconds (usually 15 seconds).
			delay := time.Duration(rand.Intn(int(dd))) * time.Millisecond

			logger.InfoContext(ctx, "introduced ambiguous delay", slog.Duration("delay", delay))

			res, err := next(ctx, req)

			duration := time.Since(start)
			if duration < delay {
				// If the request was completed before the delay, wait for the remaining time.
				select {
				case <-time.After(delay - duration):
					// Return the response after the delay
					return res, err
				case <-ctx.Done():
					// Handle context cancellation
					logger.ErrorContext(ctx, "terminated request", slog.Any("err", ctx.Err()))
					return nil, connect.NewError(connect.CodeAborted, fmt.Errorf("terminated request"))
				}
			} else {
				// If the request took longer than the delay, return the response immediately.
				// This generally shouldn't happen, but it's possible if the delay is set too low.
				logger.WarnContext(ctx, "processed for longer than the minimum delay")
				return res, err
			}
		}
	}
}
