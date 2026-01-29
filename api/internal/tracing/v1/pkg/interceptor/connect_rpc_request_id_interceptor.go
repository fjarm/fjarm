package interceptor

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"

	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
)

const connectRPCRequestIDInterceptorTag = "connect_rpc_request_id_interceptor"

// NewConnectRPCRequestIDLoggingInterceptor intercepts ConnectRPC requests and verifies that a key named `request-id` is
// in the request headers with a non-null value. If the key can't be found, the request is automatically rejected.
// Otherwise, the corresponding value is added to the context before completing the request.
func NewConnectRPCRequestIDLoggingInterceptor(l *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			clientIP := req.Peer().Addr

			logger := l.With(
				slog.String(logkeys.Tag, connectRPCRequestIDInterceptorTag),
				slog.String(logkeys.Addr, clientIP),
				slog.String(logkeys.Rpc, req.Spec().Procedure),
			)

			reqID, err := getRequestIDFromReqHeaders(req)
			if err != nil {
				logger.WarnContext(
					ctx,
					"failed to verify request-id header",
					slog.Any(logkeys.Err, err),
				)
				return nil, err
			}

			logger.InfoContext(
				ctx,
				"verified request contains request-id header",
				slog.String(tracing.RequestIDKey, reqID),
			)

			res, err := next(ctx, req)
			return res, err
		}
	}
}

func getRequestIDFromReqHeaders(req connect.AnyRequest) (string, error) {
	header := req.Header().Get(tracing.RequestIDKey)
	if header == "" {
		return "", tracing.ErrRequestIDNotFound
	}
	return header, nil
}
