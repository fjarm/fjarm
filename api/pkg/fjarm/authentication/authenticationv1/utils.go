package authenticationv1

import (
	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
	"buf.build/go/protovalidate"
)

func ValidateSessionID(id *authenticationpb.SessionId) error {
	if id == nil {
		return ErrValidationError
	}
	return protovalidate.Validate(id)
}

func ValidateAccessToken(token *authenticationpb.AccessToken) error {
	if token == nil {
		return ErrValidationError
	}
	return protovalidate.Validate(token)
}

func ValidateRefreshToken(token *authenticationpb.RefreshToken) error {
	if token == nil {
		return ErrValidationError
	}
	return protovalidate.Validate(token)
}

func ValidateSession(session *authenticationpb.Session) error {
	if session == nil {
		return ErrValidationError
	}
	return protovalidate.Validate(session)
}

func ValidateCreateSessionRequest(req *authenticationpb.CreateSessionRequest) error {
	if req == nil {
		return ErrValidationError
	}
	return protovalidate.Validate(req)
}

func ValidateCreateSessionResponse(res *authenticationpb.CreateSessionResponse) error {
	if res == nil {
		return ErrValidationError
	}
	return protovalidate.Validate(res)
}
