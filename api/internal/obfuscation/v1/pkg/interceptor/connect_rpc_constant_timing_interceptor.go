package interceptor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"connectrpc.com/connect"

	"github.com/fjarm/fjarm/api/internal/logkeys"
)

const connectRPCAmbiguousTimingInterceptorTag = "connect_rpc_ambiguous_timing_interceptor"

// NewConnectRPCConstantTimingInterceptor creates an interceptor that introduces a constant delay to requests.
func NewConnectRPCConstantTimingInterceptor(l *slog.Logger, dd DelayDuration) connect.UnaryInterceptorFunc {
	if dd <= 0 {
		l.Warn("invalid delay duration", slog.Any("delay", dd))
		dd = DelayDuration(1000)
	}
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			logger := l.With(
				slog.String(logkeys.Tag, connectRPCAmbiguousTimingInterceptorTag),
				slog.String(logkeys.Rpc, req.Spec().Procedure),
			)

			start := time.Now()

			// Introduce a constant delay to mask execution time.
			delay := time.Duration(dd) * time.Millisecond

			logger.InfoContext(
				ctx,
				"introduced timing delay",
				slog.Duration("delay", delay),
				slog.Int64(logkeys.StartTime, start.Unix()),
			)

			res, err := next(ctx, req)

			duration := time.Since(start)
			if duration < delay {
				// If the request was completed before the delay, wait for the remaining time.
				select {
				case <-time.After(delay - duration):
					// Return the response after the delay
					end := time.Now()
					logger.InfoContext(
						ctx,
						"completed request with timing delay",
						slog.Int64(logkeys.EndTime, end.Unix()),
						slog.Duration(logkeys.Duration, duration),
						slog.Any(logkeys.Err, err),
					)
					return res, err
				case <-ctx.Done():
					// Handle context cancellation
					end := time.Now()
					logger.ErrorContext(
						ctx,
						"terminated request with cancelled context",
						slog.Any("err", ctx.Err()),
						slog.Int64(logkeys.EndTime, end.Unix()),
						slog.Duration(logkeys.Duration, duration),
					)
					return nil, connect.NewError(connect.CodeAborted, fmt.Errorf("terminated request"))
				}
			} else {
				// If the request took longer than the delay, return the response immediately.
				// This generally shouldn't happen, but it's possible if the delay is set too low.
				end := time.Now()
				logger.WarnContext(
					ctx,
					"processed for longer than the minimum delay",
					slog.Int64(logkeys.EndTime, end.Unix()),
					slog.Duration(logkeys.Duration, duration),
				)
				return res, err
			}
		}
	}
}
