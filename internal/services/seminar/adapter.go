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

package seminar

import (
	"context"
	"fmt"

	seminarrepo "github.com/mikhail5545/product-service-go/internal/database/seminar"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	seminarmodel "github.com/mikhail5545/product-service-go/internal/models/seminar"
	imageowner "github.com/mikhail5545/product-service-go/internal/types/image_owner"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/services/seminar_mock/adapter_mock.go -package=seminar_mock . OwnerRepoAdapter
type OwnerRepoAdapter interface {
	GetWithUnpublished(ctx context.Context, id string) (imageowner.Owner, error)
	ListWithUnpublishedByIDs(ctx context.Context, ids ...string) ([]imageowner.Owner, error)
	AddImage(ctx context.Context, owner imageowner.Owner, image *imagemodel.Image) error
	DeleteImage(ctx context.Context, owner imageowner.Owner, mediaSvcID string) error
	AddImageBatch(ctx context.Context, owners []imageowner.Owner, image *imagemodel.Image) error
	BatchUpdate(ctx context.Context, owners []imageowner.Owner, opt uint) (int64, error)
	FindOwnerIDsByImageID(ctx context.Context, mediaSvcID string, ownerIDs []string) ([]string, error)
	DecrementImageCount(ctx context.Context, ownerIDs []string) (int64, error)
	DB() *gorm.DB
	WithTx(tx *gorm.DB) imageowner.OwnerRepo[imageowner.Owner]
}

// ownerRepoAdapter adapts a seminarrepo.Repository to the generic
// imageowner.OwnerRepo[imageowner.Owner] interface.
type ownerRepoAdapter struct {
	repo seminarrepo.Repository
}

// NewOwnerRepoAdapter creates a new adapter.
func NewOwnerRepoAdapter(repo seminarrepo.Repository) imageowner.OwnerRepo[imageowner.Owner] {
	return &ownerRepoAdapter{repo: repo}
}

func (a *ownerRepoAdapter) GetWithUnpublished(ctx context.Context, id string) (imageowner.Owner, error) {
	seminar, err := a.repo.GetWithUnpublished(ctx, id)
	if err != nil {
		return nil, err
	}
	var owner imageowner.Owner = *seminar
	return owner, nil
}

func (a *ownerRepoAdapter) ListWithUnpublishedByIDs(ctx context.Context, ids ...string) ([]imageowner.Owner, error) {
	seminars, err := a.repo.ListWithUnpublishedByIDs(ctx, ids...)
	if err != nil {
		return nil, err
	}
	owners := make([]imageowner.Owner, len(seminars))
	for i := range seminars {
		owners[i] = &seminars[i]
	}
	return owners, nil
}

// AddImage adds an image to a single owner by converting it back to a seminar.
func (a *ownerRepoAdapter) AddImage(ctx context.Context, owner imageowner.Owner, image *imagemodel.Image) error {
	if c, ok := owner.(*seminarmodel.Seminar); ok {
		return a.repo.AddImage(ctx, c, image)
	}
	return fmt.Errorf("incorrect owner type")
}

// DeleteImage deletes an image from a single owner by converting it back to a seminar.
func (a *ownerRepoAdapter) DeleteImage(ctx context.Context, owner imageowner.Owner, mediaSvcID string) error {
	if c, ok := owner.(*seminarmodel.Seminar); ok {
		return a.repo.DeleteImage(ctx, c, mediaSvcID)
	}
	return fmt.Errorf("incorrect owner type")
}

// AddImageBatch adds an image to a batch of owners by converting them back to seminars.
func (a *ownerRepoAdapter) AddImageBatch(ctx context.Context, owners []imageowner.Owner, image *imagemodel.Image) error {
	seminars := make([]seminarmodel.Seminar, len(owners))
	for i, owner := range owners {
		if c, ok := owner.(*seminarmodel.Seminar); ok {
			seminars[i] = *c
		}
	}
	return a.repo.AddImageBatch(ctx, seminars, image)
}

// DeleteImageBatch deletes an image from a batch of owners by converting them back to seminars.
func (a *ownerRepoAdapter) DeleteImageBatch(ctx context.Context, owners []imageowner.Owner, image *imagemodel.Image) error {
	seminars := make([]seminarmodel.Seminar, len(owners))
	for i, owner := range owners {
		if c, ok := owner.(*seminarmodel.Seminar); ok {
			seminars[i] = *c
		}
	}
	return a.repo.DeleteImageBatch(ctx, seminars, image)
}

func (a *ownerRepoAdapter) BatchUpdate(ctx context.Context, owners []imageowner.Owner, opt uint) (int64, error) {
	seminars := make([]seminarmodel.Seminar, len(owners))
	for i, owner := range owners {
		if c, ok := owner.(*seminarmodel.Seminar); ok {
			seminars[i] = *c
		}
	}
	return a.repo.BatchUpdate(ctx, seminars, opt)
}

func (a *ownerRepoAdapter) FindOwnerIDsByImageID(ctx context.Context, mediaSvcID string, ownerIDs []string) ([]string, error) {
	return a.repo.FindOwnerIDsByImageID(ctx, mediaSvcID, ownerIDs)
}

func (a *ownerRepoAdapter) DecrementImageCount(ctx context.Context, ownerIDs []string) (int64, error) {
	return a.repo.DecrementImageCount(ctx, ownerIDs)
}

func (a *ownerRepoAdapter) DB() *gorm.DB {
	return a.repo.DB()
}

func (a *ownerRepoAdapter) WithTx(tx *gorm.DB) imageowner.OwnerRepo[imageowner.Owner] {
	return &ownerRepoAdapter{repo: a.repo.WithTx(tx)}
}
