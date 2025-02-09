package v1

import (
	userservicepb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"github.com/bufbuild/protovalidate-go"
	"testing"
)

func TestUserPassword_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&userservicepb.UserPassword{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		password string
		wantErr  bool
	}{
		"valid_non_empty_password": {
			password: "what a cool password",
			wantErr:  false,
		},
		"valid_min_length_string_password": {
			password: "1",
			wantErr:  false,
		},
		"invalid_empty_string_password": {
			password: "",
			wantErr:  true,
		},
		"invalid_no_value_password": {
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			password := &userservicepb.UserPassword{
				Password: &tc.password,
			}
			err = validator.Validate(password)
			if err != nil && !tc.wantErr {
				t.Errorf("got error = %v, wantErr = %v, input = %v", err, tc.wantErr, tc.password)
			}
		})
	}
}
