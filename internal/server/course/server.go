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

/*
Package course provides the implementation of the gRPC
[coursepb.CourseServiceServer] interface and provides
various operations for Course models.
*/
package course

import (
	"context"

	"github.com/mikhail5545/product-service-go/internal/models"
	"github.com/mikhail5545/product-service-go/internal/services/course"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	coursepb "github.com/mikhail5545/proto-go/proto/course/v0"
)

// Server implements the gRPC [coursepb.CourseServiceServer] interface and provides
// operations for Course models. It acts as an adapter between the gRPC transport layer
// and the server-layer buusiness logic of microservice, defined in the [course.Service].
//
// For more information about underlying gRPC server, see [github.com/mikhail5545/proto-go].
type Server struct {
	coursepb.UnimplementedCourseServiceServer
	service *course.Service
}

// New creates a new Server instance.
func New(s *course.Service) *Server {
	return &Server{service: s}
}

// Get retrieves a course by their ID.
// It returns the full course object.
// If the course is not found, it returns a `NotFound` gRPC error.
func (s *Server) Get(ctx context.Context, req *coursepb.GetRequest) (*coursepb.GetResponse, error) {
	course, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return &coursepb.GetResponse{Course: types.CourseToProtobuf(course)}, nil
}

// GetReduced retrieves a course by their ID.
// It returns the reduced course object (not all fields are presented, especially it does not provide
// list of[models.CoursePart] for this course).
//
// If the course is not found, it returns a `NotFound` gRPC error.
func (s *Server) GetReduced(ctx context.Context, req *coursepb.GetReducedRequest) (*coursepb.GetReducedResponse, error) {
	course, err := s.service.GetReduced(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return &coursepb.GetReducedResponse{Course: types.CourseToProtobufListItem(course)}, nil
}

// List retrieves a paginated list of all courses.
// The response contains a list of courses
// and the total number of courses in the system.
func (s *Server) List(ctx context.Context, req *coursepb.ListRequest) (*coursepb.ListResponse, error) {
	courses, total, err := s.service.List(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbCourses []*coursepb.CourseListItem
	for _, course := range courses {
		pbCourses = append(pbCourses, types.CourseToProtobufListItem(&course))
	}

	return &coursepb.ListResponse{Courses: pbCourses, Total: total}, nil
}

// Create creates a new course record, typically in the process of direct course
// creation. It automatically creates all underlying products and populdates they're `name` and `description`
// fields from [models.Course.Name] and [models.Course.Description] if not provided.
//
// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
// It returns newly created course model with all fields.
func (s *Server) Create(ctx context.Context, req *coursepb.CreateRequest) (*coursepb.CreateResponse, error) {
	course := &models.Course{
		Name:           req.GetName(),
		Description:    req.GetDescription(),
		Topic:          req.GetTopic(),
		AccessDuration: int(req.GetAccessDuration()),
		Product: &models.Product{
			Price:       req.GetProduct().GetPrice(),
			Name:        req.GetProduct().GetName(),
			Description: req.GetProduct().GetDescription(),
		},
	}

	course, err := s.service.Create(ctx, course)
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepb.CreateResponse{Course: types.CourseToProtobuf(course)}, nil
}

// Update updates course fields that have been acually changed. All request fields
// except ID are optional, so service will update course only if at least one field
// has been updated.
//
// It populates only updated fields in the response along with the `fieldmaskpb.UpdateMask` which contains
// paths to updated fields.
func (s *Server) Update(ctx context.Context, req *coursepb.UpdateRequest) (*coursepb.UpdateResponse, error) {
	course := &models.Course{
		ID:             req.GetId(),
		Name:           req.GetName(),
		Description:    req.GetDescription(),
		Topic:          req.GetTopic(),
		AccessDuration: int(req.GetAccessDuration()),
		Product: &models.Product{
			Price: req.GetProduct().GetPrice(),
		},
	}

	updates, productUpdates, err := s.service.Update(ctx, course, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return types.CourseToProtobufUpdate(updates, productUpdates), nil
}
