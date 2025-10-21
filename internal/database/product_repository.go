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

package database

import (
	"context"

	"github.com/mikhail5545/product-service-go/internal/models"
	"gorm.io/gorm"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	// Read operations
	Find(ctx context.Context, id string) (*models.Product, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.Product, error)
	FindByType(ctx context.Context, productType string, limit, offset int) ([]models.Product, error)
	Count(ctx context.Context) (int64, error)
	CountByType(ctx context.Context, productType string) (int64, error)

	// Write operations
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id string) error

	DB() *gorm.DB
	WithTx(tx *gorm.DB) ProductRepository
}

type gormProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new GORM-based product repository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &gormProductRepository{db: db}
}

// DB returns the underlying gorm.DB instance.
func (r *gormProductRepository) DB() *gorm.DB {
	return r.db
}

// WithTx returns a new repository instance with the given transaction.
func (r *gormProductRepository) WithTx(tx *gorm.DB) ProductRepository {
	return &gormProductRepository{db: tx}
}

// Find retrieves a single product by it's ID.
func (r *gormProductRepository) Find(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	return &product, err
}

func (r *gormProductRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at desc").Find(&products).Error
	return products, err
}

func (r *gormProductRepository) FindByType(ctx context.Context, productType string, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).Where("product_type = ?", productType).Limit(limit).Offset(offset).Order("created_at desc").Find(&products).Error
	return products, err
}

func (r *gormProductRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&count).Error
	return count, err
}

func (r *gormProductRepository) CountByType(ctx context.Context, productType string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Where("product_type = ?", productType).Count(&count).Error
	return count, err
}

func (r *gormProductRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *gormProductRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *gormProductRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}
