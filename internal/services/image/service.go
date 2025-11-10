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

// Package image provides a reusable service for managing images for different owner types.
package image

import (
	"context"
	"errors"
	"fmt"

	imagerepo "github.com/mikhail5545/product-service-go/internal/database/image"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	imageowner "github.com/mikhail5545/product-service-go/internal/types/image_owner"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/services/image_mock/service_mock.go -package=image_mock . Service

// Service provides generic business logic for image management.
type Service interface {
	// AddImage adds an image for a single owner.
	// The owner must implement the Owner interface, and its repository
	// must implement the OwnerRepo interface.
	//
	// Returns an error if the request payload is invalid (ErrInvalidArgument),
	// the owner is not found (ErrOwnerNotFound), the image limit is exceeded (ErrImageLimitExceeded),
	// or a database/internal error occurs.
	AddImage(ctx context.Context, req *imagemodel.AddRequest, ownerRepo imageowner.OwnerRepo[imageowner.Owner]) error
	// DeleteImage removes an image from a single owner.
	// The owner must implement the Owner interface, and its repository
	// must implement the OwnerRepo interface.
	//
	// Returns an error if the request payload is invalid (ErrInvalidArgument), the owner is not found (ErrOwnerNotFound),
	// or a database/internal error occurs.
	DeleteImage(ctx context.Context, req *imagemodel.DeleteRequest, ownerRepo imageowner.OwnerRepo[imageowner.Owner]) error
	// AddImageBatch adds an image for a batch of owners.
	// Owners must implement Owner methods and they're repository
	// must implement OwnerRepo methods.
	//
	// It returns the number of affected owners.
	// Returns an error if no owners are found in the database (ErrOwnersNotFound), request payload is
	// invalid (ErrInvalidArgument), or a databsae/internal error occures.
	AddImageBatch(ctx context.Context, req *imagemodel.AddBatchRequest, ownerRepo imageowner.OwnerRepo[imageowner.Owner]) (int, error)
	// DeleteImageBatch removes an image from a batch of owners.
	// Owners must implement Owner methods and they're repository
	// must implement OwnerRepo methods.
	//
	// It returns the number of affected owners.
	// Returns an error if no owners are found in the database (ErrOwnersNotFound), no associations between owners and image
	// was found (ErrAssociationsNotFound), request payload is invalid (ErrInvalidArgument), or a databsae/internal error occures.
	DeleteImageBatch(ctx context.Context, req *imagemodel.DeleteBatchRequst, ownerRepo imageowner.OwnerRepo[imageowner.Owner]) (int, error)
}

// service holds [imagerepo.Repository] to perform database operations.
type service struct {
	ImageRepo imagerepo.Repository
}

// New creates a new image service instance.
func New(imageRepo imagerepo.Repository) Service {
	return &service{ImageRepo: imageRepo}
}

// AddImage adds an image for a single owner.
// The owner must implement the Owner interface, and its repository
// must implement the OwnerRepo interface.
//
// Returns an error if the request payload is invalid (ErrInvalidArgument),
// the owner is not found (ErrOwnerNotFound), the image limit is exceeded (ErrImageLimitExceeded),
// or a database/internal error occurs.
func (s *service) AddImage(ctx context.Context, req *imagemodel.AddRequest, ownerRepo imageowner.OwnerRepo[imageowner.Owner]) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	return ownerRepo.DB().Transaction(func(tx *gorm.DB) error {
		txOwnerRepo := ownerRepo.WithTx(tx)

		owner, err := txOwnerRepo.GetWithUnpublished(ctx, req.OwnerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrOwnerNotFound, err)
			}
			return fmt.Errorf("failed to retrieve owner: %w", err)
		}

		if owner.GetUploadedImageAmount() >= 5 {
			return ErrImageLimitExceeded
		}

		newImage := &imagemodel.Image{
			URL:            req.URL,
			SecureURL:      req.SecureURL,
			PublicID:       req.PublicID,
			MediaServiceID: req.MediaServiceID,
		}

		if err := txOwnerRepo.AddImage(ctx, owner, newImage); err != nil {
			return fmt.Errorf("failed to add image for owner: %w", err)
		}

		owner.SetUploadedImageAmount(owner.GetUploadedImageAmount() + 1)
		if _, err := txOwnerRepo.BatchUpdate(ctx, []imageowner.Owner{owner}, 2); err != nil {
			return fmt.Errorf("failed to update owner uploaded image count: %w", err)
		}
		return nil
	})
}

// DeleteImage removes an image from a single owner.
// The owner must implement the Owner interface, and its repository
// must implement the OwnerRepo interface.
//
// Returns an error if the request payload is invalid (ErrInvalidArgument), the owner is not found (ErrOwnerNotFound),
// or a database/internal error occurs.
func (s *service) DeleteImage(ctx context.Context, req *imagemodel.DeleteRequest, ownerRepo imageowner.OwnerRepo[imageowner.Owner]) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	return ownerRepo.DB().Transaction(func(tx *gorm.DB) error {
		txOwnerRepo := ownerRepo.WithTx(tx)

		owner, err := txOwnerRepo.GetWithUnpublished(ctx, req.OwnerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrOwnerNotFound, err)
			}
			return fmt.Errorf("failed to retrieve owner: %w", err)
		}

		if err := txOwnerRepo.DeleteImage(ctx, owner, req.MediaServiceID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrImageNotFoundOnOwner, err)
			}
			return fmt.Errorf("failed to delete image from owner: %w", err)
		}

		if _, err := txOwnerRepo.DecrementImageCount(ctx, []string{req.OwnerID}); err != nil {
			return fmt.Errorf("failed to decrement owner uploaded image count: %w", err)
		}

		return nil
	})
}

// AddImageBatch adds an image for a batch of owners.
// Owners must implement Owner methods and they're repository
// must implement OwnerRepo methods.
//
// It returns the number of affected owners.
// Returns an error if no owners are found in the database (ErrOwnersNotFound), request payload is
// invalid (ErrInvalidArgument), or a databsae/internal error occures.
func (s *service) AddImageBatch(ctx context.Context, req *imagemodel.AddBatchRequest, ownerRepo imageowner.OwnerRepo[imageowner.Owner]) (int, error) {
	affectedOwners := 0
	if err := req.Validate(); err != nil {
		return affectedOwners, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	owners, err := ownerRepo.ListWithUnpublishedByIDs(ctx, req.OwnerIDs...)
	if err != nil {
		return affectedOwners, fmt.Errorf("failed to retrieve owners: %w", err)
	}
	if len(owners) == 0 {
		return affectedOwners, fmt.Errorf("%w: %w", ErrOwnersNotFound, err)
	}

	newImage := &imagemodel.Image{
		URL:            req.URL,
		SecureURL:      req.SecureURL,
		PublicID:       req.PublicID,
		MediaServiceID: req.MediaServiceID,
	}

	var validOwners []imageowner.Owner
	for _, owner := range owners {
		if owner.GetUploadedImageAmount() < 5 {
			validOwners = append(validOwners, owner)
		}
	}

	if len(validOwners) == 0 {
		return affectedOwners, nil // No owners to update, but not an error.
	}

	err = s.ImageRepo.DB().Transaction(func(tx *gorm.DB) error {
		txOwnerRepo := ownerRepo.WithTx(tx)

		if err := txOwnerRepo.AddImageBatch(ctx, owners, newImage); err != nil {
			return fmt.Errorf("failed to batch add images for owners: %w", err)
		}

		for _, o := range validOwners {
			o.SetUploadedImageAmount(o.GetUploadedImageAmount() + 1)
		}

		if _, err := txOwnerRepo.BatchUpdate(ctx, owners, 2); err != nil {
			return fmt.Errorf("failed to batch update owners: %w", err)
		}
		affectedOwners = len(validOwners)
		return nil
	})
	if err != nil {
		return affectedOwners, err
	}
	return affectedOwners, nil
}

// DeleteImageBatch removes an image from a batch of owners.
// Owners must implement Owner methods and they're repository
// must implement OwnerRepo methods.
//
// It returns the number of affected owners.
// Returns an error if no owners are found in the database (ErrOwnersNotFound), no associations between owners and image
// was found (ErrAssociationsNotFound), request payload is invalid (ErrInvalidArgument), or a databsae/internal error occures.
func (s *service) DeleteImageBatch(ctx context.Context, req *imagemodel.DeleteBatchRequst, ownerRepo imageowner.OwnerRepo[imageowner.Owner]) (int, error) {
	affectedOwners := 0
	if err := req.Validate(); err != nil {
		return affectedOwners, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	err := s.ImageRepo.DB().Transaction(func(tx *gorm.DB) error {
		txOwnerRepo := ownerRepo.WithTx(tx)
		owners, err := txOwnerRepo.ListWithUnpublishedByIDs(ctx, req.OwnerIDs...)
		if err != nil {
			return fmt.Errorf("failed to retrieve owners: %w", err)
		}
		if len(owners) == 0 {
			return fmt.Errorf("%w: %w", ErrOwnersNotFound, err)
		}

		// Find which of the requested owners are really associated with an image.
		affectectedOwnerIDs, err := txOwnerRepo.FindOwnerIDsByImageID(ctx, req.MediaServiceID, req.OwnerIDs)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrAssociationsNotFound, err)
		}

		if err := txOwnerRepo.DeleteImageBatch(ctx, owners, &imagemodel.Image{MediaServiceID: req.MediaServiceID}); err != nil {
			return fmt.Errorf("failed to batch delete image from owners: %w", err)
		}

		// Decrement the counter if the owners was really affected.
		if len(affectectedOwnerIDs) > 0 {
			if _, err := txOwnerRepo.DecrementImageCount(ctx, affectectedOwnerIDs); err != nil {
				return fmt.Errorf("failed to decrement uploaded image count from owners: %w", err)
			}
		}
		affectedOwners = len(affectectedOwnerIDs)
		return nil
	})
	if err != nil {
		return affectedOwners, err
	}
	return affectedOwners, nil
}
