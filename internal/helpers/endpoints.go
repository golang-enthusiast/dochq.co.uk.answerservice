package helpers

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

// LoggingEndpointMiddleware - loggin for mw.
func LoggingEndpointMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				_ = logger.Log("transport_error", err, "took", time.Since(begin))
			}(time.Now())

			return next(ctx, request)
		}
	}
}

// MethodLogger - methods logger.
func MethodLogger(logger log.Logger, s string) endpoint.Middleware {
	return LoggingEndpointMiddleware(log.With(logger, "method", s))
}

// SetupEndpoint - setup endpoint.
func SetupEndpoint(handler endpoint.Endpoint, logger log.Logger, serviceName, methodName string) endpoint.Endpoint {
	result := createMiddleware()(handler)
	result = MethodLogger(logger, methodName)(result)
	return result
}

func createMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			return next(ctx, request)
		}
	}
}
