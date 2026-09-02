package authenticationv1

import (
	"testing"

	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

func TestSession_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&authenticationpb.Session{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		session *authenticationpb.Session
		wantErr bool
	}{
		"valid_session": {
			session: &authenticationpb.Session{
				SessionId: &authenticationpb.SessionId{
					SessionId: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				},
				AccessToken: &authenticationpb.AccessToken{
					AccessToken: proto.String("valid-access-token"),
				},
				RefreshToken: &authenticationpb.RefreshToken{
					RefreshToken: proto.String("valid-refresh-token"),
				},
			},
			wantErr: false,
		},
		"invalid_empty_session": {
			session: &authenticationpb.Session{},
			wantErr: true,
		},
		"invalid_missing_session_id": {
			session: &authenticationpb.Session{
				AccessToken: &authenticationpb.AccessToken{
					AccessToken: proto.String("valid-access-token"),
				},
				RefreshToken: &authenticationpb.RefreshToken{
					RefreshToken: proto.String("valid-refresh-token"),
				},
			},
			wantErr: true,
		},
		"invalid_missing_access_token": {
			session: &authenticationpb.Session{
				SessionId: &authenticationpb.SessionId{
					SessionId: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				},
				RefreshToken: &authenticationpb.RefreshToken{
					RefreshToken: proto.String("valid-refresh-token"),
				},
			},
			wantErr: true,
		},
		"invalid_missing_refresh_token": {
			session: &authenticationpb.Session{
				SessionId: &authenticationpb.SessionId{
					SessionId: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				},
				AccessToken: &authenticationpb.AccessToken{
					AccessToken: proto.String("valid-access-token"),
				},
			},
			wantErr: true,
		},
		"invalid_unset_inner_session_id": {
			session: &authenticationpb.Session{
				SessionId: &authenticationpb.SessionId{},
				AccessToken: &authenticationpb.AccessToken{
					AccessToken: proto.String("valid-access-token"),
				},
				RefreshToken: &authenticationpb.RefreshToken{
					RefreshToken: proto.String("valid-refresh-token"),
				},
			},
			wantErr: true,
		},
		"invalid_unset_inner_access_token": {
			session: &authenticationpb.Session{
				SessionId: &authenticationpb.SessionId{
					SessionId: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				},
				AccessToken: &authenticationpb.AccessToken{},
				RefreshToken: &authenticationpb.RefreshToken{
					RefreshToken: proto.String("valid-refresh-token"),
				},
			},
			wantErr: true,
		},
		"invalid_unset_inner_refresh_token": {
			session: &authenticationpb.Session{
				SessionId: &authenticationpb.SessionId{
					SessionId: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				},
				AccessToken: &authenticationpb.AccessToken{
					AccessToken: proto.String("valid-access-token"),
				},
				RefreshToken: &authenticationpb.RefreshToken{},
			},
			wantErr: true,
		},
		"invalid_nil_session": {
			session: nil,
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err = validator.Validate(tc.session)
			if err != nil && !tc.wantErr {
				t.Errorf("Validate got an unexpected error: %v", err)
			}
			if err == nil && tc.wantErr {
				t.Error("Validate expected an error but got nil")
			}
		})
	}
}
