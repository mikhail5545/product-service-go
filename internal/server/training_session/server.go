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

	"github.com/mikhail5545/product-service-go/internal/models"
	trainingsession "github.com/mikhail5545/product-service-go/internal/services/training_session"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	trainingsessionpb "github.com/mikhail5545/proto-go/proto/training_session/v0"
)

// Server implements the gRPC [trainingsessionpb.TrainingSessionServiceServer] interface and provides
// operations for TrainingSession models. It acts as an adapter between the gRPC transport layer
// and the server-layer buusiness logic of microservice, defined in the [trainingsession.Service].
//
// For more information about underlying gRPC server, see [github.com/mikhail5545/proto-go].
type Server struct {
	trainingsessionpb.UnimplementedTrainingSessionServiceServer
	service *trainingsession.Service
}

// New creates a new [trainingsession.Server].
func New(s *trainingsession.Service) *Server {
	return &Server{service: s}
}

// Get retrieves a training session by their ID.
// It returns the full training session object.
// If the training session is not found, it returns a `NotFound` gRPC error.
func (s *Server) Get(ctx context.Context, req *trainingsessionpb.GetRequest) (*trainingsessionpb.GetResponse, error) {
	ts, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return &trainingsessionpb.GetResponse{TrainingSession: types.TrainingSessionToProtobuf(ts)}, nil
}

// List retrieves a paginated list of all training sessions.
// The response contains a list of full training session objects.
// and the total number of training sessions in the system.
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

// Create creates a new training session record, typically in the process of direct training session
// creation. It automatically creates an underlying product.
//
// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
// It returns newly created course training session with all fields.
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

// Update updates training session fields that have been acually changed. All request fields
// except ID are optional, so service will update training session only if at least one field
// has been updated.
//
// It populates only updated fields in the response along with the `fieldmaskpb.UpdateMask` which contains
// paths to updated fields.
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

// Delete completely deletes training session from the system.
func (s *Server) Delete(ctx context.Context, req *trainingsessionpb.DeleteRequest) (*trainingsessionpb.DeleteResponse, error) {
	err := s.service.Delete(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return &trainingsessionpb.DeleteResponse{Id: req.GetId()}, nil
}
