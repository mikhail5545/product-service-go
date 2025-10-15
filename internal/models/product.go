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
	Image            string    `json:"image"` //Cloudinary secure URL
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
