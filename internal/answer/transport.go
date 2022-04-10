package answer

import (
	"context"

	apiv1 "dochq.co.uk.answerservice/api/generated/dochq.co.uk/answerserviceapi/v1"
	"dochq.co.uk.answerservice/internal/domain"
	errors "dochq.co.uk.answerservice/internal/error"
	"dochq.co.uk.answerservice/internal/helpers"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
)

type grpcServer struct {
	createAnswer grpctransport.Handler
	updateAnswer grpctransport.Handler
	deleteAnswer grpctransport.Handler
	getAnswer    grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(endpoints Endpoints, logger log.Logger) apiv1.AnswerServiceServer {
	options := helpers.SetupServerOptions(logger)
	return &grpcServer{
		createAnswer: grpctransport.NewServer(
			endpoints.CreateAnswerEndpoint,
			decodeCreateAnswerRequest,
			encodeCreateAnswerResponse,
			options...,
		),
		updateAnswer: grpctransport.NewServer(
			endpoints.UpdateAnswerEndpoint,
			decodeUpdateAnswerRequest,
			encodeUpdateAnswerResponse,
			options...,
		),
		deleteAnswer: grpctransport.NewServer(
			endpoints.DeleteAnswerEndpoint,
			decodeDeleteAnswerRequest,
			encodeDeleteAnswerResponse,
			options...,
		),
		getAnswer: grpctransport.NewServer(
			endpoints.GetAnswerEndpoint,
			decodeGetAnswerRequest,
			encodeGetAnswerResponse,
			options...,
		),
	}
}

// CreateAnswer Impl.
func (s *grpcServer) CreateAnswer(ctx context.Context, req *apiv1.CreateAnswerRequest) (*apiv1.CreateAnswerResponse, error) {
	rep, err := helpers.ServeGrpc(ctx, req, s.createAnswer)
	if err != nil {
		return nil, err
	}
	return rep.(*apiv1.CreateAnswerResponse), nil
}

func decodeCreateAnswerRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*apiv1.CreateAnswerRequest)
	decodedAnswer, err := decodeAnswer(req.Answer)
	if err != nil {
		return CreateAnswerRequest{}, err
	}
	return CreateAnswerRequest{
		Answer: decodedAnswer,
	}, nil
}

func encodeCreateAnswerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(CreateAnswerResponse)
	if resp.Err != nil {
		return &apiv1.CreateAnswerResponse{}, resp.Err
	}
	return &apiv1.CreateAnswerResponse{}, nil
}

// UpdateAnswer Impl.
func (s *grpcServer) UpdateAnswer(ctx context.Context, req *apiv1.UpdateAnswerRequest) (*apiv1.UpdateAnswerResponse, error) {
	rep, err := helpers.ServeGrpc(ctx, req, s.updateAnswer)
	if err != nil {
		return nil, err
	}
	return rep.(*apiv1.UpdateAnswerResponse), nil
}

func decodeUpdateAnswerRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*apiv1.UpdateAnswerRequest)
	decodedAnswer, err := decodeAnswer(req.Answer)
	if err != nil {
		return UpdateAnswerRequest{}, err
	}
	return UpdateAnswerRequest{
		Answer: decodedAnswer,
	}, nil
}

func encodeUpdateAnswerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(UpdateAnswerResponse)
	if resp.Err != nil {
		return &apiv1.UpdateAnswerResponse{}, resp.Err
	}
	return &apiv1.UpdateAnswerResponse{}, nil
}

// DeleteAnswer Impl.
func (s *grpcServer) DeleteAnswer(ctx context.Context, req *apiv1.DeleteAnswerRequest) (*apiv1.DeleteAnswerResponse, error) {
	rep, err := helpers.ServeGrpc(ctx, req, s.deleteAnswer)
	if err != nil {
		return nil, err
	}
	return rep.(*apiv1.DeleteAnswerResponse), nil
}

func decodeDeleteAnswerRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*apiv1.DeleteAnswerRequest)
	return DeleteAnswerRequest{
		Key: domain.AnswerKey(req.Key),
	}, nil
}

func encodeDeleteAnswerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(DeleteAnswerResponse)
	if resp.Err != nil {
		return &apiv1.DeleteAnswerResponse{}, resp.Err
	}
	return &apiv1.DeleteAnswerResponse{}, nil
}

// GetAnswer Impl.
func (s *grpcServer) GetAnswer(ctx context.Context, req *apiv1.GetAnswerRequest) (*apiv1.GetAnswerResponse, error) {
	rep, err := helpers.ServeGrpc(ctx, req, s.getAnswer)
	if err != nil {
		return nil, err
	}
	return rep.(*apiv1.GetAnswerResponse), nil
}

func decodeGetAnswerRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*apiv1.GetAnswerRequest)
	return GetAnswerRequest{
		Key: domain.AnswerKey(req.Key),
	}, nil
}

func encodeGetAnswerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(GetAnswerResponse)
	if resp.Err != nil {
		return &apiv1.GetAnswerResponse{}, resp.Err
	}
	encodedAnswer, err := encodeAnswer(resp.Answer)
	if err != nil {
		return &apiv1.GetAnswerResponse{}, err
	}
	return &apiv1.GetAnswerResponse{
		Answer: encodedAnswer,
	}, nil
}

func decodeAnswer(answer *apiv1.Answer) (*domain.Answer, error) {
	if answer == nil {
		return nil, errors.NewErrInternal("Cannot decode nil value")
	}
	return &domain.Answer{
		Key:   domain.AnswerKey(answer.Key),
		Value: domain.AnswerValue(answer.Value),
	}, nil
}

func encodeAnswer(answer *domain.Answer) (*apiv1.Answer, error) {
	if answer == nil {
		return nil, errors.NewErrInternal("Cannot encode nil value")
	}
	return &apiv1.Answer{
		Key:   string(answer.Key),
		Value: string(answer.Value),
	}, nil
}
