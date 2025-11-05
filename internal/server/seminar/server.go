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
Package seminar provides the implementation of the gRPC
[seminarpb.SeminarServiceServer] interface and provides
various operations for Seminar models.
*/
package seminar

import (
	"context"

	seminarmodel "github.com/mikhail5545/product-service-go/internal/models/seminar"
	seminarservice "github.com/mikhail5545/product-service-go/internal/services/seminar"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	seminarpb "github.com/mikhail5545/proto-go/proto/seminar/v0"
	"google.golang.org/grpc"
)

// Server implements the gRPC [seminarpb.SeminarServiceServer] interface and provides
// operations for Seminar models. It acts as an adapter between the gRPC transport layer
// and the server-layer buusiness logic of microservice, defined in the [seminar.Service].
//
// For more information about underlying gRPC server, see [github.com/mikhail5545/proto-go].
type Server struct {
	seminarpb.UnimplementedSeminarServiceServer
	service seminarservice.Service
}

// New creates a new Server instance.
func New(svc seminarservice.Service) *Server {
	return &Server{service: svc}
}

// Register registers the course server with a gRPC server instance.
func Register(s *grpc.Server, svc seminarservice.Service) {
	seminarpb.RegisterSeminarServiceServer(s, New(svc))
}

// Get retrieves a single published and not soft-deleted seminar record,
// along with all of its associated products details.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Get(ctx context.Context, req *seminarpb.GetRequest) (*seminarpb.GetResponse, error) {
	details, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.GetResponse{SeminarDetails: types.SeminarDetailsToProtobuf(details)}, nil
}

// GetWithDeleted retrieves a single seminar record, including soft-deleted ones,
// along with all of its associated products details.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetWithDeleted(ctx context.Context, req *seminarpb.GetWithDeletedRequest) (*seminarpb.GetWithDeletedResponse, error) {
	details, err := s.service.GetWithDeleted(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.GetWithDeletedResponse{SeminarDetails: types.SeminarDetailsToProtobuf(details)}, nil
}

// GetWithUnpublished retrieves a single seminar record, including unpublished ones,
// along with all of its associated products details.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) GetWithUnpublished(ctx context.Context, req *seminarpb.GetWithUnpublishedRequest) (*seminarpb.GetWithUnpublishedResponse, error) {
	details, err := s.service.GetWithUnpublished(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.GetWithUnpublishedResponse{SeminarDetails: types.SeminarDetailsToProtobuf(details)}, nil
}

// List retrieves a paginated list of all published and not soft-deleted seminar records.
// Each record includes all of its associated products details.
// The response also contains the total count of such records.
func (s *Server) List(ctx context.Context, req *seminarpb.ListRequest) (*seminarpb.ListResponse, error) {
	seminars, total, err := s.service.List(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbSeminars []*seminarpb.SeminarDetails
	for _, seminar := range seminars {
		pbSeminars = append(pbSeminars, types.SeminarDetailsToProtobuf(&seminar))
	}
	return &seminarpb.ListResponse{SeminarDetails: pbSeminars, Total: total}, nil
}

// ListDeleted retrieves a paginated list of all soft-deleted seminar records.
// Each record includes all of its associated products details.
// The response also contains the total count of such records.
func (s *Server) ListDeleted(ctx context.Context, req *seminarpb.ListDeletedRequest) (*seminarpb.ListDeletedResponse, error) {
	seminars, total, err := s.service.ListDeleted(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbSeminars []*seminarpb.SeminarDetails
	for _, seminar := range seminars {
		pbSeminars = append(pbSeminars, types.SeminarDetailsToProtobuf(&seminar))
	}
	return &seminarpb.ListDeletedResponse{SeminarDetails: pbSeminars, Total: total}, nil
}

// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) seminar records.
// Each record includes all of its associated products details.
// The response also contains the total count of such records.
func (s *Server) ListUnpublished(ctx context.Context, req *seminarpb.ListUnpublishedRequest) (*seminarpb.ListUnpublishedResponse, error) {
	seminars, total, err := s.service.ListUnpublished(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbSeminars []*seminarpb.SeminarDetails
	for _, seminar := range seminars {
		pbSeminars = append(pbSeminars, types.SeminarDetailsToProtobuf(&seminar))
	}
	return &seminarpb.ListUnpublishedResponse{SeminarDetails: pbSeminars, Total: total}, nil
}

// Create creates a new Course and all of its associated products (5 total product records).
// All of them are created in an unpublished state.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (s *Server) Create(ctx context.Context, req *seminarpb.CreateRequest) (*seminarpb.CreateResponse, error) {
	createReq := &seminarmodel.CreateRequest{
		Name:                req.GetName(),
		ShortDescription:    req.GetShortDescription(),
		ReservationPrice:    req.GetReservationPrice(),
		EarlyPrice:          req.GetEarlyPrice(),
		LatePrice:           req.GetLatePrice(),
		EarlySurchargePrice: req.GetEarlySurchargePrice(),
		LateSurchargePrice:  req.GetLateSurchargePrice(),
		Date:                req.GetDate().AsTime(),
		EndingDate:          req.GetDate().AsTime(),
		LatePaymentDate:     req.GetDate().AsTime(),
		Place:               req.GetPlace(),
	}
	res, err := s.service.Create(ctx, createReq)
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.CreateResponse{
		Id:                      res.ID,
		ReservationProductId:    res.ReservationProductID,
		EarlyProductId:          res.EarlyProductID,
		LateProductId:           res.LateProductID,
		EarlySurchargeProductId: res.EarlySurchargeProductID,
		LateSurchargeProductId:  res.LateSurchargeProductID,
	}, nil
}

// Publish makes a seminar and all of its associated products available in the catalog by setting `InStock` to true.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Publish(ctx context.Context, req *seminarpb.PublishRequest) (*seminarpb.PublishResponse, error) {
	if err := s.service.Publish(ctx, req.GetId()); err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.PublishResponse{Id: req.GetId()}, nil
}

// Unpublish archives a seminar and all of its associated products from the catalog
// by setting their `InStock` to false.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Unpublish(ctx context.Context, req *seminarpb.UnpublishRequest) (*seminarpb.UnpublishResponse, error) {
	if err := s.service.Unpublish(ctx, req.GetId()); err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.UnpublishResponse{Id: req.GetId()}, nil
}

// Update performs a partial update of a seminar and its related products.
// At least one field must be provided for an update to occur.
// The response contains only the fields that were actually changed.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (s *Server) Update(ctx context.Context, req *seminarpb.UpdateRequest) (*seminarpb.UpdateResponse, error) {
	updateReq := &seminarmodel.UpdateRequest{
		ID:                  req.GetId(),
		Name:                req.Name,
		ShortDescription:    req.ShortDescription,
		LongDescription:     req.LongDescription,
		ReservationPrice:    req.ReservationPrice,
		EarlyPrice:          req.EarlyPrice,
		LatePrice:           req.LatePrice,
		EarlySurchargePrice: req.EarlySurchargePrice,
		LateSurchargePrice:  req.LateSurchargePrice,
		Place:               req.Place,
		Tags:                req.Tags,
	}
	date := req.Date.AsTime()
	edate := req.EndingDate.AsTime()
	lpdate := req.LatePaymentDate.AsTime()
	updateReq.Date = &date
	updateReq.EndingDate = &edate
	updateReq.LatePaymentDate = &lpdate
	res, err := s.service.Update(ctx, updateReq)
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return types.SeminarToProtobufUpdate(&seminarpb.UpdateResponse{Id: req.GetId()}, res), nil
}

// Delete performs a soft-delete on a seminar and all of its associated products.
// It also unpublishes them, requiring manual re-publishing after restoration.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Delete(ctx context.Context, req *seminarpb.DeleteRequest) (*seminarpb.DeleteResponse, error) {
	err := s.service.Delete(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.DeleteResponse{Id: req.GetId()}, nil
}

// DeletePermanent permanently deletes a seminar and all of its associated products from the database.
// This action is irreversible.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) DeletePermanent(ctx context.Context, req *seminarpb.DeletePermanentRequest) (*seminarpb.DeletePermanentResponse, error) {
	err := s.service.DeletePermanent(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.DeletePermanentResponse{Id: req.GetId()}, nil
}

// Restore restores a soft-deleted seminar and all of its associated products.
// The restored records are not automatically published and must be published manually.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (s *Server) Restore(ctx context.Context, req *seminarpb.RestoreRequest) (*seminarpb.RestoreResponse, error) {
	err := s.service.Restore(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.RestoreResponse{Id: req.GetId()}, nil
}
