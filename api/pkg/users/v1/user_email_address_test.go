package v1

import (
	userservicepb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/users/v1"
	"github.com/bufbuild/protovalidate-go"
	"testing"
)

func TestUserEmailAddress_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&userservicepb.UserEmailAddress{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		emailAddress string
		wantErr      bool
	}{
		"valid_common_email_address": {
			emailAddress: "gleep@gmail.com",
			wantErr:      false,
		},
		"valid_uncommon_tld_email_address": {
			emailAddress: "bleep@b.xyz",
			wantErr:      false,
		},
		"valid_uncommon_domain_email_address": {
			emailAddress: "gleep@glop.com",
			wantErr:      false,
		},
		"valid_non_existent_tld_email_address": {
			emailAddress: "hello@cool.thisisnotavalidtld",
			wantErr:      false,
		},
		"invalid_empty_string_email_address": {
			emailAddress: "",
			wantErr:      true,
		},
		"invalid_non_email_address_string_email_address": {
			emailAddress: "bleep",
			wantErr:      true,
		},
		"invalid_handle_email_address": {
			emailAddress: "@gmail.com",
			wantErr:      true,
		},
		"invalid_space_as_handle_email_address": {
			emailAddress: " @gmail.com",
			wantErr:      true,
		},
		"invalid_non_alphabet_in_username_email_address": {
			emailAddress: "-@gmail.com",
			wantErr:      true,
		},
		"invalid_www_in_tld_email_address": {
			emailAddress: "gloop@www.gmail.com",
			wantErr:      true,
		},
		"invalid_no_value_email_address": {
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			emailAddress := &userservicepb.UserEmailAddress{
				EmailAddress: &tc.emailAddress,
			}
			err = validator.Validate(emailAddress)
			if err != nil && !tc.wantErr {
				t.Errorf("got error = %v, wantErr = %v, input = %v", err, tc.wantErr, tc.emailAddress)
			}
		})
	}
}
