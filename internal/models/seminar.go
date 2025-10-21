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

type Seminar struct {
	ID                      string    `gorm:"primaryKey;size:36" json:"id"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	ReservationProductID    string    `gorm:"size:36;index" json:"reservation_product_id"` // Внешний ключ
	ReservationProduct      *Product  `gorm:"foreignKey:ReservationProductID" json:"reservation_product"`
	EarlyProductID          string    `gorm:"size:36;index" json:"early_product_id"` // Внешний ключ
	EarlyProduct            *Product  `gorm:"foreignKey:EarlyProductID" json:"early_product"`
	LateProductID           string    `gorm:"size:36;index" json:"late_product_id"` // Внешний ключ
	LateProduct             *Product  `gorm:"foreignKey:LateProductID" json:"late_product"`
	EarlySurchargeProductID string    `gorm:"size:36;index" json:"early_surcharge_product_id"` // Внешний ключ
	EarlySurchargeProduct   *Product  `gorm:"foreignKey:EarlySurchargeProductID" json:"early_surcharge_product"`
	LateSurchargeProductID  string    `gorm:"size:36;index" json:"late_surcharge_product_id"` // Внешний ключ
	LateSurchargeProduct    *Product  `gorm:"foreignKey:LateSurchargeProductID" json:"late_surcharge_product"`
	Date                    time.Time `gorm:"type:timestamptz" json:"date"`
	EndingDate              time.Time `gorm:"type:timestamptz" json:"ending_date"`
	Place                   string    `json:"place"`
	LatePaymentDate         time.Time `gorm:"type:timestamptz" json:"late_payment_date"`
	Details                 string    `json:"details"`
	Price                   float32   `gorm:"-" json:"price"`                    // Игнорировать в DB
	CurrentPriceProductID   string    `gorm:"-" json:"current_price_product_id"` // Игнорировать в DB
}

// Returns early price and early product Id or late price and late product Id based on LatePymentDate
func (seminar *Seminar) GetPrice() (float32, string) {
	now := time.Now()
	if now.Before(seminar.LatePaymentDate) {
		return seminar.EarlyProduct.Price, seminar.EarlyProduct.ID
	} else {
		return seminar.LateProduct.Price, seminar.LateProduct.ID
	}
}

// DTO Models
type SeminarProductInfo struct {
	Price       float32 `json:"price" validate:"required,gt=0"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type AddSeminarRequest struct {
	Name                  string              `json:"name" validate:"required"`
	Description           string              `json:"description" validate:"required"`
	Date                  time.Time           `json:"date" validate:"required"`
	EndingDate            time.Time           `json:"ending_date" validate:"required"`
	Place                 string              `json:"place" validate:"required"`
	LatePaymentDate       time.Time           `json:"late_payment_date" validate:"required"`
	Details               string              `json:"details"`
	ReservationProduct    SeminarProductInfo  `json:"reservation_product" validate:"required"`
	EarlyProduct          SeminarProductInfo  `json:"early_product" validate:"required"`
	LateProduct           SeminarProductInfo  `json:"late_product" validate:"required"`
	EarlySurchargeProduct *SeminarProductInfo `json:"early_surcharge_product"`
	LateSurchargeProduct  *SeminarProductInfo `json:"late_surcharge_product"`
}

type UpdateSeminarRequest struct {
	Name                  *string             `json:"name"`
	Description           *string             `json:"description"`
	Date                  *time.Time          `json:"date"`
	EndingDate            *time.Time          `json:"ending_date"`
	Place                 *string             `json:"place"`
	LatePaymentDate       *time.Time          `json:"late_payment_date"`
	Details               *string             `json:"details"`
	ReservationProduct    *SeminarProductInfo `json:"reservation_product"`
	EarlyProduct          *SeminarProductInfo `json:"early_product"`
	LateProduct           *SeminarProductInfo `json:"late_product"`
	EarlySurchargeProduct *SeminarProductInfo `json:"early_surcharge_product"`
	LateSurchargeProduct  *SeminarProductInfo `json:"late_surcharge_product"`
}
