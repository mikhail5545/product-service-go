package server

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"vitainmove.com/product-service-go/internal/services"
)

// toGRPCError converts a service layer error into a gRPC status error.
func toGRPCError(err error) error {
	if err == nil {
		return nil
	}

	var serviceErr *services.ProductServiceError
	if errors.As(err, &serviceErr) {
		switch serviceErr.GetCode() {
		case http.StatusBadRequest:
			return status.Errorf(codes.InvalidArgument, serviceErr.Error())
		case http.StatusNotFound:
			return status.Errorf(codes.NotFound, serviceErr.Error())
		}
	}

	return status.Errorf(codes.Internal, "an internal error occurred: %v", err)
}
