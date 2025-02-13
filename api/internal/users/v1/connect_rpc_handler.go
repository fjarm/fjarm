package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"connectrpc.com/connect"
	"context"
	"github.com/bufbuild/protovalidate-go"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"log/slog"
)

const connectRPCHandlerTag = "connect_rpc_handler"

type userDomain interface {
	createUser(ctx context.Context, user *userspb.User) (*userspb.User, error)
	getUserWithID(ctx context.Context, id *userspb.UserId) (*userspb.User, error)
	updateUser(ctx context.Context, user *userspb.User) (*userspb.User, error)
	deleteUser(ctx context.Context, user *userspb.User) error
}

type ConnectRPCHandler struct {
	domain    userDomain
	logger    *slog.Logger
	validator protovalidate.Validator
}

// CreateUser handles gRPC requests to create a `User` entity.
func (h *ConnectRPCHandler) CreateUser(
	ctx context.Context,
	req *connect.Request[userspb.CreateUserRequest],
) (*connect.Response[userspb.CreateUserResponse], error) {
	return nil, ErrUnimplemented
}

// GetUser handles gRPC requests to retrieve a `User` entity.
func (h *ConnectRPCHandler) GetUser(
	ctx context.Context,
	req *connect.Request[userspb.GetUserRequest],
) (*connect.Response[userspb.GetUserResponse], error) {
	return nil, nil
}

// UpdateUser handles gRPC requests to modify a field in a `User` entity.
func (h *ConnectRPCHandler) UpdateUser(
	ctx context.Context,
	req *connect.Request[userspb.UpdateUserRequest],
) (*connect.Response[userspb.UpdateUserResponse], error) {
	return nil, ErrUnimplemented
}

// DeleteUser handles gRPC requests to delete an instance of a `User` entity.
func (h *ConnectRPCHandler) DeleteUser(
	ctx context.Context,
	req *connect.Request[userspb.DeleteUserRequest],
) (*connect.Response[userspb.DeleteUserResponse], error) {
	return nil, ErrUnimplemented
}

// NewConnectRPCHandler creates a concrete users ConnectRPC service with logging and business/domain logic.
func NewConnectRPCHandler(l *slog.Logger) *ConnectRPCHandler {
	logger := l.With(
		slog.String(logkeys.Tag, connectRPCHandlerTag),
	)

	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithFailFast(),
		protovalidate.WithMessages(
			&userspb.CreateUserRequest{},
			&userspb.CreateUserResponse{},
			&userspb.GetUserRequest{},
			&userspb.GetUserResponse{},
		),
	)

	if err != nil {
		logger.Error(
			"failed to create message validator",
			slog.Any(logkeys.Err, err),
		)
		return nil
	}

	dom := newUserDomain()
	han := ConnectRPCHandler{
		domain:    dom,
		logger:    logger,
		validator: validator,
	}
	return &han
}
