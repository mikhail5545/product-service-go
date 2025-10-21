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

package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mikhail5545/product-service-go/internal/database"
	"github.com/mikhail5545/product-service-go/internal/models"
	"gorm.io/gorm"
)

type TrainingSessionService struct {
	TrainingSessionRepo database.TrainingSessionRepository
	ProductRepo         database.ProductRepository
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

func NewTrainingSessionService(tsr database.TrainingSessionRepository, pr database.ProductRepository) *TrainingSessionService {
	return &TrainingSessionService{
		TrainingSessionRepo: tsr,
		ProductRepo:         pr,
	}
}

func (s *TrainingSessionService) GetTrainingSession(ctx context.Context, id string) (*models.TrainingSession, error) {
	trainingSession, err := s.TrainingSessionRepo.Find(ctx, id)
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
	trainingSessions, err := s.TrainingSessionRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, &TrainingSessionServiceError{
			Msg:  "Failed to get training sessions",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	total, err := s.TrainingSessionRepo.Count(ctx)
	if err != nil {
		return nil, 0, &TrainingSessionServiceError{
			Msg:  "Failed to count training sessions",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return trainingSessions, total, nil
}

func (s *TrainingSessionService) CreateTrainingSession(ctx context.Context, ts *models.TrainingSession) (*models.TrainingSession, error) {
	err := s.TrainingSessionRepo.DB().Transaction(func(tx *gorm.DB) error {
		txTSRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		ts.Product.ID = uuid.New().String()
		ts.Product.CreatedAt = time.Now()
		ts.Product.UpdatedAt = time.Now()
		ts.Product.ProductType = "training_session"

		if err := txProductRepo.Create(ctx, ts.Product); err != nil {
			return &TrainingSessionServiceError{
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

func (s *TrainingSessionService) UpdateTrainingSession(ctx context.Context, ts *models.TrainingSession, id string) (*models.TrainingSession, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &TrainingSessionServiceError{
			Msg:  "Invalid training session id",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	var updatedTs *models.TrainingSession
	err := s.TrainingSessionRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txTSRepo := s.TrainingSessionRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		tsToUpdate, findErr := txTSRepo.Find(ctx, id)
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

		if tsToUpdate.DurationMinutes != ts.DurationMinutes && ts.DurationMinutes >= 30 {
			tsToUpdate.DurationMinutes = ts.DurationMinutes
			tsUpdated = true
		}
		if ts.Format != "" && tsToUpdate.Format != ts.Format {
			tsToUpdate.Format = ts.Format
			tsUpdated = true
		}
		if ts.Product.Name != "" && tsToUpdate.Product.Name != ts.Product.Name {
			tsToUpdate.Product.Name = ts.Product.Name
			tsProductUpdated = true
		}
		if ts.Product.Description != "" && tsToUpdate.Product.Description != ts.Product.Description {
			tsToUpdate.Product.Description = ts.Product.Description
			tsProductUpdated = true
		}
		if tsToUpdate.Product.Price != ts.Product.Price && ts.Product.Price >= 0 {
			tsToUpdate.Product.Price = ts.Product.Price
			tsProductUpdated = true
		}

		if tsProductUpdated {
			if err := txProductRepo.Update(ctx, ts.Product); err != nil {
				return &TrainingSessionServiceError{
					Msg:  "Failed to update training session product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
			ts.Product.UpdatedAt = time.Now()
		}

		if tsUpdated || tsProductUpdated {
			ts.UpdatedAt = time.Now()
			if err := txTSRepo.Update(ctx, ts); err != nil {
				return &TrainingSessionServiceError{
					Msg:  "Failed to update training session",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		updatedTs = tsToUpdate
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

	return s.TrainingSessionRepo.Delete(ctx, id)
}
