package v1

import (
	usersrpc "buf.build/gen/go/fjarm/fjarm/grpc/go/fjarm/users/v1/usersv1grpc"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"context"
	"github.com/bufbuild/protovalidate-go"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

const handlerTag = "grpc_handler"

// ErrUnimplemented is returned when a service method is called before it's been implemented.
var ErrUnimplemented = status.Error(codes.Unimplemented, "rpc is not implemented")

type userDomain interface {
	createUser(ctx context.Context, user *userspb.User) (*userspb.User, error)
	getUserWithID(ctx context.Context, id *userspb.UserId) (*userspb.User, error)
	updateUser(ctx context.Context, user *userspb.User) (*userspb.User, error)
	deleteUser(ctx context.Context, user *userspb.User) error
}

// GrpcHandler implements the gRPC service found in `user_service.proto`.
type GrpcHandler struct {
	usersrpc.UnimplementedUserServiceServer

	domain    userDomain
	logger    *slog.Logger
	validator protovalidate.Validator
}

// CreateUser handles gRPC requests to create a `User` entity.
func (h *GrpcHandler) CreateUser(
	ctx context.Context,
	req *userspb.CreateUserRequest,
) (*userspb.CreateUserResponse, error) {
	return nil, ErrUnimplemented
}

// GetUser handles gRPC requests to retrieve a `User` entity.
func (h *GrpcHandler) GetUser(ctx context.Context, req *userspb.GetUserRequest) (*userspb.GetUserResponse, error) {
	return nil, ErrUnimplemented
}

// UpdateUser handles gRPC requests to modify a field in a `User` entity.
func (h *GrpcHandler) UpdateUser(
	ctx context.Context,
	req *userspb.UpdateUserRequest,
) (*userspb.UpdateUserResponse, error) {
	return nil, ErrUnimplemented
}

// DeleteUser handles gRPC requests to delete an instance of a `User` entity.
func (h *GrpcHandler) DeleteUser(
	ctx context.Context,
	req *userspb.DeleteUserRequest,
) (*userspb.DeleteUserResponse, error) {
	return nil, ErrUnimplemented
}

// NewGrpcHandler creates a concrete users gRPC service with logging and business/domain logic.
func NewGrpcHandler(l *slog.Logger) *GrpcHandler {
	logger := l.With(
		slog.String(logkeys.Tag, handlerTag),
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
	}

	dom := newUserDomain()

	handler := GrpcHandler{
		domain:    dom,
		logger:    logger,
		validator: validator,
	}

	return &handler
}
