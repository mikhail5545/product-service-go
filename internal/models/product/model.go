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

// Package product provides models, DTO models for [product.Service] requests and validation tools.
package product

import (
	"time"

	"gorm.io/gorm"
)

// Product holds essential data for order and cart operations.
// It acts as polymorphic model, holding ID of structure, representing detailed information.
// It can hold ID for:
type Product struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Price     float32        `json:"price"`
	// This field flags is the product available in the catalogue or is it archived.
	//
	// 	- InStock = true -> available in the catalogue
	// 	- InStock = false -> not available in the catalogue, archived
	InStock bool `json:"in_stock"`
	// ID to the details struct. It can be [models.course.Course], [models.seminar.Seminar], [models.trainingsession.TrainingSession]
	// [models.physicalgood.PhysicalGood].
	DetailsID string `gorm:"size:36;index" json:"details_id"`
	// Type of the details struct. It can be 'course', 'seminar', 'training_session', 'physical_good'.
	DetailsType string `gorm:"size:50;index" json:"details_type"`
}

type GetProductsResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
}
