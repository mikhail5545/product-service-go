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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	productrepo "github.com/mikhail5545/product-service-go/internal/database/product"
	trainingsessionrepo "github.com/mikhail5545/product-service-go/internal/database/training_session"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
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
	// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Get(ctx context.Context, id string) (*trainingsessionmodel.TrainingSessionDetails, error)
	// GetWithDeleted retrieves a single training session record from the database, including soft-deleted ones,
	// along with its associated product details (price and product ID).
	//
	// Returns a TrainingSessionDetails struct containing the combined information.
	// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	GetWithDeleted(ctx context.Context, id string) (*trainingsessionmodel.TrainingSessionDetails, error)
	// GetWithUnpublished retrieves a single training session record from the database, including unpublished ones (but not soft-deleted),
	// along with its associated product details (price and product ID).
	//
	// Returns a TrainingSessionDetails struct containing the combined information.
	// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	GetWithUnpublished(ctx context.Context, id string) (*trainingsessionmodel.TrainingSessionDetails, error)
	// List retrieves a paginated list of all published and not soft-deleted training session records.
	// Each record is returned with its associated product details.
	//
	// Returns a slice of TrainingSessionDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs (http.StatusInternalServerError).
	List(ctx context.Context, limit, offset int) ([]trainingsessionmodel.TrainingSessionDetails, int64, error)
	// ListDeleted retrieves a paginated list of all soft-deleted physical training session.
	// Each record is returned with its associated product details.
	//
	// Returns a slice of TrainingSessionDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs (http.StatusInternalServerError).
	ListDeleted(ctx context.Context, limit, offset int) ([]trainingsessionmodel.TrainingSessionDetails, int64, error)
	// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) training session records.
	// Each record is returned with its associated product details.
	//
	// Returns a slice of TrainingSessionDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs (http.StatusInternalServerError).
	ListUnpublished(ctx context.Context, limit, offset int) ([]trainingsessionmodel.TrainingSessionDetails, int64, error)
	// Create creates a new TrainingSession record and its associated Product record in the database.
	// It validates the request payload to ensure all required fields are present.
	// Both the training session and the product are created in an unpublished state (`InStock: false`).
	//
	// Returns a CreateResponse containing the newly created TrainingSessionID and ProductID.
	// Returns an error if the request payload is invalid (http.StatusBadRequest) or a database/internal error occurs (http.StatusInternalServerError).
	Create(ctx context.Context, req *trainingsessionmodel.CreateRequest) (*trainingsessionmodel.CreateResponse, error)
	// Publish sets the `InStock` field to true for a training session and its associated product,
	// making it available in the catalog.
	//
	// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Publish(ctx context.Context, id string) error
	// Unpublish sets the `InStock` field to false for a training session and its associated product,
	// archiving it from the catalog.
	//
	// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Unpublish(ctx context.Context, id string) error
	// Update performs a partial update of a training session and its related product.
	// The request should contain the training session's ID and the fields to be updated.
	// At least one field must be provided for an update to occur.
	//
	// Returns a map containing the fields that were actually changed, nested under "training_session" and "product" keys.
	// Example: `{"training_session": {"name": "new name"}, "product": {"price": 99.99}}`
	// Returns an error if the request payload is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Update(ctx context.Context, req *trainingsessionmodel.UpdateRequest) (map[string]any, error)
	// Delete performs a soft-delete of a training session and its related product record.
	// It also unpublishes both records, meaning they must be manually published again after restoration.
	//
	// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Delete(ctx context.Context, id string) error
	// DeletePermanent performs a complete delete of a training session and its related product record.
	//
	// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	DeletePermanent(ctx context.Context, id string) error
	// Restore performs a restore of a training session and its related product record.
	// Training session and its related product record are not being published. This should be
	// done manually.
	//
	// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Restore(ctx context.Context, id string) error
	// AddImage adds a new image to a training session. It's called by media-service-go upon successful image upload.
	// The function validates the request, checks the image limit, and appends the new image information.
	//
	// Returns an AddResponse with the MediaServiceID on success.
	// Returns an error if:
	// - The request payload is invalid (http.StatusBadRequest).
	// - The training session (owner) is not found (http.StatusNotFound).
	// - The image limit (5) is exceeded (http.StatusBadRequest).
	// - A database/internal error occurs (http.StatusInternalServerError).
	AddImage(ctx context.Context, req *imagemodel.AddRequest) (*imagemodel.AddResponse, error)
	// DeleteImage removes an image from a training session. It's called by media-service-go upon successful image deletion.
	// The function validates the request and removes the image information from the training session.
	//
	// Returns an error if:
	// - The request payload is invalid (http.StatusBadRequest).
	// - The training session (owner) or image is not found (http.StatusNotFound).
	// - A database/internal error occurs (http.StatusInternalServerError).
	DeleteImage(ctx context.Context, req *imagemodel.DeleteRequest) error
}

// service provides service-layer business logic for training session models.
// It holds [trainingsessionrepo.Repository] and [productrepo.Repository] instances
// to perform database operations.
type service struct {
	TrainingSessionRepo trainingsessionrepo.Repository
	ProductRepo         productrepo.Repository
}

// Error represents training session service error.
type Error struct {
	Msg  string
	Err  error
	Code int
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) GetCode() int {
	return e.Code
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
// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Get(ctx context.Context, id string) (*trainingsessionmodel.TrainingSessionDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{Msg: "Invalid training session ID", Err: err, Code: http.StatusBadRequest}
	}
	trainingSession, err := s.TrainingSessionRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{Msg: "Training session not found", Err: err, Code: http.StatusNotFound}
		}
		return nil, &Error{Msg: "Failed to get training session", Err: err, Code: http.StatusInternalServerError}
	}
	product, err := s.ProductRepo.SelectByDetailsID(ctx, trainingSession.ID, "id", "price")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{Msg: "Product not found", Err: err, Code: http.StatusNotFound}
		}
		return nil, &Error{Msg: "Failed to get product", Err: err, Code: http.StatusInternalServerError}
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
// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) GetWithDeleted(ctx context.Context, id string) (*trainingsessionmodel.TrainingSessionDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{Msg: "Invalid training session ID", Err: err, Code: http.StatusBadRequest}
	}
	trainingSession, err := s.TrainingSessionRepo.GetWithDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{Msg: "Training session not found", Err: err, Code: http.StatusNotFound}
		}
		return nil, &Error{Msg: "Failed to get training session", Err: err, Code: http.StatusInternalServerError}
	}
	product, err := s.ProductRepo.SelectWithDeletedByDetailsID(ctx, trainingSession.ID, "id", "price")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{Msg: "Product not found", Err: err, Code: http.StatusNotFound}
		}
		return nil, &Error{Msg: "Failed to get product", Err: err, Code: http.StatusInternalServerError}
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
// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) GetWithUnpublished(ctx context.Context, id string) (*trainingsessionmodel.TrainingSessionDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{Msg: "Invalid training session ID", Err: err, Code: http.StatusBadRequest}
	}
	trainingSession, err := s.TrainingSessionRepo.GetWithUnpublished(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{Msg: "Training session not found", Err: err, Code: http.StatusNotFound}
		}
		return nil, &Error{Msg: "Failed to get training session", Err: err, Code: http.StatusInternalServerError}
	}
	product, err := s.ProductRepo.SelectWithUnpublishedByDetailsID(ctx, trainingSession.ID, "id", "price")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{Msg: "Product not found", Err: err, Code: http.StatusNotFound}
		}
		return nil, &Error{Msg: "Failed to get product", Err: err, Code: http.StatusInternalServerError}
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
// Returns an error if a database/internal error occurs (http.StatusInternalServerError).
func (s *service) List(ctx context.Context, limit, offset int) ([]trainingsessionmodel.TrainingSessionDetails, int64, error) {
	trainingSessions, err := s.TrainingSessionRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, &Error{Msg: "Failed to get training sessions", Err: err, Code: http.StatusInternalServerError}
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
		return nil, 0, &Error{Msg: "Failed to get products", Err: err, Code: http.StatusInternalServerError}
	}

	total, err := s.TrainingSessionRepo.Count(ctx)
	if err != nil {
		return nil, 0, &Error{Msg: "Failed to count training sessions", Err: err, Code: http.StatusInternalServerError}
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
// Returns an error if a database/internal error occurs (http.StatusInternalServerError).
func (s *service) ListUnpublished(ctx context.Context, limit, offset int) ([]trainingsessionmodel.TrainingSessionDetails, int64, error) {
	trainingSessions, err := s.TrainingSessionRepo.ListUnpublished(ctx, limit, offset)
	if err != nil {
		return nil, 0, &Error{Msg: "Failed to get training sessions", Err: err, Code: http.StatusInternalServerError}
	}
	total, err := s.TrainingSessionRepo.CountUnpublished(ctx)
	if err != nil {
		return nil, 0, &Error{Msg: "Failed to count training sessions", Err: err, Code: http.StatusInternalServerError}
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
		return nil, 0, &Error{Msg: "Failed to get products", Err: err, Code: http.StatusInternalServerError}
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
// Returns an error if a database/internal error occurs (http.StatusInternalServerError).
func (s *service) ListDeleted(ctx context.Context, limit, offset int) ([]trainingsessionmodel.TrainingSessionDetails, int64, error) {
	trainingSessions, err := s.TrainingSessionRepo.ListDeleted(ctx, limit, offset)
	if err != nil {
		return nil, 0, &Error{Msg: "Failed to get training sessions", Err: err, Code: http.StatusInternalServerError}
	}
	total, err := s.TrainingSessionRepo.CountDeleted(ctx)
	if err != nil {
		return nil, 0, &Error{Msg: "Failed to count training sessions", Err: err, Code: http.StatusInternalServerError}
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
		return nil, 0, &Error{Msg: "Failed to get products", Err: err, Code: http.StatusInternalServerError}
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
// Returns an error if the request payload is invalid (http.StatusBadRequest) or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Create(ctx context.Context, req *trainingsessionmodel.CreateRequest) (*trainingsessionmodel.CreateResponse, error) {
	var tsID, productID string
	err := s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txTSRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		if err := req.Validate(); err != nil {
			return &Error{Msg: "Invalid request payload", Err: err, Code: http.StatusBadRequest}
		}

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
			return &Error{Msg: "Failed to create training session", Err: err, Code: http.StatusInternalServerError}
		}

		if err := txProductRepo.Create(ctx, product); err != nil {
			return &Error{Msg: "Failed to create underlying product for training session", Err: err, Code: http.StatusInternalServerError}
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
// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Publish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &Error{Msg: "Invalid training session id", Err: err, Code: http.StatusBadRequest}
	}
	return s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txTrainingSessionRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txTrainingSessionRepo.SetInStock(ctx, id, true)
		if err != nil {
			return &Error{Msg: "Failed to publish training session", Err: err, Code: http.StatusInternalServerError}
		} else if ra == 0 {
			return &Error{Msg: "Training session not found", Err: err, Code: http.StatusNotFound}
		}
		ra, err = txProductRepo.SetInStockByDetailsID(ctx, id, true)
		if err != nil {
			return &Error{Msg: "Failed to publish training session product", Err: err, Code: http.StatusInternalServerError}
		} else if ra == 0 {
			return &Error{Msg: "Training session product not found", Err: err, Code: http.StatusNotFound}
		}
		return nil
	})
}

// Unpublish sets the `InStock` field to false for a training session and its associated product,
// archiving it from the catalog.
//
// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Unpublish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &Error{Msg: "Invalid training session id", Err: err, Code: http.StatusBadRequest}
	}
	return s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txTrainingSessionRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txTrainingSessionRepo.SetInStock(ctx, id, false)
		if err != nil {
			return &Error{Msg: "Failed to unpublish training session", Err: err, Code: http.StatusInternalServerError}
		} else if ra == 0 {
			return &Error{Msg: "Training session not found", Err: err, Code: http.StatusNotFound}
		}
		ra, err = txProductRepo.SetInStockByDetailsID(ctx, id, false)
		if err != nil {
			return &Error{Msg: "Failed to unpublish training session product", Err: err, Code: http.StatusInternalServerError}
		} else if ra == 0 {
			return &Error{Msg: "Training session product not found", Err: err, Code: http.StatusNotFound}
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
// Returns an error if the request payload is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Update(ctx context.Context, req *trainingsessionmodel.UpdateRequest) (map[string]any, error) {
	updates := make(map[string]any)
	err := s.TrainingSessionRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txTSRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		if err := req.Validate(); err != nil {
			return &Error{Msg: "Invalid request payload", Err: err, Code: http.StatusBadRequest}
		}

		ts, err := txTSRepo.Select(ctx, req.ID, "name", "short_description", "long_description", "duration_minutes", "format", "tags")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{Msg: "Training session not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "Failed to get training session", Err: err, Code: http.StatusInternalServerError}
		}
		product, err := txProductRepo.SelectByDetailsID(ctx, ts.ID, "id", "price")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{Msg: "Product not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "Failed to get product", Err: err, Code: http.StatusInternalServerError}
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
				return &Error{Msg: "Failed to update training session product", Err: err, Code: http.StatusInternalServerError}
			}
		}

		if len(tsUpdates) > 0 {
			if _, err := txTSRepo.Update(ctx, ts, tsUpdates); err != nil {
				return &Error{Msg: "Failed to update training session", Err: err, Code: http.StatusInternalServerError}
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

// AddImage adds a new image to a training session. It's called by media-service-go upon successful image upload.
// The function validates the request, checks the image limit, and appends the new image information.
//
// Returns an AddResponse with the MediaServiceID on success.
// Returns an error if:
// - The request payload is invalid (http.StatusBadRequest).
// - The training session (owner) is not found (http.StatusNotFound).
// - The image limit (5) is exceeded (http.StatusBadRequest).
// - A database/internal error occurs (http.StatusInternalServerError).
func (s *service) AddImage(ctx context.Context, req *imagemodel.AddRequest) (*imagemodel.AddResponse, error) {
	if err := req.Validate(); err != nil {
		validationMsg, _ := json.Marshal(err)
		return nil, &Error{Msg: string(validationMsg), Err: err, Code: http.StatusBadRequest}
	}

	err := s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txTrainingSessionRepo := s.TrainingSessionRepo.WithTx(tx)

		tsRec, err := txTrainingSessionRepo.GetWithUnpublished(ctx, req.OwnerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{Msg: "Training session not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "Failed to get training session", Err: err, Code: http.StatusInternalServerError}
		}

		if tsRec.UploadedImageAmount >= 5 {
			return &Error{Msg: "Maximum number of uploaded images is 5 per item", Err: nil, Code: http.StatusBadRequest}
		}

		newImage := &imagemodel.Image{
			PublicID:       req.PublicID,
			URL:            req.URL,
			SecureURL:      req.SecureURL,
			MediaServiceID: req.MediaServiceID,
		}

		if err := txTrainingSessionRepo.AddImage(ctx, tsRec, newImage); err != nil {
			return &Error{Msg: "Failed to add image to training session", Err: err, Code: http.StatusInternalServerError}
		}

		// Increment the image count and save
		tsRec.UploadedImageAmount++
		if _, err := txTrainingSessionRepo.Update(ctx, tsRec, map[string]any{"uploaded_image_amount": tsRec.UploadedImageAmount}); err != nil {
			return &Error{Msg: "Failed to update training session", Err: err, Code: http.StatusInternalServerError}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &imagemodel.AddResponse{MediaServiceID: req.MediaServiceID}, nil
}

// DeleteImage removes an image from a training session. It's called by media-service-go upon successful image deletion.
// The function validates the request and removes the image information from the training session.
//
// Returns an error if:
// - The request payload is invalid (http.StatusBadRequest).
// - The training session (owner) or image is not found (http.StatusNotFound).
// - A database/internal error occurs (http.StatusInternalServerError).
func (s *service) DeleteImage(ctx context.Context, req *imagemodel.DeleteRequest) error {
	if err := req.Validate(); err != nil {
		validationMsg, _ := json.Marshal(err)
		return &Error{Msg: string(validationMsg), Err: err, Code: http.StatusBadRequest}
	}

	return s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txTrainingSessionRepo := s.TrainingSessionRepo.WithTx(tx)

		tsRec, err := txTrainingSessionRepo.GetWithUnpublished(ctx, req.OwnerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{Msg: "Training session not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "Failed to get training session", Err: err, Code: http.StatusInternalServerError}
		}

		var imageFound bool
		for i, img := range tsRec.Images {
			if img.MediaServiceID == req.MediaServiceID {
				imageFound = true
				// Remove the image from the slice
				tsRec.Images = append(tsRec.Images[:i], tsRec.Images[i+1:]...)
				break
			}
		}

		if !imageFound {
			return &Error{Msg: "Image not found on training session", Err: gorm.ErrRecordNotFound, Code: http.StatusNotFound}
		}

		if err := txTrainingSessionRepo.DeleteImage(ctx, tsRec, req.MediaServiceID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{Msg: "Image not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "Failed to delete training session image", Err: err, Code: http.StatusInternalServerError}
		}

		tsRec.UploadedImageAmount--

		// Save the updated image list and count
		if _, err := txTrainingSessionRepo.Update(ctx, tsRec, map[string]any{"images": tsRec.Images, "uploaded_image_amount": tsRec.UploadedImageAmount}); err != nil {
			return &Error{Msg: "Failed to update training session", Err: err, Code: http.StatusInternalServerError}
		}
		return nil
	})
}

// Delete performs a soft-delete of a training session and its related product record.
// It also unpublishes both records, meaning they must be manually published again after restoration.
//
// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &Error{Msg: "Invalid training session id", Err: err, Code: http.StatusBadRequest}
	}
	return s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSessionRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		// Check if training session exists
		if _, err := txSessionRepo.GetWithUnpublished(ctx, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{Msg: "Training session not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "Failed to get training session", Err: err, Code: http.StatusInternalServerError}
		}

		// Unpublish all instances
		if _, err := txSessionRepo.SetInStock(ctx, id, false); err != nil {
			return &Error{Msg: "Failed to unpublish training session", Err: err, Code: http.StatusInternalServerError}
		}
		ra, err := txProductRepo.SetInStockByDetailsID(ctx, id, false)
		if err != nil {
			return &Error{Msg: "Failed to unpublish training session product", Err: err, Code: http.StatusInternalServerError}
		} else if ra == 0 {
			return &Error{Msg: "Training session product not found", Err: err, Code: http.StatusNotFound}
		}
		if _, err = txSessionRepo.Delete(ctx, id); err != nil {
			return &Error{Msg: "Failed to delete training session", Err: err, Code: http.StatusInternalServerError}
		}
		if _, err = txProductRepo.DeleteByDetailsID(ctx, id); err != nil {
			return &Error{Msg: "Failed to delete training session product", Err: err, Code: http.StatusInternalServerError}
		}
		return nil
	})
}

// DeletePermanent performs a complete delete of a training session and its related product record.
//
// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) DeletePermanent(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &Error{Msg: "Invalid training session id", Err: err, Code: http.StatusBadRequest}
	}
	return s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSessionRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		ra, err := txSessionRepo.DeletePermanent(ctx, id)
		if err != nil {
			return &Error{Msg: "Failed to delete training session", Err: err, Code: http.StatusInternalServerError}
		} else if ra == 0 {
			return &Error{Msg: "Training session not found", Err: err, Code: http.StatusNotFound}
		}
		ra, err = txProductRepo.DeletePermanentByDetailsID(ctx, id)
		if err != nil {
			return &Error{Msg: "Failed to delete training session product", Err: err, Code: http.StatusInternalServerError}
		} else if ra == 0 {
			return &Error{Msg: "Training session product not found", Err: err, Code: http.StatusNotFound}
		}
		return nil
	})
}

// Restore performs a restore of a training session and its related product record.
// Training session and its related product record are not being published. This should be
// done manually.
//
// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Restore(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &Error{Msg: "Invalid training session id", Err: err, Code: http.StatusBadRequest}
	}
	return s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSessionRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		if ra, err := txSessionRepo.Restore(ctx, id); err != nil || ra == 0 {
			if ra == 0 {
				return &Error{Msg: "Training session not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "Failed to delete training session", Err: err, Code: http.StatusInternalServerError}
		}
		if ra, err := txProductRepo.RestoreByDetailsID(ctx, id); err != nil || ra == 0 {
			if ra == 0 {
				return &Error{Msg: "Training session product not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "Failed to delete training session product", Err: err, Code: http.StatusInternalServerError}
		}
		return nil
	})
}
