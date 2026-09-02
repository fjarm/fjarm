package authenticationv1

import (
	"testing"

	authenticationpb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/authentication/v1"
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

func TestAccessToken_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&authenticationpb.AccessToken{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		accessToken *authenticationpb.AccessToken
		wantErr     bool
	}{
		"valid_access_token": {
			accessToken: &authenticationpb.AccessToken{
				AccessToken: proto.String("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.t-ae7H8M3AzbKlCfgUpjJWCGbKaBEvpAcNJBK9l406s"),
			},
			wantErr: false,
		},
		"valid_plain_string_access_token": {
			accessToken: &authenticationpb.AccessToken{
				AccessToken: proto.String("sample-valid-access-token"),
			},
			wantErr: false,
		},
		"invalid_unset_access_token": {
			accessToken: &authenticationpb.AccessToken{},
			wantErr:     true,
		},
		"invalid_nil_access_token": {
			accessToken: nil,
			wantErr:     true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err = validator.Validate(tc.accessToken)
			if err != nil && !tc.wantErr {
				t.Errorf("Validate got an unexpected error: %v", err)
			}
			if err == nil && tc.wantErr {
				t.Error("Validate expected an error but got nil")
			}
		})
	}
}
