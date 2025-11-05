// github.com/mikhail5545/product-service-go
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

package errors

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mikhail5545/product-service-go/internal/services/course"
	coursepart "github.com/mikhail5545/product-service-go/internal/services/course_part"
	physicalgood "github.com/mikhail5545/product-service-go/internal/services/physical_good"
	"github.com/mikhail5545/product-service-go/internal/services/product"
	"github.com/mikhail5545/product-service-go/internal/services/seminar"
	trainingsession "github.com/mikhail5545/product-service-go/internal/services/training_session"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceError interface {
	GetCode() int
	Error() string
}

// HTTPErrorHandler is a custom error handler for Echo.
func HTTPErrorHandler(err error, c echo.Context) {
	var se ServiceError
	if errors.As(err, &se) {
		c.JSON(se.GetCode(), map[string]string{"error": se.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

// ToGRPCError converts a service layer error into a gRPC status error.
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	var productErr *product.Error
	if errors.As(err, &productErr) {
		switch productErr.GetCode() {
		case http.StatusBadRequest:
			return status.Errorf(codes.InvalidArgument, "Product service error occurred: %s", productErr.Error())
		case http.StatusNotFound:
			return status.Errorf(codes.NotFound, "Product service error occurred: %s", productErr.Error())
		}
	}

	var tsError *trainingsession.Error
	if errors.As(err, &tsError) {
		switch tsError.GetCode() {
		case http.StatusBadRequest:
			return status.Errorf(codes.InvalidArgument, "Training session service error occurred: %s", tsError.Error())
		case http.StatusNotFound:
			return status.Errorf(codes.NotFound, "Training session service error occurred: %s", tsError.Error())
		}
	}

	var courseErr *course.Error
	if errors.As(err, &courseErr) {
		switch courseErr.GetCode() {
		case http.StatusBadRequest:
			return status.Errorf(codes.InvalidArgument, "Course service error occurred: %s", courseErr.Error())
		case http.StatusNotFound:
			return status.Errorf(codes.NotFound, "Course service error occurred: %s", courseErr.Error())
		}
	}

	var semianrErr *seminar.Error
	if errors.As(err, &semianrErr) {
		switch semianrErr.GetCode() {
		case http.StatusBadRequest:
			return status.Errorf(codes.InvalidArgument, "Seminar service error occurred: %s", semianrErr.Error())
		case http.StatusNotFound:
			return status.Errorf(codes.NotFound, "Seminar service error occurred: %s", semianrErr.Error())
		}
	}

	var coursePartErr *coursepart.Error
	if errors.As(err, &coursePartErr) {
		switch coursePartErr.GetCode() {
		case http.StatusBadRequest:
			return status.Errorf(codes.InvalidArgument, "Course part service error occurred: %s", coursePartErr.Error())
		case http.StatusNotFound:
			return status.Errorf(codes.NotFound, "Course part service error occurred: %s", coursePartErr.Error())
		}
	}

	var physicalGoodErr *physicalgood.Error
	if errors.As(err, &physicalGoodErr) {
		switch physicalGoodErr.GetCode() {
		case http.StatusBadRequest:
			return status.Errorf(codes.InvalidArgument, "Physical good service error occurred: %s", physicalGoodErr.Error())
		case http.StatusNotFound:
			return status.Errorf(codes.NotFound, "Physical good service error occurred: %s", physicalGoodErr.Error())
		}
	}

	return status.Errorf(codes.Internal, "an internal error occurred: %v", err)
}
