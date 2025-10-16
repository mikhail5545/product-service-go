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

// DB models
type Product struct {
	ID               string    `gorm:"primaryKey;size:36" json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Name             string    `json:"name"`
	Description      string    `gorm:"size:512" json:"description"`
	Price            float32   `json:"price"`
	ImageUrl         string    `json:"image"` //Cloudinary secure URL
	Amount           int       `json:"amount"`
	ProductType      string    `gorm:"size:50;column:product_type" json:"product_type"`
	ShippingRequired bool      `json:"shipping_required"`
}

// DTO models
type AddProductRequest struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Price            float32 `json:"price" validate:"required,gt=0"`
	Amount           int     `json:"amount"`
	ShippingRequired bool    `json:"shipping_required"`
}

type EditProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Amount      int     `json:"amount"`
}

type GetProductsResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
}
