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

	"github.com/mikhail5545/product-service-go/internal/models"
	"github.com/mikhail5545/product-service-go/internal/services/seminar"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	seminarpb "github.com/mikhail5545/proto-go/proto/seminar/v0"
)

// Server implements the gRPC [seminarpb.SeminarServiceServer] interface and provides
// operations for Seminar models. It acts as an adapter between the gRPC transport layer
// and the server-layer buusiness logic of microservice, defined in the [seminar.Service].
//
// For more information about underlying gRPC server, see [github.com/mikhail5545/proto-go].
type Server struct {
	seminarpb.UnimplementedSeminarServiceServer
	service *seminar.Service
}

// New creates a new Server instance.
func New(s *seminar.Service) *Server {
	return &Server{service: s}
}

// Get retrieves a seminar by their ID.
// It returns the full seminar object.
// If the seminar is not found, it returns a `NotFound` gRPC error.
func (s *Server) Get(ctx context.Context, req *seminarpb.GetRequest) (*seminarpb.GetResponse, error) {
	seminar, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.GetResponse{Seminar: types.SeminarToProtobuf(seminar)}, nil
}

// List retrieves a paginated list of all seminars.
// The response contains a list of seminars
// and the total number of seminars in the system.
func (s *Server) List(ctx context.Context, req *seminarpb.ListRequest) (*seminarpb.ListResponse, error) {
	seminars, total, err := s.service.List(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbSeminars []*seminarpb.Seminar
	for _, seminar := range seminars {
		pbSeminars = append(pbSeminars, types.SeminarToProtobuf(&seminar))
	}
	return &seminarpb.ListResponse{Seminars: pbSeminars, Total: total}, nil
}

// Create creates a new seminar record, typically in the process of direct seminar
// creation. It automatically creates all underlying products and populdates they're `name` and `description`
// fields from [models.Seminar.Name] and [models.Seminar.Description] if not provided.
//
// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
// It returns newly created seminar model with all fields.
func (s *Server) Create(ctx context.Context, req *seminarpb.CreateRequest) (*seminarpb.CreateResponse, error) {
	seminar := &models.Seminar{
		Name:            req.GetName(),
		Description:     req.GetDescription(),
		Place:           req.GetPlace(),
		Date:            req.GetDate().AsTime(),
		EndingDate:      req.GetEndingDate().AsTime(),
		Details:         req.GetDetails(),
		LatePaymentDate: req.GetLatePaymentDate().AsTime(),
		ReservationProduct: &models.Product{
			Price:       req.GetReservationProduct().GetPrice(),
			Name:        req.GetReservationProduct().GetName(),
			Description: req.GetReservationProduct().GetDescription(),
		},
		EarlyProduct: &models.Product{
			Price:       req.GetEarlyProduct().GetPrice(),
			Name:        req.GetEarlyProduct().GetName(),
			Description: req.GetEarlyProduct().GetDescription(),
		},
		LateProduct: &models.Product{
			Price:       req.GetLateProduct().GetPrice(),
			Name:        req.GetLateProduct().GetName(),
			Description: req.GetLateProduct().GetDescription(),
		},
		EarlySurchargeProduct: &models.Product{
			Price:       req.GetEarlySurchargeProduct().GetPrice(),
			Name:        req.GetEarlySurchargeProduct().GetName(),
			Description: req.GetEarlySurchargeProduct().GetDescription(),
		},
		LateSurchargeProduct: &models.Product{
			Price:       req.GetLateSurchargeProduct().GetPrice(),
			Name:        req.GetLateSurchargeProduct().GetName(),
			Description: req.GetLateSurchargeProduct().GetDescription(),
		},
	}
	seminar, err := s.service.Create(ctx, seminar)
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.CreateResponse{Seminar: types.SeminarToProtobuf(seminar)}, nil
}

// Update updates seminar fields that have been acually changed. All request fields
// except ID are optional, so service will update seminar only if at least one field
// has been updated.
//
// It populates only updated fields in the response along with the `fieldmaskpb.UpdateMask` which contains
// paths to updated fields.
func (s *Server) Update(ctx context.Context, req *seminarpb.UpdateRequest) (*seminarpb.UpdateResponse, error) {
	seminar := &models.Seminar{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Place:       req.GetPlace(),
		Date:        req.GetDate().AsTime(),
		EndingDate:  req.GetEndingDate().AsTime(),
		Details:     req.GetDetails(),
		ReservationProduct: &models.Product{
			Price:       req.GetReservationProduct().GetPrice(),
			Name:        req.GetReservationProduct().GetName(),
			Description: req.GetReservationProduct().GetDescription(),
		},
		EarlyProduct: &models.Product{
			Price:       req.GetEarlyProduct().GetPrice(),
			Name:        req.GetEarlyProduct().GetName(),
			Description: req.GetEarlyProduct().GetDescription(),
		},
		LateProduct: &models.Product{
			Price:       req.GetLateProduct().GetPrice(),
			Name:        req.GetLateProduct().GetName(),
			Description: req.GetLateProduct().GetDescription(),
		},
		EarlySurchargeProduct: &models.Product{
			Price:       req.GetEarlySurchargeProduct().GetPrice(),
			Name:        req.GetEarlySurchargeProduct().GetName(),
			Description: req.GetEarlySurchargeProduct().GetDescription(),
		},
		LateSurchargeProduct: &models.Product{
			Price:       req.GetLateSurchargeProduct().GetPrice(),
			Name:        req.GetLateSurchargeProduct().GetName(),
			Description: req.GetLateSurchargeProduct().GetDescription(),
		},
	}
	updates, err := s.service.Update(ctx, seminar, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return types.SeminarToProtobufUpdate(updates), nil
}

// Delete completely deletes Seminar record from the system.
func (s *Server) Delete(ctx context.Context, req *seminarpb.DeleteRequest) (*seminarpb.DeleteResponse, error) {
	err := s.service.Delete(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.DeleteResponse{Id: req.GetId()}, nil
}
