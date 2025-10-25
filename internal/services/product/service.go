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

package product

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mikhail5545/product-service-go/internal/database/product"
	"github.com/mikhail5545/product-service-go/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	Repo product.Repository
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

func New(pr product.Repository) *Service {
	return &Service{Repo: pr}
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]models.Product, int64, error) {
	products, err := s.Repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to get products",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	total, err := s.Repo.Count(ctx)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to get products count",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return products, total, nil
}

func (s *Service) Get(ctx context.Context, id string) (*models.Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{
			Msg:  "Invalid product id",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	product, err := s.Repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{
				Msg:  "Product not found",
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return nil, &Error{
			Msg:  "Failed to get product",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return product, nil
}

func (s *Service) ListByType(ctx context.Context, productType string, limit, offset int) ([]models.Product, int64, error) {
	products, err := s.Repo.ListByType(ctx, productType, limit, offset)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to get product",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	total, err := s.Repo.CountByType(ctx, productType)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to count products",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return products, total, nil
}

func (s *Service) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	err := s.Repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.Repo.WithTx(tx)

		if product.Name == "" {
			return &Error{Msg: "product name cannot be empty", Err: nil, Code: http.StatusBadRequest}
		}
		if product.Description == "" {
			return &Error{Msg: "product description cannot be empty", Err: nil, Code: http.StatusBadRequest}
		}
		if product.Amount < 0 {
			return &Error{Msg: "amount cannot be negative", Err: nil, Code: http.StatusBadRequest}
		}
		if product.Price <= 0 {
			return &Error{Msg: "price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		product.ID = uuid.New().String()
		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()
		product.ProductType = "physical"

		if err := txRepo.Create(ctx, product); err != nil {
			return &Error{Msg: "failed to create product", Err: err, Code: http.StatusInternalServerError}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *Service) Update(ctx context.Context, req *models.Product, id string) (map[string]any, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{
			Msg:  "invalid product id",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	updates := make(map[string]any)
	err := s.Repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.Repo.WithTx(tx)

		product, err := txRepo.Get(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{
					Msg:  "product not found",
					Err:  err,
					Code: http.StatusNotFound,
				}
			}
			return &Error{
				Msg:  "failed to find product",
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}

		if req.Name != "" && req.Name != product.Name {
			updates["name"] = req.Name
		}
		if req.Price != 0 && req.Price != product.Price {
			updates["price"] = req.Price
		}
		if req.Amount >= 0 && req.Amount != product.Amount { // amount can be null
			updates["amount"] = req.Amount
		}
		if req.ShippingRequired != product.ShippingRequired {
			updates["shipping_required"] = req.ShippingRequired
		}
		//TODO: implement image upload via media-service-go
		if len(updates) > 0 {
			product.UpdatedAt = time.Now()
			if _, err := txRepo.Update(ctx, product, updates); err != nil {
				return &Error{Msg: "failed to update product", Err: err, Code: http.StatusInternalServerError}
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
		return &Error{Msg: "invalid product id", Err: err, Code: http.StatusBadRequest}
	}
	return s.Repo.Delete(ctx, id)
}
