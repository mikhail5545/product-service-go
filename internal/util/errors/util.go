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

// Package errors provides utility handlers for internal errors and their convertions to external ones (http, gRPC).
package errors

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mikhail5545/product-service-go/internal/services/course"
	coursepart "github.com/mikhail5545/product-service-go/internal/services/course_part"
	imageservice "github.com/mikhail5545/product-service-go/internal/services/image"
	imagemanager "github.com/mikhail5545/product-service-go/internal/services/image_manager"
	physicalgood "github.com/mikhail5545/product-service-go/internal/services/physical_good"
	"github.com/mikhail5545/product-service-go/internal/services/product"
	"github.com/mikhail5545/product-service-go/internal/services/seminar"
	trainingsession "github.com/mikhail5545/product-service-go/internal/services/training_session"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceError interface {
	// GetCode is deprecated. Use errors.Is with sentinel errors instead.
	GetCode() int
	Error() string
}

// HTTPErrorHandler is a custom error handler for Echo.
func HTTPErrorHandler(err error, c echo.Context) {
	// Handle specific sentinel errors first
	if errors.Is(err, seminar.ErrInvalidArgument) || errors.Is(err, course.ErrInvalidArgument) || errors.Is(err, trainingsession.ErrInvalidArgument) || errors.Is(err, physicalgood.ErrInvalidArgument) {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if errors.Is(err, seminar.ErrNotFound) || errors.Is(err, course.ErrNotFound) || errors.Is(err, trainingsession.ErrNotFound) || errors.Is(err, physicalgood.ErrNotFound) {
		c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	if errors.Is(err, seminar.ErrImageLimitExceeded) {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Fallback for older error types
	var se ServiceError
	if errors.As(err, &se) {
		c.JSON(se.GetCode(), map[string]string{"error": se.Error()})
		return
	}

	// Default to internal server error
	c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

// HandleServiceError converts a service layer error into a gRPC status error.
func HandleServiceError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, seminar.ErrInvalidArgument) ||
		errors.Is(err, course.ErrInvalidArgument) ||
		errors.Is(err, trainingsession.ErrInvalidArgument) ||
		errors.Is(err, physicalgood.ErrInvalidArgument) ||
		errors.Is(err, product.ErrInvalidArgument) ||
		errors.Is(err, coursepart.ErrInvalidArgument) ||
		errors.Is(err, imageservice.ErrUnknownOwner) ||
		errors.Is(err, imagemanager.ErrImageLimitExceeded) ||
		errors.Is(err, imagemanager.ErrInvalidArgument) {
		return status.Errorf(codes.InvalidArgument, "Invalid argument: %s", err.Error())
	}
	if errors.Is(err, seminar.ErrNotFound) ||
		errors.Is(err, course.ErrNotFound) ||
		errors.Is(err, trainingsession.ErrNotFound) ||
		errors.Is(err, physicalgood.ErrNotFound) ||
		errors.Is(err, product.ErrNotFound) ||
		errors.Is(err, coursepart.ErrNotFound) ||
		errors.Is(err, imagemanager.ErrOwnerNotFound) ||
		errors.Is(err, imagemanager.ErrOwnersNotFound) ||
		errors.Is(err, imagemanager.ErrImageNotFoundOnOwner) {
		return status.Errorf(codes.NotFound, "Not found: %s", err.Error())
	}

	return status.Errorf(codes.Internal, "an internal error occurred: %v", err)
}
