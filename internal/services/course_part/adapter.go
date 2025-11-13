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

package coursepart

import (
	"context"

	coursepartrepo "github.com/mikhail5545/product-service-go/internal/database/course_part"
	videoowner "github.com/mikhail5545/product-service-go/internal/types/video_owner"
	"gorm.io/gorm"
)

type OwnerRepoAdapter interface {
	GetWithUnpublished(ctx context.Context, id string) (videoowner.Owner, error)
	// UpdateVideoID updates the video ID for a given owner.
	// It should handle both setting a new ID and clearing it (by passing nil).
	UpdateVideoID(ctx context.Context, ownerID string, videoID *string) error
	DB() *gorm.DB
	WithTx(tx *gorm.DB) videoowner.OwnerRepo[videoowner.Owner]
}

// ownerRepoAdapter adapts a coursepartrepo.Repository to the generic
// videoowner.OwnerRepo[videoowner.Owner] interface.
type ownerRepoAdapter struct {
	repo coursepartrepo.Repository
}

func NewOwnerRepoAdapter(repo coursepartrepo.Repository) OwnerRepoAdapter {
	return &ownerRepoAdapter{repo: repo}
}

func (a *ownerRepoAdapter) GetWithUnpublished(ctx context.Context, id string) (videoowner.Owner, error) {
	part, err := a.repo.GetWithUnpublished(ctx, id)
	if err != nil {
		return nil, err
	}
	var owner videoowner.Owner = *part
	return owner, nil
}

func (a *ownerRepoAdapter) UpdateVideoID(ctx context.Context, ownerID string, videoID *string) error {
	return a.repo.UpdateVideoID(ctx, ownerID, videoID)
}

func (a *ownerRepoAdapter) DB() *gorm.DB {
	return a.DB()
}

func (a *ownerRepoAdapter) WithTx(tx *gorm.DB) videoowner.OwnerRepo[videoowner.Owner] {
	return &ownerRepoAdapter{repo: a.repo.WithTx(tx)}
}
