package authenticationv1

import (
	"testing"

	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

func TestCreateSessionResponse_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&authenticationpb.CreateSessionResponse{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		response *authenticationpb.CreateSessionResponse
		wantErr  bool
	}{
		"valid_create_session_response": {
			response: &authenticationpb.CreateSessionResponse{
				Session: &authenticationpb.Session{
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
			},
			wantErr: false,
		},
		"valid_empty_create_session_response": {
			response: &authenticationpb.CreateSessionResponse{},
			wantErr:  false,
		},
		"invalid_with_invalid_session": {
			response: &authenticationpb.CreateSessionResponse{
				Session: &authenticationpb.Session{},
			},
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err = validator.Validate(tc.response)
			if err != nil && !tc.wantErr {
				t.Errorf("Validate got an unexpected error: %v", err)
			}
			if err == nil && tc.wantErr {
				t.Error("Validate expected an error but got nil")
			}
		})
	}
}
