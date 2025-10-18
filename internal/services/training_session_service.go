// vitainmove.com/product-service-go
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

package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"vitainmove.com/product-service-go/internal/database"
	"vitainmove.com/product-service-go/internal/models"
)

type TrainingSessionService struct {
	Repo database.TrainingSessionRepository
}

type TrainingSessionServiceError struct {
	Msg  string
	Err  error
	Code int
}

func (e *TrainingSessionServiceError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func (e *TrainingSessionServiceError) Unwrap() error {
	return e.Err
}

func (e *TrainingSessionServiceError) GetCode() int {
	return e.Code
}

func NewTrainingSessionService(tsr database.TrainingSessionRepository) *TrainingSessionService {
	return &TrainingSessionService{Repo: tsr}
}

func (s *TrainingSessionService) GetTrainingSession(ctx context.Context, id string) (*models.TrainingSession, error) {
	trainingSession, err := s.Repo.Find(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &TrainingSessionServiceError{
				Msg:  "Training session not found",
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return nil, &TrainingSessionServiceError{
			Msg:  "Failed to get training session",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return trainingSession, nil
}

func (s *TrainingSessionService) GetTrainingSessions(ctx context.Context, limit, offset int) ([]models.TrainingSession, int64, error) {
	trainingSessions, err := s.Repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, &TrainingSessionServiceError{
			Msg:  "Failed to get training sessions",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	total, err := s.Repo.Count(ctx)
	if err != nil {
		return nil, 0, &TrainingSessionServiceError{
			Msg:  "Failed to count training sessions",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return trainingSessions, total, nil
}

func (s *TrainingSessionService) CreateTrainingSession(ctx context.Context, req *models.AddTrainingSessionRequest) (*models.TrainingSession, error) {
	var ts *models.TrainingSession

	err := s.Repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.Repo.WithTx(tx)

		ts = &models.TrainingSession{
			ID: uuid.New().String(),
			Product: &models.Product{
				ID:               uuid.New().String(),
				Name:             req.Product.Name,
				Description:      req.Product.Description,
				Price:            req.Product.Price,
				Amount:           0, // Training sessions are not stock-limited this way
				ShippingRequired: false,
				ProductType:      "training_session",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			DurationMinutes: req.DurationMinutes,
			Format:          req.Format,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		ts.ProductID = ts.Product.ID

		if err := txRepo.Create(ctx, ts); err != nil {
			return &TrainingSessionServiceError{
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

func (s *TrainingSessionService) UpdateTrainingSession(ctx context.Context, req *models.EditTrainingSessionRequest, id string) (*models.TrainingSession, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &TrainingSessionServiceError{
			Msg:  "Invalid training session id",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	var updatedTs *models.TrainingSession
	err := s.Repo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := s.Repo.WithTx(tx)

		var ts *models.TrainingSession
		ts, findErr := txRepo.Find(ctx, id)
		if findErr != nil {
			if errors.Is(findErr, gorm.ErrRecordNotFound) {
				return &TrainingSessionServiceError{
					Msg:  "Training session not found",
					Err:  findErr,
					Code: http.StatusNotFound,
				}
			}
			return &TrainingSessionServiceError{
				Msg:  "Failed to get training session",
				Err:  findErr,
				Code: http.StatusInternalServerError,
			}
		}

		var tsUpdated bool
		var tsProductUpdated bool

		if ts.DurationMinutes != req.DurationMinutes && req.DurationMinutes >= 30 {
			ts.DurationMinutes = req.DurationMinutes
			tsUpdated = true
		}
		if req.Format != "" && ts.Format != req.Format {
			ts.Format = req.Format
			tsUpdated = true
		}
		if req.Name != "" && ts.Product.Name != req.Name {
			ts.Product.Name = req.Name
			tsProductUpdated = true
		}
		if req.Description != "" && ts.Product.Description != req.Description {
			ts.Product.Description = req.Description
			tsProductUpdated = true
		}
		if ts.Product.Price != req.Price && req.Price >= 0 {
			ts.Product.Price = req.Price
			tsProductUpdated = true
		}

		if tsProductUpdated {
			ts.UpdatedAt = time.Now()
			ts.Product.UpdatedAt = time.Now()
			if err := txRepo.UpdateProduct(ctx, ts.Product); err != nil {
				return &TrainingSessionServiceError{
					Msg:  "Failed to update training session product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}

		if tsUpdated {
			ts.UpdatedAt = time.Now()
			if err := txRepo.Update(ctx, ts); err != nil {
				return &TrainingSessionServiceError{
					Msg:  "Failed to update training session",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		updatedTs = ts
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedTs, nil
}

func (s *TrainingSessionService) DeleteTrainingSession(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &TrainingSessionServiceError{
			Msg:  "Invalid training session id",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	return s.Repo.Delete(ctx, id)
}
