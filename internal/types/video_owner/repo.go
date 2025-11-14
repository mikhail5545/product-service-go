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

// Package video owner defines generic Owner and OwnerRepo interfaces for video owners.
package video_owner

import (
	"context"

	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/types/video_owner_mock/repo_mock.go -package=video_owner_mock . OwnerRepo

// Owner defines the interface for any model that can own a video.
type Owner interface {
	GetID() string
	GetVideoID() *string
	SetVideoID(videoID *string)
}

// OwnerRepo defines the interface for an owner's repository.
type OwnerRepo[T Owner] interface {
	GetWithUnpublished(ctx context.Context, id string) (T, error)
	// UpdateVideoID updates the video ID for a given owner.
	// It should handle both setting a new ID and clearing it (by passing nil).
	UpdateVideoID(ctx context.Context, ownerID string, videoID *string) error
	DB() *gorm.DB
	WithTx(tx *gorm.DB) OwnerRepo[T]
}
