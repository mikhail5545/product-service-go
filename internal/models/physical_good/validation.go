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
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/mikhail5545/product-service-go/internal/models/common"
)

// Validate validates fields of [physicalgood.CreateRequest].
// All request fields are required for creation.
// Validation rules:
//
//   - Name: required, 3-255 characters, Alpha only.
//   - ShortDescription: required, 3-255 characters.
//   - Price: required, >= 1.
//   - ShippingRequired: required, boolean.
//   - Amount: required, >= 0, >= 1 if ShippingRequired is true.
func (req CreateRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(
			&req.Name,
			validation.Required,
			validation.Length(3, 255),
			validation.By(common.ValidateName),
		),
		validation.Field(
			&req.ShortDescription,
			validation.Required,
			validation.Length(3, 255),
		),
		validation.Field(
			&req.Price,
			validation.Required,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.Amount,
			validation.Required,
			validation.Min(0),
			validation.When(
				req.ShippingRequired,
				validation.Min(1),
			),
		),
	)
}

// Validate validates fields of [physicalgood.UpdateRequest].
// All request fields except ID are optional.
// Validation rules:
//
//   - ID: required, UUID.
//   - Name: optional, 3-255 characters, Alpha only.
//   - ShortDescription: optional, 3-255 characters.
//   - ShortDescription: optional, 3-3000 characters.
//   - Price: optional, >= 1.
//   - ShippingRequired: optional, boolean.
//   - Amount: optional, >= 0, >= 1 if ShippingRequired is true.
//   - Tags: optional, 1-10 items, 3-20 characters each.
func (req UpdateRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(
			&req.ID,
			validation.Required,
			is.UUID,
		),
		validation.Field(
			&req.Name,
			validation.Length(3, 255),
			validation.By(common.ValidateName),
		),
		validation.Field(
			&req.ShortDescription,
			validation.Length(3, 255),
		),
		validation.Field(
			&req.LongDescription,
			validation.Length(3, 3000),
		),
		validation.Field(
			&req.Price,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.Amount,
			validation.Min(0),
			validation.By(func(value interface{}) error {
				if req.ShippingRequired != nil && *req.ShippingRequired {
					if amount, ok := value.(*int); ok {
						if *amount < 1 {
							return errors.New("must be greater then 0 if shipping is required")
						}
					}
				}
				return nil
			}),
		),
		validation.Field(
			&req.Tags,
			validation.Length(1, 10),
			validation.Each(validation.Length(3, 20), is.Alphanumeric),
		),
	)
}
