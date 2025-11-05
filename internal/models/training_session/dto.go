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

type CreateRequest struct {
	Name             string  `json:"name"`
	ShortDescription string  `json:"short_description"`
	DurationMinutes  int     `json:"duration_minutes"`
	Format           string  `json:"format"`
	Price            float32 `json:"price"`
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
	DurationMinutes  *int     `json:"duration_minutes,omitempty"`
	Format           *string  `json:"format,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	Price            *float32 `json:"price,omitempty"`
}

type TrainingSessionDetails struct {
	*TrainingSession
	Price     float32 `json:"price"`
	ProductID string  `json:"product_id"`
}
