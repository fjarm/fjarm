package authentication

import (
	"log/slog"

	"buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/authentication/v1/authenticationv1connect"
	"buf.build/go/protovalidate"
)

type AuthenticationServiceConnectRPCHandler struct {
	authenticationv1connect.UnimplementedAuthenticationServiceHandler
	logger    *slog.Logger
	validator protovalidate.Validator
}
