package authentication

import (
	"context"

	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
)

func (h *ConnectRPCHandler) CreateSession(
	ctx context.Context,
	req *authenticationpb.CreateSessionRequest,
) (*authenticationpb.CreateSessionResponse, error) {
	return h.UnimplementedAuthenticationServiceHandler.CreateSession(ctx, req)
}
