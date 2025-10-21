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

type ProductService struct {
	Repo database.ProductRepository
}

type ProductServiceError struct {
	Msg  string
	Err  error
	Code int
}

func (e *ProductServiceError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func (e *ProductServiceError) Unwrap() error {
	return e.Err
}

func (e *ProductServiceError) GetCode() int {
	return e.Code
}

func NewProductService(pr database.ProductRepository) *ProductService {
	return &ProductService{Repo: pr}
}

func (s *ProductService) GetProducts(ctx context.Context, limit, offset int) ([]models.Product, int64, error) {
	products, err := s.Repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, &ProductServiceError{
			Msg:  "Failed to get products",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	total, err := s.Repo.Count(ctx)
	if err != nil {
		return nil, 0, &ProductServiceError{
			Msg:  "Failed to get products count",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return products, total, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &ProductServiceError{
			Msg:  "Invalid product id",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	product, err := s.Repo.Find(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &ProductServiceError{
				Msg:  "Product not found",
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return nil, &ProductServiceError{
			Msg:  "Failed to get product",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return product, nil
}

func (s *ProductService) GetProductByType(ctx context.Context, productType string, limit, offset int) ([]models.Product, int64, error) {
	products, err := s.Repo.FindByType(ctx, productType, limit, offset)
	if err != nil {
		return nil, 0, &ProductServiceError{
			Msg:  "Failed to get product",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	total, err := s.Repo.CountByType(ctx, productType)
	if err != nil {
		return nil, 0, &ProductServiceError{
			Msg:  "Failed to count products",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return products, total, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	err := s.Repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.Repo.WithTx(tx)

		if product.Name == "" {
			return &ProductServiceError{Msg: "product name cannot be empty", Err: nil, Code: http.StatusBadRequest}
		}
		if product.Description == "" {
			return &ProductServiceError{Msg: "product description cannot be empty", Err: nil, Code: http.StatusBadRequest}
		}
		if product.Amount < 0 {
			return &ProductServiceError{Msg: "amount cannot be negative", Err: nil, Code: http.StatusBadRequest}
		}
		if product.Price <= 0 {
			return &ProductServiceError{Msg: "price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		product.ID = uuid.New().String()
		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()
		product.ProductType = "physical"

		if err := txRepo.Create(ctx, product); err != nil {
			return &ProductServiceError{Msg: "failed to create product", Err: err, Code: http.StatusInternalServerError}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *models.Product, id string) (*models.Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &ProductServiceError{Msg: "invalid product id", Err: err, Code: http.StatusBadRequest}
	}

	var productToUpdate *models.Product
	err := s.Repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.Repo.WithTx(tx)

		var findErr error
		productToUpdate, findErr = txRepo.Find(ctx, id)
		if findErr != nil {
			if errors.Is(findErr, gorm.ErrRecordNotFound) {
				return &ProductServiceError{Msg: "product not found", Err: findErr, Code: http.StatusNotFound}
			}
			return &ProductServiceError{Msg: "failed to find product", Err: findErr, Code: http.StatusInternalServerError}
		}

		var updated bool
		if product.Name != "" && product.Name != productToUpdate.Name {
			productToUpdate.Name = product.Name
			updated = true
		}
		if product.Price != 0 && product.Price != productToUpdate.Price {
			productToUpdate.Price = product.Price
			updated = true
		}
		if product.Amount >= 0 && product.Amount != productToUpdate.Amount { // amount can be null
			productToUpdate.Amount = product.Amount
			updated = true
		}
		// Only update if a field has changed to avoid unnecessary DB calls.
		if updated {
			productToUpdate.UpdatedAt = time.Now()
			if err := txRepo.Update(ctx, productToUpdate); err != nil {
				return &ProductServiceError{Msg: "failed to update product", Err: err, Code: http.StatusInternalServerError}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return productToUpdate, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &ProductServiceError{Msg: "invalid product id", Err: err, Code: http.StatusBadRequest}
	}
	return s.Repo.Delete(ctx, id)
}
