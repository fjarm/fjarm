package helloworld

import (
	rpc "buf.build/gen/go/fjarm/helloworld/grpc/go/helloworld/v1/helloworldv1grpc"
	pb "buf.build/gen/go/fjarm/helloworld/protocolbuffers/go/helloworld/v1"
	"context"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type Service struct {
	rpc.UnimplementedHelloWorldServiceServer
	repo      repository
	validator *protovalidate.Validator
}

func (svc *Service) GetHelloWorld(
	ctx context.Context,
	req *emptypb.Empty,
) (*pb.GetHelloWorldResponse, error) {
	slog.Info("received request", "rpc", "GetHelloWorld")
	res, err := svc.repo.getHelloWorld(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "request failed", "err", err, "rpc", "GetHelloWorld")
		return nil, status.Error(codes.Unknown, "failed to complete request")
	}
	err = svc.validator.Validate(res)
	if err != nil {
		slog.ErrorContext(ctx, "validation failed", "err", err, "rpc", "GetHelloWorld")
		return nil, status.Error(codes.DataLoss, "failed to validate response")
	}
	return res, nil
}

func NewService() *Service {
	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(
			&pb.GetHelloWorldResponse{},
		),
	)
	if err != nil {
		slog.Error("failed to initialize validator", "err", err)
		return nil
	}
	return &Service{
		repo:      repository{},
		validator: validator,
	}
}
