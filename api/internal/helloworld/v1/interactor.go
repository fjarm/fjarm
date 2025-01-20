package v1

import (
	"context"
	"fmt"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"log/slog"
)

const interactorTag = "interactor"

type getHelloWorldMessageer interface {
	getHelloWorldMessage(ctx context.Context) (string, error)
}

type interactor struct {
	repo getHelloWorldMessageer
}

func newInteractor(repo getHelloWorldMessageer) *interactor {
	return &interactor{
		repo: repo,
	}
}

func (svc *interactor) getHelloWorld(ctx context.Context, input string) (string, error) {
	slog.InfoContext(
		ctx,
		"called getHelloWorld",
		slog.String(logkeys.Tag, interactorTag),
		slog.String("input", input),
	)

	msg, err := svc.repo.getHelloWorldMessage(ctx)
	if err != nil {
		slog.WarnContext(
			ctx,
			"getHelloWorld failed to request message from repository",
			slog.String(logkeys.Tag, interactorTag),
			slog.String(logkeys.Err, err.Error()),
		)
		return "", err
	}

	if input == "" {
		return msg, nil
	}

	res := fmt.Sprintf("%s, %s", msg, input)
	return res, nil
}
