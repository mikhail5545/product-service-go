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

// Package trainingsession provides service-layer business logic for training sessions.
package trainingsession

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	productrepo "github.com/mikhail5545/product-service-go/internal/database/product"
	trainingsessionrepo "github.com/mikhail5545/product-service-go/internal/database/training_session"
	productmodel "github.com/mikhail5545/product-service-go/internal/models/product"
	trainingsessionmodel "github.com/mikhail5545/product-service-go/internal/models/training_session"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/services/training_session_mock/service_mock.go -package=training_session_mock . Service

// Service provides service-layer business logic for training session models.
type Service interface {
	// Get retrieves a single published and not soft-deleted training session record from the database,
	// along with its associated product details (price and product ID).
	//
	// Returns a TrainingSessionDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	Get(ctx context.Context, id string) (*trainingsessionmodel.TrainingSessionDetails, error)
	// GetWithDeleted retrieves a single training session record from the database, including soft-deleted ones,
	// along with its associated product details (price and product ID).
	//
	// Returns a TrainingSessionDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	GetWithDeleted(ctx context.Context, id string) (*trainingsessionmodel.TrainingSessionDetails, error)
	// GetWithUnpublished retrieves a single training session record from the database, including unpublished ones (but not soft-deleted),
	// along with its associated product details (price and product ID).
	//
	// Returns a TrainingSessionDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	GetWithUnpublished(ctx context.Context, id string) (*trainingsessionmodel.TrainingSessionDetails, error)
	// List retrieves a paginated list of all published and not soft-deleted training session records.
	// Each record is returned with its associated product details.
	//
	// Returns a slice of TrainingSessionDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs.
	List(ctx context.Context, limit, offset int) ([]trainingsessionmodel.TrainingSessionDetails, int64, error)
	// ListDeleted retrieves a paginated list of all soft-deleted physical training session.
	// Each record is returned with its associated product details.
	//
	// Returns a slice of TrainingSessionDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs.
	ListDeleted(ctx context.Context, limit, offset int) ([]trainingsessionmodel.TrainingSessionDetails, int64, error)
	// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) training session records.
	// Each record is returned with its associated product details.
	//
	// Returns a slice of TrainingSessionDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs.
	ListUnpublished(ctx context.Context, limit, offset int) ([]trainingsessionmodel.TrainingSessionDetails, int64, error)
	// Create creates a new TrainingSession record and its associated Product record in the database.
	// It validates the request payload to ensure all required fields are present.
	// Both the training session and the product are created in an unpublished state (`InStock: false`).
	//
	// Returns a CreateResponse containing the newly created TrainingSessionID and ProductID.
	// Returns an error if the request payload is invalid (ErrInvalidArgument) or a database/internal error occurs.
	Create(ctx context.Context, req *trainingsessionmodel.CreateRequest) (*trainingsessionmodel.CreateResponse, error)
	// Publish sets the `InStock` field to true for a training session and its associated product,
	// making it available in the catalog.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Publish(ctx context.Context, id string) error
	// Unpublish sets the `InStock` field to false for a training session and its associated product,
	// archiving it from the catalog.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Unpublish(ctx context.Context, id string) error
	// Update performs a partial update of a training session and its related product.
	// The request should contain the training session's ID and the fields to be updated.
	// At least one field must be provided for an update to occur.
	//
	// Returns a map containing the fields that were actually changed, nested under "training_session" and "product" keys.
	// Example: `{"training_session": {"name": "new name"}, "product": {"price": 99.99}}`
	// Returns an error if the request payload is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Update(ctx context.Context, req *trainingsessionmodel.UpdateRequest) (map[string]any, error)
	// Delete performs a soft-delete of a training session and its related product record.
	// It also unpublishes both records, meaning they must be manually published again after restoration.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Delete(ctx context.Context, id string) error
	// DeletePermanent performs a complete delete of a training session and its related product record.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	DeletePermanent(ctx context.Context, id string) error
	// Restore performs a restore of a training session and its related product record.
	// Training session and its related product record are not being published. This should be
	// done manually.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Restore(ctx context.Context, id string) error
}

// service provides service-layer business logic for training session models.
// It holds [trainingsessionrepo.Repository] and [productrepo.Repository] instances
// to perform database operations.
type service struct {
	TrainingSessionRepo trainingsessionrepo.Repository
	ProductRepo         productrepo.Repository
}

// New creates a new service instance with provided training session and product repositories.
func New(tsr trainingsessionrepo.Repository, pr productrepo.Repository) Service {
	return &service{
		TrainingSessionRepo: tsr,
		ProductRepo:         pr,
	}
}

// Get retrieves a single published and not soft-deleted training session record from the database,
// along with its associated product details (price and product ID).
//
// Returns a TrainingSessionDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Get(ctx context.Context, id string) (*trainingsessionmodel.TrainingSessionDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	trainingSession, err := s.TrainingSessionRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to get training session: %w", err)
	}
	product, err := s.ProductRepo.SelectByDetailsID(ctx, trainingSession.ID, "id", "price")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to get training session product: %w", err)
	}
	return &trainingsessionmodel.TrainingSessionDetails{
		TrainingSession: trainingSession,
		Price:           product.Price,
		ProductID:       product.ID,
	}, nil
}

// GetWithDeleted retrieves a single training session record from the database, including soft-deleted ones,
// along with its associated product details (price and product ID).
//
// Returns a TrainingSessionDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) GetWithDeleted(ctx context.Context, id string) (*trainingsessionmodel.TrainingSessionDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	trainingSession, err := s.TrainingSessionRepo.GetWithDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to get training session: %w", err)
	}
	product, err := s.ProductRepo.SelectWithDeletedByDetailsID(ctx, trainingSession.ID, "id", "price")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to get training session product: %w", err)
	}
	return &trainingsessionmodel.TrainingSessionDetails{
		TrainingSession: trainingSession,
		Price:           product.Price,
		ProductID:       product.ID,
	}, nil
}

// GetWithUnpublished retrieves a single training session record from the database, including unpublished ones (but not soft-deleted),
// along with its associated product details (price and product ID).
//
// Returns a TrainingSessionDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) GetWithUnpublished(ctx context.Context, id string) (*trainingsessionmodel.TrainingSessionDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	trainingSession, err := s.TrainingSessionRepo.GetWithUnpublished(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to get training session: %w", err)
	}
	product, err := s.ProductRepo.SelectWithUnpublishedByDetailsID(ctx, trainingSession.ID, "id", "price")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to get training session product: %w", err)
	}
	return &trainingsessionmodel.TrainingSessionDetails{
		TrainingSession: trainingSession,
		Price:           product.Price,
		ProductID:       product.ID,
	}, nil
}

// List retrieves a paginated list of all published and not soft-deleted training session records.
// Each record is returned with its associated product details.
//
// Returns a slice of TrainingSessionDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs.
func (s *service) List(ctx context.Context, limit, offset int) ([]trainingsessionmodel.TrainingSessionDetails, int64, error) {
	trainingSessions, err := s.TrainingSessionRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get training sessions: %w", err)
	}

	var tsIDs []string
	// Create a map for quick product lookup by ID
	sessionMap := make(map[string]*trainingsessionmodel.TrainingSession, len(trainingSessions))
	for i := range trainingSessions {
		sessionMap[trainingSessions[i].ID] = &trainingSessions[i]
		tsIDs = append(tsIDs, trainingSessions[i].ID)
	}

	products, err := s.ProductRepo.SelectByDetailsIDs(ctx, tsIDs, "id", "price", "details_id")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products: %w", err)
	}

	total, err := s.TrainingSessionRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count training sessions: %w", err)
	}

	var allDetails []trainingsessionmodel.TrainingSessionDetails
	for _, p := range products {
		allDetails = append(allDetails, trainingsessionmodel.TrainingSessionDetails{
			TrainingSession: sessionMap[p.DetailsID],
			Price:           p.Price,
			ProductID:       p.ID,
		})
	}
	return allDetails, total, nil
}

// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) training session records.
// Each record is returned with its associated product details.
//
// Returns a slice of TrainingSessionDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs.
func (s *service) ListUnpublished(ctx context.Context, limit, offset int) ([]trainingsessionmodel.TrainingSessionDetails, int64, error) {
	trainingSessions, err := s.TrainingSessionRepo.ListUnpublished(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get training sessions: %w", err)
	}
	total, err := s.TrainingSessionRepo.CountUnpublished(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count training sessions: %w", err)
	}

	var tsIDs []string
	// Create a map for quick product lookup by ID
	sessionMap := make(map[string]*trainingsessionmodel.TrainingSession, len(trainingSessions))
	for i := range trainingSessions {
		sessionMap[trainingSessions[i].ID] = &trainingSessions[i]
		tsIDs = append(tsIDs, trainingSessions[i].ID)
	}

	products, err := s.ProductRepo.SelectWithUnpublishedByDetailsIDs(ctx, tsIDs, "id", "price", "details_id")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products: %w", err)
	}
	var allDetails []trainingsessionmodel.TrainingSessionDetails
	for _, p := range products {
		allDetails = append(allDetails, trainingsessionmodel.TrainingSessionDetails{
			TrainingSession: sessionMap[p.DetailsID],
			Price:           p.Price,
			ProductID:       p.ID,
		})
	}
	return allDetails, total, nil
}

// ListDeleted retrieves a paginated list of all soft-deleted physical training session.
// Each record is returned with its associated product details.
//
// Returns a slice of TrainingSessionDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs.
func (s *service) ListDeleted(ctx context.Context, limit, offset int) ([]trainingsessionmodel.TrainingSessionDetails, int64, error) {
	trainingSessions, err := s.TrainingSessionRepo.ListDeleted(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get training sessions: %w", err)
	}
	total, err := s.TrainingSessionRepo.CountDeleted(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count training sessions: %w", err)
	}

	var tsIDs []string
	// Create a map for quick product lookup by ID
	sessionMap := make(map[string]*trainingsessionmodel.TrainingSession, len(trainingSessions))
	for i := range trainingSessions {
		sessionMap[trainingSessions[i].ID] = &trainingSessions[i]
		tsIDs = append(tsIDs, trainingSessions[i].ID)
	}

	products, err := s.ProductRepo.SelectWithDeletedByDetailsIDs(ctx, tsIDs, "id", "price", "details_id")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products: %w", err)
	}
	var allDetails []trainingsessionmodel.TrainingSessionDetails
	for _, p := range products {
		allDetails = append(allDetails, trainingsessionmodel.TrainingSessionDetails{
			TrainingSession: sessionMap[p.DetailsID],
			Price:           p.Price,
			ProductID:       p.ID,
		})
	}
	return allDetails, total, nil
}

// Create creates a new TrainingSession record and its associated Product record in the database.
// It validates the request payload to ensure all required fields are present.
// Both the training session and the product are created in an unpublished state (`InStock: false`).
//
// Returns a CreateResponse containing the newly created TrainingSessionID and ProductID.
// Returns an error if the request payload is invalid (ErrInvalidArgument) or a database/internal error occurs.
func (s *service) Create(ctx context.Context, req *trainingsessionmodel.CreateRequest) (*trainingsessionmodel.CreateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w, %w", ErrInvalidArgument, err)
	}

	var tsID, productID string
	err := s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txTSRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		ts := &trainingsessionmodel.TrainingSession{
			ID:               uuid.New().String(),
			Name:             req.Name,
			ShortDescription: req.ShortDescription,
			DurationMinutes:  req.DurationMinutes,
			Format:           req.Format,
			InStock:          false,
		}

		product := &productmodel.Product{
			ID:          uuid.New().String(),
			Price:       req.Price,
			DetailsID:   ts.ID,
			DetailsType: "training_session",
			InStock:     false,
		}

		if err := txTSRepo.Create(ctx, ts); err != nil {
			return fmt.Errorf("failed to create training session: %w", err)
		}

		if err := txProductRepo.Create(ctx, product); err != nil {
			return fmt.Errorf("failed to create training session product: %w", err)
		}
		tsID = ts.ID
		productID = product.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &trainingsessionmodel.CreateResponse{
		ID:        tsID,
		ProductID: productID,
	}, nil
}

// Publish sets the `InStock` field to true for a training session and its associated product,
// making it available in the catalog.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Publish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txTrainingSessionRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txTrainingSessionRepo.SetInStock(ctx, id, true)
		if err != nil {
			return fmt.Errorf("failed to publish training session: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		ra, err = txProductRepo.SetInStockByDetailsID(ctx, id, true)
		if err != nil {
			return fmt.Errorf("failed to publich training session product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil
	})
}

// Unpublish sets the `InStock` field to false for a training session and its associated product,
// archiving it from the catalog.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Unpublish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txTrainingSessionRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txTrainingSessionRepo.SetInStock(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to unpublish training session: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		ra, err = txProductRepo.SetInStockByDetailsID(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to unpublich training session product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil
	})
}

// Update performs a partial update of a training session and its related product.
// The request should contain the training session's ID and the fields to be updated.
// At least one field must be provided for an update to occur.
//
// Returns a map containing the fields that were actually changed, nested under "training_session" and "product" keys.
// Example: `{"training_session": {"name": "new name"}, "product": {"price": 99.99}}`
// Returns an error if the request payload is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Update(ctx context.Context, req *trainingsessionmodel.UpdateRequest) (map[string]any, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	updates := make(map[string]any)
	err := s.TrainingSessionRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txTSRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		ts, err := txTSRepo.Select(ctx, req.ID, "name", "short_description", "long_description", "duration_minutes", "format", "tags")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to get training session: %w", err)
		}
		product, err := txProductRepo.SelectByDetailsID(ctx, ts.ID, "id", "price")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to get training session product: %w", err)
		}

		tsUpdates := make(map[string]any)
		productUpdates := make(map[string]any)
		if req.Name != nil && *req.Name != ts.Name {
			tsUpdates["name"] = *req.Name
		}
		if req.ShortDescription != nil && *req.ShortDescription != ts.ShortDescription {
			tsUpdates["short_description"] = *req.ShortDescription
		}
		if req.LongDescription != nil && *req.LongDescription != ts.LongDescription {
			tsUpdates["long_description"] = *req.LongDescription
		}
		if req.DurationMinutes != nil && *req.DurationMinutes != ts.DurationMinutes {
			tsUpdates["duration_minutes"] = *req.DurationMinutes
		}
		if req.Format != nil && *req.Format != ts.Format {
			tsUpdates["format"] = *req.Format
		}
		if req.Price != nil && *req.Price != product.Price {
			productUpdates["price"] = *req.Price
		}
		if len(req.Tags) > 0 {
			tsUpdates["tags"] = req.Tags
		}

		if len(productUpdates) > 0 {
			if _, err := txProductRepo.Update(ctx, product, productUpdates); err != nil {
				return fmt.Errorf("failed to update training session: %w", err)
			}
		}

		if len(tsUpdates) > 0 {
			if _, err := txTSRepo.Update(ctx, ts, tsUpdates); err != nil {
				return fmt.Errorf("failed to update training session product: %w", err)
			}
		}
		updates["training_session"] = tsUpdates
		updates["product"] = productUpdates
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updates, nil
}

// Delete performs a soft-delete of a training session and its related product record.
// It also unpublishes both records, meaning they must be manually published again after restoration.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSessionRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		// Check if training session exists
		if _, err := txSessionRepo.GetWithUnpublished(ctx, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to get training session: %w", err)
		}

		// Unpublish all instances
		if _, err := txSessionRepo.SetInStock(ctx, id, false); err != nil {
			return fmt.Errorf("failed to unpublish training session: %w", err)
		}
		ra, err := txProductRepo.SetInStockByDetailsID(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to unpublish training session product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		if _, err := txSessionRepo.Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete training session: %w", err)
		}
		if _, err := txProductRepo.DeleteByDetailsID(ctx, id); err != nil {
			return fmt.Errorf("failed to delete training session product: %w", err)
		}
		return nil
	})
}

// DeletePermanent performs a complete delete of a training session and its related product record.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) DeletePermanent(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSessionRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		ra, err := txSessionRepo.DeletePermanent(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete training session: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		ra, err = txProductRepo.DeletePermanentByDetailsID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete training session product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil
	})
}

// Restore performs a restore of a training session and its related product record.
// Training session and its related product record are not being published. This should be
// done manually.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Restore(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSessionRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txSessionRepo.Restore(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to restore training session: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		ra, err = txProductRepo.RestoreByDetailsID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to restore training session product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil
	})
}
