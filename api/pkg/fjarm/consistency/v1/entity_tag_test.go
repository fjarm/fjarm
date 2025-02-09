package v1

import (
	consistencypb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/consistency/v1"
	"github.com/bufbuild/protovalidate-go"
	"testing"
)

func TestEntityTag_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithFailFast(),
		protovalidate.WithMessages(
			&consistencypb.EntityTag{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		entityTag string
		wantErr   bool
	}{
		"valid_non_empty_etag": {
			entityTag: "blah",
			wantErr:   false,
		},
		"invalid_empty_string_etag": {
			entityTag: "",
			wantErr:   true,
		},
		"invalid_missing_value_etag": {
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			eTag := &consistencypb.EntityTag{
				EntityTag: &tc.entityTag,
			}
			err = validator.Validate(eTag)
			if err != nil && !tc.wantErr {
				t.Errorf("got error = %v, wantErr = %v, input = %v", err, tc.wantErr, tc.entityTag)
			}
		})
	}
}
