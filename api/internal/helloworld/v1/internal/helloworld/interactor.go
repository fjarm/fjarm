package helloworld

import (
	"context"
	"fmt"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
	"log/slog"
)

const interactorTag = "interactor"

type getHelloWorldMessageer interface {
	getHelloWorldMessage(ctx context.Context) (string, error)
}

type interactor struct {
	logger *slog.Logger
	repo   getHelloWorldMessageer
}

func newInteractor(logger *slog.Logger, repo getHelloWorldMessageer) *interactor {
	return &interactor{
		logger: logger,
		repo:   repo,
	}
}

func (svc *interactor) getHelloWorld(ctx context.Context, input string) (string, error) {
	logger := svc.logger.With(
		slog.String(logkeys.Tag, interactorTag),
		slog.Any(tracing.RequestIDKey, ctx.Value(tracing.RequestIDKey)),
	)

	msg, err := svc.repo.getHelloWorldMessage(ctx)
	if err != nil {
		logger.WarnContext(
			ctx,
			"getHelloWorld failed to request message from repository",
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
