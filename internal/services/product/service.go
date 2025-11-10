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

// Package product provides service-layer business logic for products.
package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	productrepo "github.com/mikhail5545/product-service-go/internal/database/product"
	productmodel "github.com/mikhail5545/product-service-go/internal/models/product"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/services/product_mock/service_mock.go -package=product_mock . Service

// Service provides service-layer business logic for product models.
type Service interface {
	// Get retrieves a single published and not soft-deleted product record from the database.
	//
	// Returns a Product struct containing the information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occures.
	Get(ctx context.Context, id string) (*productmodel.Product, error)
	// GetWithDeleted retrieves a single product record from the database, including soft-deleted ones.
	//
	// Returns a Product struct containing the information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occures.
	GetWithDeleted(ctx context.Context, id string) (*productmodel.Product, error)
	// GetWithUnpublished retrieves a single product record from the database, including unpublished ones (but not soft-deleted).
	//
	// Returns a Product struct containing the information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occures.
	GetWithUnpublished(ctx context.Context, id string) (*productmodel.Product, error)
	// GetByDetailsID retrieves a single published and not soft-deleted product record from the database by it's DetailsID.
	// Cannot be used for products with seminar `DetailsType`.
	//
	// Returns a Product struct containing the information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occures.
	GetByDetailsID(ctx context.Context, detailsID string) (*productmodel.Product, error)
	// GetWithDeletedByDetailsID retrieves a single product record from the database by it's DetailsID, including soft-deleted ones.
	// Cannot be used for products with seminar `DetailsType`.
	//
	// Returns a Product struct containing the information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occures.
	GetWithDeletedByDetailsID(ctx context.Context, detailsID string) (*productmodel.Product, error)
	// GetWithUnpublishedByDetailsID retrieves a single product record from the database by it's DetailsID,
	// including unpublished ones (but not soft-deleted).
	// Cannot be used for products with seminar `DetailsType`.
	//
	// Returns a Product struct containing the information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occures.
	GetWithUnpublishedByDetailsID(ctx context.Context, detailsID string) (*productmodel.Product, error)
	// List retrieves a paginated list of all published and not soft-deleted product records.
	//
	// Returns a slice of ProductDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occures.
	List(ctx context.Context, limit, offset int) ([]productmodel.Product, int64, error)
	// ListDeleted retrieves a paginated list of all soft-deleted product records.
	//
	// Returns a slice of ProductDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occures.
	ListDeleted(ctx context.Context, limit, offset int) ([]productmodel.Product, int64, error)
	// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) product records.
	//
	// Returns a slice of ProductDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occures.
	ListUnpublished(ctx context.Context, limit, offset int) ([]productmodel.Product, int64, error)
	// List retrieves a paginated list of all published and not soft-deleted product records with specified DetailsType.
	//
	// Returns a slice of ProductDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occures.
	ListByDetailsType(ctx context.Context, detailsType string, limit, offset int) ([]productmodel.Product, int64, error)
}

// service provides service-layer business logic for product models.
// It holds [productrepo.Repository] instance
// to perform database operations.
type service struct {
	Repo productrepo.Repository
}

// New creates a new service instance with provided product repository.
func New(pr productrepo.Repository) Service {
	return &service{Repo: pr}
}

// Get retrieves a single published and not soft-deleted product record from the database.
//
// Returns a Product struct containing the information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occures.
func (s *service) Get(ctx context.Context, id string) (*productmodel.Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	product, err := s.Repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}
	return product, nil
}

// GetWithDeleted retrieves a single product record from the database, including soft-deleted ones.
//
// Returns a Product struct containing the information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occures.
func (s *service) GetWithDeleted(ctx context.Context, id string) (*productmodel.Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	product, err := s.Repo.GetWithDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}
	return product, nil
}

// GetWithUnpublished retrieves a single product record from the database, including unpublished ones (but not soft-deleted).
//
// Returns a Product struct containing the information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occures.
func (s *service) GetWithUnpublished(ctx context.Context, id string) (*productmodel.Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	product, err := s.Repo.GetWithUnpublished(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}
	return product, nil
}

// GetByDetailsID retrieves a single published and not soft-deleted product record from the database by it's DetailsID.
// Cannot be used for products with seminar `DetailsType`.
//
// Returns a Product struct containing the information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occures.
func (s *service) GetByDetailsID(ctx context.Context, detailsID string) (*productmodel.Product, error) {
	if _, err := uuid.Parse(detailsID); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	product, err := s.Repo.GetByDetailsID(ctx, detailsID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}
	return product, nil
}

// GetWithDeletedByDetailsID retrieves a single product record from the database by it's DetailsID, including soft-deleted ones.
// Cannot be used for products with seminar `DetailsType`.
//
// Returns a Product struct containing the information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occures.
func (s *service) GetWithDeletedByDetailsID(ctx context.Context, detailsID string) (*productmodel.Product, error) {
	if _, err := uuid.Parse(detailsID); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	product, err := s.Repo.GetWithDeletedByDetailsID(ctx, detailsID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}
	return product, nil
}

// GetWithUnpublishedByDetailsID retrieves a single product record from the database by it's DetailsID,
// including unpublished ones (but not soft-deleted).
// Cannot be used for products with seminar `DetailsType`.
//
// Returns a Product struct containing the information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occures.
func (s *service) GetWithUnpublishedByDetailsID(ctx context.Context, detailsID string) (*productmodel.Product, error) {
	if _, err := uuid.Parse(detailsID); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	product, err := s.Repo.GetWithUnpublishedByDetailsID(ctx, detailsID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}
	return product, nil
}

// List retrieves a paginated list of all published and not soft-deleted product records.
//
// Returns a slice of ProductDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occures.
func (s *service) List(ctx context.Context, limit, offset int) ([]productmodel.Product, int64, error) {
	products, err := s.Repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}
	total, err := s.Repo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}
	return products, total, nil
}

// ListDeleted retrieves a paginated list of all soft-deleted product records.
//
// Returns a slice of ProductDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occures.
func (s *service) ListDeleted(ctx context.Context, limit, offset int) ([]productmodel.Product, int64, error) {
	products, err := s.Repo.ListDeleted(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}
	total, err := s.Repo.CountDeleted(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}
	return products, total, nil
}

// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) product records.
//
// Returns a slice of ProductDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occures.
func (s *service) ListUnpublished(ctx context.Context, limit, offset int) ([]productmodel.Product, int64, error) {
	products, err := s.Repo.ListUnpublished(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}
	total, err := s.Repo.CountUnpublished(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}
	return products, total, nil
}

// List retrieves a paginated list of all published and not soft-deleted product records with specified DetailsType.
//
// Returns a slice of ProductDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occures.
func (s *service) ListByDetailsType(ctx context.Context, detailsType string, limit, offset int) ([]productmodel.Product, int64, error) {
	products, err := s.Repo.ListByDetailsType(ctx, detailsType, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}
	total, err := s.Repo.CountByDetailsType(ctx, detailsType)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}
	return products, total, nil
}
