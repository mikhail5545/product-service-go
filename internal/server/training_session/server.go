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
Package trainingsession provides the implementation of the gRPC
[trainingsessionpb.TrainingSessionServiceServer] interface and provides
various operations for TrainingSession models.
*/
package trainingsession

import (
	"context"

	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	trainingsessionmodel "github.com/mikhail5545/product-service-go/internal/models/training_session"
	trainingsessionservice "github.com/mikhail5545/product-service-go/internal/services/training_session"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	trainingsessionpb "github.com/mikhail5545/proto-go/proto/training_session/v0"
	"google.golang.org/grpc"
)

// Server implements the gRPC [trainingsessionpb.TrainingSessionServiceServer] interface and provides
// operations for TrainingSession models. It acts as an adapter between the gRPC transport layer
// and the server-layer buusiness logic of microservice, defined in the [trainingsession.Service].
//
// For more information about underlying gRPC server, see [github.com/mikhail5545/proto-go].
type Server struct {
	trainingsessionpb.UnimplementedTrainingSessionServiceServer
	service trainingsessionservice.Service
}

// New creates a new [trainingsession.Server].
func New(svc trainingsessionservice.Service) *Server {
	return &Server{service: svc}
}

// Register registers the course server with a gRPC server instance.
func Register(s *grpc.Server, svc trainingsessionservice.Service) {
	trainingsessionpb.RegisterTrainingSessionServiceServer(s, New(svc))
}

// Get retrieves a single published and not soft-deleted training session record,
// along with its associated product details.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Get(ctx context.Context, req *trainingsessionpb.GetRequest) (*trainingsessionpb.GetResponse, error) {
	details, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return &trainingsessionpb.GetResponse{TrainingSessionDetails: types.TrainingSessionDetailsToProtobuf(details)}, nil
}

// GetWithDeleted retrieves a single training session record, including soft-deleted ones,
// along with its associated product details.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetWithDeleted(ctx context.Context, req *trainingsessionpb.GetWithDeletedRequest) (*trainingsessionpb.GetWithDeletedResponse, error) {
	details, err := s.service.GetWithDeleted(ctx, req.GetId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return &trainingsessionpb.GetWithDeletedResponse{TrainingSessionDetails: types.TrainingSessionDetailsToProtobuf(details)}, nil
}

// GetWithUnpublished retrieves a single training session record, including unpublished ones,
// along with its associated product details.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetWithUnpublished(ctx context.Context, req *trainingsessionpb.GetWithUnpublishedRequest) (*trainingsessionpb.GetWithUnpublishedResponse, error) {
	details, err := s.service.GetWithUnpublished(ctx, req.GetId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return &trainingsessionpb.GetWithUnpublishedResponse{TrainingSessionDetails: types.TrainingSessionDetailsToProtobuf(details)}, nil
}

// List retrieves a paginated list of all published and not soft-deleted training session records.
// Each record includes its associated product details.
// The response also contains the total count of such records.
func (s *Server) List(ctx context.Context, req *trainingsessionpb.ListRequest) (*trainingsessionpb.ListResponse, error) {
	ts, total, err := s.service.List(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	var pbdetails []*trainingsessionpb.TrainingSessionDetails
	for _, ts := range ts {
		pbdetails = append(pbdetails, types.TrainingSessionDetailsToProtobuf(&ts))
	}

	return &trainingsessionpb.ListResponse{TrainingSessionsDetails: pbdetails, Total: total}, nil
}

// ListDeleted retrieves a paginated list of all soft-deleted training session records.
// Each record includes its associated product details.
// The response also contains the total count of such records.
func (s *Server) ListDeleted(ctx context.Context, req *trainingsessionpb.ListDeletedRequest) (*trainingsessionpb.ListDeletedResponse, error) {
	ts, total, err := s.service.ListDeleted(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	var pbdetails []*trainingsessionpb.TrainingSessionDetails
	for _, ts := range ts {
		pbdetails = append(pbdetails, types.TrainingSessionDetailsToProtobuf(&ts))
	}

	return &trainingsessionpb.ListDeletedResponse{TrainingSessionsDetails: pbdetails, Total: total}, nil
}

// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) training session records.
// Each record includes its associated product details.
// The response also contains the total count of such records.
func (s *Server) ListUnpublished(ctx context.Context, req *trainingsessionpb.ListUnpublishedRequest) (*trainingsessionpb.ListUnpublishedResponse, error) {
	ts, total, err := s.service.ListUnpublished(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	var pbdetails []*trainingsessionpb.TrainingSessionDetails
	for _, ts := range ts {
		pbdetails = append(pbdetails, types.TrainingSessionDetailsToProtobuf(&ts))
	}

	return &trainingsessionpb.ListUnpublishedResponse{TrainingSessionsDetails: pbdetails, Total: total}, nil
}

// Publish makes a training session and its associated product available in the catalog by setting `InStock` to true.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Publish(ctx context.Context, req *trainingsessionpb.PublishRequest) (*trainingsessionpb.PublishResponse, error) {
	if err := s.service.Publish(ctx, req.GetId()); err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &trainingsessionpb.PublishResponse{Id: req.GetId()}, nil
}

// Unpublish archives a training session and its associated product from the catalog
// by setting their `InStock` field to false.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Unpublish(ctx context.Context, req *trainingsessionpb.UnpublishRequest) (*trainingsessionpb.UnpublishResponse, error) {
	if err := s.service.Unpublish(ctx, req.GetId()); err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &trainingsessionpb.UnpublishResponse{Id: req.GetId()}, nil
}

// Create creates a new TrainingSession and its associated Product.
// Both are created in an unpublished state.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (s *Server) Create(ctx context.Context, req *trainingsessionpb.CreateRequest) (*trainingsessionpb.CreateResponse, error) {
	createReq := &trainingsessionmodel.CreateRequest{
		Name:             req.GetName(),
		ShortDescription: req.GetShortDescription(),
		Format:           req.GetFormat(),
		Price:            req.GetPrice(),
		DurationMinutes:  int(req.GetDurationMinutes()),
	}
	res, err := s.service.Create(ctx, createReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &trainingsessionpb.CreateResponse{Id: res.ID, ProductId: res.ProductID}, nil
}

// Update performs a partial update of a training session and its related product.
// At least one field must be provided for an update to occur.
// The response contains only the fields that were actually changed.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (s *Server) Update(ctx context.Context, req *trainingsessionpb.UpdateRequest) (*trainingsessionpb.UpdateResponse, error) {
	updateReq := &trainingsessionmodel.UpdateRequest{
		ID:               req.GetId(),
		Name:             req.Name,
		ShortDescription: req.ShortDescription,
		LongDescription:  req.LongDescription,
		Format:           req.Format,
		Price:            req.Price,
		Tags:             req.Tags,
	}
	dm := int(req.GetDurationMinutes())
	updateReq.DurationMinutes = &dm
	res, err := s.service.Update(ctx, updateReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return types.TrainingSessionToProtobufUpdate(&trainingsessionpb.UpdateResponse{Id: req.GetId()}, res), nil
}

// AddImage adds a new image to a training session. It's called by media-service-go upon successful image upload.
// It validates the request, checks the image limit and appends the new information.
//
// Returns `InvalidArgument` gRPC error if the request payload is invalid/image limit is exceeded.
// Returns `NotFound` gRPC error if the record is not found.
func (s *Server) AddImage(ctx context.Context, req *trainingsessionpb.AddImageRequest) (*trainingsessionpb.AddImageResponse, error) {
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
	return &trainingsessionpb.AddImageResponse{MediaServiceId: req.MediaServiceId, OwnerId: req.OwnerId}, nil
}

// DeleteImage deletes an image from a training session. It's called by media-service-go upon successful image deletion.
// The function validates the request and removes the image information from the training session.
// This action is irreversable.
//
// Returns `InvalidArgument` gRPC error if the request payload is invalid.
// Returns `NotFound` gRPC error if any of records is not found.
func (s *Server) DeleteImage(ctx context.Context, req *trainingsessionpb.DeleteImageRequest) (*trainingsessionpb.DeleteImageResponse, error) {
	deleteReq := &imagemodel.DeleteRequest{
		OwnerID:        req.GetOwnerId(),
		MediaServiceID: req.GetMediaServiceId(),
	}
	err := s.service.DeleteImage(ctx, deleteReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &trainingsessionpb.DeleteImageResponse{OwnerId: req.GetOwnerId(), MediaServiceId: req.GetMediaServiceId()}, nil
}

// AddImageBatch adds an image for a batch of training sessions. It's called by media-service-go
// upon successful image uplaod.
//
// Returns `InvalidArgument` gRPC error if the request payload is invalid.
// Returns `NotFound` gRPC error none of the training sessions were found.
func (s *Server) AddImageBatch(ctx context.Context, req *trainingsessionpb.AddImageBatchRequest) (*trainingsessionpb.AddImageBatchResponse, error) {
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
	return &trainingsessionpb.AddImageBatchResponse{OwnersAffected: int32(affectedOwners)}, nil
}

// DeleteImageBatch deletes an image from a batch of training sessions. It's called by media-service-go
// upon successful image deletion.
//
// Returns `InvalidArgument` gRPC error if the request payload is invalid.
// Returns `NotFound` gRPC error none of the training sessions were found or the image was not found.
func (s *Server) DeleteImageBatch(ctx context.Context, req *trainingsessionpb.DeleteImageBatchRequest) (*trainingsessionpb.DeleteImageBatchResponse, error) {
	deleteReq := &imagemodel.DeleteBatchRequst{
		MediaServiceID: req.GetMediaServiceId(),
		OwnerIDs:       req.GetOwnerIds(),
	}
	affectedOwners, err := s.service.DeleteImageBatch(ctx, deleteReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &trainingsessionpb.DeleteImageBatchResponse{OwnersAffected: int32(affectedOwners)}, nil
}

// Delete performs a soft-delete on a training session and its associated product.
// It also unpublishes them, requiring manual re-publishing after restoration.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Delete(ctx context.Context, req *trainingsessionpb.DeleteRequest) (*trainingsessionpb.DeleteResponse, error) {
	err := s.service.Delete(ctx, req.GetId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return &trainingsessionpb.DeleteResponse{Id: req.GetId()}, nil
}

// DeletePermanent permanently deletes a trainign session and its associated product from the database.
// This action is irreversible.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) DeletePermanent(ctx context.Context, req *trainingsessionpb.DeletePermanentRequest) (*trainingsessionpb.DeletePermanentResponse, error) {
	err := s.service.DeletePermanent(ctx, req.GetId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return &trainingsessionpb.DeletePermanentResponse{Id: req.GetId()}, nil
}

// Restore restores a soft-deleted training session and its associated product.
// The restored records are not automatically published and must be published manually.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Restore(ctx context.Context, req *trainingsessionpb.RestoreRequest) (*trainingsessionpb.RestoreResponse, error) {
	err := s.service.Restore(ctx, req.GetId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}

	return &trainingsessionpb.RestoreResponse{Id: req.GetId()}, nil
}
