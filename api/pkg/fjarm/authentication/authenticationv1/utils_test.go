package authenticationv1

import (
	"context"
	"testing"

	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"google.golang.org/protobuf/proto"
)

func TestValidateSessionID(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		id  *authenticationpb.SessionId
		err bool
	}{
		"valid_session_id": {
			id: &authenticationpb.SessionId{
				SessionId: proto.String("sess-123"),
			},
			err: false,
		},
		"invalid_nil_session_id": {
			id:  nil,
			err: true,
		},
		"invalid_unset_session_id": {
			id:  &authenticationpb.SessionId{},
			err: true,
		},
		"invalid_empty_session_id": {
			id: &authenticationpb.SessionId{
				SessionId: proto.String(""),
			},
			err: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateSessionID(ctx, tc.id)
			if err != nil && !tc.err {
				t.Errorf("ValidateSessionID got an unexpected error = %v", err)
			}
			if err == nil && tc.err {
				t.Error("ValidateSessionID expected an error but got nil")
			}
		})
	}
}

func TestValidateAccessToken(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		token *authenticationpb.AccessToken
		err   bool
	}{
		"valid_access_token": {
			token: &authenticationpb.AccessToken{
				AccessToken: proto.String("valid-token"),
			},
			err: false,
		},
		"invalid_nil_access_token": {
			token: nil,
			err:   true,
		},
		"invalid_unset_access_token": {
			token: &authenticationpb.AccessToken{},
			err:   true,
		},
		"invalid_empty_access_token": {
			token: &authenticationpb.AccessToken{
				AccessToken: proto.String(""),
			},
			err: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateAccessToken(ctx, tc.token)
			if err != nil && !tc.err {
				t.Errorf("ValidateAccessToken got an unexpected error = %v", err)
			}
			if err == nil && tc.err {
				t.Error("ValidateAccessToken expected an error but got nil")
			}
		})
	}
}

func TestValidateRefreshToken(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		token *authenticationpb.RefreshToken
		err   bool
	}{
		"valid_refresh_token": {
			token: &authenticationpb.RefreshToken{
				RefreshToken: proto.String("valid-refresh"),
			},
			err: false,
		},
		"invalid_nil_refresh_token": {
			token: nil,
			err:   true,
		},
		"invalid_unset_refresh_token": {
			token: &authenticationpb.RefreshToken{},
			err:   true,
		},
		"invalid_empty_refresh_token": {
			token: &authenticationpb.RefreshToken{
				RefreshToken: proto.String(""),
			},
			err: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateRefreshToken(ctx, tc.token)
			if err != nil && !tc.err {
				t.Errorf("ValidateRefreshToken got an unexpected error = %v", err)
			}
			if err == nil && tc.err {
				t.Error("ValidateRefreshToken expected an error but got nil")
			}
		})
	}
}

func TestValidateSession(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		session *authenticationpb.Session
		err     bool
	}{
		"valid_session": {
			session: &authenticationpb.Session{
				SessionId:    &authenticationpb.SessionId{SessionId: proto.String("sess-123")},
				AccessToken:  &authenticationpb.AccessToken{AccessToken: proto.String("access-123")},
				RefreshToken: &authenticationpb.RefreshToken{RefreshToken: proto.String("refresh-123")},
			},
			err: false,
		},
		"invalid_nil_session": {
			session: nil,
			err:     true,
		},
		"invalid_empty_session": {
			session: &authenticationpb.Session{},
			err:     true,
		},
		"invalid_missing_session_id": {
			session: &authenticationpb.Session{
				AccessToken:  &authenticationpb.AccessToken{AccessToken: proto.String("access-123")},
				RefreshToken: &authenticationpb.RefreshToken{RefreshToken: proto.String("refresh-123")},
			},
			err: true,
		},
		"invalid_missing_access_token": {
			session: &authenticationpb.Session{
				SessionId:    &authenticationpb.SessionId{SessionId: proto.String("sess-123")},
				RefreshToken: &authenticationpb.RefreshToken{RefreshToken: proto.String("refresh-123")},
			},
			err: true,
		},
		"invalid_missing_refresh_token": {
			session: &authenticationpb.Session{
				SessionId:   &authenticationpb.SessionId{SessionId: proto.String("sess-123")},
				AccessToken: &authenticationpb.AccessToken{AccessToken: proto.String("access-123")},
			},
			err: true,
		},
		"invalid_empty_inner_session_id": {
			session: &authenticationpb.Session{
				SessionId:    &authenticationpb.SessionId{SessionId: proto.String("")},
				AccessToken:  &authenticationpb.AccessToken{AccessToken: proto.String("access-123")},
				RefreshToken: &authenticationpb.RefreshToken{RefreshToken: proto.String("refresh-123")},
			},
			err: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateSession(ctx, tc.session)
			if err != nil && !tc.err {
				t.Errorf("ValidateSession got an unexpected error = %v", err)
			}
			if err == nil && tc.err {
				t.Error("ValidateSession expected an error but got nil")
			}
		})
	}
}

func TestValidateCreateSessionRequest(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		req *authenticationpb.CreateSessionRequest
		err bool
	}{
		"valid_request": {
			req: &authenticationpb.CreateSessionRequest{
				IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				EmailAddress:   &userspb.UserEmailAddress{EmailAddress: proto.String("user@example.com")},
				Password:       &userspb.UserPassword{Password: proto.String("password123")},
			},
			err: false,
		},
		"invalid_nil_request": {
			req: nil,
			err: true,
		},
		"invalid_empty_request": {
			req: &authenticationpb.CreateSessionRequest{},
			err: true,
		},
		"invalid_empty_password": {
			req: &authenticationpb.CreateSessionRequest{
				IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				EmailAddress:   &userspb.UserEmailAddress{EmailAddress: proto.String("user@example.com")},
				Password:       &userspb.UserPassword{Password: proto.String("")},
			},
			err: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateCreateSessionRequest(ctx, tc.req)
			if err != nil && !tc.err {
				t.Errorf("ValidateCreateSessionRequest got an unexpected error = %v", err)
			}
			if err == nil && tc.err {
				t.Error("ValidateCreateSessionRequest expected an error but got nil")
			}
		})
	}
}

func TestValidateCreateSessionResponse(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		res *authenticationpb.CreateSessionResponse
		err bool
	}{
		"valid_empty_response": {
			res: &authenticationpb.CreateSessionResponse{},
			err: false,
		},
		"valid_response_with_session": {
			res: &authenticationpb.CreateSessionResponse{
				Session: &authenticationpb.Session{
					SessionId:    &authenticationpb.SessionId{SessionId: proto.String("sess-123")},
					AccessToken:  &authenticationpb.AccessToken{AccessToken: proto.String("access-123")},
					RefreshToken: &authenticationpb.RefreshToken{RefreshToken: proto.String("refresh-123")},
				},
			},
			err: false,
		},
		"invalid_nil_response": {
			res: nil,
			err: true,
		},
		"invalid_response_with_invalid_session": {
			res: &authenticationpb.CreateSessionResponse{
				Session: &authenticationpb.Session{},
			},
			err: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ValidateCreateSessionResponse(ctx, tc.res)
			if err != nil && !tc.err {
				t.Errorf("ValidateCreateSessionResponse got an unexpected error = %v", err)
			}
			if err == nil && tc.err {
				t.Error("ValidateCreateSessionResponse expected an error but got nil")
			}
		})
	}
}
