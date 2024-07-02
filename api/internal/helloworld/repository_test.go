package helloworld

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestRepository_GetHelloWorld(t *testing.T) {
	tests := map[string]struct {
		wantCode   string
		wantOutput string
	}{
		"valid": {
			wantCode:   codes.OK.String(),
			wantOutput: "Hello World",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			repo := NewRepository()
			actual, err := repo.GetHelloWorld(context.Background(), &emptypb.Empty{})
			if err != nil {
				t.Errorf("GetHelloWorld got unexpected error: %v", err)
			}
			if tc.wantCode != actual.Status {
				t.Errorf("want %s, got %s", tc.wantCode, actual.Status)
			}
		})
	}
}
