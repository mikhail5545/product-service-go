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

// Package physicalgood provides repository-layer logic for physical good models.
package physicalgood

import (
	"context"

	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	physicalgoodmodel "github.com/mikhail5545/product-service-go/internal/models/physical_good"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/database/physical_good_mock/repo_mock.go -package=physical_good_mock github.com/mikhail5545/product-service-go/internal/database/physical_good Repository

// Repository defines the interface for physical good data operations.
type Repository interface {
	// --- Only published and not soft-deleted ---

	// Get retrieves a single physical good record from the database.
	Get(ctx context.Context, id string) (*physicalgoodmodel.PhysicalGood, error)
	// Select retrieves specified fields of a physical good record from the database.
	Select(ctx context.Context, id string, fields ...string) (*physicalgoodmodel.PhysicalGood, error)
	// List retrieves a paginated list of all physical good records int the database.
	List(ctx context.Context, limit, offset int) ([]physicalgoodmodel.PhysicalGood, error)
	// Count counts the total number of all the physical good records in the database.
	Count(ctx context.Context) (int64, error)

	// --- With soft-deleted, if soft-deleted then also unpublished ---

	// GetWithDeleted retrieves a single physical good record from the database including soft-deleted physial goods.
	GetWithDeleted(ctx context.Context, id string) (*physicalgoodmodel.PhysicalGood, error)
	// ListDeleted retrieves a paginated list of all soft-deleted physical good records in the database.
	ListDeleted(ctx context.Context, limit, offset int) ([]physicalgoodmodel.PhysicalGood, error)
	// CountDeleted counts the total number of all soft-deleted physical good records in the database.
	CountDeleted(ctx context.Context) (int64, error)

	// --- With unpublished, but not soft-deleted ---

	// GetWithUnpublished retrieves a single physical good record from the database including unpublished physial goods.
	GetWithUnpublished(ctx context.Context, id string) (*physicalgoodmodel.PhysicalGood, error)
	// ListUnpublished retrieves a paginated list of all unpublished physical good records in the database.
	ListUnpublished(ctx context.Context, limit, offset int) ([]physicalgoodmodel.PhysicalGood, error)
	// CountUnpublished counts the total number of all unpublished physical good records in the database.
	CountUnpublished(ctx context.Context) (int64, error)

	// --- Common ---

	// Create creates a new physical good record in the database.
	Create(ctx context.Context, ts *physicalgoodmodel.PhysicalGood) error
	// SetInStock sets a new value for physical good's InStock field.
	SetInStock(ctx context.Context, id string, inStock bool) (int64, error)
	// Update performs partial update of a physical good record using updates.
	Update(ctx context.Context, ts *physicalgoodmodel.PhysicalGood, updates any) (int64, error)
	// AddImage adds a new image for the physical good record in the database.
	AddImage(ctx context.Context, good *physicalgoodmodel.PhysicalGood, image *imagemodel.Image) error
	// DeleteImage deletes an image from the physical good record.
	DeleteImage(ctx context.Context, good *physicalgoodmodel.PhysicalGood, mediaSvcID string) error
	// Delete performs soft-delete of a physical good record.
	Delete(ctx context.Context, id string) (int64, error)
	// DeletePermanent performs permanent delete of a physical good record.
	DeletePermanent(ctx context.Context, id string) (int64, error)
	// Restore restores soft-deleted physical good record.
	Restore(ctx context.Context, id string) (int64, error)
	// DB returns the underlying gorm.DB instance.
	DB() *gorm.DB
	// WithTx returns a new repository instance with the given transaction.
	WithTx(tx *gorm.DB) Repository
}

// gormRepository holds gorm.DB for GORM-based database operations.
type gormRepository struct {
	db *gorm.DB
}

// New creates a new GORM-based physical good repository.
func New(db *gorm.DB) Repository {
	return &gormRepository{
		db: db,
	}
}

// DB returns the underlying gorm.DB instance.
func (r *gormRepository) DB() *gorm.DB {
	return r.db
}

// WithTx returns a new repository instance with the given transaction.
func (r *gormRepository) WithTx(tx *gorm.DB) Repository {
	return &gormRepository{
		db: tx,
	}
}

// --- Only published and not soft-deleted ---

// Get retrieves a single physical good record from the database.
func (r *gormRepository) Get(ctx context.Context, id string) (*physicalgoodmodel.PhysicalGood, error) {
	var good physicalgoodmodel.PhysicalGood
	err := r.db.WithContext(ctx).Preload("Images").Where("in_stock = ?", true).First(&good, id).Error
	return &good, err
}

// Select retrieves specified fields of a physical good record from the database.
func (r *gormRepository) Select(ctx context.Context, id string, fields ...string) (*physicalgoodmodel.PhysicalGood, error) {
	var good physicalgoodmodel.PhysicalGood
	err := r.db.WithContext(ctx).Model(&physicalgoodmodel.PhysicalGood{}).Select(fields).Where("id = ?", id).First(&good).Error
	return &good, err
}

// List retrieves a paginated list of all physical good records int the database.
func (r *gormRepository) List(ctx context.Context, limit, offset int) ([]physicalgoodmodel.PhysicalGood, error) {
	var goods []physicalgoodmodel.PhysicalGood
	err := r.db.WithContext(ctx).Where("in_stock = ?", true).Preload("Images").Limit(limit).Offset(offset).Order("created_at desc").Find(&goods).Error
	return goods, err
}

// Count counts the total number of all the physical good records in the database.
func (r *gormRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&physicalgoodmodel.PhysicalGood{}).Where("in_stock = ?", true).Count(&count).Error
	return count, err
}

// --- With soft-deleted, if soft-deleted then also unpublished ---

// GetWithDeleted retrieves a single physical good record from the database including soft-deleted physial goods.
func (r *gormRepository) GetWithDeleted(ctx context.Context, id string) (*physicalgoodmodel.PhysicalGood, error) {
	var good physicalgoodmodel.PhysicalGood
	err := r.db.WithContext(ctx).Unscoped().Preload("Images").First(&good, id).Error
	return &good, err
}

// ListDeleted retrieves a paginated list of all soft-deleted physical good records in the database.
func (r *gormRepository) ListDeleted(ctx context.Context, limit, offset int) ([]physicalgoodmodel.PhysicalGood, error) {
	var goods []physicalgoodmodel.PhysicalGood
	err := r.db.WithContext(ctx).Unscoped().Preload("Images").Where("deleted_at IS NOT NULL").Limit(limit).Offset(offset).Order("created_at desc").Find(&goods).Error
	return goods, err
}

// ListDeleted retrieves a paginated list of all soft-deleted physical good records in the database.
func (r *gormRepository) CountDeleted(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Unscoped().
		Model(&physicalgoodmodel.PhysicalGood{}).
		Where("deleted_at IS NOT NULL").
		Count(&count).Error
	return count, err
}

// --- With unpublished, but not soft-deleted ---

// GetWithUnpublished retrieves a single physical good record from the database including unpublished physial goods.
func (r *gormRepository) GetWithUnpublished(ctx context.Context, id string) (*physicalgoodmodel.PhysicalGood, error) {
	var good physicalgoodmodel.PhysicalGood
	err := r.db.WithContext(ctx).Preload("Images").First(&good, id).Error
	return &good, err
}

// ListUnpublished retrieves a paginated list of all unpublished physical good records in the database.
func (r *gormRepository) ListUnpublished(ctx context.Context, limit, offset int) ([]physicalgoodmodel.PhysicalGood, error) {
	var goods []physicalgoodmodel.PhysicalGood
	err := r.db.WithContext(ctx).
		Model(&physicalgoodmodel.PhysicalGood{}).
		Preload("Images").
		Where("in_stock = ?", false).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&goods).Error
	return goods, err
}

// CountUnpublished counts the total number of all unpublished physical good records in the database.
func (r *gormRepository) CountUnpublished(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&physicalgoodmodel.PhysicalGood{}).Where("in_stock = ?", false).Count(&count).Error
	return count, err
}

// --- Common ---

// Create creates a new physical good record in the database.
func (r *gormRepository) Create(ctx context.Context, good *physicalgoodmodel.PhysicalGood) error {
	return r.db.WithContext(ctx).Create(good).Error
}

// SetInStock sets a new value for physical good's InStock field.
func (r *gormRepository) SetInStock(ctx context.Context, id string, inStock bool) (int64, error) {
	res := r.db.WithContext(ctx).Model(&physicalgoodmodel.PhysicalGood{}).Where("id = ?", id).Update("in_stock", inStock)
	return res.RowsAffected, res.Error
}

// Update performs partial update of a physical good record using updates.
func (r *gormRepository) Update(ctx context.Context, good *physicalgoodmodel.PhysicalGood, updates any) (int64, error) {
	res := r.db.WithContext(ctx).Model(good).Updates(updates)
	return res.RowsAffected, res.Error
}

// AddImage adds a new image for the physical good record in the database.
func (r *gormRepository) AddImage(ctx context.Context, good *physicalgoodmodel.PhysicalGood, image *imagemodel.Image) error {
	return r.db.WithContext(ctx).Model(good).Association("Images").Append(image)
}

// DeleteImage deletes an image from the physical good record.
func (r *gormRepository) DeleteImage(ctx context.Context, good *physicalgoodmodel.PhysicalGood, mediaSvcID string) error {
	return r.db.WithContext(ctx).Model(good).Association("Images").Delete(&imagemodel.Image{MediaServiceID: mediaSvcID})
}

// Delete performs soft-delete of a physical good record.
func (r *gormRepository) Delete(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Delete(&physicalgoodmodel.PhysicalGood{}, id)
	return res.RowsAffected, res.Error
}

// DeletePermanent performs permanent delete of a physical good record.
func (r *gormRepository) DeletePermanent(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Delete(&physicalgoodmodel.PhysicalGood{}, id)
	return res.RowsAffected, res.Error
}

// Restore restores soft-deleted physical good record.
func (r *gormRepository) Restore(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Model(&physicalgoodmodel.PhysicalGood{}).Where("id = ?", id).Update("deleted_at", nil)
	return res.RowsAffected, res.Error
}
