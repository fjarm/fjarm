package interceptor

import (
	"connectrpc.com/connect"
	"context"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
	"log/slog"
	"time"
)

func NewConnectRPCRequestIDLoggingInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()

			reqID, err := getRequestIDFromReqHeaders(req)
			lvl := slog.LevelInfo
			if err != nil {
				lvl = slog.LevelWarn
			}

			clientIP := req.Peer().Addr

			logger.Log(
				ctx,
				lvl,
				"intercepted request",
				slog.String(tracing.RequestIDKey, reqID),
				slog.Time(logkeys.StartTime, start),
				slog.String(logkeys.Addr, clientIP),
				slog.String(logkeys.Rpc, req.Spec().Procedure),
				slog.Any(logkeys.Err, err),
			)

			var res connect.AnyResponse = nil
			if err == nil {
				res, err = next(ctx, req)
			}

			duration := time.Since(start)

			logger.InfoContext(
				ctx,
				"completed request",
				slog.String(tracing.RequestIDKey, reqID),
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
