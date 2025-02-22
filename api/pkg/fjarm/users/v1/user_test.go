package v1

import (
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestUserTest_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&userspb.User{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		request *userspb.User
		err     bool
	}{
		"valid_user": {
			request: &userspb.User{
				UserId: &userspb.UserId{
					UserId: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				},
			},
			err: false,
		},
		"invalid_empty_user": {
			request: &userspb.User{},
			err:     false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err = validator.Validate(tc.request)
			if err != nil && !tc.err {
				t.Errorf("Validate got an unexpected error: %v", err)
			}
			if err == nil && tc.err {
				t.Error("Validate expected an error but got nil")
			}
		})
	}
}
