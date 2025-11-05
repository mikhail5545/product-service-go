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
Package physicalgood provides the implementation of the gRPC
[physicalgoodpb.PhysicalGoodServiceServer] interface and provides
various operations for PhysicalGood models.
*/
package physicalgood

import (
	"context"

	physicalgoodmodel "github.com/mikhail5545/product-service-go/internal/models/physical_good"
	physicalgoodservice "github.com/mikhail5545/product-service-go/internal/services/physical_good"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	physicalgoodpb "github.com/mikhail5545/proto-go/proto/physical_good/v0"
	"google.golang.org/grpc"
)

// Server implements the gRPC [physicalgoodpb.PhysicalGoodServiceServer] interface and provides
// operations for PhysicalGood models. It acts as an adapter between the gRPC transport layer
// and the server-layer buusiness logic of microservice, defined in the [physicalgoodservice.Service].
//
// For more information about underlying gRPC server, see [github.com/mikhail5545/proto-go].
type Server struct {
	physicalgoodpb.UnimplementedPhysicalGoodServiceServer
	service physicalgoodservice.Service
}

// New creates a new [physicalgood.Server].
func New(svc physicalgoodservice.Service) *Server {
	return &Server{service: svc}
}

// Register registers the course server with a gRPC server instance.
func Register(s *grpc.Server, svc physicalgoodservice.Service) {
	physicalgoodpb.RegisterPhysicalGoodServiceServer(s, New(svc))
}

// Get retrieves a single published and not soft-deleted physical good record,
// along with its associated product details.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Get(ctx context.Context, req *physicalgoodpb.GetRequest) (*physicalgoodpb.GetResponse, error) {
	details, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &physicalgoodpb.GetResponse{PhysicalGoodDetails: types.PhysicalGoodDetailsToProtobuf(details)}, nil
}

// GetWithDeleted retrieves a single physical good record, including soft-deleted ones,
// along with its associated product details.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetWithDeleted(ctx context.Context, req *physicalgoodpb.GetWithDeletedRequest) (*physicalgoodpb.GetWithDeletedResponse, error) {
	details, err := s.service.GetWithDeleted(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &physicalgoodpb.GetWithDeletedResponse{PhysicalGoodDetails: types.PhysicalGoodDetailsToProtobuf(details)}, nil
}

// GetWithUnpublished retrieves a single physical good record, including unpublished ones,
// along with its associated product details.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetWithUnpublished(ctx context.Context, req *physicalgoodpb.GetWithUnpublishedRequest) (*physicalgoodpb.GetWithUnpublishedResponse, error) {
	details, err := s.service.GetWithUnpublished(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &physicalgoodpb.GetWithUnpublishedResponse{PhysicalGoodDetails: types.PhysicalGoodDetailsToProtobuf(details)}, nil
}

// List retrieves a paginated list of all published and not soft-deleted physical good records.
// Each record includes its associated product details.
// The response also contains the total count of such records.
func (s *Server) List(ctx context.Context, req *physicalgoodpb.ListRequest) (*physicalgoodpb.ListResponse, error) {
	goods, total, err := s.service.List(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbgoods []*physicalgoodpb.PhysicalGoodDetails
	for _, g := range goods {
		pbgoods = append(pbgoods, types.PhysicalGoodDetailsToProtobuf(&g))
	}
	return &physicalgoodpb.ListResponse{PhysicalGoodDetails: pbgoods, Total: total}, nil
}

// ListDeleted retrieves a paginated list of all soft-deleted physical good records.
// Each record includes its associated product details.
// The response also contains the total count of such records.
func (s *Server) ListDeleted(ctx context.Context, req *physicalgoodpb.ListDeletedRequest) (*physicalgoodpb.ListDeletedResponse, error) {
	goods, total, err := s.service.ListDeleted(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbgoods []*physicalgoodpb.PhysicalGoodDetails
	for _, g := range goods {
		pbgoods = append(pbgoods, types.PhysicalGoodDetailsToProtobuf(&g))
	}
	return &physicalgoodpb.ListDeletedResponse{PhysicalGoodDetails: pbgoods, Total: total}, nil
}

// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) physical good records.
// Each record includes its associated product details.
// The response also contains the total count of such records.
func (s *Server) ListUnpublished(ctx context.Context, req *physicalgoodpb.ListUnpublishedRequest) (*physicalgoodpb.ListUnpublishedResponse, error) {
	goods, total, err := s.service.ListUnpublished(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbgoods []*physicalgoodpb.PhysicalGoodDetails
	for _, g := range goods {
		pbgoods = append(pbgoods, types.PhysicalGoodDetailsToProtobuf(&g))
	}
	return &physicalgoodpb.ListUnpublishedResponse{PhysicalGoodDetails: pbgoods, Total: total}, nil
}

// Publish makes a physical good and its associated product available in the catalog by setting `InStock` to true.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Publish(ctx context.Context, req *physicalgoodpb.PublishRequest) (*physicalgoodpb.PublishResponse, error) {
	if err := s.service.Publish(ctx, req.GetId()); err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &physicalgoodpb.PublishResponse{Id: req.GetId()}, nil
}

// Unpublish archives a physical good and its associated product from the catalog
// by setting their `InStock` field to false.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Unpublish(ctx context.Context, req *physicalgoodpb.UnpublishRequest) (*physicalgoodpb.UnpublishResponse, error) {
	if err := s.service.Unpublish(ctx, req.GetId()); err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &physicalgoodpb.UnpublishResponse{Id: req.GetId()}, nil
}

// Create creates a new PhysicalGood and its associated Product.
// Both are created in an unpublished state.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (s *Server) Create(ctx context.Context, req *physicalgoodpb.CreateRequest) (*physicalgoodpb.CreateResponse, error) {
	createReq := &physicalgoodmodel.CreateRequest{
		Name:             req.GetName(),
		ShortDescription: req.GetShortDescription(),
		Price:            req.GetPrice(),
		Amount:           int(req.GetAmount()),
		ShippingRequired: req.GetShippingRequired(),
	}
	res, err := s.service.Create(ctx, createReq)
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &physicalgoodpb.CreateResponse{Id: res.ID, ProductId: res.ProductID}, nil
}

// Update performs a partial update of a physical good and its related product.
// At least one field must be provided for an update to occur.
// The response contains only the fields that were actually changed.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (s *Server) Update(ctx context.Context, req *physicalgoodpb.UpdateRequest) (*physicalgoodpb.UpdateResponse, error) {
	updateReq := &physicalgoodmodel.UpdateRequest{
		ID:               req.GetId(),
		Name:             req.Name,
		ShortDescription: req.ShortDescription,
		LongDescription:  req.LongDescription,
		Price:            req.Price,
		ShippingRequired: req.ShippingRequired,
		Tags:             req.Tags,
	}
	a := int(req.GetAmount())
	updateReq.Amount = &a
	res, err := s.service.Update(ctx, updateReq)
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return types.PhysicalGoodToProtobufUpdate(&physicalgoodpb.UpdateResponse{Id: req.GetId()}, res), nil
}

// Delete performs a soft-delete on a physical good and its associated product.
// It also unpublishes them, requiring manual re-publishing after restoration.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Delete(ctx context.Context, req *physicalgoodpb.DeleteRequest) (*physicalgoodpb.DeleteResponse, error) {
	if err := s.service.Delete(ctx, req.GetId()); err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &physicalgoodpb.DeleteResponse{Id: req.GetId()}, nil
}

// DeletePermanent permanently deletes a physical good and its associated product from the database.
// This action is irreversible.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) DeletePermanent(ctx context.Context, req *physicalgoodpb.DeletePermanentRequest) (*physicalgoodpb.DeletePermanentResponse, error) {
	if err := s.service.DeletePermanent(ctx, req.GetId()); err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &physicalgoodpb.DeletePermanentResponse{Id: req.GetId()}, nil
}

// Restore restores a soft-deleted physical good and its associated product.
// The restored records are not automatically published and must be published manually.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Restore(ctx context.Context, req *physicalgoodpb.RestoreRequest) (*physicalgoodpb.RestoreResponse, error) {
	if err := s.service.Restore(ctx, req.GetId()); err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &physicalgoodpb.RestoreResponse{Id: req.GetId()}, nil
}
