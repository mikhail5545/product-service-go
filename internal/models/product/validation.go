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
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Validate validates fields of [product.CreateRequest].
// All request fields are required for product creation.
// Validation rules:
//
//   - DetailsID: required, UUID
//   - Price: required, min 1.
//   - Name: required, "physical_good", "training_session", seminar or "course".
func (req *AddRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(
			&req.Price,
			validation.Required,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.DetailsID,
			validation.Required,
			is.UUID,
		),
		validation.Field(
			&req.DetailsType,
			validation.Required,
			validation.In("physical_good", "training_session", "seminar", "course"),
		),
	)
}
