package error

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCErrorEncoder error encoder
func GRPCErrorEncoder(err error) error {
	if err == nil {
		return err
	}
	switch err.(type) {
	case *ErrInvalidArgument:
		return status.Error(codes.InvalidArgument, err.Error())
	case *ErrAlreadyExist:
		return status.Error(codes.AlreadyExists, err.Error())
	case *ErrNotFound:
		return status.Error(codes.NotFound, err.Error())
	case *ErrFailedPrecondition:
		return status.Error(codes.FailedPrecondition, err.Error())
	case *ErrInternal:
		return status.Error(codes.Internal, err.Error())
	case *ErrUnauthorized:
		return status.Error(codes.Unauthenticated, err.Error())
	case *ErrPermissionDenied:
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}
