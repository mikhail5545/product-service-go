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
Package image provides the implementation of the gRPC
[imagepb.ImageServiceServer] interface and provides
various operations for images.
*/
package image

import (
	"context"

	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	imageservice "github.com/mikhail5545/product-service-go/internal/services/image"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	imagepb "github.com/mikhail5545/proto-go/proto/product_service/image/v0"
)

type Server struct {
	imagepb.UnimplementedImageServiceServer
	service imageservice.Service
}

func (s *Server) Add(ctx context.Context, req *imagepb.AddRequest) (*imagepb.AddResponse, error) {
	addReq := &imagemodel.AddRequest{
		PublicID:       req.GetPublicId(),
		URL:            req.GetUrl(),
		SecureURL:      req.GetSecureUrl(),
		MediaServiceID: req.GetMediaServiceId(),
		OwnerID:        req.GetOwnerId(),
	}
	err := s.service.Add(ctx, req.GetOwnerType(), addReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &imagepb.AddResponse{MediaServiceId: req.GetMediaServiceId(), OwnerId: req.GetOwnerId()}, nil
}

func (s *Server) Delete(ctx context.Context, req *imagepb.DeleteRequest) (*imagepb.DeleteResponse, error) {
	deleteReq := &imagemodel.DeleteRequest{
		OwnerID:        req.GetOwnerId(),
		MediaServiceID: req.GetMediaServiceId(),
	}
	err := s.service.Delete(ctx, req.GetOwnerType(), deleteReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &imagepb.DeleteResponse{MediaServiceId: req.GetMediaServiceId(), OwnerId: req.GetOwnerId()}, nil
}

func (s *Server) AddBatch(ctx context.Context, req *imagepb.AddBatchRequest) (*imagepb.AddBatchResponse, error) {
	addReq := &imagemodel.AddBatchRequest{
		PublicID:       req.GetPublicId(),
		URL:            req.GetUrl(),
		SecureURL:      req.GetSecureUrl(),
		MediaServiceID: req.GetMediaServiceId(),
		OwnerIDs:       req.GetOwnerIds(),
	}
	ownersAffected, err := s.service.AddBatch(ctx, req.GetOwnerType(), addReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &imagepb.AddBatchResponse{OwnersAffected: int64(ownersAffected)}, nil
}

func (s *Server) DeleteBatch(ctx context.Context, req *imagepb.DeleteBatchRequest) (*imagepb.DeleteBatchResponse, error) {
	deleteReq := &imagemodel.DeleteBatchRequst{
		MediaServiceID: req.GetMediaServiceId(),
		OwnerIDs:       req.GetOwnerIds(),
	}
	ownersAffected, err := s.service.DeleteBatch(ctx, req.GetOwnerType(), deleteReq)
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &imagepb.DeleteBatchResponse{OwnersAffected: int64(ownersAffected)}, nil
}
