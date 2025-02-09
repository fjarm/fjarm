package v1

import (
	pb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/helloworld/v1"
	"github.com/bufbuild/protovalidate-go"
	"testing"
)

func TestGetHelloWorldResponse_OutputValidation(t *testing.T) {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&pb.GetHelloWorldResponse{},
		),
	)
	if err != nil {
		t.Errorf("failed to initialize validator: %v", err)
	}
	tests := map[string]struct {
		output  string
		wantErr bool
	}{
		"valid_one_character_string": {
			output:  "a",
			wantErr: false,
		},
		"valid_multi_character_string": {
			output:  "hello world",
			wantErr: false,
		},
		"invalid_zero_character_string": {
			output:  "",
			wantErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp := pb.GetHelloWorldResponse{
				Status: nil,
				Output: &pb.HelloWorldOutput{
					Output: &tc.output,
				},
			}
			err = validator.Validate(&resp)
			if (err != nil) != tc.wantErr {
				t.Errorf("got error = %v, wantErr %v, input = %v", err, tc.wantErr, tc.output)
			}
		})
	}
}
