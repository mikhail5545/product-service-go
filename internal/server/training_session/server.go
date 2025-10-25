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
	trainingsession "github.com/mikhail5545/product-service-go/internal/services/training_session"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	trainingsessionpb "github.com/mikhail5545/proto-go/proto/training_session/v0"
)

type Server struct {
	trainingsessionpb.UnimplementedTrainingSessionServiceServer
	service *trainingsession.Service
}

func NewServer(s *trainingsession.Service) *Server {
	return &Server{service: s}
}

func (s *Server) Get(ctx context.Context, req *trainingsessionpb.GetRequest) (*trainingsessionpb.GetResponse, error) {
	ts, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return &trainingsessionpb.GetResponse{TrainingSession: types.TrainingSessionToProtobuf(ts)}, nil
}

func (s *Server) List(ctx context.Context, req *trainingsessionpb.ListRequest) (*trainingsessionpb.ListResponse, error) {
	ts, total, err := s.service.List(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbTs []*trainingsessionpb.TrainingSession
	for _, ts := range ts {
		pbTs = append(pbTs, types.TrainingSessionToProtobuf(&ts))
	}

	return &trainingsessionpb.ListResponse{TrainingSessions: pbTs, Total: total}, nil
}

func (s *Server) Create(ctx context.Context, req *trainingsessionpb.CreateRequest) (*trainingsessionpb.CreateResponse, error) {
	ts := &models.TrainingSession{
		DurationMinutes: int(req.GetDurationMinutes()),
		Format:          req.GetFormat(),
		Product: &models.Product{
			Name:        req.GetProduct().GetName(),
			Description: req.GetProduct().GetDescription(),
			Price:       req.GetProduct().GetPrice(),
		},
	}
	ts, err := s.service.Create(ctx, ts)
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return &trainingsessionpb.CreateResponse{TrainingSession: types.TrainingSessionToProtobuf(ts)}, nil
}

func (s *Server) Update(ctx context.Context, req *trainingsessionpb.UpdateRequest) (*trainingsessionpb.UpdateResponse, error) {
	ts := &models.TrainingSession{
		DurationMinutes: int(req.GetDurationMinutes()),
		Format:          req.GetFormat(),
		Product: &models.Product{
			Name:        req.GetProduct().GetName(),
			Description: req.GetProduct().GetDescription(),
			Price:       req.GetProduct().GetPrice(),
		},
	}
	updates, productUpdates, err := s.service.Update(ctx, ts, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return types.TrainingSessionToProtobufUpdate(updates, productUpdates), nil
}

func (s *Server) Delete(ctx context.Context, req *trainingsessionpb.DeleteRequest) (*trainingsessionpb.DeleteResponse, error) {
	err := s.service.Delete(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return &trainingsessionpb.DeleteResponse{Id: req.GetId()}, nil
}
