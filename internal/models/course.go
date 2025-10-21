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

package models

import (
	"time"
)

type Course struct {
	ID             string        `gorm:"primaryKey;size:36" json:"id"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	Name           string        `json:"name"`
	ImageUrl       string        `json:"image_url"`
	Topic          string        `json:"topic"`
	Description    string        `json:"description"`
	ProductID      string        `gorm:"size:36;index" json:"product_id"` // Внешний ключ
	Product        *Product      `gorm:"foreignKey:ProductID" json:"product"`
	AccessDuration int           `json:"access_duration"`
	CourseParts    []*CoursePart `gorm:"foreignKey:CourseID" json:"course_parts"` // Обратная связь
}

type CoursePart struct {
	ID          string     `gorm:"primaryKey;size:36" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Number      int        `json:"number"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CourseID    string     `gorm:"size:36;index" json:"course_id"` // Внешний ключ
	Course      *Course    `gorm:"foreignKey:CourseID" json:"course"`
	MUXVideoID  *string    `gorm:"size:36;index" json:"mux_video_id,omitempty"` // Внешний ключ
	MUXVideo    *MUXUpload `gorm:"-" json:"mux_video,omitempty"`
}

type Attachment struct {
	ID           string      `gorm:"primaryKey;size:36" json:"id"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Name         string      `json:"name"`
	Path         string      `json:"path"`
	Description  string      `json:"description"`
	CoursePartID string      `gorm:"size:36;index" json:"course_part_id"` // Внешний ключ
	CoursePart   *CoursePart `gorm:"foreignKey:CoursePartID" json:"-"`
}

// DTO models
type CourseProductInfo struct {
	Price float32 `json:"price" validate:"required,gt=0"`
}

type AddCourseRequest struct {
	Name           string            `json:"name" validate:"required"`
	Description    string            `json:"description" validate:"required"`
	Topic          string            `json:"topic" validate:"required"`
	Price          float32           `json:"price" validate:"required,gt=0"`
	AccessDuration int               `json:"access_duration"  validate:"required,gt=0"`
	Product        CourseProductInfo `json:"product" validate:"required"`
}

type EditCourseRequest struct {
	Name           *string            `json:"name"`
	Description    *string            `json:"description"`
	Topic          *string            `json:"topic"`
	AccessDuration *int               `json:"access_duration"`
	Product        *CourseProductInfo `json:"product"`
}

type AddCoursePartRequest struct {
	Number      int    `json:"number"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PublicCoursePart struct {
	ID          string    `json:"id"`
	Number      int       `json:"number"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Course      *Course   `json:"course"`
	MUXVideo    *MUXVideo `json:"mux_video,omitempty"`
}

type GetCoursesResponse struct {
	Courses []Course `json:"courses"`
	Total   int64    `json:"total"`
}
