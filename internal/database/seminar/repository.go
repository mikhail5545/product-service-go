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

// Package seminar provides repository-layer logic for seminar models.
package seminar

import (
	"context"
	"fmt"
	"strings"

	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	seminarmodel "github.com/mikhail5545/product-service-go/internal/models/seminar"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/database/seminar_mock/repo_mock.go -package=seminar_mock github.com/mikhail5545/product-service-go/internal/database/seminar Repository

// Repository defines the interface for seminar data operations.
type Repository interface {
	// --- Only published and not soft-deleted ---

	// Get retrieves single seminar record from the database.
	Get(ctx context.Context, id string) (*seminarmodel.Seminar, error)
	// Select retrieves specidied seminar fields from the database.
	Select(ctx context.Context, id string, fields ...string) (*seminarmodel.Seminar, error)
	// List retrieves a paginated list of all seminar records in the database.
	List(ctx context.Context, limit, offset int) ([]seminarmodel.Seminar, error)
	// Count counts the total number of all seminar records in the database.
	Count(ctx context.Context) (int64, error)

	// --- With soft-deleted, if soft-deleted then also unpublished ---

	// GetWithDeleted retrieves single seminar record from the database including soft-deleted ones.
	GetWithDeleted(ctx context.Context, id string) (*seminarmodel.Seminar, error)
	// ListDeleted retrieves a paginated list of all soft-deleted seminar records from database.
	ListDeleted(ctx context.Context, limit, offset int) ([]seminarmodel.Seminar, error)
	// CountDeleted counts the total number of all soft-deleted seminar records in the database.
	CountDeleted(ctx context.Context) (int64, error)

	// --- With unpublished, but not soft-deleted ---

	// GetWithUnpublished retrieves single seminar record from the database including unpublished seminars.
	GetWithUnpublished(ctx context.Context, id string) (*seminarmodel.Seminar, error)
	// ListWithUnpublished retrieves paginated list of all unpublished seminar records from the database.
	ListUnpublished(ctx context.Context, limit, offset int) ([]seminarmodel.Seminar, error)
	// ListWithUnpublishedByIDs retrieves seminar records by ids from database including unpublished ones.
	ListWithUnpublishedByIDs(ctx context.Context, ids ...string) ([]seminarmodel.Seminar, error)
	// CountWithUnpublished counts the total number of all unpublished seminar records in the database.
	CountUnpublished(ctx context.Context) (int64, error)

	// --- Common ---

	// Create creates a new seminar record in the database.
	Create(ctx context.Context, seminar *seminarmodel.Seminar) error
	// SetInStock sets a new value for seminar's InStock field.
	SetInStock(ctx context.Context, id string, inStock bool) (int64, error)
	// Update performs partial update of a seminar record using updates.
	Update(ctx context.Context, seminar *seminarmodel.Seminar, updates any) (int64, error)
	// BatchUpdate performs partial update for a batch of seminar records in the database.
	// Field that needs to be updated must be populated in all seminar records.
	// Opt param indicates which field needs to be updated:
	//
	//   - 0: Name
	//   - 1: ShortDescription
	//   - 2: UploadedImageCount
	BatchUpdate(ctx context.Context, updates []seminarmodel.Seminar, opt uint) (int64, error)
	// FindOwnerIDsByImageID finds all seminar IDs associated with a given image media service ID within a specific set of owners.
	FindOwnerIDsByImageID(ctx context.Context, mediaSvcID string, ownerIDs []string) ([]string, error)
	// DecrementImageCount decrements the uploaded_image_amount for the given seminar IDs.
	DecrementImageCount(ctx context.Context, seminarIDs []string) (int64, error)
	// AddImage adds a new image for the Seminar record in the database.
	AddImage(ctx context.Context, seminar *seminarmodel.Seminar, image *imagemodel.Image) error
	// AddImageBatch adds a new image (single) for the many seminar records in the database.
	AddImageBatch(ctx context.Context, seminars []seminarmodel.Seminar, image *imagemodel.Image) error
	// DeleteImageBatch deletes an image (single) from many seminar records in the database.
	// Note: This only removes the association. The caller is responsible for updating any related counters
	// within the same transaction to ensure data consistency.
	DeleteImageBatch(ctx context.Context, seminars []seminarmodel.Seminar, image *imagemodel.Image) error
	// DeleteImage deletes an image from the Seminar record.
	DeleteImage(ctx context.Context, seminar *seminarmodel.Seminar, mediaSvcID string) error
	// Delete performs soft-delete of a seminar record.
	Delete(ctx context.Context, id string) (int64, error)
	// DeletePermanent performs permanent delete of a seminar record.
	DeletePermanent(ctx context.Context, id string) (int64, error)
	// Restore restores soft-deleted seminar record.
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

// New creates a new GORM-based seminar repository.
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

// Get retrieves single seminar record from the database.
func (r *gormRepository) Get(ctx context.Context, id string) (*seminarmodel.Seminar, error) {
	var seminar *seminarmodel.Seminar
	err := r.db.WithContext(ctx).Where("in_stock = ?", true).Preload("Images").First(&seminar, "id = ?", id).Error
	return seminar, err
}

// Select retrieves specidied seminar fields from the database.
func (r *gormRepository) Select(ctx context.Context, id string, fields ...string) (*seminarmodel.Seminar, error) {
	var seminar *seminarmodel.Seminar
	err := r.db.WithContext(ctx).Model(&seminarmodel.Seminar{}).Where("in_stock = ?", true).Select(fields).Where("id = ?", id).First(&seminar).Error
	return seminar, err
}

// List retrieves a paginated list of all seminar records in the database.
func (r *gormRepository) List(ctx context.Context, limit, offset int) ([]seminarmodel.Seminar, error) {
	var seminars []seminarmodel.Seminar
	err := r.db.WithContext(ctx).Model(&seminarmodel.Seminar{}).Preload("Images").Where("in_stock = ?", true).Order("created_at desc").Limit(limit).Offset(offset).Find(&seminars).Error
	return seminars, err
}

// Count counts the total number of all seminar records in the database.
func (r *gormRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&seminarmodel.Seminar{}).Where("in_stock = ?", true).Count(&count).Error
	return count, err
}

// --- With soft-deleted, if soft-deleted then also unpublished ---

// GetWithDeleted retrieves single seminar record from the database including soft-deleted ones.
func (r *gormRepository) GetWithDeleted(ctx context.Context, id string) (*seminarmodel.Seminar, error) {
	var seminar *seminarmodel.Seminar
	err := r.db.WithContext(ctx).Unscoped().Preload("Images").First(&seminar, "id = ?", id).Error
	return seminar, err
}

// ListDeleted retrieves a paginated list of all soft-deleted seminar records from database.
func (r *gormRepository) ListDeleted(ctx context.Context, limit, offset int) ([]seminarmodel.Seminar, error) {
	var seminars []seminarmodel.Seminar
	err := r.db.WithContext(ctx).Unscoped().Preload("Images").Where("deleted_at IS NOT NULL").Order("created_at desc").Limit(limit).Offset(offset).Find(&seminars).Error
	return seminars, err
}

// CountDeleted counts the total number of all soft-deleted seminar records in the database.
func (r *gormRepository) CountDeleted(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Unscoped().
		Model(&seminarmodel.Seminar{}).
		Where("deleted_at IS NOT NULL").
		Count(&count).Error
	return count, err
}

// --- With unpublished, but not soft-deleted ---

// GetWithUnpublished retrieves single seminar record from the database including unpublished seminars.
func (r *gormRepository) GetWithUnpublished(ctx context.Context, id string) (*seminarmodel.Seminar, error) {
	var seminar seminarmodel.Seminar
	err := r.db.WithContext(ctx).Preload("Images").First(&seminar, id).Error
	return &seminar, err
}

// ListUnpublished retrieves paginated list of all unpublished seminar records from the database.
func (r *gormRepository) ListUnpublished(ctx context.Context, limit, offset int) ([]seminarmodel.Seminar, error) {
	var seminars []seminarmodel.Seminar
	err := r.db.WithContext(ctx).
		Model(&seminarmodel.Seminar{}).
		Preload("Images").
		Where("in_stock = ?", false).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&seminars).Error
	return seminars, err
}

// ListWithUnpublishedByIDs retrieves seminar records by ids from database including unpublished ones.
func (r *gormRepository) ListWithUnpublishedByIDs(ctx context.Context, ids ...string) ([]seminarmodel.Seminar, error) {
	var seminars []seminarmodel.Seminar
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&seminars).Error
	return seminars, err
}

// CountUnpublished counts the total number of all unpublished seminar records in the database.
func (r *gormRepository) CountUnpublished(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&seminarmodel.Seminar{}).Where("in_stock = ?", false).Count(&count).Error
	return count, err
}

// --- Common ---

// Create creates a new SeminarPart record in the database.
func (r *gormRepository) Create(ctx context.Context, seminar *seminarmodel.Seminar) error {
	return r.db.WithContext(ctx).Create(seminar).Error
}

// SetInStock sets a new value for seminar's InStock field.
func (r *gormRepository) SetInStock(ctx context.Context, id string, inStock bool) (int64, error) {
	res := r.db.WithContext(ctx).Model(&seminarmodel.Seminar{}).Where("id = ?", id).Update("in_stock", inStock)
	return res.RowsAffected, res.Error
}

// Update performs partial update of a seminar record using updates.
func (r *gormRepository) Update(ctx context.Context, seminar *seminarmodel.Seminar, updates any) (int64, error) {
	res := r.db.WithContext(ctx).Model(seminar).Updates(updates)
	return res.RowsAffected, res.Error
}

// BatchUpdate performs partial update for a batch of seminar records in the database.
// Field that needs to be updated must be populated in all seminar records.
// Opt param indicates which field needs to be updated:
//
//   - 0: Name
//   - 1: ShortDescription
//   - 2: UploadedImageCount
func (r *gormRepository) BatchUpdate(ctx context.Context, updates []seminarmodel.Seminar, opt uint) (int64, error) {
	if len(updates) == 0 {
		return 0, nil
	}

	var fieldName string
	var caseClauses []string
	var ids []string

	switch opt {
	case 0:
		fieldName = "name"
		for _, u := range updates {
			ids = append(ids, u.ID)
			caseClauses = append(caseClauses, fmt.Sprintf("WHEN '%s' THEN '%s'", u.ID, u.Name))
		}
	case 1:
		fieldName = "short_description"
		for _, u := range updates {
			ids = append(ids, u.ID)
			caseClauses = append(caseClauses, fmt.Sprintf("WHEN '%s' THEN '%s'", u.ID, u.ShortDescription))
		}
	case 2:
		fieldName = "uploaded_image_amount"
		for _, u := range updates {
			ids = append(ids, u.ID)
			caseClauses = append(caseClauses, fmt.Sprintf("WHEN '%s' THEN %d", u.ID, u.UploadedImageAmount))
		}
	default:
		return 0, nil
	}

	query := fmt.Sprintf(
		`UPDATE seminars SET %s = (CASE id %s END) WHERE id IN (%s)`,
		fieldName,
		strings.Join(caseClauses, " "),
		"'"+strings.Join(ids, "','")+"'",
	)

	res := r.db.WithContext(ctx).Exec(query)
	return res.RowsAffected, res.Error
}

// FindOwnerIDsByImageID finds all seminar IDs associated with a given image media service ID within a specific set of owners.
func (r *gormRepository) FindOwnerIDsByImageID(ctx context.Context, mediaSvcID string, ownerIDs []string) ([]string, error) {
	var affectedSeminarIDs []string
	joinTable := r.db.WithContext(ctx).Model(&seminarmodel.Seminar{}).Association("Images").Relationship.JoinTable
	err := r.db.WithContext(ctx).Table(joinTable.Table).
		Where("image_media_service_id = ?", mediaSvcID).
		Where("seminar_id IN ?", ownerIDs).
		Pluck("seminar_id", &affectedSeminarIDs).Error
	return affectedSeminarIDs, err
}

// DecrementImageCount decrements the uploaded_image_amount for the given seminar IDs.
func (r *gormRepository) DecrementImageCount(ctx context.Context, seminarIDs []string) (int64, error) {
	res := r.db.WithContext(ctx).
		Model(&seminarmodel.Seminar{}).
		Where("id IN ?", seminarIDs).
		UpdateColumn("uploaded_image_amount", gorm.Expr("uploaded_image_amount - 1"))
	return res.RowsAffected, res.Error
}

// AddImage adds a new image for the Seminar record in the database.
func (r *gormRepository) AddImage(ctx context.Context, seminar *seminarmodel.Seminar, image *imagemodel.Image) error {
	return r.db.WithContext(ctx).Model(seminar).Association("Images").Append(image)
}

// AddImageBatch adds a new image (single) for the many seminar records in the database.
func (r *gormRepository) AddImageBatch(ctx context.Context, seminars []seminarmodel.Seminar, image *imagemodel.Image) error {
	return r.db.WithContext(ctx).Model(&seminars).Association("Images").Append(image)
}

// DeleteImage deletes an image from the Seminar record.
func (r *gormRepository) DeleteImage(ctx context.Context, seminar *seminarmodel.Seminar, mediaSvcID string) error {
	return r.db.WithContext(ctx).Model(seminar).Association("Images").Delete(&imagemodel.Image{MediaServiceID: mediaSvcID})
}

// DeleteImageBatch deletes an image (single) from many seminar records in the database.
func (r *gormRepository) DeleteImageBatch(ctx context.Context, seminars []seminarmodel.Seminar, image *imagemodel.Image) error {
	return r.db.WithContext(ctx).Model(&seminars).Association("Images").Delete(image)
}

// Delete performs soft-delete of a seminar record.
func (r *gormRepository) Delete(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Delete(&seminarmodel.Seminar{}, id)
	return res.RowsAffected, res.Error
}

// DeletePermanent performs permanent delete of a seminar record.
func (r *gormRepository) DeletePermanent(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Delete(&seminarmodel.Seminar{}, id)
	return res.RowsAffected, res.Error
}

// Restore restores soft-deleted seminar record.
func (r *gormRepository) Restore(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Model(&seminarmodel.Seminar{}).Where("id = ?", id).Update("deleted_at", nil)
	return res.RowsAffected, res.Error
}
