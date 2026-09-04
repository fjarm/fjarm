package users

import (
	"context"
	"errors"
	"log/slog"

	"buf.build/gen/go/fjarm/fjarm/connectrpc/gosimple/fjarm/users/v1/usersv1connect"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"buf.build/go/protovalidate"
	"connectrpc.com/connect"

	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
)

const connectRPCHandlerTag = "connect_rpc_handler"

type userDomain interface {
	createUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.User, error)
	getUser(ctx context.Context, req *userspb.GetUserRequest) (*userspb.User, error)
	updateUser(ctx context.Context, req *userspb.UpdateUserRequest) (*userspb.User, error)
	deleteUser(ctx context.Context, req *userspb.DeleteUserRequest) error
}

// ConnectRPCHandler defines a ConnectRPC handler for the `fjarm.users.v1.UserService` service.
type ConnectRPCHandler struct {
	usersv1connect.UnimplementedUserServiceHandler
	domain    userDomain
	logger    *slog.Logger
	validator protovalidate.Validator
}

// CreateUser handles ConnectRPC requests to create a `User` entity.
func (h *ConnectRPCHandler) CreateUser(
	ctx context.Context,
	req *userspb.CreateUserRequest,
) (*userspb.CreateUserResponse, error) {
	logger := h.logger.With(
		slog.String(logkeys.Tag, connectRPCHandlerTag),
		slog.Any(tracing.RequestIDKey, ctx.Value(tracing.RequestIDKey)),
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

// NewConnectRPCHandler creates a concrete users ConnectRPC service with logging and business/domain logic.
func NewConnectRPCHandler(
	logger *slog.Logger,
	domain userDomain,
	validator protovalidate.Validator,
) *ConnectRPCHandler {
	han := ConnectRPCHandler{
		domain:    domain,
		logger:    logger,
		validator: validator,
	}
	return &han
}
