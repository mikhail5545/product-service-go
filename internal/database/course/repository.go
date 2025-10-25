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

package course

import (
	"context"

	"github.com/mikhail5545/product-service-go/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	// Read operations
	// Find retrieves Course data from the database
	Get(ctx context.Context, id string) (*models.Course, error)
	// FindWithParts retrieves Course data withound any Course Parts.
	GetReduced(ctx context.Context, id string) (*models.Course, error)
	// List retrieves all Course records from database. (Without course parts)
	List(ctx context.Context, limit, offset int) ([]models.Course, error)
	// Count counts the number of Course records in the database.
	Count(ctx context.Context) (int64, error)

	// Write operations
	Create(ctx context.Context, course *models.Course) error
	Update(ctx context.Context, course *models.Course, updates any) (int64, error)
	Delete(ctx context.Context, id string) error

	DB() *gorm.DB
	WithTx(tx *gorm.DB) Repository
}

type gormRepository struct {
	db *gorm.DB
}

// New creates a new GORM-based course repository.
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

func (r *gormRepository) Get(ctx context.Context, id string) (*models.Course, error) {
	var course models.Course
	err := r.db.WithContext(ctx).Preload("Product").Preload("CourseParts").First(&course, "id = ?", id).Error
	return &course, err
}

func (r *gormRepository) GetReduced(ctx context.Context, id string) (*models.Course, error) {
	var course models.Course
	err := r.db.WithContext(ctx).Preload("Product").First(&course, "id = ?", id).Error
	return &course, err
}

func (r *gormRepository) List(ctx context.Context, limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.WithContext(ctx).Preload("Product").Limit(limit).Offset(offset).Order("created_at desc").Find(&courses).Error
	return courses, err
}

func (r *gormRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Course{}).Count(&count).Error
	return count, err
}

func (r *gormRepository) Create(ctx context.Context, course *models.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *gormRepository) Update(ctx context.Context, course *models.Course, updates any) (int64, error) {
	res := r.db.WithContext(ctx).Model(course).Updates(updates)
	return res.RowsAffected, res.Error
}

func (r *gormRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Course{}, id).Error
}
