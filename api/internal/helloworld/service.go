package helloworld

import (
	"context"
	"github.com/bufbuild/protovalidate-go"
	pb "github.com/fjarm/fjarm/api/pkg/helloworld/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type Service struct {
	pb.UnimplementedHelloWorldServiceServer
	repo      repository
	validator *protovalidate.Validator
}

func (svc *Service) GetHelloWorld(
	ctx context.Context,
	req *emptypb.Empty,
) (*pb.GetHelloWorldResponse, error) {
	res, err := svc.repo.getHelloWorld(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "request failed", "err", err, "rpc", "GetHelloWorld")
		return nil, nil
	}
	err = svc.validator.Validate(res)
	if err != nil {
		slog.ErrorContext(ctx, "validation failed", "err", err, "rpc", "GetHelloWorld")
	}
	return res, nil
}

func NewService() *Service {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(
			&pb.GetHelloWorldRequest{},
		),
	)
	if err != nil {
		slog.Error("Failed to initialize validator", "err", err)
	}
	return &Service{
		repo:      repository{},
		validator: validator,
	}
}
