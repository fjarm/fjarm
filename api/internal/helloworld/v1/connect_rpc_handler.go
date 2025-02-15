package v1

import (
	"buf.build/gen/go/fjarm/fjarm/connectrpc/go/fjarm/helloworld/v1/helloworldv1connect"
	pb "buf.build/gen/go/fjarm/fjarm/protocolbuffers/go/fjarm/helloworld/v1"
	"connectrpc.com/connect"
	"context"
	"github.com/bufbuild/protovalidate-go"
	"github.com/fjarm/fjarm/api/internal/logkeys"
	"github.com/fjarm/fjarm/api/internal/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

const connectRPCHandlerTag = "connect_rpc_handler"

// ConnectRPCHandler defines a ConnectRPC handler for the `fjarm.helloworld.v1.HelloWorldService` service.
type ConnectRPCHandler struct {
	domain    getHelloWorlder
	logger    *slog.Logger
	validator protovalidate.Validator
}

// GetHelloWorld implements the similarly named RPC defined in the `fjarm.helloworld.v1.HelloWorldService` service.
func (h *ConnectRPCHandler) GetHelloWorld(
	ctx context.Context,
	req *connect.Request[pb.GetHelloWorldRequest],
) (*connect.Response[pb.GetHelloWorldResponse], error) {
	logger := h.logger.With(
		slog.String(logkeys.Rpc, helloworldv1connect.HelloWorldServiceGetHelloWorldProcedure),
		slog.String(tracing.RequestIDKey, req.Header().Get(tracing.RequestIDKey)),
	)

	logger.InfoContext(
		ctx,
		"received request",
		slog.String(logkeys.Request, req.Msg.String()),
	)

	err := h.validator.Validate(req.Msg)
	if err != nil {
		logger.WarnContext(ctx, "failed to validate request", slog.Any(logkeys.Raw, req.Msg))
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	msg, err := h.domain.getHelloWorld(ctx, req.Msg.GetInput().GetInput())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	res := connect.NewResponse(&pb.GetHelloWorldResponse{
		Status: status.New(codes.OK, codes.OK.String()).Proto(),
		Output: &pb.HelloWorldOutput{
			Output: &msg,
		},
	})
	return res, nil
}

// NewConnectRPCHandler creates a concrete helloworld
func NewConnectRPCHandler(l *slog.Logger) *ConnectRPCHandler {
	logger := l.With(
		slog.String(logkeys.Tag, connectRPCHandlerTag),
	)

	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithFailFast(),
		protovalidate.WithMessages(
			&pb.GetHelloWorldRequest{},
			&pb.GetHelloWorldResponse{},
		),
	)

	if err != nil {
		logger.Error("failed to create message validator", slog.Any(logkeys.Err, err))
		return nil
	}

	repo := newInMemoryRepository()
	dom := newInteractor(repo)
	han := ConnectRPCHandler{
		domain:    dom,
		logger:    logger,
		validator: validator,
	}
	return &han
}
