package usersv1

import (
	userservicepb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"buf.build/go/protovalidate"
	"testing"
)

func TestUserHandle_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&userservicepb.UserHandle{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		handle  string
		wantErr bool
	}{
		"valid_single_character": {
			handle:  "a",
			wantErr: false,
		},
		"valid_from_example": {
			handle:  "therealgleepglop",
			wantErr: false,
		},
		"invalid_includes_space": {
			handle:  "the real",
			wantErr: true,
		},
		"invalid_incorrect_length": {
			handle:  "",
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			handle := &userservicepb.UserHandle{
				Handle: &tc.handle,
			}
			err = validator.Validate(handle)
			if err != nil && !tc.wantErr {
				t.Errorf("got error = %v, wantErr = %v, input = %v", err, tc.wantErr, tc.handle)
			}
		})
	}
}
