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

package trainingsession

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mikhail5545/product-service-go/internal/database/product"
	trainingsession "github.com/mikhail5545/product-service-go/internal/database/training_session"
	"github.com/mikhail5545/product-service-go/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	TrainingSessionRepo trainingsession.Repository
	ProductRepo         product.Repository
}

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

func New(tsr trainingsession.Repository, pr product.Repository) *Service {
	return &Service{
		TrainingSessionRepo: tsr,
		ProductRepo:         pr,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*models.TrainingSession, error) {
	trainingSession, err := s.TrainingSessionRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{
				Msg:  "Training session not found",
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return nil, &Error{
			Msg:  "Failed to get training session",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return trainingSession, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]models.TrainingSession, int64, error) {
	trainingSessions, err := s.TrainingSessionRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to get training sessions",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	total, err := s.TrainingSessionRepo.Count(ctx)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to count training sessions",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return trainingSessions, total, nil
}

func (s *Service) Create(ctx context.Context, ts *models.TrainingSession) (*models.TrainingSession, error) {
	err := s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txTSRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		ts.Product.ID = uuid.New().String()
		ts.Product.CreatedAt = time.Now()
		ts.Product.UpdatedAt = time.Now()
		ts.Product.ProductType = "training_session"

		if err := txProductRepo.Create(ctx, ts.Product); err != nil {
			return &Error{
				Msg:  "Failed to create underlying product for training session",
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}

		ts.ProductID = ts.Product.ID
		ts.CreatedAt = time.Now()
		ts.UpdatedAt = time.Now()
		ts.ID = uuid.New().String()

		if err := txTSRepo.Create(ctx, ts); err != nil {
			return &Error{
				Msg:  "Failed to create training session",
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *Service) Update(ctx context.Context, req *models.TrainingSession, id string) (map[string]any, map[string]any, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, nil, &Error{
			Msg:  "Invalid training session id",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	updates := make(map[string]any)
	productUpdates := make(map[string]any)
	err := s.TrainingSessionRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txTSRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		ts, err := txTSRepo.Get(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{
					Msg:  "Training session not found",
					Err:  err,
					Code: http.StatusNotFound,
				}
			}
			return &Error{
				Msg:  "Failed to get training session",
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}

		if ts.DurationMinutes != req.DurationMinutes && req.DurationMinutes >= 30 {
			updates["duration_minutes"] = req.DurationMinutes
		}
		if req.Format != "" && (req.Format == "offline" || req.Format == "online") && ts.Format != req.Format {
			updates["format"] = req.Format
		}
		if req.Product.Name != "" && ts.Product.Name != req.Product.Name {
			productUpdates["name"] = req.Product.Name
		}
		if req.Product.Description != "" && ts.Product.Description != req.Product.Description {
			productUpdates["description"] = req.Product.Description
		}
		if ts.Product.Price != req.Product.Price && req.Product.Price >= 0 {
			productUpdates["price"] = req.Product.Price
		}
		// TODO: Add implementation to add new image with media-service-go

		if len(productUpdates) > 0 {
			if _, err := txProductRepo.Update(ctx, ts.Product, productUpdates); err != nil {
				return &Error{
					Msg:  "Failed to update training session product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}

		if len(updates) > 0 {
			if _, err := txTSRepo.Update(ctx, ts, updates); err != nil {
				return &Error{
					Msg:  "Failed to update training session",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return updates, productUpdates, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &Error{
			Msg:  "Invalid training session id",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	return s.TrainingSessionRepo.Delete(ctx, id)
}
