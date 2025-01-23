package v1

import (
	"context"
	"testing"
)

func TestInMemoryRepository_GetHelloWorldMessage(t *testing.T) {
	tests := map[string]struct {
		want string
		err  bool
	}{
		"valid": {
			want: "Hello World",
			err:  false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			repo := newInMemoryRepository()
			actual, err := repo.getHelloWorldMessage(context.Background())
			if err != nil && !tc.err {
				t.Errorf("getHelloWorldMesssage got an unexpected error: %v", err)
			}

			if actual != tc.want {
				t.Errorf("getHelloWorldMessage got: %v, want: %v", actual, tc.want)
			}
		})
	}
}
