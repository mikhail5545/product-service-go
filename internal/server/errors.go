// vitainmove.com/product-service-go
// microservice for vitianmove project family
// Copyright (C) 2025  Mikhail Kulik

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
