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

package seminar

import "time"

type CreateRequest struct {
	Name                string    `json:"name"`
	ShortDescription    string    `json:"short_description"`
	ReservationPrice    float32   `json:"reservation_price"`
	EarlyPrice          float32   `json:"early_price"`
	LatePrice           float32   `json:"late_price"`
	EarlySurchargePrice float32   `json:"early_surcharge_price"`
	LateSurchargePrice  float32   `json:"late_surcharge_price"`
	Date                time.Time `json:"date"`
	EndingDate          time.Time `json:"ending_date"`
	Place               string    `json:"place"`
	LatePaymentDate     time.Time `json:"late_payment_date"`
}

type CreateResponse struct {
	ID                      string `json:"id"`
	ReservationProductID    string `json:"reservation_product_id"`
	EarlyProductID          string `json:"early_product_id"`
	LateProductID           string `json:"late_product_id"`
	EarlySurchargeProductID string `json:"early_surcharge_product_id"`
	LateSurchargeProductID  string `json:"late_surcharge_product_id"`
}

type UpdateRequest struct {
	ID                  string     `json:"id"`
	Name                *string    `json:"name,omitempty"`
	ShortDescription    *string    `json:"short_description,omitempty"`
	LongDescription     *string    `json:"long_description,omitempty"`
	ReservationPrice    *float32   `json:"reservation_price,omitempty"`
	EarlyPrice          *float32   `json:"early_price,omitempty"`
	LatePrice           *float32   `json:"late_price,omitempty"`
	EarlySurchargePrice *float32   `json:"early_surcharge_price,omitempty"`
	LateSurchargePrice  *float32   `json:"late_surcharge_price,omitempty"`
	Date                *time.Time `json:"date,omitempty"`
	EndingDate          *time.Time `json:"ending_date,omitempty"`
	Place               *string    `json:"place,omitempty"`
	Tags                []string   `json:"tags,omitempty"`
	LatePaymentDate     *time.Time `json:"late_payment_date,omitempty"`
}

type SeminarDetails struct {
	*Seminar                       `json:"id"`
	ReservationPrice               float32 `json:"reservation_price"`
	EarlyPrice                     float32 `json:"early_price"`
	LatePrice                      float32 `json:"late_price"`
	EarlySurchargePrice            float32 `json:"early_surcharge_price"`
	LateSurchargePrice             float32 `json:"late_surcharge_price"`
	CurrentPrice                   float32 `json:"current_price"`
	CurrentPriceProductID          string  `json:"current_price_product_id"`
	CurrentSurchargePrice          float32 `json:"current_surcharge_price"`
	CurrentSurchargePriceProductID string  `json:"current_surcharge_price_product_id"`
}

// Current populates the following fields in the [seminar.SeminarDetails] struct
// depnding on Seminar.LatePaymentDate value:
//
//   - CurrentPrice: early or late price
//   - CurrentPriceID: Seminar.EarlyProductID or Seminar.LateProductID
//   - CurrentSurchargePrice: early or late surcharge price
//   - CurrentPriceID: Seminar.EarlySurchargeProductID or Seminar.LateSurchargeProductID
func (d *SeminarDetails) Current() {
	if d.Seminar == nil {
		return
	}

	if d.LatePaymentDate.After(time.Now()) {
		d.CurrentPrice = d.EarlyPrice
		if d.EarlyProductID != nil {
			d.CurrentPriceProductID = *d.EarlyProductID
		}
		d.CurrentSurchargePrice = d.EarlySurchargePrice
		if d.EarlySurchargeProductID != nil {
			d.CurrentSurchargePriceProductID = *d.EarlySurchargeProductID
		}
	} else {
		d.CurrentPrice = d.LatePrice
		if d.LateProductID != nil {
			d.CurrentPriceProductID = *d.LateProductID
		}
		d.CurrentSurchargePrice = d.LateSurchargePrice
		if d.LateSurchargeProductID != nil {
			d.CurrentSurchargePriceProductID = *d.LateSurchargeProductID
		}
	}
}
