// vitainmove.com/product-service-go
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

	"gorm.io/gorm"
	"vitainmove.com/product-service-go/internal/models"
)

type SeminarRepository interface {
	// Read operations
	Find(ctx context.Context, id string) (*models.Seminar, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.Seminar, error)
	Count(ctx context.Context) (int64, error)

	// Write operations
	Create(ctx context.Context, seminar *models.Seminar) error
	Update(ctx context.Context, seminar *models.Seminar) error
	Delete(ctx context.Context, id string) error

	DB() *gorm.DB
	WithTx(tx *gorm.DB) SeminarRepository
}

type gormSeminarRepository struct {
	db *gorm.DB
}

// NewSeminarRepository creates a new GORM-based seminar repository.
func NewSeminarRepository(db *gorm.DB) SeminarRepository {
	return &gormSeminarRepository{db: db}
}

// DB returns the underlying gorm.DB instance.
func (r *gormSeminarRepository) DB() *gorm.DB {
	return r.db
}

// WithTx returns a new repository instance with the given transaction.
func (r *gormSeminarRepository) WithTx(tx *gorm.DB) SeminarRepository {
	return &gormSeminarRepository{db: tx}
}

func (r *gormSeminarRepository) Find(ctx context.Context, id string) (*models.Seminar, error) {
	var seminar *models.Seminar
	err := r.db.WithContext(ctx).Preload("ReservationProduct").Preload("EarlyProduct").Preload("LateProduct").Preload("EarlySurchargeProduct").Preload("LateSurchargeProduct").First(&seminar, "id = ?", id).Error
	return seminar, err
}

func (r *gormSeminarRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Seminar, error) {
	var seminars []models.Seminar
	err := r.db.WithContext(ctx).Preload("ReservationProduct").Preload("EarlyProduct").Preload("LateProduct").Preload("EarlySurchargeProduct").Preload("LateSurchargeProduct").Order("created_at desc").Limit(limit).Offset(offset).Find(&seminars).Error
	return seminars, err
}

func (r *gormSeminarRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Seminar{}).Count(&count).Error
	return count, err
}

func (r *gormSeminarRepository) Create(ctx context.Context, seminar *models.Seminar) error {
	return r.db.WithContext(ctx).Create(seminar).Error
}

func (r *gormSeminarRepository) Update(ctx context.Context, seminar *models.Seminar) error {
	return r.db.WithContext(ctx).Save(seminar).Error
}

func (r *gormSeminarRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Seminar{}, id).Error
}
