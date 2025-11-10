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

// Package physicalgood provides models, DTO models for [physicalgood.Service] requests and validation tools.
package physicalgood

import (
	"time"

	"github.com/mikhail5545/product-service-go/internal/models/image"
	"gorm.io/gorm"
)

type PhysicalGood struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Tags      []string       `gorm:"type:varchar(128)[]" json:"tags"`
	Name      string         `gorm:"type:varchar(255)" json:"name"`
	// For concise, limited text. Brief description
	ShortDescription string `gorm:"type:varchar(255)" json:"short_description"`
	// For large text\Markdown content. Detailed description
	LongDescription string  `gorm:"type:text" json:"long_description"`
	Price           float32 `json:"price"`
	Amount          int     `json:"amount"`
	// This field flags is the product available in the catalogue or is it archived.
	//
	// 	- InStock = true -> available in the catalogue
	// 	- InStock = false -> not available in the catalogue, archived
	InStock             bool          `json:"in_stock"`
	UploadedImageAmount int           `json:"uploaded_image_amount"`
	Images              []image.Image `gorm:"polymorphic:Owner;" json:"images"`
	ShippingRequired    bool          `json:"shipping_required"`
}

func (g PhysicalGood) GetUploadedImageAmount() int {
	return g.UploadedImageAmount
}

func (g PhysicalGood) SetUploadedImageAmount(amount int) {
	g.UploadedImageAmount = amount
}
