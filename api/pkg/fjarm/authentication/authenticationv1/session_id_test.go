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
		err       bool
	}{
		"valid_session_id": {
			sessionId: &authenticationpb.SessionId{
				SessionId: proto.String("123e4567-e89b-12d3-a456-426614174000"),
			},
			err: false,
		},
		"invalid_unset_session_id": {
			sessionId: &authenticationpb.SessionId{},
			err:       true,
		},
		"invalid_nil_session_id": {
			sessionId: nil,
			err:       true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err = validator.Validate(tc.sessionId)
			if err != nil && !tc.err {
				t.Errorf("Validate got an unexpected error: %v", err)
			}
			if err == nil && tc.err {
				t.Error("Validate expected an error but got nil")
			}
		})
	}
}
