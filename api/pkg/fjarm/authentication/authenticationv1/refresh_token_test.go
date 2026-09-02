package authenticationv1

import (
	"testing"

	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

func TestRefreshToken_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&authenticationpb.RefreshToken{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		refreshToken *authenticationpb.RefreshToken
		err          bool
	}{
		"valid_refresh_token": {
			refreshToken: &authenticationpb.RefreshToken{
				RefreshToken: proto.String("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.t-ae7H8M3AzbKlCfgUpjJWCGbKaBEvpAcNJBK9l406s"),
			},
			err: false,
		},
		"invalid_unset_refresh_token": {
			refreshToken: &authenticationpb.RefreshToken{},
			err:         true,
		},
		"invalid_nil_refresh_token": {
			refreshToken: nil,
			err:         true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err = validator.Validate(tc.refreshToken)
			if err != nil && !tc.err {
				t.Errorf("Validate got an unexpected error: %v", err)
			}
			if err == nil && tc.err {
				t.Error("Validate expected an error but got nil")
			}
		})
	}
}
