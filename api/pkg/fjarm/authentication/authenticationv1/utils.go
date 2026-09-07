package authenticationv1

import (
	"context"

	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
	"buf.build/go/protovalidate"
)

func ValidateSessionID(_ context.Context, id *authenticationpb.SessionId) error {
	if id == nil {
		return ErrValidationError
	}
	if !id.HasSessionId() {
		return ErrValidationError
	}
	if id.GetSessionId() == "" {
		return ErrValidationError
	}
	return protovalidate.Validate(id)
}

func ValidateAccessToken(_ context.Context, token *authenticationpb.AccessToken) error {
	if token == nil {
		return ErrValidationError
	}
	if !token.HasAccessToken() {
		return ErrValidationError
	}
	if token.GetAccessToken() == "" {
		return ErrValidationError
	}
	return protovalidate.Validate(token)
}

func ValidateRefreshToken(_ context.Context, token *authenticationpb.RefreshToken) error {
	if token == nil {
		return ErrValidationError
	}
	if !token.HasRefreshToken() {
		return ErrValidationError
	}
	if token.GetRefreshToken() == "" {
		return ErrValidationError
	}
	return protovalidate.Validate(token)
}

func ValidateSession(ctx context.Context, session *authenticationpb.Session) error {
	if session == nil {
		return ErrValidationError
	}
	err := ValidateSessionID(ctx, session.GetSessionId())
	if err != nil {
		return err
	}
	err = ValidateAccessToken(ctx, session.GetAccessToken())
	if err != nil {
		return err
	}
	err = ValidateRefreshToken(ctx, session.GetRefreshToken())
	if err != nil {
		return err
	}
	return protovalidate.Validate(session)
}

func ValidateCreateSessionRequest(_ context.Context, req *authenticationpb.CreateSessionRequest) error {
	if req == nil {
		return ErrValidationError
	}
	if req.GetPassword().GetPassword() == "" {
		return ErrValidationError
	}
	return protovalidate.Validate(req)
}

func ValidateCreateSessionResponse(ctx context.Context, res *authenticationpb.CreateSessionResponse) error {
	if res == nil {
		return ErrValidationError
	}
	if res.HasSession() {
		err := ValidateSession(ctx, res.GetSession())
		if err != nil {
			return err
		}
	}
	return protovalidate.Validate(res)
}
