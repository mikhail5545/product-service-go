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
	"github.com/mikhail5545/product-service-go/internal/services"
	"github.com/mikhail5545/product-service-go/internal/utils"
	seminarpb "github.com/mikhail5545/proto-go/proto/seminar/v0"
)

type SeminarServer struct {
	seminarpb.UnimplementedSeminarServiceServer
	seminarService *services.SeminarService
}

func NewSeminarServer(ss *services.SeminarService) *SeminarServer {
	return &SeminarServer{seminarService: ss}
}

func (s *SeminarServer) GetSeminar(ctx context.Context, req *seminarpb.GetSeminarRequest) (*seminarpb.GetSeminarResponse, error) {
	seminar, err := s.seminarService.GetSeminar(ctx, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &seminarpb.GetSeminarResponse{Seminar: utils.ConvertToProtobufSeminar(seminar)}, nil
}

func (s *SeminarServer) ListSeminars(ctx context.Context, req *seminarpb.ListSeminarsRequest) (*seminarpb.ListSeminarsResponse, error) {
	seminars, total, err := s.seminarService.GetSeminars(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, toGRPCError(err)
	}
	var pbSeminars []*seminarpb.Seminar
	for _, seminar := range seminars {
		pbSeminars = append(pbSeminars, utils.ConvertToProtobufSeminar(&seminar))
	}
	return &seminarpb.ListSeminarsResponse{Seminars: pbSeminars, Total: total}, nil
}

func (s *SeminarServer) CreateSeminar(ctx context.Context, req *seminarpb.CreateSeminarRequest) (*seminarpb.CreateSeminarResponse, error) {
	seminar := &models.Seminar{
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
	seminar, err := s.seminarService.CreateSeminar(ctx, seminar)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &seminarpb.CreateSeminarResponse{Seminar: utils.ConvertToProtobufSeminar(seminar)}, nil
}

func (s *SeminarServer) UpdateSeminar(ctx context.Context, req *seminarpb.UpdateSeminarRequest) (*seminarpb.UpdateSeminarResponse, error) {
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
	seminar, err := s.seminarService.UpdateSeminar(ctx, seminar, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &seminarpb.UpdateSeminarResponse{Seminar: utils.ConvertToProtobufSeminar(seminar)}, nil
}

func (s *SeminarServer) DeleteSeminar(ctx context.Context, req *seminarpb.DeleteSeminarRequest) (*seminarpb.DeleteSeminarResponse, error) {
	err := s.seminarService.DeleteSeminar(ctx, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &seminarpb.DeleteSeminarResponse{Id: req.GetId()}, nil
}
