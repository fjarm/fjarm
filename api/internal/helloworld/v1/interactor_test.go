package v1

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestInteractor_GetHelloWorld(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	tests := map[string]struct {
		given string
		want  string
		err   bool
	}{
		"valid_non_empty_input": {
			given: "gleep",
			want:  "Hello World, gleep",
			err:   false,
		},
		"valid_empty_input": {
			given: "",
			want:  "Hello World",
			err:   false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			repository := newInMemoryRepository()
			domain := newInteractor(logger, repository)

			actual, err := domain.getHelloWorld(context.Background(), tc.given)
			if err != nil && !tc.err {
				t.Errorf("getHelloWorld got an unexpected error: %v", err)
			}

			if actual != tc.want {
				t.Errorf("getHelloWorld got: %v, want: %v", actual, tc.want)
			}
		})
	}
}
