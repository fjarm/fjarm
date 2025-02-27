package interceptor

import (
	"connectrpc.com/connect"
	"context"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
	"log/slog"
	"time"
)

// NewConnectRPCRequestIDLoggingInterceptor intercepts ConnectRPC requests and verifies that a key named `request-id` is
// in the request headers with a non-null value. If the key can't be found, the request is automatically rejected.
// Otherwise, the corresponding value is added to the context before completing the request.
func NewConnectRPCRequestIDLoggingInterceptor(l *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()

			reqID, err := getRequestIDFromReqHeaders(req)
			lvl := slog.LevelInfo
			if err != nil {
				lvl = slog.LevelWarn
			}

			clientIP := req.Peer().Addr

			logger := l.With(
				slog.String(tracing.RequestIDKey, reqID),
				slog.String(logkeys.Addr, clientIP),
				slog.String(logkeys.Rpc, req.Spec().Procedure),
			)

			logger.Log(
				ctx,
				lvl,
				"intercepted request",
				slog.Time(logkeys.StartTime, start),
				slog.Any(logkeys.Err, err),
			)

			var res connect.AnyResponse = nil
			if err == nil {
				res, err = next(context.WithValue(ctx, tracing.RequestIDKey, reqID), req)
			}

			duration := time.Since(start)

			logger.InfoContext(
				ctx,
				"completed request",
				slog.Duration(logkeys.Duration, duration),
				slog.Any(logkeys.Err, err),
			)

			return res, err
		}
	}
}

func getRequestIDFromReqHeaders(req connect.AnyRequest) (string, error) {
	header := req.Header().Get(tracing.RequestIDKey)
	if header == "" {
		return "", ErrRequestIDNotFound
	}
	return header, nil
}
