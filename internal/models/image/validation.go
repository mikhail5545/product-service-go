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

// Package image provides models, DTO models for image related requests and validation tools.
package image

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Validate validates fields of [image.AddRequest].
// All request fields are required for image creation.
// Validation rules:
//
//   - URL: required, valid URL.
//   - SecureURL: required, valid URL.
//   - PublicID: required, string.
//   - MediaServiceID: required, valid UUID.
//   - OwnerID: required, valid UUID.
func (req AddRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(
			&req.URL,
			validation.Required,
			is.URL,
		),
		validation.Field(
			&req.SecureURL,
			validation.Required,
			is.URL,
		),
		validation.Field(
			&req.PublicID,
			validation.Required,
		),
		validation.Field(
			&req.MediaServiceID,
			validation.Required,
			is.UUID,
		),
		validation.Field(
			&req.OwnerID,
			validation.Required,
			is.UUID,
		),
	)
}

// Validate validates fields of [image.DeleteRequest].
// All request fields are required for image deletion.
// Validation rules:
//
//   - MediaServiceID: required, valid UUID.
//   - OwnerID: required, valid UUID.
func (req DeleteRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(
			&req.MediaServiceID,
			validation.Required,
			is.UUID,
		),
		validation.Field(
			&req.OwnerID,
			validation.Required,
			is.UUID,
		),
	)
}
