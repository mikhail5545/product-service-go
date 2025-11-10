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

// Package course provides models, DTO models for [course.Service] requests and validation tools.
package course

// CreateCourseRequest provides essential fields to create new [database.Course] model.
// Other fields should be added later with update request.
type CreateRequest struct {
	Name             string  `json:"name" validate:"required"`
	ShortDescription string  `json:"short_description" validate:"required"`
	Topic            string  `json:"topic" validate:"required"`
	Price            float32 `json:"price" validate:"required,gt=0"`
	AccessDuration   int     `json:"access_duration"  validate:"required,gt=0"`
}

type CreateResponse struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
}

type UpdateRequest struct {
	ID               string   `json:"id" validate:"required"`
	Name             *string  `json:"name"`
	ShortDescription *string  `json:"short_description"`
	LongDescription  *string  `json:"long_description"`
	Topic            *string  `json:"topic"`
	AccessDuration   *int     `json:"access_duration"`
	Tags             []string `json:"tags"`
	Price            *float32 `json:"price"`
}

// CourseDetails is a DTO that combines the Course model with its associated Product price.
type CourseDetails struct {
	*Course
	Price     float32 `json:"price"`
	ProductID string  `json:"product_id"`
}
