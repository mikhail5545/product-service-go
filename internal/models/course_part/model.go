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

import (
	"time"

	video "github.com/mikhail5545/product-service-go/internal/models/video"
	"gorm.io/gorm"
)

type CoursePart struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Tags      []string       `gorm:"type:varchar(128)[]" json:"tags"`
	// Order of a part in the course
	Number int    `json:"number"`
	Name   string `gorm:"type:varchar(255)" json:"name"`
	// For concise, limited text. Brief description
	ShortDescription string `gorm:"type:varchar(255)" json:"short_description"`
	// For large text\Markdown content. Detailed description
	LongDescription string `gorm:"type:text" json:"long_description"`
	// This field flags is the non sellable item available for the users or is it archived.
	//
	// 	- Published = true -> available for the users
	// 	- Published = false -> not available for the users, archived
	Published bool   `json:"published"`
	CourseID  string `gorm:"size:36;index" json:"course_id"`
	// Unique identifier for the associated Video instance. May be nil. It represents the association with the [media-service-go] MUX Asset.
	//
	// [media-service-go]: https://github.com/mikhail5545/media-service-go
	VideoID *string `gorm:"size:36;index" json:"video_id,omitempty"`
	// This object represents associated Video. May be nil. It represents the [media-service-go] MUX Asset.
	//
	// [media-service-go]: https://github.com/mikhail5545/media-service-go
	Video *video.Video `gorm:"-" json:"video,omitempty"`
}

func (p CoursePart) GetVideoID() *string {
	return p.VideoID
}

func (p CoursePart) SetVideoID(id *string) {
	p.VideoID = id
}
