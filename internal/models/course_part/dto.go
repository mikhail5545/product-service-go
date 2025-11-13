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

// Package coursepart provides models, DTO models for [coursepart.Service] requests and validation tools.
package coursepart

// CreateRequest holds neccessary fields to create new Course Part.
// All fields are required.
type CreateRequest struct {
	CourseID         string `json:"course_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Number           int    `json:"number"`
}

type CreateResponse struct {
	ID       string `json:"id"`
	CourseID string `json:"course_id"`
}

// CreateRequest holds neccessary fields to update new Course Part.
// All fields are optional except ID and CourseID.
type UpdateRequest struct {
	ID               string   `json:"id"`
	CourseID         string   `json:"course_id"`
	Name             *string  `json:"name"`
	LongDescription  *string  `json:"long_description"`
	ShortDescription *string  `json:"short_description"`
	Number           *int     `json:"number"`
	Tags             []string `json:"tags"`
}
