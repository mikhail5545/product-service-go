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
	trainingsessionpb "vitainmove.com/product-service-go/proto/training_session/v0"
)

type TrainingSessionServer struct {
	trainingsessionpb.UnimplementedTrainingSessionServiceServer
	trainingSessionService *services.TrainingSessionService
}

func NewTrainingSessionServer(tsService *services.TrainingSessionService) *TrainingSessionServer {
	return &TrainingSessionServer{trainingSessionService: tsService}
}

func (s *TrainingSessionServer) GetTrainingSession(ctx context.Context, req *trainingsessionpb.GetTrainingSessionRequest) (*trainingsessionpb.GetTrainingSessionResponse, error) {
	ts, err := s.trainingSessionService.GetTrainingSession(ctx, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &trainingsessionpb.GetTrainingSessionResponse{TrainingSession: utils.ConvertToProtobufTrainingSession(ts)}, nil
}

func (s *TrainingSessionServer) ListTrainingSessions(ctx context.Context, req *trainingsessionpb.ListTrainingSessionsRequest) (*trainingsessionpb.ListTrainingSessionsResponse, error) {
	ts, total, err := s.trainingSessionService.GetTrainingSessions(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, toGRPCError(err)
	}
	var pbTs []*trainingsessionpb.TrainingSession
	for _, ts := range ts {
		pbTs = append(pbTs, utils.ConvertToProtobufTrainingSession(&ts))
	}

	return &trainingsessionpb.ListTrainingSessionsResponse{TrainingSessions: pbTs, Total: total}, nil
}

func (s *TrainingSessionServer) CreateTrainingSession(ctx context.Context, req *trainingsessionpb.CreateTrainingSessionRequest) (*trainingsessionpb.CreateTrainingSessionResponse, error) {
	ts := &models.TrainingSession{
		DurationMinutes: int(req.GetDurationMinutes()),
		Format:          req.GetFormat(),
		Product: &models.Product{
			Name:        req.GetProductInfo().GetName(),
			Description: req.GetProductInfo().GetDescription(),
			Price:       req.GetProductInfo().GetPrice(),
		},
	}
	ts, err := s.trainingSessionService.CreateTrainingSession(ctx, ts)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &trainingsessionpb.CreateTrainingSessionResponse{TrainingSession: utils.ConvertToProtobufTrainingSession(ts)}, nil
}

func (s *TrainingSessionServer) UpdateTrainingSession(ctx context.Context, req *trainingsessionpb.UpdateTrainingSessionRequest) (*trainingsessionpb.UpdateTrainingSessionResponse, error) {
	ts := &models.TrainingSession{
		DurationMinutes: int(req.GetDurationMinutes()),
		Format:          req.GetFormat(),
		Product: &models.Product{
			Name:        req.GetProductInfo().GetName(),
			Description: req.GetProductInfo().GetDescription(),
			Price:       req.GetProductInfo().GetPrice(),
		},
	}
	ts, err := s.trainingSessionService.UpdateTrainingSession(ctx, ts, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &trainingsessionpb.UpdateTrainingSessionResponse{TrainingSession: utils.ConvertToProtobufTrainingSession(ts)}, nil
}

func (s *TrainingSessionServer) DeleteTrainingSession(ctx context.Context, req *trainingsessionpb.DeleteTrainingSessionRequest) (*trainingsessionpb.DeleteTrainingSessionResponse, error) {
	err := s.trainingSessionService.DeleteTrainingSession(ctx, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &trainingsessionpb.DeleteTrainingSessionResponse{Id: req.GetId()}, nil
}
