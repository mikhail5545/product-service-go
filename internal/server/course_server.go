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
	"context"

	"vitainmove.com/product-service-go/internal/models"
	"vitainmove.com/product-service-go/internal/services"
	"vitainmove.com/product-service-go/internal/utils"
	coursepb "vitainmove.com/product-service-go/proto/course/v0"
)

type CourseServer struct {
	coursepb.UnimplementedCourseServiceServer
	courseService *services.CourseService
}

func NewCourseServer(cs *services.CourseService) *CourseServer {
	return &CourseServer{courseService: cs}
}

func (s *CourseServer) GetCourse(ctx context.Context, req *coursepb.GetCourseRequest) (*coursepb.GetCourseResponse, error) {
	course, err := s.courseService.GetCourse(ctx, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &coursepb.GetCourseResponse{Course: utils.ConvertToProtobufCourse(course)}, nil
}

func (s *CourseServer) ListCourses(ctx context.Context, req *coursepb.ListCoursesRequest) (*coursepb.ListCoursesResponse, error) {
	courses, total, err := s.courseService.GetCourses(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, toGRPCError(err)
	}
	var pbCourses []*coursepb.CourseListItem
	for _, course := range courses {
		pbCourses = append(pbCourses, utils.ConvertToProtobufListCourseItem(&course))
	}

	return &coursepb.ListCoursesResponse{Courses: pbCourses, Total: total}, nil
}

func (s *CourseServer) CreateCourse(ctx context.Context, req *coursepb.CreateCourseRequest) (*coursepb.CreateCourseResponse, error) {
	// The request is now flat and specific, not a nested object
	course := &models.Course{
		Name:           req.GetName(),
		Description:    req.GetDescription(),
		Topic:          req.GetTopic(),
		AccessDuration: int(req.GetAccessDuration()),
		Product: &models.Product{
			// Name and Description for the product can be derived in the service layer
			Price: req.GetProductInfo().GetPrice(),
		},
	}

	course, err := s.courseService.CreateCourse(ctx, course)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &coursepb.CreateCourseResponse{Course: utils.ConvertToProtobufCourse(course)}, nil
}

func (s *CourseServer) UpdateCourse(ctx context.Context, req *coursepb.UpdateCourseRequest) (*coursepb.UpdateCourseResponse, error) {
	// The request now contains optional fields for partial updates.
	// The service layer will handle which fields to update.
	course := &models.Course{
		ID:             req.GetId(),
		Name:           req.GetName(),
		Description:    req.GetDescription(),
		Topic:          req.GetTopic(),
		AccessDuration: int(req.GetAccessDuration()),
		Product: &models.Product{
			Price: req.GetProductInfo().GetPrice(),
		},
	}

	course, err := s.courseService.UpdateCourse(ctx, course)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &coursepb.UpdateCourseResponse{Course: utils.ConvertToProtobufCourse(course)}, nil
}
