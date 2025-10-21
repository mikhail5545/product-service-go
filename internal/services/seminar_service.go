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

type SeminarService struct {
	SeminarRepo database.SeminarRepository
	ProductRepo database.ProductRepository
}

type SeminarServiceError struct {
	Msg  string
	Err  error
	Code int
}

func (e *SeminarServiceError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func (e *SeminarServiceError) Unwrap() error {
	return e.Err
}

func (e *SeminarServiceError) GetCode() int {
	return e.Code
}

func NewSeminarService(sr database.SeminarRepository, pr database.ProductRepository) *SeminarService {
	return &SeminarService{
		SeminarRepo: sr,
		ProductRepo: pr,
	}
}

func (s *SeminarService) GetSeminar(ctx context.Context, id string) (*models.Seminar, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &SeminarServiceError{
			Msg:  "Invalid seminar ID",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	seminar, err := s.SeminarRepo.Find(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &SeminarServiceError{
				Msg:  "Seminar not found",
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return nil, &SeminarServiceError{
			Msg:  "Failed to retrieve seminar",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return seminar, nil
}

func (s *SeminarService) GetSeminars(ctx context.Context, limit, offset int) ([]models.Seminar, int64, error) {
	seminars, err := s.SeminarRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, &SeminarServiceError{
			Msg:  "Failed to retrieve seminars",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	total, err := s.SeminarRepo.Count(ctx)
	if err != nil {
		return nil, 0, &SeminarServiceError{
			Msg:  "Failed to count seminars",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return seminars, total, nil
}

func (s *SeminarService) CreateSeminar(ctx context.Context, seminar *models.Seminar) (*models.Seminar, error) {
	err := s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		if seminar.Name == "" {
			return &SeminarServiceError{Msg: "seminar name cannot be empty", Err: nil, Code: http.StatusBadRequest}
		}

		seminar.ReservationProduct.ID = uuid.New().String()
		if seminar.ReservationProduct.Name == "" {
			seminar.ReservationProduct.Name = fmt.Sprintf("Закрепить место на семинар %s", seminar.Name)
		}
		if seminar.ReservationProduct.Description == "" {
			seminar.ReservationProduct.Description = seminar.Description
		}
		if seminar.ReservationProduct.Price <= 0 {
			return &SeminarServiceError{Msg: "seminar reservation product price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		seminar.EarlyProduct.ID = uuid.New().String()
		if seminar.EarlyProduct.Name == "" {
			seminar.EarlyProduct.Name = fmt.Sprintf("Ранняя оплата места на семинар %s", seminar.Name)
		}
		if seminar.EarlyProduct.Description == "" {
			seminar.EarlyProduct.Description = seminar.Description
		}
		if seminar.EarlyProduct.Price <= 0 {
			return &SeminarServiceError{Msg: "seminar early product price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		seminar.LateProduct.ID = uuid.New().String()
		if seminar.LateProduct.Name == "" {
			seminar.LateProduct.Name = fmt.Sprintf("Поздняя оплата места на семинар %s", seminar.Name)
		}
		if seminar.LateProduct.Description == "" {
			seminar.LateProduct.Description = seminar.Description
		}
		if seminar.LateProduct.Price <= 0 {
			return &SeminarServiceError{Msg: "seminar late product price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		seminar.EarlySurchargeProduct.ID = uuid.New().String()
		if seminar.EarlySurchargeProduct.Name == "" {
			seminar.EarlySurchargeProduct.Name = fmt.Sprintf("Ранняя доплата места на семинар %s", seminar.Name)
		}
		if seminar.EarlySurchargeProduct.Description == "" {
			seminar.EarlySurchargeProduct.Description = seminar.Description
		}
		if seminar.EarlySurchargeProduct.Price <= 0 {
			return &SeminarServiceError{Msg: "seminar early surcharge product price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		seminar.LateSurchargeProduct.ID = uuid.New().String()
		if seminar.LateSurchargeProduct.Name == "" {
			seminar.LateSurchargeProduct.Name = fmt.Sprintf("Поздняя доплата места на семинар %s", seminar.Name)
		}
		if seminar.LateSurchargeProduct.Description == "" {
			seminar.LateSurchargeProduct.Description = seminar.Description
		}
		if seminar.LateSurchargeProduct.Price <= 0 {
			return &SeminarServiceError{Msg: "seminar late surcharge product price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		if err := txProductRepo.Create(ctx, seminar.ReservationProduct); err != nil {
			return &SeminarServiceError{Msg: "failed to create seminar reservation product", Err: err, Code: http.StatusInternalServerError}
		}
		if err := txProductRepo.Create(ctx, seminar.EarlyProduct); err != nil {
			return &SeminarServiceError{Msg: "failed to create seminar early product", Err: err, Code: http.StatusInternalServerError}
		}
		if err := txProductRepo.Create(ctx, seminar.LateProduct); err != nil {
			return &SeminarServiceError{Msg: "failed to create seminar late product", Err: err, Code: http.StatusInternalServerError}
		}
		if err := txProductRepo.Create(ctx, seminar.EarlySurchargeProduct); err != nil {
			return &SeminarServiceError{Msg: "failed to create seminar early surcharge product", Err: err, Code: http.StatusInternalServerError}
		}
		if err := txProductRepo.Create(ctx, seminar.LateSurchargeProduct); err != nil {
			return &SeminarServiceError{Msg: "failed to create seminar late surcharge product", Err: err, Code: http.StatusInternalServerError}
		}

		seminar.ReservationProductID = seminar.ReservationProduct.ID
		seminar.EarlyProductID = seminar.EarlyProduct.ID
		seminar.LateProductID = seminar.LateProduct.ID
		seminar.EarlySurchargeProductID = seminar.EarlySurchargeProduct.ID
		seminar.LateSurchargeProductID = seminar.LateSurchargeProduct.ID

		seminar.ID = uuid.New().String()

		if err := txSeminarRepo.Create(ctx, seminar); err != nil {
			return &SeminarServiceError{Msg: "failed to create seminar", Err: err, Code: http.StatusInternalServerError}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return seminar, nil
}

func (s *SeminarService) UpdateSeminar(ctx context.Context, seminar *models.Seminar, id string) (*models.Seminar, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &SeminarServiceError{Msg: "invalid seminar id", Err: err, Code: http.StatusBadRequest}
	}

	var seminarToUpdate *models.Seminar
	err := s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		var findErr error
		seminarToUpdate, findErr = txSeminarRepo.Find(ctx, id)
		if findErr != nil {
			if errors.Is(findErr, gorm.ErrRecordNotFound) {
				return &SeminarServiceError{Msg: "seminar not found", Err: findErr, Code: http.StatusNotFound}
			}
			return &SeminarServiceError{Msg: "failed to find seminar", Err: findErr, Code: http.StatusInternalServerError}
		}

		var seminarUpdated bool
		var reservationProductUpdated bool
		var earlyProductUpdated bool
		var lateProductUpdated bool
		var earlySurchargeProductUpdated bool
		var lateSurchargeProductUpdated bool

		if seminar.Name != "" && seminar.Name != seminarToUpdate.Name {
			seminarToUpdate.Name = seminar.Name
			seminarUpdated = true
		}
		if seminar.Description != "" && seminar.Description != seminarToUpdate.Description {
			seminarToUpdate.Description = seminar.Description
			seminarUpdated = true
		}
		if seminar.Place != "" && seminar.Place != seminarToUpdate.Place {
			seminarToUpdate.Place = seminar.Place
			seminarUpdated = true
		}
		if time.Now().Before(seminar.Date) && seminar.Date != seminarToUpdate.Date {
			seminarToUpdate.Date = seminar.Date
			seminarUpdated = true
		}
		if seminar.EndingDate.After(time.Now()) && seminar.EndingDate != seminarToUpdate.EndingDate {
			seminarToUpdate.EndingDate = seminar.EndingDate
			seminarUpdated = true
		}
		if seminar.Details != "" && seminar.Details != seminarToUpdate.Details {
			seminarToUpdate.Details = seminar.Details
			seminarUpdated = true
		}
		if seminar.ReservationProduct.Price != 0 && seminar.ReservationProduct.Price != seminarToUpdate.ReservationProduct.Price {
			seminarToUpdate.ReservationProduct.Price = seminar.ReservationProduct.Price
			reservationProductUpdated = true
		}
		if seminar.EarlyProduct.Price != 0 && seminar.EarlyProduct.Price != seminarToUpdate.EarlyProduct.Price {
			seminarToUpdate.EarlyProduct.Price = seminar.EarlyProduct.Price
			earlyProductUpdated = true
		}
		if seminar.LateProduct.Price != 0 && seminar.LateProduct.Price != seminarToUpdate.LateProduct.Price {
			seminarToUpdate.LateProduct.Price = seminar.LateProduct.Price
			lateProductUpdated = true
		}
		if seminar.EarlySurchargeProduct.Price != 0 && seminar.EarlySurchargeProduct.Price != seminarToUpdate.EarlySurchargeProduct.Price {
			seminarToUpdate.EarlySurchargeProduct.Price = seminar.EarlySurchargeProduct.Price
			earlySurchargeProductUpdated = true
		}
		if seminar.LateSurchargeProduct.Price != 0 && seminar.LateSurchargeProduct.Price != seminarToUpdate.LateSurchargeProduct.Price {
			seminarToUpdate.LateSurchargeProduct.Price = seminar.LateSurchargeProduct.Price
			lateSurchargeProductUpdated = true
		}

		if reservationProductUpdated {
			if err := txProductRepo.Update(ctx, seminarToUpdate.ReservationProduct); err != nil {
				return &SeminarServiceError{
					Msg:  "Failed to update seminar reservation product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		if earlyProductUpdated {
			if err := txProductRepo.Update(ctx, seminarToUpdate.EarlyProduct); err != nil {
				return &SeminarServiceError{
					Msg:  "Failed to update seminar early product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		if lateProductUpdated {
			if err := txProductRepo.Update(ctx, seminarToUpdate.LateProduct); err != nil {
				return &SeminarServiceError{
					Msg:  "Failed to update seminar late product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		if earlySurchargeProductUpdated {
			if err := txProductRepo.Update(ctx, seminarToUpdate.EarlySurchargeProduct); err != nil {
				return &SeminarServiceError{
					Msg:  "Failed to update seminar early surcharge product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		if lateSurchargeProductUpdated {
			if err := txProductRepo.Update(ctx, seminarToUpdate.LateSurchargeProduct); err != nil {
				return &SeminarServiceError{
					Msg:  "Failed to update seminar late surcharge product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		if seminarUpdated {
			if err := txSeminarRepo.Update(ctx, seminarToUpdate); err != nil {
				return &SeminarServiceError{
					Msg:  "Failed to update seminar",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return seminarToUpdate, nil
}

func (s *SeminarService) DeleteSeminar(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &SeminarServiceError{Msg: "invalid seminar id", Err: err, Code: http.StatusBadRequest}
	}
	return s.SeminarRepo.Delete(ctx, id)
}
