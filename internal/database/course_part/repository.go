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

// Package coursepart provides repository-layer logic for course part models.
package coursepart

import (
	"context"

	coursepartmodel "github.com/mikhail5545/product-service-go/internal/models/course_part"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/database/course_part_mock/repo_mock.go -package=course_part_mock github.com/mikhail5545/product-service-go/internal/database/course_part Repository

// Repository defines the interface for course part data operations.
type Repository interface {
	// --- Only published and not soft-deleted ---

	// Get retrieves single course part record from the database.
	Get(ctx context.Context, id string) (*coursepartmodel.CoursePart, error)
	// Select retrieves specified fields of the course part record from the database.
	Select(ctx context.Context, id string, fields ...string) (*coursepartmodel.CoursePart, error)
	// List retrieves a paginated list of all course part records in the database.
	List(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, error)
	// Count counts the total number of all course part records by courseID in the database.
	Count(ctx context.Context, courseID string) (int64, error)
	// CountQuery counts the total number of course part racords in the database by query.
	CountQuery(ctx context.Context, query any, args ...any) (int64, error)

	// --- With soft-deleted, if soft-deleted then also unpublished ---

	// GetWithDeleted retrieves single course part record from the database including soft-deleted course parts.
	GetWithDeleted(ctx context.Context, id string) (*coursepartmodel.CoursePart, error)
	// ListDeleted retrieves a paginated list of all soft-deleted course part records in database for the specific course.
	ListDeleted(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, error)
	// CountDeleted counts the total number of all soft-deleted course part records in the database for the specitic course.
	CountDeleted(ctx context.Context, courseID string) (int64, error)

	// --- With unpublished, but not soft-deleted ---

	// GetWithUnpublished retrieves single course part record record from the database including unpublished course parts.
	GetWithUnpublished(ctx context.Context, id string) (*coursepartmodel.CoursePart, error)
	// ListUnpublished retrieves a paginated list of all unpublished course part records in database for the specific course.
	ListUnpublished(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, error)
	// CountUnpublished counts the total number of all unpublished course part records in the database for the specitic course.
	CountUnpublished(ctx context.Context, courseID string) (int64, error)

	// --- Common ---

	// Create creates a new CoursePart record in the database.
	Create(ctx context.Context, coursePart *coursepartmodel.CoursePart) error
	// SetPublished sets a new value for course part's Published field.
	SetPublished(ctx context.Context, id string, published bool) (int64, error)
	// SetPublishedByCourseID sets a new value for Published field in all course parts with specified courseID.
	SetPublishedByCourseID(ctx context.Context, courseID string, published bool) (int64, error)
	// Update performs partial update of a course part record using updates.
	Update(ctx context.Context, coursePart *coursepartmodel.CoursePart, updates any) (int64, error)
	// Delete performs soft-delete of a course part record.
	Delete(ctx context.Context, id string) (int64, error)
	// DeleteByCourseID performs soft-delete for all course parts related to a course.
	DeleteByCourseID(ctx context.Context, courseID string) (int64, error)
	// DeletePermanent performs permanent delete of a course part record.
	DeletePermanent(ctx context.Context, id string) (int64, error)
	// DeletePermanentByCourseID performs permanent delete for all course parts related to a course.
	DeletePermanentByCourseID(ctx context.Context, courseID string) (int64, error)
	// Restore restores soft-deleted course part record.
	Restore(ctx context.Context, id string) (int64, error)
	// RestoreByCourseID restores all soft-deleted course parts for a given course.
	RestoreByCourseID(ctx context.Context, courseID string) (int64, error)

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

// Get retrieves single course part record from the database.
func (r *gormRepository) Get(ctx context.Context, id string) (*coursepartmodel.CoursePart, error) {
	var coursePart coursepartmodel.CoursePart
	err := r.db.WithContext(ctx).Where("published = ?", true).First(&coursePart, "id = ?", id).Error
	return &coursePart, err
}

// Select retrieves specified fields of the course part record from the database.
func (r *gormRepository) Select(ctx context.Context, id string, fields ...string) (*coursepartmodel.CoursePart, error) {
	var coursePart coursepartmodel.CoursePart
	err := r.db.WithContext(ctx).Model(&coursepartmodel.CoursePart{}).Where("published = ?", true).Select(fields).Where("id = ?", id).First(&coursePart).Error
	return &coursePart, err
}

// List retrieves a paginated list of all course part records in the database.
func (r *gormRepository) List(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, error) {
	var courseParts []coursepartmodel.CoursePart
	err := r.db.WithContext(ctx).Where("published = ?", true).Order("created_at desc").Limit(limit).Offset(offset).Find(&courseParts, "course_id = ?", courseID).Error
	return courseParts, err
}

// Count counts the total number of all course part records by courseID in the database.
func (r *gormRepository) Count(ctx context.Context, courseID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&coursepartmodel.CoursePart{}).Where("published = ?", true).Where("course_id = ?", courseID).Count(&count).Error
	return count, err
}

// CountQuery counts the total number of course part racords in the database by query.
func (r *gormRepository) CountQuery(ctx context.Context, query any, args ...any) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&coursepartmodel.CoursePart{}).Where("published = ?", true).Where(query, args).Count(&count).Error
	return count, err
}

// --- With soft-deleted, if soft-deleted then also unpublished ---

// GetWithDeleted retrieves single course part record from the database including soft-deleted course parts.
func (r *gormRepository) GetWithDeleted(ctx context.Context, id string) (*coursepartmodel.CoursePart, error) {
	var coursePart coursepartmodel.CoursePart
	err := r.db.WithContext(ctx).Unscoped().First(&coursePart, "id = ?", id).Error
	return &coursePart, err
}

// ListDeleted retrieves a paginated list of all soft-deleted course part records in database for the specific course.
func (r *gormRepository) ListDeleted(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, error) {
	var courseParts []coursepartmodel.CoursePart
	err := r.db.WithContext(ctx).Unscoped().
		Model(&coursepartmodel.CoursePart{}).
		Where("course_id = ?", courseID).
		Where("deleted_at IS NOT NULL").
		Order("created_at desc").Limit(limit).Offset(offset).
		Find(&courseParts).Error
	return courseParts, err
}

// CountDeleted counts the total number of all soft-deleted course part records in the database for the specitic course.
func (r *gormRepository) CountDeleted(ctx context.Context, courseID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Unscoped().
		Model(&coursepartmodel.CoursePart{}).
		Where("course_id = ?", courseID).
		Where("deleted_at IS NOT NULL").
		Count(&count).Error
	return count, err
}

// --- With unpublished, but not soft-deleted ---

// GetWithUnpublished retrieves single course part record record from the database including unpublished course parts.
func (r *gormRepository) GetWithUnpublished(ctx context.Context, id string) (*coursepartmodel.CoursePart, error) {
	var coursePart coursepartmodel.CoursePart
	err := r.db.WithContext(ctx).First(&coursePart, id).Error
	return &coursePart, err
}

// ListUnpublished retrieves a paginated list of all unpublished course part records in database for the specific course.
func (r *gormRepository) ListUnpublished(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, error) {
	var courseParts []coursepartmodel.CoursePart
	err := r.db.WithContext(ctx).
		Model(&coursepartmodel.CoursePart{}).
		Where("published = ?", false).
		Where("course_id = ?", courseID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&courseParts).Error
	return courseParts, err

}

// CountUnpublished counts the total number of all unpublished course part records in the database for the specitic course.
func (r *gormRepository) CountUnpublished(ctx context.Context, courseID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&coursepartmodel.CoursePart{}).
		Where("published = ?", false).
		Where("course_id = ?", courseID).
		Count(&count).Error
	return count, err
}

// --- Common ---

// Create creates a new CoursePart record in the database.
func (r *gormRepository) Create(ctx context.Context, coursePart *coursepartmodel.CoursePart) error {
	return r.db.WithContext(ctx).Create(coursePart).Error
}

// SetPublished sets a new value for course part's Published field.
func (r *gormRepository) SetPublished(ctx context.Context, id string, published bool) (int64, error) {
	res := r.db.WithContext(ctx).Model(&coursepartmodel.CoursePart{}).Where("id = ?", id).Update("published", published)
	return res.RowsAffected, res.Error
}

// SetPublishedByCourseID sets a new value for Published field in all course parts with specified courseID.
func (r *gormRepository) SetPublishedByCourseID(ctx context.Context, courseID string, published bool) (int64, error) {
	res := r.db.WithContext(ctx).Model(&coursepartmodel.CoursePart{}).Where("course_id = ?", courseID).Update("published", published)
	return res.RowsAffected, res.Error
}

// Update performs partial update of a course part record using updates.
func (r *gormRepository) Update(ctx context.Context, coursePart *coursepartmodel.CoursePart, updates any) (int64, error) {
	res := r.db.WithContext(ctx).Model(&coursepartmodel.CoursePart{}).Where("id = ?", coursePart.ID).Updates(updates)
	return res.RowsAffected, res.Error
}

// Delete performs soft-delete of a course part record.
func (r *gormRepository) Delete(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Delete(&coursepartmodel.CoursePart{}, id)
	return res.RowsAffected, res.Error
}

// Delete performs soft-delete of a course part records by course id.
func (r *gormRepository) DeleteByCourseID(ctx context.Context, courseID string) (int64, error) {
	res := r.db.WithContext(ctx).Where("course_id = ?", courseID).Delete(&coursepartmodel.CoursePart{})
	return res.RowsAffected, res.Error
}

// DeletePermanent performs permanent delete of a course part record.
func (r *gormRepository) DeletePermanent(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Delete(&coursepartmodel.CoursePart{}, id)
	return res.RowsAffected, res.Error
}

// DeletePermanent performs permanent delete of a course part records by course id.
func (r *gormRepository) DeletePermanentByCourseID(ctx context.Context, courseID string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Where("course_id = ?", courseID).Delete(&coursepartmodel.CoursePart{})
	return res.RowsAffected, res.Error
}

// Restore restores soft-deleted course part record.
func (r *gormRepository) Restore(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Model(&coursepartmodel.CoursePart{}).Where("id = ?", id).Update("deleted_at", nil)
	return res.RowsAffected, res.Error
}

// Restore restores soft-deleted course part records.
func (r *gormRepository) RestoreByCourseID(ctx context.Context, courseID string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Model(&coursepartmodel.CoursePart{}).Where("course_id = ?", courseID).Update("deleted_at", nil)
	return res.RowsAffected, res.Error
}
