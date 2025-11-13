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

// Package video provides service-layer business logic for videos. It represents a generic
// entry point for all video operations across all types of services (video owners). It routes
// video requests through generic [videomanager.Service] to perform low-level logic.
package video

import (
	"context"
	"fmt"

	coursepartrepo "github.com/mikhail5545/product-service-go/internal/database/course_part"
	videomodel "github.com/mikhail5545/product-service-go/internal/models/video"
	coursepart "github.com/mikhail5545/product-service-go/internal/services/course_part"
	videomanager "github.com/mikhail5545/product-service-go/internal/services/video_manager"
	videoowner "github.com/mikhail5545/product-service-go/internal/types/video_owner"
)

// Service provides service-layer logic for videos.
// It acts as the router for video operations to create generic entry point
// for all types of video owners (services) like 'course part', etc.
type Service interface {
	// Add associates a video with a single owner using [videomanager.AddVideo] for specified owner type.
	// If there was another video, associated with this owner, it will be replaced with the new one. It also
	// should be deassociated in the corresponding service separately. This function handles only local owner-video relations.
	// It first validates that the video exists in the media service.
	//
	// It returns an error if provided video instance is already associated with provided owner (videomanager.ErrVideoInUse),
	// owner/video are not found (videomanager.ErrOwnerNotFound/videomanager.ErrVideoNotFound), the request payload
	// is invalid (videomanager.ErrInvalidArgument) or a database/internal error occurres.
	Add(ctx context.Context, ownerType string, req *videomodel.AddRequest) error
	// Remove disassociates a video from a single owner using [videomanager.AddVideo] for specified owner type.
	// This function handles only local owner-video relations.
	// Owner should be also deassociated from the video in the corresponding service.
	//
	//
	// It returns an error if owner/video are not found (videomanager.ErrOwnerNotFound/videomanager.ErrVideoNotFound),
	// the request payload is invalid (videomanager.ErrInvalidArgument) or a database/internal error occurres.
	Remove(ctx context.Context, ownerType string, req *videomodel.RemoveRequest) error
}

// service holds an instance of [coursepartrepo.Repository] to perform database operations for all
// services and generic [videomanager.Service] to
// perform generic video operations.
type service struct {
	coursePartRepo coursepartrepo.Repository
	manager        videomanager.Service
}

// New creates a new Service instance.
func New(cpr coursepartrepo.Repository, m videomanager.Service) Service {
	return &service{coursePartRepo: cpr, manager: m}
}

// getOwnerRepoAdapter returns an adapter for service "ownerType". ownerType should be 'course_part', etc.
//
// Returns ErrUnknownOwner if ownerType is invalid.
func (s *service) getOwnerRepoAdapter(ownerType string) (videoowner.OwnerRepo[videoowner.Owner], error) {
	switch ownerType {
	case "course_part":
		return coursepart.NewOwnerRepoAdapter(s.coursePartRepo), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownOwner, ownerType)
	}
}

// Add associates a video with a single owner using [videomanager.AddVideo] for specified owner type.
// If there was another video, associated with this owner, it will be replaced with the new one. It also
// should be deassociated in the corresponding service separately. This function handles only local owner-video relations.
// It first validates that the video exists in the media service.
func (s *service) Add(ctx context.Context, ownerType string, req *videomodel.AddRequest) error {
	adapter, err := s.getOwnerRepoAdapter(ownerType)
	if err != nil {
		return err
	}
	return s.manager.Add(ctx, req, adapter)
}

// Remove disassociates a video from a single owner using [videomanager.AddVideo] for specified owner type.
// This function handles only local owner-video relations.
// Owner should be also deassociated from the video in the corresponding service.
func (s *service) Remove(ctx context.Context, ownerType string, req *videomodel.RemoveRequest) error {
	adapter, err := s.getOwnerRepoAdapter(ownerType)
	if err != nil {
		return err
	}
	return s.manager.Remove(ctx, req, adapter)
}
