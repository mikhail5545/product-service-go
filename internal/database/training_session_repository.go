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

type TrainingSessionRepository interface {
	// Read operations
	Find(ctx context.Context, id string) (*models.TrainingSession, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.TrainingSession, error)
	Count(ctx context.Context) (int64, error)

	// Write operations
	Create(ctx context.Context, trainingSession *models.TrainingSession) error
	Update(ctx context.Context, trainingSession *models.TrainingSession) error
	Delete(ctx context.Context, id string) error

	DB() *gorm.DB
	WithTx(tx *gorm.DB) TrainingSessionRepository
}

type gormTrainingSessionRepository struct {
	db *gorm.DB
}

// NewTrainingSessionRepository creates a new GORM-based training_session repository.
func NewTrainingSessionRepository(db *gorm.DB) TrainingSessionRepository {
	return &gormTrainingSessionRepository{db: db}
}

// DB returns the underlying gorm.DB instance.
func (r *gormTrainingSessionRepository) DB() *gorm.DB {
	return r.db
}

// WithTx returns a new repository instance with the given transaction.
func (r *gormTrainingSessionRepository) WithTx(tx *gorm.DB) TrainingSessionRepository {
	return &gormTrainingSessionRepository{db: tx}
}

func (r *gormTrainingSessionRepository) Find(ctx context.Context, id string) (*models.TrainingSession, error) {
	var trainingSession models.TrainingSession
	err := r.db.WithContext(ctx).Preload("Product").First(&trainingSession, id).Error
	return &trainingSession, err
}

func (r *gormTrainingSessionRepository) FindAll(ctx context.Context, limit, offset int) ([]models.TrainingSession, error) {
	var trainingSessions []models.TrainingSession
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at desc").Find(&trainingSessions).Error
	return trainingSessions, err
}

func (r *gormTrainingSessionRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.TrainingSession{}).Count(&count).Error
	return count, err
}

func (r *gormTrainingSessionRepository) Create(ctx context.Context, trainingSession *models.TrainingSession) error {
	return r.db.WithContext(ctx).Create(trainingSession).Error
}

func (r *gormTrainingSessionRepository) Update(ctx context.Context, trainingSession *models.TrainingSession) error {
	return r.db.WithContext(ctx).Save(trainingSession).Error
}

func (r *gormTrainingSessionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.TrainingSession{}, id).Error
}
