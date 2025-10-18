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

package models

import (
	"time"
)

// DB model
type TrainingSession struct {
	ID              string    `gorm:"primaryKey;size:36" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ProductID       string    `gorm:"size:36;index" json:"product_id"` // Внешний ключ, не null
	Product         *Product  `gorm:"foreignKey:ProductID" json:"product"`
	DurationMinutes int       `json:"duration_minutes"`
	Format          string    `gorm:"size:50" json:"format,omitempty"`
}

// DTO models
type AddTrainingSessionRequest struct {
	DurationMinutes int               `json:"duration_minutes"  validate:"required,gt=0"`
	Format          string            `json:"format" validate:"required"`
	Product         AddProductRequest `json:"product"`
}

type EditTrainingSessionRequest struct {
	DurationMinutes int     `json:"duration_minutes" validate:"required,gt=0"`
	Format          string  `json:"format" validate:"required"`
	Name            string  `json:"name" validate:"required"`
	Description     string  `json:"description" validate:"required"`
	Price           float32 `json:"price" validate:"required,gt=0"`
}

type EditTrainingSessionImageRequest struct {
	Image string `json:"image" validate:"required"`
}

type GetTrainingSessionsResponse struct {
	TrainingSessions []TrainingSession `json:"training_sessions"`
	Total            int64             `json:"total"`
}
