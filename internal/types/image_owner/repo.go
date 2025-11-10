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

package image_owner

import (
	"context"

	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/types/image_owner_mock/repo_mock.go -package=image_owner_mock . OwnerRepo

// Owner defines the interface for any model that can own images.
type Owner interface {
	GetUploadedImageAmount() int
	SetUploadedImageAmount(amount int)
}

// OwnerRepo defines the interface for an owner's repository.
type OwnerRepo[T Owner] interface {
	ListWithUnpublishedByIDs(ctx context.Context, ids ...string) ([]T, error)
	GetWithUnpublished(ctx context.Context, id string) (T, error)
	AddImage(ctx context.Context, owner T, image *imagemodel.Image) error
	AddImageBatch(ctx context.Context, owners []T, image *imagemodel.Image) error
	DeleteImage(ctx context.Context, owner T, mediaSvcID string) error
	DeleteImageBatch(ctx context.Context, owners []T, image *imagemodel.Image) error
	BatchUpdate(ctx context.Context, owners []T, opt uint) (int64, error)
	FindOwnerIDsByImageID(ctx context.Context, mediaSvcID string, ownerIDs []string) ([]string, error)
	DecrementImageCount(ctx context.Context, ownerIDs []string) (int64, error)
	DB() *gorm.DB
	WithTx(tx *gorm.DB) OwnerRepo[T]
}
