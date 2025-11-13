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

// Package videomanager provides a reusable service for managing videos for different owner types.
package videomanager

import (
	"context"
	"errors"
	"fmt"

	mediaservice "github.com/mikhail5545/product-service-go/internal/clients/mediaservice"
	videomodel "github.com/mikhail5545/product-service-go/internal/models/video"
	videoowner "github.com/mikhail5545/product-service-go/internal/types/video_owner"
	"gorm.io/gorm"
)

// Service provides generic business logic for video management.
type Service interface {
	// Add associates a video with a single owner.
	// If there was another video, associated with this owner, it will be replaced with the new one. It also
	// should be deassociated in the corresponding service separately. This function handles only local owner-video relations.
	// It first validates that the video exists in the media service.
	//
	// It returns an error if provided video instance is already associated with provided owner (ErrVideoInUse),
	// owner/video are not found (ErrOwnerNotFound/ErrVideoNotFound), the request payload is invalid (ErrInvalidArgument) or
	// a database/internal error occurres.
	Add(ctx context.Context, req *videomodel.AddRequest, ownerRepo videoowner.OwnerRepo[videoowner.Owner]) error
	// Remove disassociates a video from a single owner.
	// This function handles only local owner-video relations.
	// Owner should be also deassociated from the video in the corresponding service.
	//
	// It returns an error if owner/video are not found (ErrOwnerNotFound/ErrVideoNotFound),
	// the request payload is invalid (ErrInvalidArgument) or a database/internal error occurres.
	Remove(ctx context.Context, req *videomodel.RemoveRequest, ownerRepo videoowner.OwnerRepo[videoowner.Owner]) error
}

// service holds a media service client to interact with video data.
type service struct {
	mediaClient mediaservice.Client // TODO: replace with protobuf client
}

// New creates a new video manager service instance.
func New(mediaClient mediaservice.Client) Service {
	return &service{mediaClient: mediaClient}
}

// Add associates a video with a single owner.
// If there was another video, associated with this owner, it will be replaced with the new one. It also
// should be deassociated in the corresponding service separately. This function handles only local owner-video relations.
// It first validates that the video exists in the media service.
//
// It returns an error if provided video instance is already associated with provided owner (ErrVideoInUse),
// owner/video are not found (ErrOwnerNotFound/ErrVideoNotFound), the request payload is invalid (ErrInvalidArgument) or
// a database/internal error occurres.
func (s *service) Add(ctx context.Context, req *videomodel.AddRequest, ownerRepo videoowner.OwnerRepo[videoowner.Owner]) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	// TODO: make gRPC call to the media-service-go to fetch the mux asset

	return ownerRepo.DB().Transaction(func(tx *gorm.DB) error {
		txOwnerRepo := ownerRepo.WithTx(tx)

		owner, err := txOwnerRepo.GetWithUnpublished(ctx, req.OwnerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrOwnerNotFound, err)
			}
			return fmt.Errorf("failed to retreive video owner: %w", err)
		}

		if owner.GetVideoID() != nil && *owner.GetVideoID() == req.MediaServiceID {
			return ErrVideoInUse
		}

		if err := txOwnerRepo.UpdateVideoID(ctx, req.OwnerID, &req.MediaServiceID); err != nil {
			return fmt.Errorf("failed to add video to owner: %w", err)
		}
		return nil
	})
}

// Remove disassociates a video from a single owner.
// This function handles only local owner-video relations.
// Owner should be also deassociated from the video in the corresponding service.
//
// It returns an error if owner/video are not found (ErrOwnerNotFound/ErrVideoNotFound),
// the request payload is invalid (ErrInvalidArgument) or a database/internal error occurres.
func (s *service) Remove(ctx context.Context, req *videomodel.RemoveRequest, ownerRepo videoowner.OwnerRepo[videoowner.Owner]) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	return ownerRepo.DB().Transaction(func(tx *gorm.DB) error {
		txOwnerRepo := ownerRepo.WithTx(tx)

		_, err := txOwnerRepo.GetWithUnpublished(ctx, req.OwnerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrOwnerNotFound, err)
			}
			return fmt.Errorf("failed to retrieve owner: %w", err)
		}

		if err := txOwnerRepo.UpdateVideoID(ctx, req.OwnerID, nil); err != nil {
			return fmt.Errorf("failed to remove video from owner: %w", err)
		}
		return nil
	})
}
