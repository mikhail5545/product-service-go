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

// Package trainingsession provides repository-layer logic for training session models.
package trainingsession

import (
	"context"

	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	tsmodel "github.com/mikhail5545/product-service-go/internal/models/training_session"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/database/training_session_mock/repo_mock.go -package=training_session_mock github.com/mikhail5545/product-service-go/internal/database/training_session Repository

// Repository defines the interface for training session data operations.
type Repository interface {
	// --- Only published and not soft-deleted ---

	// Get retrieves a single published and not soft-deleted training session record from the database.
	Get(ctx context.Context, id string) (*tsmodel.TrainingSession, error)
	// Select retrieves specified fields of a published and not soft-deleted training session record from the database.
	Select(ctx context.Context, id string, fields ...string) (*tsmodel.TrainingSession, error)
	// List retrieves a paginated list of all published and not soft-deleted training session records in the database.
	List(ctx context.Context, limit, offset int) ([]tsmodel.TrainingSession, error)
	// Count counts the total number of all published and not soft-deleted training session records in the database.
	Count(ctx context.Context) (int64, error)

	// --- With soft-deleted, if soft-deleted then also unpublished ---

	// GetWithDeleted retrieves a single training session record from the database, including soft-deleted ones.
	GetWithDeleted(ctx context.Context, id string) (*tsmodel.TrainingSession, error)
	// ListDeleted retrieves a paginated list of all soft-deleted training session records from the database.
	ListDeleted(ctx context.Context, limit, offset int) ([]tsmodel.TrainingSession, error)
	// CountDeleted counts the total number of all soft-deleted training session records in the database.
	CountDeleted(ctx context.Context) (int64, error)

	// --- With unpublished, but not soft-deleted ---

	// GetWithUnpublished retrieves a single training session record from the database, including unpublished ones (but not soft-deleted).
	GetWithUnpublished(ctx context.Context, id string) (*tsmodel.TrainingSession, error)
	// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) training session records from the database.
	ListUnpublished(ctx context.Context, limit, offset int) ([]tsmodel.TrainingSession, error)
	// CountUnpublished counts the total number of all unpublished (but not soft-deleted) training session records in the database.
	CountUnpublished(ctx context.Context) (int64, error)

	// --- Common ---

	// Create creates a new training session record in the database.
	Create(ctx context.Context, ts *tsmodel.TrainingSession) error
	// SetInStock sets a new value for the training session's InStock field.
	SetInStock(ctx context.Context, id string, inStock bool) (int64, error)
	// Update performs a partial update of a training session record using the provided updates map.
	Update(ctx context.Context, ts *tsmodel.TrainingSession, updates any) (int64, error)
	// AddImage adds a new image for the training session record in the database.
	AddImage(ctx context.Context, ts *tsmodel.TrainingSession, image *imagemodel.Image) error
	// DeleteImage deletes an image from the training session record.
	DeleteImage(ctx context.Context, ts *tsmodel.TrainingSession, mediaSvcID string) error
	// Delete performs a soft-delete of a training session record.
	Delete(ctx context.Context, id string) (int64, error)
	// DeletePermanent performs a permanent delete of a training session record.
	DeletePermanent(ctx context.Context, id string) (int64, error)
	// Restore restores a soft-deleted training session record.
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

// New creates a new GORM-based training session repository.
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

// Get retrieves a single published and not soft-deleted training session record from the database.
func (r *gormRepository) Get(ctx context.Context, id string) (*tsmodel.TrainingSession, error) {
	var ts tsmodel.TrainingSession
	err := r.db.WithContext(ctx).Preload("Images").Where("in_stock = ?", true).First(&ts, "id = ?", id).Error
	return &ts, err
}

// Select retrieves specified fields of a published and not soft-deleted training session record from the database.
func (r *gormRepository) Select(ctx context.Context, id string, fields ...string) (*tsmodel.TrainingSession, error) {
	var ts tsmodel.TrainingSession
	err := r.db.WithContext(ctx).Model(&tsmodel.TrainingSession{}).Where("in_stock = ?", true).Select(fields).Where("id = ?", id).First(&ts).Error
	return &ts, err
}

// List retrieves a paginated list of all published and not soft-deleted training session records in the database.
func (r *gormRepository) List(ctx context.Context, limit, offset int) ([]tsmodel.TrainingSession, error) {
	var ts []tsmodel.TrainingSession
	err := r.db.WithContext(ctx).Where("in_stock = ?", true).Preload("Images").Limit(limit).Offset(offset).Order("created_at desc").Find(&ts).Error
	return ts, err
}

// Count counts the total number of all published and not soft-deleted training session records in the database.
func (r *gormRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&tsmodel.TrainingSession{}).Where("in_stock = ?", true).Count(&count).Error
	return count, err
}

// --- With soft-deleted, if soft-deleted then also unpublished ---

// GetWithDeleted retrieves a single training session record from the database, including soft-deleted ones.
func (r *gormRepository) GetWithDeleted(ctx context.Context, id string) (*tsmodel.TrainingSession, error) {
	var ts tsmodel.TrainingSession
	err := r.db.WithContext(ctx).Unscoped().Preload("Images").First(&ts, id).Error
	return &ts, err
}

// ListDeleted retrieves a paginated list of all soft-deleted training session records in the database.
func (r *gormRepository) ListDeleted(ctx context.Context, limit, offset int) ([]tsmodel.TrainingSession, error) { // Corrected comment
	var ts []tsmodel.TrainingSession
	err := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Preload("Images").Limit(limit).Offset(offset).Order("created_at desc").Find(&ts).Error
	return ts, err
}

// CountDeleted counts the total number of all soft-deleted training session records in the database.
func (r *gormRepository) CountDeleted(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Unscoped().
		Model(&tsmodel.TrainingSession{}).
		Where("deleted_at IS NOT NULL").
		Count(&count).Error
	return count, err
}

// --- With unpublished, but not soft-deleted ---

// GetWithUnpublished retrieves a single training session record from the database, including unpublished ones (but not soft-deleted).
func (r *gormRepository) GetWithUnpublished(ctx context.Context, id string) (*tsmodel.TrainingSession, error) {
	var ts tsmodel.TrainingSession
	err := r.db.WithContext(ctx).Preload("Images").First(&ts, id).Error
	return &ts, err
}

// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) training session records in the database.
func (r *gormRepository) ListUnpublished(ctx context.Context, limit, offset int) ([]tsmodel.TrainingSession, error) {
	var ts []tsmodel.TrainingSession
	err := r.db.WithContext(ctx).
		Model(&tsmodel.TrainingSession{}).
		Preload("Images").
		Where("in_stock = ?", false).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&ts).Error
	return ts, err
}

// CountUnpublished counts the total number of all unpublished (but not soft-deleted) training session records in the database.
func (r *gormRepository) CountUnpublished(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&tsmodel.TrainingSession{}).Where("in_stock = ?", false).Count(&count).Error
	return count, err
}

// --- Common ---

// Create creates a new training session record in the database.
func (r *gormRepository) Create(ctx context.Context, ts *tsmodel.TrainingSession) error {
	return r.db.WithContext(ctx).Create(ts).Error
}

// Update performs a partial update of a training session record using the provided updates map.
func (r *gormRepository) Update(ctx context.Context, ts *tsmodel.TrainingSession, updates any) (int64, error) {
	res := r.db.WithContext(ctx).Model(ts).Updates(updates)
	return res.RowsAffected, res.Error
}

// AddImage adds a new image for the training session record in the database.
func (r *gormRepository) AddImage(ctx context.Context, ts *tsmodel.TrainingSession, image *imagemodel.Image) error {
	return r.db.WithContext(ctx).Model(ts).Association("Images").Append(image)
}

// DeleteImage deletes an image from the training session record.
func (r *gormRepository) DeleteImage(ctx context.Context, ts *tsmodel.TrainingSession, mediaSvcID string) error {
	return r.db.WithContext(ctx).Model(ts).Association("Images").Delete(&imagemodel.Image{MediaServiceID: mediaSvcID})
}

// SetInStock sets a new value for the training session's InStock field.
func (r *gormRepository) SetInStock(ctx context.Context, id string, inStock bool) (int64, error) {
	res := r.db.WithContext(ctx).Model(&tsmodel.TrainingSession{}).Where("id = ?", id).Update("in_stock", inStock)
	return res.RowsAffected, res.Error
}

// Delete performs soft-delete of a training session record.
func (r *gormRepository) Delete(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Delete(&tsmodel.TrainingSession{}, id)
	return res.RowsAffected, res.Error
}

// DeletePermanent performs permanent delete of a training session record.
func (r *gormRepository) DeletePermanent(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Delete(&tsmodel.TrainingSession{}, id)
	return res.RowsAffected, res.Error
}

// Restore restores soft-deleted training session record.
func (r *gormRepository) Restore(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Model(&tsmodel.TrainingSession{}).Where("id = ?", id).Update("deleted_at", nil)
	return res.RowsAffected, res.Error
}
