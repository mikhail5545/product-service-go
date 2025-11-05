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

type PhysicalGoodDetails struct {
	PhysicalGood
	Price     float32
	ProductID string
}

type CreateRequest struct {
	Name             string  `json:"name"`
	ShortDescription string  `json:"short_description"`
	Price            float32 `json:"price"`
	Amount           int     `json:"amount"`
	ShippingRequired bool    `json:"shipping_required"`
}

type CreateResponse struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
}

type UpdateRequest struct {
	ID               string   `json:"id"`
	Name             *string  `json:"name,omitempty"`
	ShortDescription *string  `json:"short_description,omitempty"`
	LongDescription  *string  `json:"long_description,omitempty"`
	Price            *float32 `json:"price,omitempty"`
	Amount           *int     `json:"amount,omitempty"`
	ShippingRequired *bool    `json:"shipping_required,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}
