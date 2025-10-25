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

package seminar

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mikhail5545/product-service-go/internal/database/product"
	"github.com/mikhail5545/product-service-go/internal/database/seminar"
	"github.com/mikhail5545/product-service-go/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	SeminarRepo seminar.Repository
	ProductRepo product.Repository
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

func New(sr seminar.Repository, pr product.Repository) *Service {
	return &Service{
		SeminarRepo: sr,
		ProductRepo: pr,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*models.Seminar, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{
			Msg:  "Invalid seminar ID",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	seminar, err := s.SeminarRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{
				Msg:  "Seminar not found",
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return nil, &Error{
			Msg:  "Failed to retrieve seminar",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return seminar, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]models.Seminar, int64, error) {
	seminars, err := s.SeminarRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to retrieve seminars",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	total, err := s.SeminarRepo.Count(ctx)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to count seminars",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return seminars, total, nil
}

func (s *Service) Create(ctx context.Context, seminar *models.Seminar) (*models.Seminar, error) {
	err := s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		if seminar.Name == "" {
			return &Error{Msg: "seminar name cannot be empty", Err: nil, Code: http.StatusBadRequest}
		}

		seminar.ReservationProduct.ID = uuid.New().String()
		if seminar.ReservationProduct.Name == "" {
			seminar.ReservationProduct.Name = fmt.Sprintf("Закрепить место на семинар %s", seminar.Name)
		}
		if seminar.ReservationProduct.Description == "" {
			seminar.ReservationProduct.Description = seminar.Description
		}
		if seminar.ReservationProduct.Price <= 0 {
			return &Error{Msg: "seminar reservation product price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		seminar.EarlyProduct.ID = uuid.New().String()
		if seminar.EarlyProduct.Name == "" {
			seminar.EarlyProduct.Name = fmt.Sprintf("Ранняя оплата места на семинар %s", seminar.Name)
		}
		if seminar.EarlyProduct.Description == "" {
			seminar.EarlyProduct.Description = seminar.Description
		}
		if seminar.EarlyProduct.Price <= 0 {
			return &Error{Msg: "seminar early product price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		seminar.LateProduct.ID = uuid.New().String()
		if seminar.LateProduct.Name == "" {
			seminar.LateProduct.Name = fmt.Sprintf("Поздняя оплата места на семинар %s", seminar.Name)
		}
		if seminar.LateProduct.Description == "" {
			seminar.LateProduct.Description = seminar.Description
		}
		if seminar.LateProduct.Price <= 0 {
			return &Error{Msg: "seminar late product price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		seminar.EarlySurchargeProduct.ID = uuid.New().String()
		if seminar.EarlySurchargeProduct.Name == "" {
			seminar.EarlySurchargeProduct.Name = fmt.Sprintf("Ранняя доплата места на семинар %s", seminar.Name)
		}
		if seminar.EarlySurchargeProduct.Description == "" {
			seminar.EarlySurchargeProduct.Description = seminar.Description
		}
		if seminar.EarlySurchargeProduct.Price <= 0 {
			return &Error{Msg: "seminar early surcharge product price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		seminar.LateSurchargeProduct.ID = uuid.New().String()
		if seminar.LateSurchargeProduct.Name == "" {
			seminar.LateSurchargeProduct.Name = fmt.Sprintf("Поздняя доплата места на семинар %s", seminar.Name)
		}
		if seminar.LateSurchargeProduct.Description == "" {
			seminar.LateSurchargeProduct.Description = seminar.Description
		}
		if seminar.LateSurchargeProduct.Price <= 0 {
			return &Error{Msg: "seminar late surcharge product price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		if err := txProductRepo.Create(ctx, seminar.ReservationProduct); err != nil {
			return &Error{Msg: "failed to create seminar reservation product", Err: err, Code: http.StatusInternalServerError}
		}
		if err := txProductRepo.Create(ctx, seminar.EarlyProduct); err != nil {
			return &Error{Msg: "failed to create seminar early product", Err: err, Code: http.StatusInternalServerError}
		}
		if err := txProductRepo.Create(ctx, seminar.LateProduct); err != nil {
			return &Error{Msg: "failed to create seminar late product", Err: err, Code: http.StatusInternalServerError}
		}
		if err := txProductRepo.Create(ctx, seminar.EarlySurchargeProduct); err != nil {
			return &Error{Msg: "failed to create seminar early surcharge product", Err: err, Code: http.StatusInternalServerError}
		}
		if err := txProductRepo.Create(ctx, seminar.LateSurchargeProduct); err != nil {
			return &Error{Msg: "failed to create seminar late surcharge product", Err: err, Code: http.StatusInternalServerError}
		}

		seminar.ReservationProductID = seminar.ReservationProduct.ID
		seminar.EarlyProductID = seminar.EarlyProduct.ID
		seminar.LateProductID = seminar.LateProduct.ID
		seminar.EarlySurchargeProductID = seminar.EarlySurchargeProduct.ID
		seminar.LateSurchargeProductID = seminar.LateSurchargeProduct.ID

		seminar.ID = uuid.New().String()

		if err := txSeminarRepo.Create(ctx, seminar); err != nil {
			return &Error{Msg: "failed to create seminar", Err: err, Code: http.StatusInternalServerError}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return seminar, nil
}

func (s *Service) Update(ctx context.Context, req *models.Seminar, id string) (map[string]map[string]any, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{Msg: "invalid seminar id", Err: err, Code: http.StatusBadRequest}
	}

	updates := make(map[string]map[string]any)
	err := s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		seminar, err := txSeminarRepo.Get(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{
					Msg:  "seminar not found",
					Err:  err,
					Code: http.StatusNotFound,
				}
			}
			return &Error{
				Msg:  "failed to find seminar",
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}

		if req.Name != "" && req.Name != seminar.Name {
			updates["seminar"]["name"] = req.Name
		}
		if req.Description != "" && req.Description != seminar.Description {
			updates["seminar"]["description"] = req.Description
		}
		if req.Place != "" && req.Place != seminar.Place {
			updates["seminar"]["place"] = req.Place
		}
		if time.Now().Before(req.Date) && req.Date != seminar.Date {
			updates["seminar"]["date"] = req.Date
		}
		if req.EndingDate.After(time.Now()) && req.EndingDate != seminar.EndingDate {
			updates["seminar"]["ending_date"] = req.EndingDate
		}
		if time.Now().Before(req.LatePaymentDate) && req.LatePaymentDate != seminar.LatePaymentDate {
			updates["seminar"]["late_payment_date"] = req.LatePaymentDate
		}
		if req.Details != "" && req.Details != seminar.Details {
			updates["seminar"]["details"] = req.Details
		}
		if req.ReservationProduct.Price != 0 && req.ReservationProduct.Price != seminar.ReservationProduct.Price {
			updates["reservation_product"]["price"] = req.ReservationProduct.Price
		}
		if req.ReservationProduct.Name != "" && req.ReservationProduct.Name != seminar.ReservationProduct.Name {
			updates["reservation_product"]["name"] = req.ReservationProduct.Name
		}
		if req.ReservationProduct.Description != "" && req.ReservationProduct.Description != seminar.ReservationProduct.Description {
			updates["reservation_product"]["description"] = req.ReservationProduct.Description
		}
		if req.EarlyProduct.Price != 0 && req.EarlyProduct.Price != seminar.EarlyProduct.Price {
			updates["early_product"]["price"] = req.EarlyProduct.Price
		}
		if req.EarlyProduct.Name != "" && req.EarlyProduct.Name != seminar.EarlyProduct.Name {
			updates["early_product"]["name"] = req.EarlyProduct.Name
		}
		if req.EarlyProduct.Description != "" && req.EarlyProduct.Description != seminar.EarlyProduct.Description {
			updates["early_product"]["description"] = req.EarlyProduct.Description
		}
		if req.LateProduct.Price != 0 && req.LateProduct.Price != seminar.LateProduct.Price {
			updates["late_product"]["price"] = req.LateProduct.Price
		}
		if req.LateProduct.Name != "" && req.LateProduct.Name != seminar.LateProduct.Name {
			updates["late_product"]["name"] = req.LateProduct.Name
		}
		if req.LateProduct.Description != "" && req.LateProduct.Description != seminar.LateProduct.Description {
			updates["late_product"]["description"] = req.LateProduct.Description
		}
		if req.EarlySurchargeProduct.Price != 0 && req.EarlySurchargeProduct.Price != seminar.EarlySurchargeProduct.Price {
			updates["early_surcharge_product"]["price"] = req.EarlySurchargeProduct.Price
		}
		if req.EarlySurchargeProduct.Name != "" && req.EarlySurchargeProduct.Name != seminar.EarlySurchargeProduct.Name {
			updates["early_surcharge_product"]["name"] = req.EarlySurchargeProduct.Name
		}
		if req.EarlySurchargeProduct.Description != "" && req.EarlySurchargeProduct.Description != seminar.EarlySurchargeProduct.Description {
			updates["early_surcharge_product"]["description"] = req.EarlySurchargeProduct.Description
		}
		if req.LateSurchargeProduct.Price != 0 && req.LateSurchargeProduct.Price != seminar.LateSurchargeProduct.Price {
			updates["late_surcharge_product"]["price"] = req.LateSurchargeProduct.Price
		}
		if req.LateSurchargeProduct.Name != "" && req.LateSurchargeProduct.Name != seminar.LateSurchargeProduct.Name {
			updates["late_surcharge_product"]["name"] = req.LateSurchargeProduct.Name
		}
		if req.LateSurchargeProduct.Description != "" && req.LateSurchargeProduct.Description != seminar.LateSurchargeProduct.Description {
			updates["late_surcharge_product"]["description"] = req.LateSurchargeProduct.Description
		}

		if val, ok := updates["seminar"]["name"].(string); ok {
			if req.ReservationProduct.Name == "" {
				updates["reservation_product"]["name"] = fmt.Sprintf("Зарезервировать место на семинар %s", val)
			}
			if req.EarlyProduct.Name == "" {
				updates["early_product"]["name"] = fmt.Sprintf("Ранняя оплата семинара %s", val)
			}
			if req.LateProduct.Name == "" {
				updates["late_product"]["name"] = fmt.Sprintf("Поздняя оплата семинара %s", val)
			}
			if req.EarlySurchargeProduct.Name == "" {
				updates["early_surcharge_product"]["name"] = fmt.Sprintf("Ранняя доплата на семинар %s", val)
			}
			if req.LateSurchargeProduct.Name == "" {
				updates["late_surcharge_product"]["name"] = fmt.Sprintf("Поздняя доплата на семинар %s", val)
			}
		}

		if val, ok := updates["seminar"]["description"].(string); ok {
			if req.ReservationProduct.Description == "" {
				updates["reservation_product"]["description"] = val
			}
			if req.EarlyProduct.Description == "" {
				updates["early_product"]["description"] = val
			}
			if req.LateProduct.Description == "" {
				updates["late_product"]["description"] = val
			}
			if req.EarlySurchargeProduct.Description == "" {
				updates["early_surcharge_product"]["description"] = val
			}
			if req.LateSurchargeProduct.Description == "" {
				updates["late_surcharge_product"]["description"] = val
			}
		}

		if len(updates["reservation_product"]) > 0 {
			if _, err := txProductRepo.Update(ctx, seminar.ReservationProduct, updates["reservation_product"]); err != nil {
				return &Error{
					Msg:  "Failed to update seminar reservation product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		if len(updates["early_product"]) > 0 {
			if _, err := txProductRepo.Update(ctx, seminar.EarlyProduct, updates["early_product"]); err != nil {
				return &Error{
					Msg:  "Failed to update seminar early product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		if len(updates["late_product"]) > 0 {
			if _, err := txProductRepo.Update(ctx, seminar.LateProduct, updates["late_product"]); err != nil {
				return &Error{
					Msg:  "Failed to update seminar late product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		if len(updates["early_surcharge_product"]) > 0 {
			if _, err := txProductRepo.Update(ctx, seminar.EarlySurchargeProduct, updates["early_surcharge_product"]); err != nil {
				return &Error{
					Msg:  "Failed to update seminar early surcharge product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		if len(updates["late_surcharge_product"]) > 0 {
			if _, err := txProductRepo.Update(ctx, seminar.LateSurchargeProduct, updates["late_surcharge_product"]); err != nil {
				return &Error{
					Msg:  "Failed to update seminar late surcharge product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		if len(updates["seminar"]) > 0 {
			if _, err := txSeminarRepo.Update(ctx, seminar, updates["seminar"]); err != nil {
				return &Error{
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
	return updates, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &Error{Msg: "invalid seminar id", Err: err, Code: http.StatusBadRequest}
	}
	return s.SeminarRepo.Delete(ctx, id)
}
