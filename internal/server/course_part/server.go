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

	coursepartmodel "github.com/mikhail5545/product-service-go/internal/models/course_part"
	coursepart "github.com/mikhail5545/product-service-go/internal/services/course_part"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	coursepartpb "github.com/mikhail5545/proto-go/proto/course_part/v0"
	"google.golang.org/grpc"
)

// Server implements the gRPC [coursepartpb.CoursePartServiceServer] interface and provides
// operations for Course part models. It acts as an adapter between the gRPC transport layer
// and the server-layer buusiness logic of microservice, defined in the [coursepart.Service].
//
// For more information about underlying gRPC server, see [github.com/mikhail5545/proto-go].
type Server struct {
	coursepartpb.UnimplementedCoursePartServiceServer
	service coursepart.Service
}

// New creates a new Server instance.
func New(svc coursepart.Service) *Server {
	return &Server{service: svc}
}

// Register registers the course server with a gRPC server instance.
func Register(s *grpc.Server, svc coursepart.Service) {
	coursepartpb.RegisterCoursePartServiceServer(s, New(svc))
}

// Get retrieves a single published and not soft-deleted course part record.
// It attemps to retrieve MUXVideo information by calling the media service.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Get(ctx context.Context, req *coursepartpb.GetRequest) (*coursepartpb.GetResponse, error) {
	part, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.GetResponse{CoursePart: types.CoursePartToProtobuf(part)}, nil
}

// GetWithDeleted retrieves a single course part record, including soft-deleted ones.
// It attemps to retrieve MUXVideo information by calling the media service.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetWithDeleted(ctx context.Context, req *coursepartpb.GetWithDeletedRequest) (*coursepartpb.GetWithDeletedResponse, error) {
	part, err := s.service.GetWithDeleted(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.GetWithDeletedResponse{CoursePart: types.CoursePartToProtobuf(part)}, nil
}

// GetWithUnpublished retrieves a single course record, including unpublished ones.
// It attemps to retrieve MUXVideo information by calling the media service.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetWithUnpublished(ctx context.Context, req *coursepartpb.GetWithUnpublishedRequest) (*coursepartpb.GetWithUnpublishedResponse, error) {
	part, err := s.service.GetWithUnpublished(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.GetWithUnpublishedResponse{CoursePart: types.CoursePartToProtobuf(part)}, nil
}

// GetReduced retrieves a single published and not soft-deleted course part record.
// It does not populate MUXVideo details; the MUXVideo field in the returned course part struct will be nil.
// This is a lighter version of the `Get` method.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetReduced(ctx context.Context, req *coursepartpb.GetReducedRequest) (*coursepartpb.GetReducedResponse, error) {
	part, err := s.service.GetReduced(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.GetReducedResponse{CoursePart: types.CoursePartToProtobuf(part)}, nil
}

// GetWithDeletedReduced retrieves a single course part record, including soft-deleted ones.
// It does not populate MUXVideo details; the MUXVideo field in the returned course part struct will be nil.
// This is a lighter version of the `GetWithDeleted` method.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetWithDeletedReduced(ctx context.Context, req *coursepartpb.GetWithDeletedReducedRequest) (*coursepartpb.GetWithDeletedReducedResponse, error) {
	part, err := s.service.GetWithDeletedReduced(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.GetWithDeletedReducedResponse{CoursePart: types.CoursePartToProtobuf(part)}, nil
}

// GetWithUnpublishedReduced retrieves a single course part record, including unpublished ones.
// It does not populate MUXVideo details; the MUXVideo field in the returned course part struct will be nil.
// This is a lighter version of the `GetWithUnpublished` method.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetWithUnpublishedReduced(ctx context.Context, req *coursepartpb.GetWithUnpublishedReducedRequest) (*coursepartpb.GetWithUnpublishedReducedResponse, error) {
	part, err := s.service.GetWithUnpublishedReduced(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.GetWithUnpublishedReducedResponse{CoursePart: types.CoursePartToProtobuf(part)}, nil
}

// List retrieves a paginated list of all published and not soft-deleted course part records.
// It attemps to retrieve MUXVideo information for each record by calling the media service.
// Each record includes MUXVideo details.
// The response also contains the total count of such records.
func (s *Server) List(ctx context.Context, req *coursepartpb.ListRequest) (*coursepartpb.ListResponse, error) {
	parts, total, err := s.service.List(ctx, req.GetCourseId(), int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbparts []*coursepartpb.CoursePart
	for _, p := range parts {
		pbparts = append(pbparts, types.CoursePartToProtobuf(&p))
	}
	return &coursepartpb.ListResponse{CourseParts: pbparts, Total: total}, nil
}

// ListDeleted retrieves a paginated list of all soft-deleted course part records.
// It does not populate MUXVideo details; the MUXVideo field in the returned course part structs will be nil.
// The response also contains the total count of such records.
func (s *Server) ListDeleted(ctx context.Context, req *coursepartpb.ListDeletedRequest) (*coursepartpb.ListDeletedResponse, error) {
	parts, total, err := s.service.ListDeleted(ctx, req.GetCourseId(), int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbparts []*coursepartpb.CoursePart
	for _, p := range parts {
		pbparts = append(pbparts, types.CoursePartToProtobuf(&p))
	}
	return &coursepartpb.ListDeletedResponse{CourseParts: pbparts, Total: total}, nil
}

// ListUnpublished retrieves a paginated list of all unpublished course part records.
// It does not populate MUXVideo details; the MUXVideo field in the returned course part structs will be nil.
// The response also contains the total count of such records.
func (s *Server) ListUnpublished(ctx context.Context, req *coursepartpb.ListUnpublishedRequest) (*coursepartpb.ListUnpublishedResponse, error) {
	parts, total, err := s.service.ListUnpublished(ctx, req.GetCourseId(), int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbparts []*coursepartpb.CoursePart
	for _, p := range parts {
		pbparts = append(pbparts, types.CoursePartToProtobuf(&p))
	}
	return &coursepartpb.ListUnpublishedResponse{CourseParts: pbparts, Total: total}, nil
}

// Create creates a new CoursePart.
// It is created in an unpublished state.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (s *Server) Create(ctx context.Context, req *coursepartpb.CreateRequest) (*coursepartpb.CreateResponse, error) {
	createReq := &coursepartmodel.CreateRequest{
		CourseID:         req.GetCourseId(),
		Name:             req.GetName(),
		ShortDescription: req.GetShortDescription(),
		Number:           int(req.GetNumber()),
	}
	res, err := s.service.Create(ctx, createReq)
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.CreateResponse{Id: res.ID, CourseId: res.CourseID}, nil
}

// Publish makes a course part available in the catalog by setting `Published` to true.
// It will fail if parent course is not published.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID or parent course is not published.
func (s *Server) Publish(ctx context.Context, req *coursepartpb.PublishRequest) (*coursepartpb.PublishResponse, error) {
	if err := s.service.Publish(ctx, req.GetId()); err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.PublishResponse{Id: req.GetId()}, nil
}

// Unpublish archives a course part from the catalog
// by setting its `Published` field to false.
//
// Returns a `NotFound` gRPC error if the course or its parts are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Unpublish(ctx context.Context, req *coursepartpb.UnpublishRequest) (*coursepartpb.UnpublishResponse, error) {
	if err := s.service.Unpublish(ctx, req.GetId()); err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.UnpublishResponse{Id: req.GetId()}, nil
}

// Update performs a partial update of a course part.
// At least one field must be provided for an update to occur.
// The response contains only the fields that were actually changed.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (s *Server) Update(ctx context.Context, req *coursepartpb.UpdateRequest) (*coursepartpb.UpdateResponse, error) {
	updateReq := &coursepartmodel.UpdateRequest{
		ID:               req.GetId(),
		CourseID:         req.GetCourseId(),
		Name:             req.Name,
		ShortDescription: req.ShortDescription,
		LongDescription:  req.LongDescription,
		Tags:             req.GetTags(),
	}
	n := int(req.GetNumber())
	updateReq.Number = &n
	res, err := s.service.Update(ctx, updateReq)
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return types.CoursePartToProtobufUpdate(&coursepartpb.UpdateResponse{Id: req.GetId(), CourseId: req.GetCourseId()}, res), nil
}

// AddVideo associates MUXVideo with the course part record by
// setting `MUXVideoID` field value in the course part record.
// It will populate/update `MUXVideoID` field only if new value is different from the previous one.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns a `InvalidArgument` gRPC error if any of the provided IDs is not a valid UUID.
func (s *Server) AddVideo(ctx context.Context, req *coursepartpb.AddVideoRequest) (*coursepartpb.AddVideoResponse, error) {
	addVideoReq := &coursepartmodel.AddVideoRequest{
		ID:         req.GetId(),
		MUXVideoID: req.GetMuxVideoId(),
	}
	_, err := s.service.AddVideo(ctx, addVideoReq)
	if err != nil {
		return nil, err
	}
	return &coursepartpb.AddVideoResponse{MuxVideoId: req.MuxVideoId}, nil
}

// Delete performs a soft-delete on a course part.
// It also unpublishes is, requiring manual re-publishing after restoration.
//
// Returns a `NotFound` gRPC error if course part is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Delete(ctx context.Context, req *coursepartpb.DeleteRequest) (*coursepartpb.DeleteResponse, error) {
	if err := s.service.Delete(ctx, req.GetId()); err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.DeleteResponse{Id: req.GetId()}, nil
}

// DeletePermanent permanently deletes a course part.
//
// Returns a `NotFound` gRPC error if course part is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) DeletePermanent(ctx context.Context, req *coursepartpb.DeletePermanentRequest) (*coursepartpb.DeletePermanentResponse, error) {
	if err := s.service.DeletePermanent(ctx, req.GetId()); err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.DeletePermanentResponse{Id: req.GetId()}, nil
}

// Restore restores a soft-deleted course part.
// The restored record is not automatically published and must be published manually.
//
// Returns a `NotFound` gRPC error if course part is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Restore(ctx context.Context, req *coursepartpb.RestoreRequest) (*coursepartpb.RestoreResponse, error) {
	if err := s.service.Restore(ctx, req.GetId()); err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &coursepartpb.RestoreResponse{Id: req.GetId()}, nil
}
