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

package server

import (
	"context"

	"github.com/mikhail5545/product-service-go/internal/models"
	"github.com/mikhail5545/product-service-go/internal/services/seminar"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	seminarpb "github.com/mikhail5545/proto-go/proto/seminar/v0"
)

type Server struct {
	seminarpb.UnimplementedSeminarServiceServer
	service *seminar.Service
}

func New(s *seminar.Service) *Server {
	return &Server{service: s}
}

func (s *Server) Get(ctx context.Context, req *seminarpb.GetRequest) (*seminarpb.GetResponse, error) {
	seminar, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.GetResponse{Seminar: types.SeminarToProtobuf(seminar)}, nil
}

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

func (s *Server) Delete(ctx context.Context, req *seminarpb.DeleteRequest) (*seminarpb.DeleteResponse, error) {
	err := s.service.Delete(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	return &seminarpb.DeleteResponse{Id: req.GetId()}, nil
}
