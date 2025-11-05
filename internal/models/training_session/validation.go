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

// Package trainingsession provides models, DTO models for [trainingsession.Service] requests and validation tools.
package trainingsession

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/mikhail5545/product-service-go/internal/models/common"
)

// Validate validates fields of [trainingsession.CreateRequest].
// All request fields are required.
// Validation rules:
//
//   - Name: required, 3-255 characters, Alpha only.
//   - ShortDescription: required, 3-255 characters.
//   - Price: required, >= 1.
//   - DurationMinutes: required, min 30, must be a multiple of 30.
//   - Format: required, "online" or "offline".
//   - AccessDuration: required, >= 1.
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
			&req.DurationMinutes,
			validation.Required,
			validation.Min(30),
			validation.MultipleOf(30),
		),
		validation.Field(
			&req.Price,
			validation.Required,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.Format,
			validation.Required,
			validation.In("online", "offline"),
		),
	)
}

// Validate validates fields of [trainingsession.UpdateRequest].
// All request fields except ID are optional.
// Validation rules:
//
//   - ID: required, UUID
//   - Name: optional, 3-255 characters, Alpha only.
//   - ShortDescription: optional, 3-255 characters.
//   - LongDescription: optional, 3-3000 characters.
//   - Price: optional, >= 1.
//   - DurationMinutes: optional, min 30, must be a multiple of 30.
//   - Format: optional, "online" or "offline".
//   - AccessDuration: optional, >= 1.
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
			&req.DurationMinutes,
			validation.By(func(value interface{}) error {
				if duration, ok := value.(int); ok {
					if duration <= 0 || duration%30 != 0 {
						return errors.New("should be greater then 0 and be multiple of 30")
					}
				}
				return nil
			}),
		),
		validation.Field(
			&req.Price,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.Format,
			validation.In("online", "offline"),
		),
		validation.Field(
			&req.Tags,
			validation.Length(1, 10),
			validation.Each(validation.Length(3, 20), is.Alphanumeric),
		),
	)
}
