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

// Package course provides repository-layer logic for course models.
package course

import (
	"context"

	coursemodel "github.com/mikhail5545/product-service-go/internal/models/course"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/database/course_mock/repo_mock.go -package=course_mock github.com/mikhail5545/product-service-go/internal/database/course Repository

// Repository defines the interface for course data operations.
type Repository interface {
	// --- Only published and not soft-deleted ---

	// Get retrieves single course record from the database.
	Get(ctx context.Context, id string) (*coursemodel.Course, error)
	// Select retrieves specified course fields from the database.
	Select(ctx context.Context, id string, fields ...string) (*coursemodel.Course, error)
	// GetReduced retrieves single course record withound any course parts.
	GetReduced(ctx context.Context, id string) (*coursemodel.Course, error)
	// List retrieves all course records from the database without any course parts.
	List(ctx context.Context, limit, offset int) ([]coursemodel.Course, error)
	// Count counts the total number of course records in the database.
	Count(ctx context.Context) (int64, error)

	// --- With soft-deleted, if soft-deleted then also unpublished ---

	// GetWithDeleted retrieves single course record from the database including soft-deleted courses.
	GetWithDeleted(ctx context.Context, id string) (*coursemodel.Course, error)
	// GetReducedWithDeleted retrieves single course record withound any course parts including soft-deleted courses.
	GetReducedWithDeleted(ctx context.Context, id string) (*coursemodel.Course, error)
	// ListDeleted retrieves all soft-deleted course records from database without any course parts.
	ListDeleted(ctx context.Context, limit, offset int) ([]coursemodel.Course, error)
	// CountDeleted counts the total number of soft-deleted Course records in the database.
	CountDeleted(ctx context.Context) (int64, error)

	// --- With unpublished, but not soft-deleted ---

	// GetWithUnpublished retrieves single course record from the database including unpublished courses.
	GetWithUnpublished(ctx context.Context, id string) (*coursemodel.Course, error)
	// GetReducedWithDeleted retrieves single course record withound any course parts including soft-deleted courses.
	GetReducedWithUnpublished(ctx context.Context, id string) (*coursemodel.Course, error)
	// ListUnpublished retrieves all unpublished course records from database without any course parts.
	ListUnpublished(ctx context.Context, limit, offset int) ([]coursemodel.Course, error)
	// CountUnpublished counts the total number of unpublished course records in the database.
	CountUnpublished(ctx context.Context) (int64, error)

	// --- Common ---

	// Create creates a new Course record in the database.
	Create(ctx context.Context, course *coursemodel.Course) error
	// SetInStock sets new value for course's InStock field.
	SetInStock(ctx context.Context, id string, inStock bool) (int64, error)
	// Update performs partial update of Course record in the database using updates.
	Update(ctx context.Context, course *coursemodel.Course, updates any) (int64, error)
	// AddImage adds a new image for the Course record in the database.
	AddImage(ctx context.Context, course *coursemodel.Course, image *imagemodel.Image) error
	// DeleteImage deletes an image from the course record.
	DeleteImage(ctx context.Context, course *coursemodel.Course, mediaSvcID string) error
	// Delete performs soft-delete of Course record.
	Delete(ctx context.Context, id string) (int64, error)
	// DeletePermanent performs permanent delete of course record.
	DeletePermanent(ctx context.Context, id string) (int64, error)
	// Restore restores soft-deleted course record.
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

// New creates a new GORM-based Course repository.
func New(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// DB returns the underlying gorm.DB instance.
func (r *gormRepository) DB() *gorm.DB {
	return r.db
}

// WithTx returns a new repository instance with the given transaction.
func (r *gormRepository) WithTx(tx *gorm.DB) Repository {
	return &gormRepository{db: tx}
}

// --- Only published and not soft-deleted ---

// Get retrieves single Course record from the database.
func (r *gormRepository) Get(ctx context.Context, id string) (*coursemodel.Course, error) {
	var course coursemodel.Course
	err := r.db.WithContext(ctx).Preload("CourseParts").Preload("Images").First(&course, "id = ?", id).Error
	return &course, err
}

// Select retrieves specified course fields from the database.
func (r *gormRepository) Select(ctx context.Context, id string, fields ...string) (*coursemodel.Course, error) {
	var course coursemodel.Course
	err := r.db.WithContext(ctx).Model(&coursemodel.Course{}).Preload("CourseParts").Select(fields).Where("id = ?", id).First(&course).Error
	return &course, err
}

// GetReduced retrieves single course record withound any course parts.
func (r *gormRepository) GetReduced(ctx context.Context, id string) (*coursemodel.Course, error) {
	var course coursemodel.Course
	err := r.db.WithContext(ctx).Preload("Images").First(&course, "id = ?", id).Error
	return &course, err
}

// List retrieves all course records from the database without any course parts.
func (r *gormRepository) List(ctx context.Context, limit, offset int) ([]coursemodel.Course, error) {
	var courses []coursemodel.Course
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at desc").Find(&courses).Error
	return courses, err
}

// Count counts the total number of course records in the database.
func (r *gormRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&coursemodel.Course{}).Count(&count).Error
	return count, err
}

// --- With soft-deleted, if soft-deleted then also unpublished ---

// GetWithDeleted retrieves single course record from the database including soft-deleted courses.
func (r *gormRepository) GetWithDeleted(ctx context.Context, id string) (*coursemodel.Course, error) {
	var course coursemodel.Course
	err := r.db.WithContext(ctx).Unscoped().Preload("CourseParts").Preload("Images").First(&course, "id = ?", id).Error
	return &course, err
}

// GetReducedWithDeleted retrieves course data withound any Course Parts including soft-deleted ones.
func (r *gormRepository) GetReducedWithDeleted(ctx context.Context, id string) (*coursemodel.Course, error) {
	var course coursemodel.Course
	err := r.db.WithContext(ctx).Unscoped().Preload("Images").First(&course, "id = ?", id).Error
	return &course, err
}

// ListDeleted retrieves all soft-deleted course records from database without any course parts.
func (r *gormRepository) ListDeleted(ctx context.Context, limit, offset int) ([]coursemodel.Course, error) {
	var courses []coursemodel.Course
	err := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Preload("Images").Limit(limit).Offset(offset).Order("created_at desc").Find(&courses).Error
	return courses, err
}

// CountDeleted counts the total number of soft-deleted Course records in the database.
func (r *gormRepository) CountDeleted(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Unscoped().Model(&coursemodel.Course{}).Where("deleted_at IS NOT NULL").Count(&count).Error
	return count, err
}

// --- With unpublished, but not soft-deleted ---

// GetWithUnpublished retrieves single course record from the database including unpublished courses.
func (r *gormRepository) GetWithUnpublished(ctx context.Context, id string) (*coursemodel.Course, error) {
	var course coursemodel.Course
	err := r.db.WithContext(ctx).Preload("CourseParts").Preload("Images").First(&course, id).Error
	return &course, err
}

// GetReducedWithDeleted retrieves single course record withound any course parts including soft-deleted courses.
func (r *gormRepository) GetReducedWithUnpublished(ctx context.Context, id string) (*coursemodel.Course, error) {
	var course coursemodel.Course
	err := r.db.WithContext(ctx).Preload("Images").First(&course, id).Error
	return &course, err
}

// ListUnpublished retrieves all unpublished course records from database without any course parts.
func (r *gormRepository) ListUnpublished(ctx context.Context, limit, offset int) ([]coursemodel.Course, error) {
	var courses []coursemodel.Course
	err := r.db.WithContext(ctx).
		Model(&coursemodel.Course{}).
		Where("in_stock = ?", false).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&courses).Error
	return courses, err
}

// CountUnpublished counts the total number of unpublished course records in the database.
func (r *gormRepository) CountUnpublished(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&coursemodel.Course{}).Where("in_stock = ?", false).Count(&count).Error
	return count, err
}

// --- Common ---

// Create creates a new Course record in the database.
func (r *gormRepository) Create(ctx context.Context, course *coursemodel.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

// SetInStock sets new value for course's InStock field.
func (r *gormRepository) SetInStock(ctx context.Context, id string, inStock bool) (int64, error) {
	res := r.db.WithContext(ctx).Model(&coursemodel.Course{}).Where("id = ?", id).Update("in_stock", inStock)
	return res.RowsAffected, res.Error
}

// Update performs partial update of Course record in the database using updates.
func (r *gormRepository) Update(ctx context.Context, course *coursemodel.Course, updates any) (int64, error) {
	res := r.db.WithContext(ctx).Model(course).Updates(updates)
	return res.RowsAffected, res.Error
}

// AddImage adds a new image for the Course record in the database.
func (r *gormRepository) AddImage(ctx context.Context, course *coursemodel.Course, image *imagemodel.Image) error {
	return r.db.WithContext(ctx).Model(course).Association("Images").Append(image)
}

// DeleteImage deletes an image from the course record.
func (r *gormRepository) DeleteImage(ctx context.Context, course *coursemodel.Course, mediaSvcID string) error {
	return r.db.WithContext(ctx).Model(course).Association("Images").Delete(&imagemodel.Image{MediaServiceID: mediaSvcID})
}

// Delete performs soft-delete of Course record.
func (r *gormRepository) Delete(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Delete(&coursemodel.Course{}, id)
	return res.RowsAffected, res.Error
}

// DeletePermanent performs permanent delete of course record.
func (r *gormRepository) DeletePermanent(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Delete(&coursemodel.Course{}, id)
	return res.RowsAffected, res.Error
}

// Restore restores soft-deleted course record.
func (r *gormRepository) Restore(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Model(&coursemodel.Course{}).Where("id = ?", id).Update("deleted_at", nil)
	return res.RowsAffected, res.Error
}
