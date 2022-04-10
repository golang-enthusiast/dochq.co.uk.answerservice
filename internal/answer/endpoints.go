package answer

import (
	"context"

	"dochq.co.uk.answerservice/internal/domain"
	"dochq.co.uk.answerservice/internal/helpers"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

// Endpoints collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Endpoints struct {
	CreateAnswerEndpoint     endpoint.Endpoint
	UpdateAnswerEndpoint     endpoint.Endpoint
	DeleteAnswerEndpoint     endpoint.Endpoint
	GetAnswerEndpoint        endpoint.Endpoint
	GetAnswerHistoryEndpoint endpoint.Endpoint
}

// NewEndpoint returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func NewEndpoint(service domain.AnswerService, logger log.Logger) Endpoints {
	factory := func(creator func(domain.AnswerService) endpoint.Endpoint, logKey string) endpoint.Endpoint {
		return helpers.SetupEndpoint(creator(service), logger, "AnswerEndpoints", logKey)
	}
	return Endpoints{
		CreateAnswerEndpoint:     factory(MakeCreateAnswerEndpoint, "CreateAnswer"),
		UpdateAnswerEndpoint:     factory(MakeUpdateAnswerEndpoint, "UpdateAnswer"),
		DeleteAnswerEndpoint:     factory(MakeDeleteAnswerEndpoint, "DeleteAnswer"),
		GetAnswerEndpoint:        factory(MakeGetAnswerEndpoint, "GetAnswer"),
		GetAnswerHistoryEndpoint: factory(MakeGetAnswerHistoryEndpoint, "GetAnswerHistory"),
	}
}

// MakeCreateAnswerEndpoint Impl.
func MakeCreateAnswerEndpoint(service domain.AnswerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateAnswerRequest)

		// Call the service.
		err = service.CreateAnswer(ctx, req.Answer)
		return CreateAnswerResponse{
			Err: err,
		}, nil
	}
}

// CreateAnswerRequest - request.
type CreateAnswerRequest struct {
	Answer *domain.Answer
}

// CreateAnswerResponse - response.
type CreateAnswerResponse struct {
	Err error
}

// MakeUpdateAnswerEndpoint Impl.
func MakeUpdateAnswerEndpoint(service domain.AnswerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateAnswerRequest)

		// Call the service.
		err = service.UpdateAnswer(ctx, req.Answer)
		return UpdateAnswerResponse{
			Err: err,
		}, nil
	}
}

// UpdateAnswerRequest - request.
type UpdateAnswerRequest struct {
	Answer *domain.Answer
}

// UpdateAnswerResponse - response.
type UpdateAnswerResponse struct {
	Err error
}

// MakeDeleteAnswerEndpoint Impl.
func MakeDeleteAnswerEndpoint(service domain.AnswerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteAnswerRequest)

		// Call the service.
		err = service.DeleteAnswer(ctx, req.Key)
		return DeleteAnswerResponse{
			Err: err,
		}, nil
	}
}

// DeleteAnswerRequest - request.
type DeleteAnswerRequest struct {
	Key domain.AnswerKey
}

// DeleteAnswerResponse - response.
type DeleteAnswerResponse struct {
	Err error
}

// MakeGetAnswerEndpoint Impl.
func MakeGetAnswerEndpoint(service domain.AnswerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetAnswerRequest)

		// Call the service.
		foundAnswer, err := service.GetAnswer(ctx, req.Key)
		if err != nil {
			return nil, err
		}
		return GetAnswerResponse{
			Answer: foundAnswer,
		}, nil
	}
}

// GetAnswerRequest request.
type GetAnswerRequest struct {
	Key domain.AnswerKey
}

// GetAnswerResponse response.
type GetAnswerResponse struct {
	Answer *domain.Answer
	Err    error
}

// MakeGetAnswerHistoryEndpoint Impl.
func MakeGetAnswerHistoryEndpoint(service domain.AnswerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetAnswerHistoryRequest)

		// Call the service.
		events, err := service.GetAnswerHistory(ctx, req.Key)
		if err != nil {
			return nil, err
		}
		return GetAnswerHistoryResponse{
			Events: events,
		}, nil
	}
}

// GetAnswerHistoryRequest request.
type GetAnswerHistoryRequest struct {
	Key domain.AnswerKey
}

// GetAnswerHistoryResponse response.
type GetAnswerHistoryResponse struct {
	Events []*domain.AnswerEvent
	Err    error
}

// //
//
// Error interceptors.
//
// //

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = CreateAnswerResponse{}
	_ endpoint.Failer = UpdateAnswerResponse{}
	_ endpoint.Failer = DeleteAnswerResponse{}
	_ endpoint.Failer = GetAnswerResponse{}
	_ endpoint.Failer = GetAnswerHistoryResponse{}
)

// Failed implements endpoint.Failer.
func (r CreateAnswerResponse) Failed() error { return r.Err }

// Failed implements endpoint.Failer.
func (r UpdateAnswerResponse) Failed() error { return r.Err }

// Failed implements endpoint.Failer.
func (r DeleteAnswerResponse) Failed() error { return r.Err }

// Failed implements endpoint.Failer.
func (r GetAnswerResponse) Failed() error { return r.Err }

// Failed implements endpoint.Failer.
func (r GetAnswerHistoryResponse) Failed() error { return r.Err }
