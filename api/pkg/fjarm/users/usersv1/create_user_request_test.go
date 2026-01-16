package usersv1

import (
	idempotencypb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/idempotency/v1"
	userspb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		"valid_create_user_request": {
			request: &userspb.CreateUserRequest{
				IdempotencyKey: &idempotencypb.IdempotencyKey{
					IdempotencyKey: proto.String("123e4567-e89b-12d3-a456-426614174000"),
					Timestamp:      timestamppb.Now(),
				},
				UserId: &userspb.UserId{
					UserId: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				},
				User: &userspb.User{
					UserId: &userspb.UserId{
						UserId: proto.String("123e4567-e89b-12d3-a456-426614174000"),
					},
				},
			},
			err: false,
		},
		"invalid_empty_idempotency_key_create_user_request": {
			request: &userspb.CreateUserRequest{
				IdempotencyKey: &idempotencypb.IdempotencyKey{},
				UserId: &userspb.UserId{
					UserId: proto.String("123e4567-e89b-12d3-a456-426614174000"),
				},
				User: &userspb.User{
					UserId: &userspb.UserId{
						UserId: proto.String("123e4567-e89b-12d3-a456-426614174000"),
					},
				},
			},
			err: true,
		},
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
