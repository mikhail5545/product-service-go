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

// Package image provides service-layer business logic for images. It represents a generic
// entry point for all image operations across all types of services (image owners). It routes
// image requests through generic [imagemanager.Service] to perform low-level logic.
package image

import (
	"context"
	"fmt"

	courserepo "github.com/mikhail5545/product-service-go/internal/database/course"
	physicalgoodrepo "github.com/mikhail5545/product-service-go/internal/database/physical_good"
	seminarrepo "github.com/mikhail5545/product-service-go/internal/database/seminar"
	trainingsessionrepo "github.com/mikhail5545/product-service-go/internal/database/training_session"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	"github.com/mikhail5545/product-service-go/internal/services/course"
	physicalgood "github.com/mikhail5545/product-service-go/internal/services/physical_good"
	"github.com/mikhail5545/product-service-go/internal/services/seminar"
	trainingsession "github.com/mikhail5545/product-service-go/internal/services/training_session"

	imagemanager "github.com/mikhail5545/product-service-go/internal/services/image_manager"
	imageowner "github.com/mikhail5545/product-service-go/internal/types/image_owner"
)

// Service provides service-layer logic for images.
// It acts as the router for image operations to create generic entry point
// for all types of image owners (services) like 'training session', 'course', etc.
type Service interface {
	Add(ctx context.Context, ownerType string, req *imagemodel.AddRequest) error
	Delete(ctx context.Context, ownerType string, req *imagemodel.DeleteRequest) error
	AddBatch(ctx context.Context, ownerType string, req *imagemodel.AddBatchRequest) (int, error)
	DeleteBatch(ctx context.Context, ownerType string, req *imagemodel.DeleteBatchRequst) (int, error)
}

// service holds instances of [courserepo.Repository], [seminarrepo.Repository], [trainingsessionrepo.Repository],
// [physicalgoodrepo.Repository] to perform database operations for all services and generic [imagemanager.Service] to
// perform generic image operations.
type service struct {
	imageSvc            imagemanager.Service
	courseRepo          courserepo.Repository
	seminarRepo         seminarrepo.Repository
	trainingSessionRepo trainingsessionrepo.Repository
	physicalGoodRepo    physicalgoodrepo.Repository
}

// New creates a new Service instance.
func New(
	imgSvc imagemanager.Service,
	cr courserepo.Repository,
	sr seminarrepo.Repository,
	tsr trainingsessionrepo.Repository,
	pgr physicalgoodrepo.Repository,
) Service {
	return &service{
		imageSvc:            imgSvc,
		courseRepo:          cr,
		seminarRepo:         sr,
		trainingSessionRepo: tsr,
		physicalGoodRepo:    pgr,
	}
}

// getOwnerRepoAdapter returns an adapter for service "ownerType". ownerType should be 'course', 'seminar', etc.
//
// Returns ErrUnknownOwner if ownerType is invalid.
func (s *service) getOwnerRepoAdapter(ownerType string) (imageowner.OwnerRepo[imageowner.Owner], error) {
	switch ownerType {
	case "course":
		return course.NewOwnerRepoAdapter(s.courseRepo), nil
	case "seminar":
		return seminar.NewOwnerRepoAdapter(s.seminarRepo), nil
	case "training_session":
		return trainingsession.NewOwnerRepoAdapter(s.trainingSessionRepo), nil
	case "physical_good":
		return physicalgood.NewOwnerRepoAdapter(s.physicalGoodRepo), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownOwner, ownerType)
	}
}

// Add adds an image for owner using [imagemanager.AddImage] for specified owner type.
func (s *service) Add(ctx context.Context, ownerType string, req *imagemodel.AddRequest) error {
	adapter, err := s.getOwnerRepoAdapter(ownerType)
	if err != nil {
		return err
	}
	return s.imageSvc.AddImage(ctx, req, adapter)
}

// Delete deletes an image from owner using [imagemanager.DeleteImage] for specified owner type.
func (s *service) Delete(ctx context.Context, ownerType string, req *imagemodel.DeleteRequest) error {
	adapter, err := s.getOwnerRepoAdapter(ownerType)
	if err != nil {
		return err
	}
	return s.imageSvc.DeleteImage(ctx, req, adapter)
}

// AddBatch adds an image for batch of owners using [imagemanager.AddImageBatch] for specified owner type.
//
// Returns the number of affected owners.
func (s *service) AddBatch(ctx context.Context, ownerType string, req *imagemodel.AddBatchRequest) (int, error) {
	adapter, err := s.getOwnerRepoAdapter(ownerType)
	if err != nil {
		return 0, err
	}
	return s.imageSvc.AddImageBatch(ctx, req, adapter)
}

// DeleteBatch deletes an image from batch of owners using [imagemanager.DeleteImageBatch] for specified owner type.
//
// Returns the number of affected owners.
func (s *service) DeleteBatch(ctx context.Context, ownerType string, req *imagemodel.DeleteBatchRequst) (int, error) {
	adapter, err := s.getOwnerRepoAdapter(ownerType)
	if err != nil {
		return 0, err
	}
	return s.imageSvc.DeleteImageBatch(ctx, req, adapter)
}
