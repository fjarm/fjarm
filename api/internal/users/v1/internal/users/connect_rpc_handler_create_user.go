package users

import (
	"context"
	"errors"
	"log/slog"

	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"

	"connectrpc.com/connect"

	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
)

const connectRPCHandlerCreateUserTag = "connect_rpc_handler_create_user"

// CreateUser handles ConnectRPC requests to create a `User` entity.
func (h *ConnectRPCHandler) CreateUser(
	ctx context.Context,
	req *userspb.CreateUserRequest,
) (*userspb.CreateUserResponse, error) {
	logger := h.logger.With(
		slog.String(logkeys.Tag, connectRPCHandlerCreateUserTag),
		slog.String(tracing.RequestIDKey, tracing.RequestIDFromContext(ctx)),
	)
	logger.InfoContext(ctx, "received request to create user")

	// Create the user entity.
	usr, err := h.domain.createUser(ctx, req)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"failed to create user entity",
			slog.String(logkeys.Raw, req.String()),
			slog.Any(logkeys.Err, err),
		)
	}
	if err != nil && errors.Is(err, ErrOperationFailed) {
		// Typically returned for authentication issues. Obscure this from the client.
		return nil, connect.NewError(connect.CodeInternal, err)
	} else if err != nil && errors.Is(err, ErrInvalidArgument) {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	} else if err != nil {
		// A true edge case. This shouldn't happen.
		return nil, connect.NewError(connect.CodeUnknown, ErrOperationFailed)
	}

	res := &userspb.CreateUserResponse{
		User: usr,
	}
	// User creation was successful.
	return res, nil
}
