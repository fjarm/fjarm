package authenticationv1

import (
	"testing"

	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

func TestSessionId_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&authenticationpb.SessionId{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		sessionId *authenticationpb.SessionId
		wantErr   bool
	}{
		"valid_session_id": {
			sessionId: &authenticationpb.SessionId{
				SessionId: proto.String("123e4567-e89b-12d3-a456-426614174000"),
			},
			wantErr: false,
		},
		"valid_string_session_id": {
			sessionId: &authenticationpb.SessionId{
				SessionId: proto.String("session-abc-123"),
			},
			wantErr: false,
		},
		"invalid_unset_session_id": {
			sessionId: &authenticationpb.SessionId{},
			wantErr:   true,
		},
		"invalid_nil_session_id": {
			sessionId: nil,
			wantErr:   true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err = validator.Validate(tc.sessionId)
			if err != nil && !tc.wantErr {
				t.Errorf("Validate got an unexpected error: %v", err)
			}
			if err == nil && tc.wantErr {
				t.Error("Validate expected an error but got nil")
			}
		})
	}
}
