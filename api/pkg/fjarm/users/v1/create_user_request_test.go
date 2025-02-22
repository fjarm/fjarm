package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"github.com/bufbuild/protovalidate-go"
	"testing"
)

func TestCreateUserRequest_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&userspb.CreateUserRequest{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		request *userspb.CreateUserRequest
		err     bool
	}{
		"invalid_empty_create_user_request": {
			request: &userspb.CreateUserRequest{},
			err:     true,
		},
		"invalid_missing_create_user_request": {
			err: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err = validator.Validate(tc.request)
			if err != nil && !tc.err {
				t.Errorf("Validate got an unexpected error = %v", err)
			}
			if err == nil && tc.err {
				t.Error("Validate expected an error but got nil")
			}
		})
	}
}
