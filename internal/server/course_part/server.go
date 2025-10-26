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
Package coursepart provides the implementation of the gRPC
[coursepartpb.CoursePartServiceServer] interface and provides
various operations for Course part models.
*/
package coursepart

import (
	"context"

	coursepart "github.com/mikhail5545/product-service-go/internal/services/course_part"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	coursepartpb "github.com/mikhail5545/proto-go/proto/course_part/v0"
)

// Server implements the gRPC [coursepartpb.CoursePartServiceServer] interface and provides
// operations for Course part models. It acts as an adapter between the gRPC transport layer
// and the server-layer buusiness logic of microservice, defined in the [coursepart.Service].
//
// For more information about underlying gRPC server, see [github.com/mikhail5545/proto-go].
type Server struct {
	coursepartpb.UnimplementedCoursePartServiceServer
	service *coursepart.Service
}

// New creates a new Server instance.
func New(s *coursepart.Service) *Server {
	return &Server{service: s}
}

// Get retrieves a course part by their ID.
// It returns the full course part object.
// If the course part is not found, it returns a `NotFound` gRPC error.
func (s *Server) Get(ctx context.Context, req *coursepartpb.GetRequest) (*coursepartpb.GetPartResponse, error) {
	part, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.GetPartResponse{CoursePart: types.ConvertToProtobufCoursePart(part)}, nil
}

// Unimplemented
func (s *Server) AddVideo(ctx context.Context, req *coursepartpb.AddVideoRequest) (*coursepartpb.AddVideoResponse, error) {
	return &coursepartpb.AddVideoResponse{}, nil
}
