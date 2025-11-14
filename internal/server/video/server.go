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
Package video provides the implementation of the gRPC
[videopb.VideoServiceServer] interface and provides
various operations for videos.
*/
package video

import (
	"context"

	videomodel "github.com/mikhail5545/product-service-go/internal/models/video"
	videoservice "github.com/mikhail5545/product-service-go/internal/services/video"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	videopb "github.com/mikhail5545/proto-go/proto/product_service/video/v0"
	"google.golang.org/grpc"
)

// Server implements [videopb.UnimplementedVideoServiceServer] and provides
// various operations for videos. It acts as an adapter between
// gRPC server and the business service-layer logic from [videoservice.Service].
// See more details about [underlying protobuf services].
//
// [underlying protobuf services]: https://github.com/mikhail5545/proto-go
type Server struct {
	videopb.UnimplementedVideoServiceServer
	service videoservice.Service
}

// New creates a new Server instance.
func New(svc videoservice.Service) *Server {
	return &Server{service: svc}
}

// Register registers the video server with a gRPC server instance.
func Register(s *grpc.Server, svc videoservice.Service) {
	videopb.RegisterVideoServiceServer(s, New(svc))
}

// Add associates a video with a single owner for specified owner type.
// If there was another video, associated with this owner, it will be replaced with the new one. It also
// should be deassociated in the corresponding service separately. This function handles only local owner-video relations.
// It first validates that the video exists in the media service.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid or video already associated with this owner.
// Returns a `NotFound` gRPC error if any of the video is not found or the owner is not found.
func (s *Server) Add(ctx context.Context, req *videopb.AddRequest) (*videopb.AddResponse, error) {
	addReq := &videomodel.AddRequest{
		OwnerID:        req.GetOwnerId(),
		MediaServiceID: req.GetMediaServiceId(),
	}
	if err := s.service.Add(ctx, req.GetOwnerType(), addReq); err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &videopb.AddResponse{OwnerId: req.GetOwnerId(), OwnerType: req.GetOwnerType()}, nil
}

// Remove disassociates a video from a single owner for specified owner type.
// This function handles only local owner-video relations.
// Owner should be also deassociated from the video in the corresponding service.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
// Returns a `NotFound` gRPC error if any of the video is not found or the owner is not found.
func (s *Server) Remove(ctx context.Context, req *videopb.RemoveRequest) (*videopb.RemoveResponse, error) {
	removeReq := &videomodel.RemoveRequest{
		OwnerID:        req.GetOwnerId(),
		MediaServiceID: req.GetMediaServiceId(),
	}
	if err := s.service.Remove(ctx, req.GetOwnerType(), removeReq); err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &videopb.RemoveResponse{OwnerId: req.GetOwnerId(), MediaServiceId: req.GetMediaServiceId()}, nil
}

// GetOwner retrieves a single owner information including unpublished ones.
// Returns minimal necessary owner information. If more owner information is needed,
// specific owner gRPC service's Get method should be called.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
// Returns a `NotFound` gRPC error if owner is not found.
func (s *Server) GetOwner(ctx context.Context, req *videopb.GetOwnerRequest) (*videopb.GetOwnerResponse, error) {
	owner, err := s.service.GetOwner(ctx, req.GetOwnerType(), req.GetOwnerId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	pbOwner := &videopb.Owner{
		Id:      owner.ID,
		VideoId: owner.VideoID,
	}
	return &videopb.GetOwnerResponse{Owner: pbOwner}, nil
}
