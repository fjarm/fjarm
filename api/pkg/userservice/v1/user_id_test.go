package v1

import (
	pb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/userservice/v1"
	"github.com/bufbuild/protovalidate-go"
	"github.com/google/uuid"
	"testing"
)

func TestUserId_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(
			&pb.UserId{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		userId  string
		wantErr bool
	}{
		"invalid_not_uuid": {
			userId:  "123",
			wantErr: true,
		},
		"invalid_not_uuid_with_dashes": {
			userId:  "123-456",
			wantErr: true,
		},
		"invalid_empty": {
			userId:  "",
			wantErr: true,
		},
		"valid_uuid": {
			userId:  uuid.New().String(),
			wantErr: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			userId := pb.UserId{UserId: &tc.userId}
			err = validator.Validate(&userId)
			if (err != nil) != tc.wantErr {
				t.Errorf("got error = %v, wantErr = %v, input = %v", err, tc.wantErr, tc.userId)
			}
		})
	}
}
