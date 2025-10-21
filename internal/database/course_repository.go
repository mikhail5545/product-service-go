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

type CourseRepository interface {
	// Read operations
	Find(ctx context.Context, id string) (*models.Course, error)
	FindWithParts(ctx context.Context, id string) (*models.Course, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.Course, error)
	Count(ctx context.Context) (int64, error)
	// TODO: Create all CRUD functions for course_part

	// Write operations
	Create(ctx context.Context, course *models.Course) error
	Update(ctx context.Context, course *models.Course) error
	Delete(ctx context.Context, id string) error

	DB() *gorm.DB
	WithTx(tx *gorm.DB) CourseRepository
}

type gormCourseRepository struct {
	db *gorm.DB
}

// NewCourseRepository creates a new GORM-based course repository.
func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &gormCourseRepository{db: db}
}

// DB returns the underlying gorm.DB instance.
func (r *gormCourseRepository) DB() *gorm.DB {
	return r.db
}

// WithTx returns a new repository instance with the given transaction.
func (r *gormCourseRepository) WithTx(tx *gorm.DB) CourseRepository {
	return &gormCourseRepository{db: tx}
}

func (r *gormCourseRepository) Find(ctx context.Context, id string) (*models.Course, error) {
	var course models.Course
	err := r.db.WithContext(ctx).Preload("Product").First(&course, "id = ?", id).Error
	return &course, err
}

func (r *gormCourseRepository) FindWithParts(ctx context.Context, id string) (*models.Course, error) {
	var course models.Course
	err := r.db.WithContext(ctx).Preload("Product").Preload("CourseParts").First(&course, "id = ?", id).Error
	return &course, err
}

func (r *gormCourseRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.WithContext(ctx).Preload("Product").Limit(limit).Offset(offset).Order("created_at desc").Find(&courses).Error
	return courses, err
}

func (r *gormCourseRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Course{}).Count(&count).Error
	return count, err
}

func (r *gormCourseRepository) Create(ctx context.Context, course *models.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *gormCourseRepository) Update(ctx context.Context, course *models.Course) error {
	return r.db.WithContext(ctx).Save(course).Error
}

func (r *gormCourseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Course{}, id).Error
}
