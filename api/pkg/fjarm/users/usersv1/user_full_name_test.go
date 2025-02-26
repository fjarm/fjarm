package usersv1

import (
	userservicepb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"github.com/bufbuild/protovalidate-go"
	"testing"
)

func TestUserFullName_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&userservicepb.UserFullName{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		first   string
		last    string
		wantErr bool
	}{
		"valid_from_example": {
			first:   "bella",
			last:    "hadid",
			wantErr: false,
		},
		"valid_first_and_last_name_one_character": {
			first:   "b",
			last:    "h",
			wantErr: false,
		},
		"invalid_first_name_empty": {
			first:   "",
			last:    "hadid",
			wantErr: true,
		},
		"invalid_last_name_empty": {
			first:   "bella",
			last:    "",
			wantErr: true,
		},
		"invalid_first_and_last_name_empty": {
			first:   "",
			last:    "",
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fullName := &userservicepb.UserFullName{
				FamilyName: &tc.first,
				GivenName:  &tc.last,
			}
			err = validator.Validate(fullName)
			if err != nil && !tc.wantErr {
				t.Errorf("got error = %v, wantErr = %v, input = %v, %v", err, tc.wantErr, tc.first, tc.last)
			}
		})
	}
}
