package v1

import (
	idempotencypb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/idempotency/v1"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestIdempotencyKey_Validation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithFailFast(),
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&idempotencypb.IdempotencyKey{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		idempotencyKey string
		timestamp      time.Time
		wantErr        bool
	}{
		"valid_non_empty_uuid_idempotency_key": {
			idempotencyKey: "123e4567-e89b-12d3-a456-426614174000",
			timestamp:      time.Now(),
			wantErr:        false,
		},
		"valid_past_timestamp_idempotency_key": {
			idempotencyKey: "123e4567-e89b-12d3-a456-426614174000",
			timestamp:      time.Now().Add(-60 * time.Second),
			wantErr:        false,
		},
		"valid_future_timestamp_idempotency_key": {
			idempotencyKey: "123e4567-e89b-12d3-a456-426614174000",
			timestamp:      time.Now().Add(60 * time.Second),
			wantErr:        false,
		},
		"invalid_non_uuid_idempotency_key": {
			idempotencyKey: "1abcd-2efgh-3ijkl",
			timestamp:      time.Now(),
			wantErr:        true,
		},
		"invalid_empty_string_idempotency_key": {
			idempotencyKey: "",
			timestamp:      time.Now(),
			wantErr:        true,
		},
		"invalid_three_character_string_idempotency_key": {
			idempotencyKey: "124",
			timestamp:      time.Now(),
			wantErr:        true,
		},
		"invalid_missing_value_idempotency_key": {
			timestamp: time.Now(),
			wantErr:   true,
		},
		"invalid_missing_value_timestamp_idempotency_key": {
			idempotencyKey: "123e4567-e89b-12d3-a456-426614174000",
			wantErr:        true,
		},
		"invalid_missing_value_idempotency_key_and_timestamp_idempotency_key": {
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ik := &idempotencypb.IdempotencyKey{
				IdempotencyKey: &tc.idempotencyKey,
				Timestamp:      timestamppb.New(tc.timestamp),
			}
			err = validator.Validate(ik)
			if err != nil && !tc.wantErr {
				t.Errorf(
					"got error = %v, wantErr = %v, input = %v, %v",
					err,
					tc.wantErr,
					tc.idempotencyKey,
					tc.timestamp,
				)
			}
		})
	}
}
