package authentication

import (
	"log/slog"

	"buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/authentication/v1/authenticationv1connect"
	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
	"buf.build/go/protovalidate"

	"github.com/fjarm/fjarm/api/internal/logkeys"
)

const connectRPCHandlerTag = "connect_rpc_handler"

type AuthenticationServiceConnectRPCHandler struct {
	authenticationv1connect.UnimplementedAuthenticationServiceHandler
	logger    *slog.Logger
	validator protovalidate.Validator
}

func NewConnectRPCHandler(l *slog.Logger) *AuthenticationServiceConnectRPCHandler {
	logger := l.With(
		slog.String(logkeys.Tag, connectRPCHandlerTag),
	)

	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithFailFast(),
		protovalidate.WithMessages(
			&authenticationpb.CreateSessionRequest{},
			&authenticationpb.CreateSessionResponse{},
		),
	)
	if err != nil {
		logger.Error("failed to create message validator", slog.Any(logkeys.Err, err))
		return nil
	}

	h := AuthenticationServiceConnectRPCHandler{
		logger:    logger,
		validator: validator,
	}
	return &h
}
