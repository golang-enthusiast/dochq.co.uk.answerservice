package helpers

import (
	"context"

	errors "dochq.co.uk.answerservice/internal/error"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
)

// ServeGrpc - wraps the error
func ServeGrpc(ctx context.Context, req interface{}, handler grpc.Handler) (interface{}, error) {
	_, resp, err := handler.ServeGRPC(ctx, req)
	return resp, errors.GRPCErrorEncoder(err)
}

// SetupServerOptions - setups server options.
func SetupServerOptions(logger log.Logger) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		grpc.ServerBefore(jwt.GRPCToContext()),
	}
}
