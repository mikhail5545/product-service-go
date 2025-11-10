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

package course

import (
	"context"
	"fmt"

	courserepo "github.com/mikhail5545/product-service-go/internal/database/course"
	coursemodel "github.com/mikhail5545/product-service-go/internal/models/course"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	imageowner "github.com/mikhail5545/product-service-go/internal/types/image_owner"
	"gorm.io/gorm"
)

// courseOwnerRepoAdapter adapts a courserepo.Repository to the generic
// imageowner.OwnerRepo[imageowner.Owner] interface.
//
//	ownerRepoAdapter := newCourseOwnerRepoAdapter(s.CourseRepo.withTx(tx))
//	// transaction based repo, which implements [imageowner.OwnerRepo]
type courseOwnerRepoAdapter struct {
	repo courserepo.Repository
}

// NewOwnerRepoAdapter creates a new adapter.
func NewOwnerRepoAdapter(repo courserepo.Repository) imageowner.OwnerRepo[imageowner.Owner] {
	return &courseOwnerRepoAdapter{repo: repo}
}

func (a *courseOwnerRepoAdapter) GetWithUnpublished(ctx context.Context, id string) (imageowner.Owner, error) {
	course, err := a.repo.GetWithUnpublished(ctx, id)
	if err != nil {
		return nil, err
	}
	var owner imageowner.Owner = *course
	return owner, nil
}

func (a *courseOwnerRepoAdapter) ListWithUnpublishedByIDs(ctx context.Context, ids ...string) ([]imageowner.Owner, error) {
	courses, err := a.repo.ListWithUnpublishedByIDs(ctx, ids...)
	if err != nil {
		return nil, err
	}
	owners := make([]imageowner.Owner, len(courses))
	for i := range courses {
		owners[i] = &courses[i]
	}
	return owners, nil
}

// AddImage adds an image to a single owner by converting it back to a course.
func (a *courseOwnerRepoAdapter) AddImage(ctx context.Context, owner imageowner.Owner, image *imagemodel.Image) error {
	if c, ok := owner.(*coursemodel.Course); ok {
		return a.repo.AddImage(ctx, c, image)
	}
	return fmt.Errorf("incorrect owner type")
}

// DeleteImage deletes an image from a single owner by converting it back to a course.
func (a *courseOwnerRepoAdapter) DeleteImage(ctx context.Context, owner imageowner.Owner, mediaSvcID string) error {
	if c, ok := owner.(*coursemodel.Course); ok {
		return a.repo.DeleteImage(ctx, c, mediaSvcID)
	}
	return fmt.Errorf("incorrect owner type")
}

// AddImageBatch adds an image to a batch of owners by converting them back to courses.
func (a *courseOwnerRepoAdapter) AddImageBatch(ctx context.Context, owners []imageowner.Owner, image *imagemodel.Image) error {
	courses := make([]coursemodel.Course, len(owners))
	for i, owner := range owners {
		if c, ok := owner.(*coursemodel.Course); ok {
			courses[i] = *c
		}
	}
	return a.repo.AddImageBatch(ctx, courses, image)
}

// DeleteImageBatch deletes an image from a batch of owners by converting them back to courses.
func (a *courseOwnerRepoAdapter) DeleteImageBatch(ctx context.Context, owners []imageowner.Owner, image *imagemodel.Image) error {
	courses := make([]coursemodel.Course, len(owners))
	for i, owner := range owners {
		if c, ok := owner.(*coursemodel.Course); ok {
			courses[i] = *c
		}
	}
	return a.repo.DeleteImageBatch(ctx, courses, image)
}

func (a *courseOwnerRepoAdapter) BatchUpdate(ctx context.Context, owners []imageowner.Owner, opt uint) (int64, error) {
	courses := make([]coursemodel.Course, len(owners))
	for i, owner := range owners {
		if c, ok := owner.(*coursemodel.Course); ok {
			courses[i] = *c
		}
	}
	return a.repo.BatchUpdate(ctx, courses, opt)
}

func (a *courseOwnerRepoAdapter) FindOwnerIDsByImageID(ctx context.Context, mediaSvcID string, ownerIDs []string) ([]string, error) {
	return a.repo.FindOwnerIDsByImageID(ctx, mediaSvcID, ownerIDs)
}

func (a *courseOwnerRepoAdapter) DecrementImageCount(ctx context.Context, ownerIDs []string) (int64, error) {
	return a.repo.DecrementImageCount(ctx, ownerIDs)
}

func (a *courseOwnerRepoAdapter) DB() *gorm.DB {
	return a.repo.DB()
}

func (a *courseOwnerRepoAdapter) WithTx(tx *gorm.DB) imageowner.OwnerRepo[imageowner.Owner] {
	return &courseOwnerRepoAdapter{repo: a.repo.WithTx(tx)}
}
