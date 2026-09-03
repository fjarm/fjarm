package authentication

import (
	"context"

	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
	"connectrpc.com/connect"
)

func (h *AuthenticationServiceConnectRPCHandler) CreateSession(
	ctx context.Context,
	req *connect.Request[authenticationpb.CreateSessionRequest],
) (*connect.Response[authenticationpb.CreateSessionResponse], error) {
	return h.UnimplementedAuthenticationServiceHandler.CreateSession(ctx, req)
}
