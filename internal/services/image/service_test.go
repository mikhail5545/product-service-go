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

// Package image provides a reusable service for managing images for different owner types.
package image

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	imagerepomock "github.com/mikhail5545/product-service-go/internal/test/database/image_mock"
	imageownermock "github.com/mikhail5545/product-service-go/internal/test/types/image_owner_mock"
	"github.com/mikhail5545/product-service-go/internal/types/image_owner"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// mockOwner implements the Owner interface for testing purposes.
type mockOwner struct {
	id                  string
	uploadedImageAmount int
}

func (m *mockOwner) GetUploadedImageAmount() int {
	return m.uploadedImageAmount
}

func (m *mockOwner) SetUploadedImageAmount(amount int) {
	m.uploadedImageAmount = amount
}

func TestService_AddImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageRepo := imagerepomock.NewMockRepository(ctrl)
	mockOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)

	testService := New(mockImageRepo)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	ownerID := uuid.New().String()
	addReq := &imagemodel.AddRequest{
		URL:            "http://example.com/image.jpg",
		SecureURL:      "https://example.com/image.jpg",
		PublicID:       "public-id",
		MediaServiceID: uuid.NewString(),
		OwnerID:        ownerID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().DB().Return(db)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		owner := &mockOwner{id: ownerID, uploadedImageAmount: 2}

		mockTxOwnerRepo.EXPECT().GetWithUnpublished(gomock.Any(), ownerID).Return(owner, nil)
		mockTxOwnerRepo.EXPECT().AddImage(gomock.Any(), owner, gomock.Any()).Return(nil)
		mockTxOwnerRepo.EXPECT().BatchUpdate(gomock.Any(), gomock.Any(), uint(2)).
			DoAndReturn(func(_ context.Context, owners []image_owner.Owner, _ uint) (int64, error) {
				assert.Equal(t, 3, owners[0].GetUploadedImageAmount())
				return 1, nil
			})

		// Act
		err := testService.AddImage(context.Background(), addReq, mockOwnerRepo)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		invalidReq := &imagemodel.AddRequest{
			OwnerID: "not-a-uuid",
		}

		// Act
		err := testService.AddImage(context.Background(), invalidReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("owner not found", func(t *testing.T) {
		// Arrange
		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().DB().Return(db)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		mockTxOwnerRepo.EXPECT().GetWithUnpublished(gomock.Any(), ownerID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		err := testService.AddImage(context.Background(), addReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrOwnerNotFound)
	})

	t.Run("image limit exceeded", func(t *testing.T) {
		// Arrange
		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().DB().Return(db)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		owner := &mockOwner{id: ownerID, uploadedImageAmount: 5}

		mockTxOwnerRepo.EXPECT().GetWithUnpublished(gomock.Any(), ownerID).Return(owner, nil)

		// Act
		err := testService.AddImage(context.Background(), addReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrImageLimitExceeded)
	})

	t.Run("add image db error", func(t *testing.T) {
		// Arrange
		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().DB().Return(db)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		owner := &mockOwner{id: ownerID, uploadedImageAmount: 2}

		mockTxOwnerRepo.EXPECT().GetWithUnpublished(gomock.Any(), ownerID).Return(owner, nil)
		dbErr := errors.New("database error")
		mockTxOwnerRepo.EXPECT().AddImage(gomock.Any(), owner, gomock.Any()).Return(dbErr)

		// Act
		err := testService.AddImage(context.Background(), addReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to add image for owner")
	})
}

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageRepo := imagerepomock.NewMockRepository(ctrl)
	mockOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)

	testService := New(mockImageRepo)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	ownerID := uuid.New().String()
	mediaSvcID := uuid.New().String()
	deleteReq := &imagemodel.DeleteRequest{
		OwnerID:        ownerID,
		MediaServiceID: mediaSvcID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().DB().Return(db)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		owner := &mockOwner{id: ownerID, uploadedImageAmount: 2}

		mockTxOwnerRepo.EXPECT().GetWithUnpublished(gomock.Any(), ownerID).Return(owner, nil)
		mockTxOwnerRepo.EXPECT().DeleteImage(gomock.Any(), owner, mediaSvcID).Return(nil)
		mockTxOwnerRepo.EXPECT().DecrementImageCount(gomock.Any(), []string{ownerID}).Return(int64(1), nil)

		// Act
		err := testService.DeleteImage(context.Background(), deleteReq, mockOwnerRepo)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		invalidReq := &imagemodel.DeleteRequest{OwnerID: "not-a-uuid"}

		// Act
		err := testService.DeleteImage(context.Background(), invalidReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("owner not found", func(t *testing.T) {
		// Arrange
		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().DB().Return(db)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		mockTxOwnerRepo.EXPECT().GetWithUnpublished(gomock.Any(), ownerID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		err := testService.DeleteImage(context.Background(), deleteReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrOwnerNotFound)
	})

	t.Run("image not found on owner", func(t *testing.T) {
		// Arrange
		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().DB().Return(db)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		owner := &mockOwner{id: ownerID, uploadedImageAmount: 2}

		mockTxOwnerRepo.EXPECT().GetWithUnpublished(gomock.Any(), ownerID).Return(owner, nil)
		mockTxOwnerRepo.EXPECT().DeleteImage(gomock.Any(), owner, mediaSvcID).Return(gorm.ErrRecordNotFound)

		// Act
		err := testService.DeleteImage(context.Background(), deleteReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrImageNotFoundOnOwner)
	})

	t.Run("decrement uploaded image amount db error", func(t *testing.T) {
		// Arrange
		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().DB().Return(db)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		owner := &mockOwner{id: ownerID, uploadedImageAmount: 2}

		mockTxOwnerRepo.EXPECT().GetWithUnpublished(gomock.Any(), ownerID).Return(owner, nil)
		mockTxOwnerRepo.EXPECT().DeleteImage(gomock.Any(), owner, mediaSvcID).Return(nil)
		dbErr := errors.New("database error")
		mockTxOwnerRepo.EXPECT().DecrementImageCount(gomock.Any(), []string{ownerID}).Return(int64(0), dbErr)

		// Act
		err := testService.DeleteImage(context.Background(), deleteReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decrement owner uploaded image count")
	})
}

func TestService_AddImageBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageRepo := imagerepomock.NewMockRepository(ctrl)
	mockOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)

	testService := New(mockImageRepo)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	ownerID_1 := uuid.New().String()
	ownerID_2 := uuid.New().String()
	ownerID_3 := uuid.New().String()
	mediaSvcID := uuid.New().String()

	addReq := &imagemodel.AddBatchRequest{
		URL:            "http://example.com/image.jpg",
		SecureURL:      "https://example.com/image.jpg",
		PublicID:       "public-id",
		MediaServiceID: mediaSvcID,
		OwnerIDs:       []string{ownerID_1, ownerID_2, ownerID_3},
	}
	owner_1 := mockOwner{
		id:                  ownerID_1,
		uploadedImageAmount: 2,
	}
	owner_2 := mockOwner{
		id:                  ownerID_2,
		uploadedImageAmount: 3,
	}
	owner_3 := mockOwner{
		id:                  ownerID_3,
		uploadedImageAmount: 5, // not a valid owner (image limit exceeded)
	}
	t.Run("success", func(t *testing.T) {
		// Arrange
		mockOwners := []mockOwner{owner_1, owner_2, owner_3}
		// Convert []mockOwner to []image_owner.Owner
		owners := make([]image_owner.Owner, len(mockOwners))
		for i := range mockOwners {
			owners[i] = &mockOwners[i]
		}

		// Convert slice of strings to slice of any for variadic mock expectation
		ownerIDsAny := make([]any, len(addReq.OwnerIDs))
		for i, v := range addReq.OwnerIDs {
			ownerIDsAny[i] = v
		}
		mockOwnerRepo.EXPECT().ListWithUnpublishedByIDs(gomock.Any(), ownerIDsAny...).Return(owners, nil)

		mockImageRepo.EXPECT().DB().Return(db)

		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		mockTxOwnerRepo.EXPECT().AddImageBatch(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ []image_owner.Owner, img *imagemodel.Image) error {
				assert.Equal(t, addReq.URL, img.URL)
				assert.Equal(t, addReq.SecureURL, img.SecureURL)
				assert.Equal(t, addReq.MediaServiceID, img.MediaServiceID)
				return nil
			})

		mockTxOwnerRepo.EXPECT().BatchUpdate(gomock.Any(), gomock.Any(), uint(2)).
			DoAndReturn(func(_ context.Context, owners []image_owner.Owner, _ uint) (int64, error) {
				assert.Equal(t, 3, owners[0].GetUploadedImageAmount())
				assert.Equal(t, 4, owners[1].GetUploadedImageAmount())
				return int64(2), nil
			})

		// Act
		affectedOwners, err := testService.AddImageBatch(context.Background(), addReq, mockOwnerRepo)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 2, affectedOwners)
	})

	t.Run("success with no valid owners", func(t *testing.T) {
		// Arrange
		mockOwners := []mockOwner{
			{id: ownerID_1, uploadedImageAmount: 5},
			{id: ownerID_2, uploadedImageAmount: 5},
		}
		owners := make([]image_owner.Owner, len(mockOwners))
		for i := range mockOwners {
			owners[i] = &mockOwners[i]
		}

		// Convert slice of strings to slice of any for variadic mock expectation
		ownerIDsAny := make([]any, len(addReq.OwnerIDs))
		for i, v := range addReq.OwnerIDs {
			ownerIDsAny[i] = v
		}
		mockOwnerRepo.EXPECT().ListWithUnpublishedByIDs(gomock.Any(), ownerIDsAny...).Return(owners, nil)

		// Act
		affectedOwners, err := testService.AddImageBatch(context.Background(), addReq, mockOwnerRepo)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0, affectedOwners)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		invalidReq := &imagemodel.AddBatchRequest{
			MediaServiceID: "not-a-uuid",
		}

		// Act
		_, err := testService.AddImageBatch(context.Background(), invalidReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("owners not found", func(t *testing.T) {
		// Arrange
		// Convert slice of strings to slice of any for variadic mock expectation
		ownerIDsAny := make([]any, len(addReq.OwnerIDs))
		for i, v := range addReq.OwnerIDs {
			ownerIDsAny[i] = v
		}
		mockOwnerRepo.EXPECT().ListWithUnpublishedByIDs(gomock.Any(), ownerIDsAny...).Return([]image_owner.Owner{}, nil)

		// Act
		_, err := testService.AddImageBatch(context.Background(), addReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrOwnersNotFound)
	})

	t.Run("add image batch db error", func(t *testing.T) {
		// Arrange
		mockOwners := []mockOwner{
			{id: ownerID_1, uploadedImageAmount: 2},
			{id: ownerID_2, uploadedImageAmount: 3},
			{id: ownerID_3, uploadedImageAmount: 5}, // not a valid owner (image limit exceeded)
		}
		owners := make([]image_owner.Owner, len(mockOwners))
		for i := range mockOwners {
			owners[i] = &mockOwners[i]
		}

		// Convert slice of strings to slice of any for variadic mock expectation
		ownerIDsAny := make([]any, len(addReq.OwnerIDs))
		for i, v := range addReq.OwnerIDs {
			ownerIDsAny[i] = v
		}
		mockOwnerRepo.EXPECT().ListWithUnpublishedByIDs(gomock.Any(), ownerIDsAny...).Return(owners, nil)

		mockImageRepo.EXPECT().DB().Return(db)

		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		dbErr := errors.New("database error")
		mockTxOwnerRepo.EXPECT().AddImageBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbErr)

		// Act
		_, err := testService.AddImageBatch(context.Background(), addReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to batch add images for owners")
	})

	t.Run("batch owner update db error", func(t *testing.T) {
		// Arrange
		mockOwners := []mockOwner{
			{id: ownerID_1, uploadedImageAmount: 2},
			{id: ownerID_2, uploadedImageAmount: 3},
			{id: ownerID_3, uploadedImageAmount: 5}, // not a valid owner (image limit exceeded)
		}
		owners := make([]image_owner.Owner, len(mockOwners))
		for i := range mockOwners {
			owners[i] = &mockOwners[i]
		}

		// Convert slice of strings to slice of any for variadic mock expectation
		ownerIDsAny := make([]any, len(addReq.OwnerIDs))
		for i, v := range addReq.OwnerIDs {
			ownerIDsAny[i] = v
		}
		mockOwnerRepo.EXPECT().ListWithUnpublishedByIDs(gomock.Any(), ownerIDsAny...).Return(owners, nil)

		mockImageRepo.EXPECT().DB().Return(db)

		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		dbErr := errors.New("database error")
		mockTxOwnerRepo.EXPECT().AddImageBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockTxOwnerRepo.EXPECT().BatchUpdate(gomock.Any(), gomock.Any(), uint(2)).Return(int64(0), dbErr)

		// Act
		_, err := testService.AddImageBatch(context.Background(), addReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to batch update owners")
	})
}

func TestService_DeleteImageBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageRepo := imagerepomock.NewMockRepository(ctrl)
	mockOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)

	testService := New(mockImageRepo)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	ownerID_1 := uuid.New().String()
	ownerID_2 := uuid.New().String()
	mediaSvcID := uuid.New().String()

	deleteReq := &imagemodel.DeleteBatchRequst{
		MediaServiceID: mediaSvcID,
		OwnerIDs:       []string{ownerID_1, ownerID_2},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockImageRepo.EXPECT().DB().Return(db)

		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		mockOwners := []mockOwner{
			{id: ownerID_1, uploadedImageAmount: 2},
			{id: ownerID_2, uploadedImageAmount: 3},
		}
		owners := make([]image_owner.Owner, len(mockOwners))
		for i := range mockOwners {
			owners[i] = &mockOwners[i]
		}

		// Convert slice of strings to slice of any for variadic mock expectation
		ownerIDsAny := make([]any, len(deleteReq.OwnerIDs))
		for i, v := range deleteReq.OwnerIDs {
			ownerIDsAny[i] = v
		}

		mockTxOwnerRepo.EXPECT().ListWithUnpublishedByIDs(gomock.Any(), ownerIDsAny...).Return(owners, nil)
		mockTxOwnerRepo.EXPECT().FindOwnerIDsByImageID(gomock.Any(), mediaSvcID, deleteReq.OwnerIDs).Return(deleteReq.OwnerIDs, nil)
		mockTxOwnerRepo.EXPECT().DeleteImageBatch(gomock.Any(), owners, &imagemodel.Image{MediaServiceID: deleteReq.MediaServiceID}).Return(nil)
		mockTxOwnerRepo.EXPECT().DecrementImageCount(gomock.Any(), deleteReq.OwnerIDs).Return(int64(2), nil)

		// Act
		affectedOwners, err := testService.DeleteImageBatch(context.Background(), deleteReq, mockOwnerRepo)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 2, affectedOwners)
	})

	t.Run("success with no affected owners", func(t *testing.T) {
		// Arrange
		mockImageRepo.EXPECT().DB().Return(db)

		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		mockOwners := []mockOwner{
			{id: ownerID_1, uploadedImageAmount: 2},
			{id: ownerID_2, uploadedImageAmount: 3},
		}
		owners := make([]image_owner.Owner, len(mockOwners))
		for i := range mockOwners {
			owners[i] = &mockOwners[i]
		}

		// Convert slice of strings to slice of any for variadic mock expectation
		ownerIDsAny := make([]any, len(deleteReq.OwnerIDs))
		for i, v := range deleteReq.OwnerIDs {
			ownerIDsAny[i] = v
		}

		mockTxOwnerRepo.EXPECT().ListWithUnpublishedByIDs(gomock.Any(), ownerIDsAny...).Return(owners, nil)
		mockTxOwnerRepo.EXPECT().FindOwnerIDsByImageID(gomock.Any(), mediaSvcID, deleteReq.OwnerIDs).Return([]string{}, nil)
		mockTxOwnerRepo.EXPECT().DeleteImageBatch(gomock.Any(), owners, &imagemodel.Image{MediaServiceID: deleteReq.MediaServiceID}).Return(nil)

		// Act
		affectedOwners, err := testService.DeleteImageBatch(context.Background(), deleteReq, mockOwnerRepo)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0, affectedOwners)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		invalidReq := &imagemodel.DeleteBatchRequst{
			MediaServiceID: "not-a-uuid",
		}

		// Act
		_, err := testService.DeleteImageBatch(context.Background(), invalidReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("owners not found", func(t *testing.T) {
		// Arrange
		mockImageRepo.EXPECT().DB().Return(db)

		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		// Convert slice of strings to slice of any for variadic mock expectation
		ownerIDsAny := make([]any, len(deleteReq.OwnerIDs))
		for i, v := range deleteReq.OwnerIDs {
			ownerIDsAny[i] = v
		}

		mockTxOwnerRepo.EXPECT().ListWithUnpublishedByIDs(gomock.Any(), ownerIDsAny...).Return([]image_owner.Owner{}, nil)

		// Act
		_, err := testService.DeleteImageBatch(context.Background(), deleteReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrOwnersNotFound)
	})

	t.Run("associations not found", func(t *testing.T) {
		// Arrange
		mockImageRepo.EXPECT().DB().Return(db)

		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		mockOwners := []mockOwner{
			{id: ownerID_1, uploadedImageAmount: 2},
			{id: ownerID_2, uploadedImageAmount: 3},
		}
		owners := make([]image_owner.Owner, len(mockOwners))
		for i := range mockOwners {
			owners[i] = &mockOwners[i]
		}

		// Convert slice of strings to slice of any for variadic mock expectation
		ownerIDsAny := make([]any, len(deleteReq.OwnerIDs))
		for i, v := range deleteReq.OwnerIDs {
			ownerIDsAny[i] = v
		}

		mockTxOwnerRepo.EXPECT().ListWithUnpublishedByIDs(gomock.Any(), ownerIDsAny...).Return(owners, nil)
		mockTxOwnerRepo.EXPECT().FindOwnerIDsByImageID(gomock.Any(), mediaSvcID, deleteReq.OwnerIDs).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.DeleteImageBatch(context.Background(), deleteReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrAssociationsNotFound)
	})

	t.Run("delete image batch db error", func(t *testing.T) {
		// Arrange
		mockImageRepo.EXPECT().DB().Return(db)

		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		mockOwners := []mockOwner{
			{id: ownerID_1, uploadedImageAmount: 2},
			{id: ownerID_2, uploadedImageAmount: 3},
		}
		owners := make([]image_owner.Owner, len(mockOwners))
		for i := range mockOwners {
			owners[i] = &mockOwners[i]
		}

		dbErr := errors.New("database error")
		// Convert slice of strings to slice of any for variadic mock expectation
		ownerIDsAny := make([]any, len(deleteReq.OwnerIDs))
		for i, v := range deleteReq.OwnerIDs {
			ownerIDsAny[i] = v
		}

		mockTxOwnerRepo.EXPECT().ListWithUnpublishedByIDs(gomock.Any(), ownerIDsAny...).Return(owners, nil)
		mockTxOwnerRepo.EXPECT().FindOwnerIDsByImageID(gomock.Any(), mediaSvcID, deleteReq.OwnerIDs).Return(deleteReq.OwnerIDs, nil)
		mockTxOwnerRepo.EXPECT().DeleteImageBatch(gomock.Any(), owners, &imagemodel.Image{MediaServiceID: deleteReq.MediaServiceID}).Return(dbErr)

		// Act
		_, err := testService.DeleteImageBatch(context.Background(), deleteReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to batch delete image from owners")
	})

	t.Run("decrement uploaded image amount db error", func(t *testing.T) {
		// Arrange
		mockImageRepo.EXPECT().DB().Return(db)

		mockTxOwnerRepo := imageownermock.NewMockOwnerRepo[image_owner.Owner](ctrl)
		mockOwnerRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxOwnerRepo)

		mockOwners := []mockOwner{
			{id: ownerID_1, uploadedImageAmount: 2},
			{id: ownerID_2, uploadedImageAmount: 3},
		}
		owners := make([]image_owner.Owner, len(mockOwners))
		for i := range mockOwners {
			owners[i] = &mockOwners[i]
		}

		dbErr := errors.New("database error")
		// Convert slice of strings to slice of any for variadic mock expectation
		ownerIDsAny := make([]any, len(deleteReq.OwnerIDs))
		for i, v := range deleteReq.OwnerIDs {
			ownerIDsAny[i] = v
		}

		mockTxOwnerRepo.EXPECT().ListWithUnpublishedByIDs(gomock.Any(), ownerIDsAny...).Return(owners, nil)
		mockTxOwnerRepo.EXPECT().FindOwnerIDsByImageID(gomock.Any(), mediaSvcID, deleteReq.OwnerIDs).Return(deleteReq.OwnerIDs, nil)
		mockTxOwnerRepo.EXPECT().DeleteImageBatch(gomock.Any(), owners, &imagemodel.Image{MediaServiceID: deleteReq.MediaServiceID}).Return(nil)
		mockTxOwnerRepo.EXPECT().DecrementImageCount(gomock.Any(), deleteReq.OwnerIDs).Return(int64(0), dbErr)

		// Act
		_, err := testService.DeleteImageBatch(context.Background(), deleteReq, mockOwnerRepo)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decrement uploaded image count from owners")
	})
}
