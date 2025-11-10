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

	coursemodel "github.com/mikhail5545/product-service-go/internal/models/course"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	courseservice "github.com/mikhail5545/product-service-go/internal/services/course"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	coursepb "github.com/mikhail5545/proto-go/proto/course/v0"
	"google.golang.org/grpc"
)

// Server implements the gRPC [coursepb.CourseServiceServer] interface and provides
// operations for Course models. It acts as an adapter between the gRPC transport layer
// and the server-layer buusiness logic of microservice, defined in the [course.Service].
//
// For more information about underlying gRPC server, see [github.com/mikhail5545/proto-go].
type Server struct {
	coursepb.UnimplementedCourseServiceServer
	service courseservice.Service
}

// New creates a new Server instance.
func New(svc courseservice.Service) *Server {
	return &Server{service: svc}
}

// Register registers the course server with a gRPC server instance.
func Register(s *grpc.Server, svc courseservice.Service) {
	coursepb.RegisterCourseServiceServer(s, New(svc))
}

// Get retrieves a single published and not soft-deleted course record,
// along with its associated product details and preloaded course parts.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Get(ctx context.Context, req *coursepb.GetRequest) (*coursepb.GetResponse, error) {
	details, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return &coursepb.GetResponse{CourseDetails: types.CourseDetaisToProtobuf(details)}, nil
}

// GetWithDeleted retrieves a single course record, including soft-deleted ones,
// along with its associated product details and preloaded course parts.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetWithDeleted(ctx context.Context, req *coursepb.GetWithDeletedRequest) (*coursepb.GetWithDeletedResponse, error) {
	details, err := s.service.GetWithDeleted(ctx, req.GetId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return &coursepb.GetWithDeletedResponse{CourseDetails: types.CourseDetaisToProtobuf(details)}, nil
}

// GetWithUnpublished retrieves a single course record, including unpublished ones,
// along with its associated product details and preloaded course parts.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetWithUnpublished(ctx context.Context, req *coursepb.GetWithUnpublishedRequest) (*coursepb.GetWithUnpublishedResponse, error) {
	details, err := s.service.GetWithUnpublished(ctx, req.GetId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return &coursepb.GetWithUnpublishedResponse{CourseDetails: types.CourseDetaisToProtobuf(details)}, nil
}

// GetReduced retrieves a single published and not soft-deleted course record,
// along with its associated product details, but without its course parts.
// This is a lighter version of the `Get` method.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetReduced(ctx context.Context, req *coursepb.GetReducedRequest) (*coursepb.GetReducedResponse, error) {
	details, err := s.service.GetReduced(ctx, req.GetId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return &coursepb.GetReducedResponse{CourseDetails: types.CourseDetaisToProtobuf(details)}, nil
}

// GetReducedWithDeleted retrieves a single course record, including soft-deleted ones,
// along with its associated product details, but without its course parts.
// This is a lighter version of the `GetWithDeleted` method.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetReducedWithDeleted(ctx context.Context, req *coursepb.GetReducedWithDeletedRequest) (*coursepb.GetReducedWithDeletedResponse, error) {
	details, err := s.service.GetReducedWithDeleted(ctx, req.GetId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return &coursepb.GetReducedWithDeletedResponse{CourseDetails: types.CourseDetaisToProtobuf(details)}, nil
}

// List retrieves a paginated list of all published and not soft-deleted course records.
// Each record includes its associated product details and preloaded course parts.
// The response also contains the total count of such records.
func (s *Server) List(ctx context.Context, req *coursepb.ListRequest) (*coursepb.ListResponse, error) {
	courses, total, err := s.service.List(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	var pbcourses []*coursepb.CourseDetails
	for _, c := range courses {
		pbcourses = append(pbcourses, types.CourseDetaisToProtobuf(&c))
	}

	return &coursepb.ListResponse{CourseDetails: pbcourses, Total: total}, nil
}

// ListDeleted retrieves a paginated list of all soft-deleted course records.
// Each record includes its associated product details.
// The response also contains the total count of such records.
func (s *Server) ListDeleted(ctx context.Context, req *coursepb.ListDeletedRequest) (*coursepb.ListDeletedResponse, error) {
	courses, total, err := s.service.ListDeleted(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	var pbcourses []*coursepb.CourseDetails
	for _, c := range courses {
		pbcourses = append(pbcourses, types.CourseDetaisToProtobuf(&c))
	}

	return &coursepb.ListDeletedResponse{CourseDetails: pbcourses, Total: total}, nil
}

// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) course records.
// Each record includes its associated product details.
// The response also contains the total count of such records.
func (s *Server) ListUnpublished(ctx context.Context, req *coursepb.ListUnpublishedRequest) (*coursepb.ListUnpublishedResponse, error) {
	courses, total, err := s.service.ListUnpublished(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	var pbcourses []*coursepb.CourseDetails
	for _, c := range courses {
		pbcourses = append(pbcourses, types.CourseDetaisToProtobuf(&c))
	}

	return &coursepb.ListUnpublishedResponse{CourseDetails: pbcourses, Total: total}, nil
}

// Publish makes a course and its associated product available in the catalog by setting `InStock` to true.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Publish(ctx context.Context, req *coursepb.PublishRequest) (*coursepb.PublishResponse, error) {
	if err := s.service.Publish(ctx, req.GetId()); err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &coursepb.PublishResponse{Id: req.GetId()}, nil
}

// Unpublish archives a course, its associated product, and all its course parts from the catalog
// by setting their `InStock` or `Published` fields to false.
//
// Returns a `NotFound` gRPC error if the course or its parts are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Unpublish(ctx context.Context, req *coursepb.UnpublishRequest) (*coursepb.UnpublishResponse, error) {
	if err := s.service.Unpublish(ctx, req.GetId()); err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &coursepb.UnpublishResponse{Id: req.GetId()}, nil
}

// Create creates a new Course and its associated Product.
// Both are created in an unpublished state.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (s *Server) Create(ctx context.Context, req *coursepb.CreateRequest) (*coursepb.CreateResponse, error) {
	createReq := &coursemodel.CreateRequest{
		Name:             req.Name,
		ShortDescription: req.ShortDescription,
		Topic:            req.Topic,
		Price:            req.Price,
		AccessDuration:   int(req.AccessDuration),
	}
	res, err := s.service.Create(ctx, createReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return &coursepb.CreateResponse{Id: res.ID, ProductId: res.ProductID}, nil
}

// Update performs a partial update of a course and its related product.
// At least one field must be provided for an update to occur.
// The response contains only the fields that were actually changed.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (s *Server) Update(ctx context.Context, req *coursepb.UpdateRequest) (*coursepb.UpdateResponse, error) {
	updateReq := &coursemodel.UpdateRequest{
		Name:             req.Name,
		ShortDescription: req.ShortDescription,
		LongDescription:  req.LongDescription,
		Topic:            req.Topic,
		Price:            req.Price,
		Tags:             req.Tags,
	}
	ad := int(req.GetAccessDuration())
	updateReq.AccessDuration = &ad

	res, err := s.service.Update(ctx, updateReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return types.CourseToProtobufUpdate(res), nil
}

// AddImage adds a new image to a course. It's called by media-service-go upon successful image upload.
// It validates the request, checks the image limit and appends the new information.
//
// Returns `InvalidArgument` gRPC error if the request payload is invalid/image limit is exceeded.
// Returns `NotFound` gRPC error if the record is not found.
func (s *Server) AddImage(ctx context.Context, req *coursepb.AddImageRequest) (*coursepb.AddImageResponse, error) {
	addRequest := &imagemodel.AddRequest{
		OwnerID:        req.GetOwnerId(),
		MediaServiceID: req.GetMediaServiceId(),
		URL:            req.GetUrl(),
		SecureURL:      req.GetSecureUrl(),
		PublicID:       req.GetPublicId(),
	}
	err := s.service.AddImage(ctx, addRequest)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &coursepb.AddImageResponse{MediaServiceId: req.MediaServiceId, OwnerId: req.OwnerId}, nil
}

// DeleteImage deletes an image from a course. It's called by media-service-go upon successful image deletion.
// The function validates the request and removes the image information from the course.
// This action is irreversable.
//
// Returns `InvalidArgument` gRPC error if the request payload is invalid.
// Returns `NotFound` gRPC error if any of records is not found.
func (s *Server) DeleteImage(ctx context.Context, req *coursepb.DeleteImageRequest) (*coursepb.DeleteImageResponse, error) {
	deleteReq := &imagemodel.DeleteRequest{
		OwnerID:        req.GetOwnerId(),
		MediaServiceID: req.GetMediaServiceId(),
	}
	err := s.service.DeleteImage(ctx, deleteReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &coursepb.DeleteImageResponse{OwnerId: req.GetOwnerId(), MediaServiceId: req.GetMediaServiceId()}, nil
}

// AddImageBatch adds an image for a batch of courses. It's called by media-service-go
// upon successful image uplaod.
//
// Returns the number of affected courses.
// Returns `InvalidArgument` gRPC error if the request payload is invalid.
// Returns `NotFound` gRPC error none of the courses were found.
func (s *Server) AddImageBatch(ctx context.Context, req *coursepb.AddImageBatchRequest) (*coursepb.AddImageBatchResponse, error) {
	addReq := &imagemodel.AddBatchRequest{
		MediaServiceID: req.GetMediaServiceId(),
		URL:            req.GetUrl(),
		SecureURL:      req.GetUrl(),
		PublicID:       req.GetPublicId(),
		OwnerIDs:       req.GetOwnerIds(),
	}
	affectedOwners, err := s.service.AddImageBatch(ctx, addReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &coursepb.AddImageBatchResponse{OwnersAffected: int32(affectedOwners)}, nil
}

// DeleteImageBatch deletes an image from a batch of courses. It's called by media-service-go
// upon successful image deletion.
//
// Returns the number of affected courses.
// Returns `InvalidArgument` gRPC error if the request payload is invalid.
// Returns `NotFound` gRPC error none of the courses were found or the image was not found.
func (s *Server) DeleteImageBatch(ctx context.Context, req *coursepb.DeleteImageBatchRequest) (*coursepb.DeleteImageBatchResponse, error) {
	deleteReq := &imagemodel.DeleteBatchRequst{
		MediaServiceID: req.GetMediaServiceId(),
		OwnerIDs:       req.GetOwnerIds(),
	}
	affectedOwners, err := s.service.DeleteImageBatch(ctx, deleteReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &coursepb.DeleteImageBatchResponse{OwnersAffected: int32(affectedOwners)}, nil
}

// Delete performs a soft-delete on a course, its associated product, and all its course parts.
// It also unpublishes them, requiring manual re-publishing after restoration.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Delete(ctx context.Context, req *coursepb.DeleteRequest) (*coursepb.DeleteResponse, error) {
	if err := s.service.Delete(ctx, req.GetId()); err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &coursepb.DeleteResponse{Id: req.GetId()}, nil
}

// DeletePermanent permanently deletes a course, its associated product, and all its course parts from the database.
// This action is irreversible.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) DeletePermanent(ctx context.Context, req *coursepb.DeletePermanentRequest) (*coursepb.DeletePermanentResponse, error) {
	if err := s.service.DeletePermanent(ctx, req.GetId()); err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &coursepb.DeletePermanentResponse{Id: req.GetId()}, nil
}

// Restore restores a soft-deleted course, its associated product, and all its course parts.
// The restored records are not automatically published and must be published manually.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Restore(ctx context.Context, req *coursepb.RestoreRequest) (*coursepb.RestoreResponse, error) {
	if err := s.service.Restore(ctx, req.GetId()); err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &coursepb.RestoreResponse{Id: req.GetId()}, nil
}
