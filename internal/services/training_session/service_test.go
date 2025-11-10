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

package trainingsession

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	"github.com/mikhail5545/product-service-go/internal/models/product"
	trainingsession "github.com/mikhail5545/product-service-go/internal/models/training_session"
	imageservice "github.com/mikhail5545/product-service-go/internal/services/image_manager"
	productmock "github.com/mikhail5545/product-service-go/internal/test/database/product_mock"
	trainingsessionmock "github.com/mikhail5545/product-service-go/internal/test/database/training_session_mock"

	imageservicemock "github.com/mikhail5545/product-service-go/internal/test/services/image_mock"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	tsID := uuid.New().String()
	productID := uuid.New().String()

	mockTrainingSession := &trainingsession.TrainingSession{
		ID:               tsID,
		DurationMinutes:  30,
		Format:           "online",
		ShortDescription: "Training session short description",
		Name:             "Training session name",
	}

	mockProduct := &product.Product{
		ID:          productID,
		DetailsID:   tsID,
		DetailsType: "training_session",
		Price:       35.55,
	}

	expectedDetails := &trainingsession.TrainingSessionDetails{
		TrainingSession: mockTrainingSession,
		Price:           mockProduct.Price,
		ProductID:       mockProduct.ID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTrainingSessionRepo.EXPECT().Get(gomock.Any(), tsID).Return(mockTrainingSession, nil)
		mockProductRepo.EXPECT().SelectByDetailsID(gomock.Any(), tsID, gomock.Any()).Return(mockProduct, nil)

		// Act
		details, err := testService.Get(context.Background(), tsID)

		// Arrange
		assert.NoError(t, err)
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("Get() got %v, want %v", details, expectedDetails)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		_, err := testService.Get(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTrainingSessionRepo.EXPECT().Get(gomock.Any(), tsID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.Get(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockTrainingSessionRepo.EXPECT().Get(gomock.Any(), tsID).Return(nil, dbErr)

		// Act
		_, err := testService.Get(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	tsID := uuid.New().String()
	productID := uuid.New().String()

	mockTrainingSession := &trainingsession.TrainingSession{
		ID:               tsID,
		DurationMinutes:  30,
		Format:           "online",
		ShortDescription: "Training session short description",
		Name:             "Training session name",
	}

	mockProduct := &product.Product{
		ID:          productID,
		DetailsID:   tsID,
		DetailsType: "training_session",
		Price:       35.55,
	}

	expectedDetails := &trainingsession.TrainingSessionDetails{
		TrainingSession: mockTrainingSession,
		Price:           mockProduct.Price,
		ProductID:       mockProduct.ID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTrainingSessionRepo.EXPECT().GetWithDeleted(gomock.Any(), tsID).Return(mockTrainingSession, nil)
		mockProductRepo.EXPECT().SelectWithDeletedByDetailsID(gomock.Any(), tsID, gomock.Any()).Return(mockProduct, nil)

		// Act
		details, err := testService.GetWithDeleted(context.Background(), tsID)

		// Arrange
		assert.NoError(t, err)
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("GetWithDeleted() got %v, want %v", details, expectedDetails)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		_, err := testService.GetWithDeleted(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTrainingSessionRepo.EXPECT().GetWithDeleted(gomock.Any(), tsID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockTrainingSessionRepo.EXPECT().GetWithDeleted(gomock.Any(), tsID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetWithUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	tsID := uuid.New().String()
	productID := uuid.New().String()

	mockTrainingSession := &trainingsession.TrainingSession{
		ID:               tsID,
		DurationMinutes:  30,
		Format:           "online",
		ShortDescription: "Training session short description",
		Name:             "Training session name",
	}

	mockProduct := &product.Product{
		ID:          productID,
		DetailsID:   tsID,
		DetailsType: "training_session",
		Price:       35.55,
	}

	expectedDetails := &trainingsession.TrainingSessionDetails{
		TrainingSession: mockTrainingSession,
		Price:           mockProduct.Price,
		ProductID:       mockProduct.ID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTrainingSessionRepo.EXPECT().GetWithUnpublished(gomock.Any(), tsID).Return(mockTrainingSession, nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByDetailsID(gomock.Any(), tsID, gomock.Any()).Return(mockProduct, nil)

		// Act
		details, err := testService.GetWithUnpublished(context.Background(), tsID)

		// Arrange
		assert.NoError(t, err)
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("GetWithUnpublished() got %v, want %v", details, expectedDetails)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTrainingSessionRepo.EXPECT().GetWithUnpublished(gomock.Any(), tsID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockTrainingSessionRepo.EXPECT().GetWithUnpublished(gomock.Any(), tsID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	tsID_1 := uuid.New().String()
	tsID_2 := uuid.New().String()
	pID_1 := uuid.New().String()
	pID_2 := uuid.New().String()

	mockTrainingSessions := []trainingsession.TrainingSession{
		{
			ID:               tsID_1,
			Name:             "Training session 1 name",
			ShortDescription: "Training session 1 short description",
			DurationMinutes:  30,
			Format:           "online",
		},
		{
			ID:               tsID_2,
			Name:             "Training session 2name",
			ShortDescription: "Training session 2 short description",
			DurationMinutes:  30,
			Format:           "online",
		},
	}

	mockProducts := []product.Product{
		{
			ID:          pID_1,
			Price:       34.44,
			DetailsID:   tsID_1,
			DetailsType: "training_session",
		},
		{
			ID:          pID_2,
			Price:       25.44,
			DetailsID:   tsID_2,
			DetailsType: "training_session",
		},
	}

	expectedDetails := []trainingsession.TrainingSessionDetails{
		{
			TrainingSession: &mockTrainingSessions[0],
			Price:           mockProducts[0].Price,
			ProductID:       mockProducts[0].ID,
		},
		{
			TrainingSession: &mockTrainingSessions[1],
			Price:           mockProducts[1].Price,
			ProductID:       mockProducts[1].ID,
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockTrainingSessionRepo.EXPECT().List(gomock.Any(), limit, offset).Return(mockTrainingSessions, nil)
		mockTrainingSessionRepo.EXPECT().Count(gomock.Any()).Return(int64(2), nil)
		mockProductRepo.EXPECT().SelectByDetailsIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)

		// Act
		details, total, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, len(expectedDetails), len(details))
	})

	t.Run("success empty list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockTrainingSessionRepo.EXPECT().List(gomock.Any(), limit, offset).Return([]trainingsession.TrainingSession{}, nil)
		mockTrainingSessionRepo.EXPECT().Count(gomock.Any()).Return(int64(0), nil)
		mockProductRepo.EXPECT().SelectByDetailsIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]product.Product{}, nil)

		// Act
		details, total, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, details, 0)
	})

	t.Run("db error on count", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("db count error")
		mockTrainingSessionRepo.EXPECT().List(gomock.Any(), limit, offset).Return(mockTrainingSessions, nil)
		mockProductRepo.EXPECT().SelectByDetailsIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)
		mockTrainingSessionRepo.EXPECT().Count(gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, _, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})

	t.Run("db error on product select", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("db product error")
		mockTrainingSessionRepo.EXPECT().List(gomock.Any(), limit, offset).Return(mockTrainingSessions, nil)
		mockProductRepo.EXPECT().SelectByDetailsIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)

		// Act
		_, _, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})

	t.Run("db error on list", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		limit, offset := 2, 0
		mockTrainingSessionRepo.EXPECT().List(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_ListDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	tsID_1 := uuid.New().String()
	tsID_2 := uuid.New().String()
	pID_1 := uuid.New().String()
	pID_2 := uuid.New().String()

	mockTrainingSessions := []trainingsession.TrainingSession{
		{
			ID:               tsID_1,
			Name:             "Training session 1 name",
			ShortDescription: "Training session 1 short description",
			DurationMinutes:  30,
			Format:           "online",
		},
		{
			ID:               tsID_2,
			Name:             "Training session 2name",
			ShortDescription: "Training session 2 short description",
			DurationMinutes:  30,
			Format:           "online",
		},
	}

	mockProducts := []product.Product{
		{
			ID:          pID_1,
			Price:       34.44,
			DetailsID:   tsID_1,
			DetailsType: "training_session",
		},
		{
			ID:          pID_2,
			Price:       25.44,
			DetailsID:   tsID_2,
			DetailsType: "training_session",
		},
	}

	expectedDetails := []trainingsession.TrainingSessionDetails{
		{
			TrainingSession: &mockTrainingSessions[0],
			Price:           mockProducts[0].Price,
			ProductID:       mockProducts[0].ID,
		},
		{
			TrainingSession: &mockTrainingSessions[1],
			Price:           mockProducts[1].Price,
			ProductID:       mockProducts[1].ID,
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockTrainingSessionRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(mockTrainingSessions, nil)
		mockTrainingSessionRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(2), nil)
		mockProductRepo.EXPECT().SelectWithDeletedByDetailsIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)

		// Act
		details, total, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, len(expectedDetails), len(details))
	})

	t.Run("success empty list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockTrainingSessionRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return([]trainingsession.TrainingSession{}, nil)
		mockTrainingSessionRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(0), nil)
		mockProductRepo.EXPECT().SelectWithDeletedByDetailsIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]product.Product{}, nil)

		// Act
		details, total, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, details, 0)
	})

	t.Run("db error on count", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("db count error")
		mockTrainingSessionRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(mockTrainingSessions, nil)
		mockTrainingSessionRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, _, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})

	t.Run("db error on list", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		limit, offset := 2, 0
		mockTrainingSessionRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_ListUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	tsID_1 := uuid.New().String()
	tsID_2 := uuid.New().String()
	pID_1 := uuid.New().String()
	pID_2 := uuid.New().String()

	mockTrainingSessions := []trainingsession.TrainingSession{
		{
			ID:               tsID_1,
			Name:             "Training session 1 name",
			ShortDescription: "Training session 1 short description",
			DurationMinutes:  30,
			Format:           "online",
		},
		{
			ID:               tsID_2,
			Name:             "Training session 2name",
			ShortDescription: "Training session 2 short description",
			DurationMinutes:  30,
			Format:           "online",
		},
	}

	mockProducts := []product.Product{
		{
			ID:          pID_1,
			Price:       34.44,
			DetailsID:   tsID_1,
			DetailsType: "training_session",
		},
		{
			ID:          pID_2,
			Price:       25.44,
			DetailsID:   tsID_2,
			DetailsType: "training_session",
		},
	}

	expectedDetails := []trainingsession.TrainingSessionDetails{
		{
			TrainingSession: &mockTrainingSessions[0],
			Price:           mockProducts[0].Price,
			ProductID:       mockProducts[0].ID,
		},
		{
			TrainingSession: &mockTrainingSessions[1],
			Price:           mockProducts[1].Price,
			ProductID:       mockProducts[1].ID,
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockTrainingSessionRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(mockTrainingSessions, nil)
		mockTrainingSessionRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(2), nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByDetailsIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)

		// Act
		details, total, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, len(expectedDetails), len(details))
	})

	t.Run("success empty list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockTrainingSessionRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return([]trainingsession.TrainingSession{}, nil)
		mockTrainingSessionRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(0), nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByDetailsIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]product.Product{}, nil)

		// Act
		details, total, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, details, 0)
	})

	t.Run("db error on count", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("db count error")
		mockTrainingSessionRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(mockTrainingSessions, nil)
		mockTrainingSessionRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, _, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})

	t.Run("db error on list", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		limit, offset := 2, 0
		mockTrainingSessionRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	createReq := &trainingsession.CreateRequest{
		Name:             "Training session name",
		ShortDescription: "Training session short description",
		DurationMinutes:  30,
		Price:            44.55,
		Format:           "online",
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		var createdTs *trainingsession.TrainingSession
		mockTxTrainingSessionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, ts *trainingsession.TrainingSession) {
				createdTs = ts
			}).Return(nil).AnyTimes()

		var createdProduct *product.Product
		mockTxProductRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, p *product.Product) {
				createdProduct = p
			}).Return(nil).AnyTimes()

		// Act
		resp, err := testService.Create(context.Background(), createReq)

		// Assert
		assert.NoError(t, err)
		if _, err := uuid.Parse(createdTs.ID); err != nil {
			t.Errorf("Expected training_session.ID to be a valid UUID, got %s", createdTs.ID)
		}
		assert.Equal(t, createReq.Name, createdTs.Name)
		assert.Equal(t, createReq.ShortDescription, createdTs.ShortDescription)
		assert.Equal(t, createReq.Format, createdTs.Format)
		assert.Equal(t, createReq.DurationMinutes, createdTs.DurationMinutes)
		assert.False(t, createdTs.InStock)

		if _, err := uuid.Parse(createdProduct.ID); err != nil {
			t.Errorf("Expected product.ID to be a valid UUID, got %s", createdProduct.ID)
		}
		assert.Equal(t, createReq.Price, createdProduct.Price)
		assert.Equal(t, createdTs.ID, createdProduct.DetailsID)
		assert.Equal(t, "training_session", createdProduct.DetailsType)
		assert.False(t, createdProduct.InStock)
		assert.Equal(t, createdTs.ID, resp.ID)
		assert.Equal(t, createdProduct.ID, resp.ProductID)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		invalidReq := &trainingsession.CreateRequest{
			Name:             "1invalidname",
			ShortDescription: "Valid description",
			DurationMinutes:  15,       // invalid
			Format:           "format", // invalid
			Price:            324.44,
		}

		// Act
		_, err := testService.Create(context.Background(), invalidReq)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxTrainingSessionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)

		// Act
		_, err := testService.Create(context.Background(), createReq)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().SetInStock(gomock.Any(), tsID, true).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), tsID, true).Return(int64(1), nil)

		// Act
		err := testService.Publish(context.Background(), tsID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		err := testService.Publish(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().SetInStock(gomock.Any(), tsID, true).Return(int64(0), nil)

		// Act
		err := testService.Publish(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxTrainingSessionRepo.EXPECT().SetInStock(gomock.Any(), tsID, true).Return(int64(0), dbErr)

		// Act
		err := testService.Publish(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Unpublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().SetInStock(gomock.Any(), tsID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), tsID, false).Return(int64(1), nil)

		// Act
		err := testService.Unpublish(context.Background(), tsID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		err := testService.Unpublish(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().SetInStock(gomock.Any(), tsID, false).Return(int64(0), nil)

		// Act
		err := testService.Unpublish(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxTrainingSessionRepo.EXPECT().SetInStock(gomock.Any(), tsID, false).Return(int64(0), dbErr)

		// Act
		err := testService.Unpublish(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tsID := uuid.New().String()
	productID := uuid.New().String()

	mockTrainingSession := &trainingsession.TrainingSession{
		ID:               tsID,
		Name:             "Old training session name",
		ShortDescription: "Old training session short description",
		DurationMinutes:  30,
		Format:           "online",
	}

	mockProduct := &product.Product{
		ID:          productID,
		Price:       45.55,
		DetailsID:   tsID,
		DetailsType: "training_session",
	}

	newName := "New training session name"
	newLongDescription := "New training session long description"
	newTags := []string{"new", "training", "tags", "session"}
	newPrice := float32(88.99)

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().Select(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockTrainingSession, nil)
		mockTxProductRepo.EXPECT().SelectByDetailsID(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProduct, nil)

		var tsUpdates map[string]any
		mockTxTrainingSessionRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, _ *trainingsession.TrainingSession, u map[string]any) {
				tsUpdates = u
			}).Return(int64(1), nil).AnyTimes()

		var productUpdates map[string]any
		mockTxProductRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, _ *product.Product, u map[string]any) {
				productUpdates = u
			}).Return(int64(1), nil).AnyTimes()

		// Act
		updates, err := testService.Update(context.Background(), &trainingsession.UpdateRequest{
			ID:              tsID,
			Name:            &newName,
			LongDescription: &newLongDescription,
			Tags:            newTags,
			Price:           &newPrice,
		})

		// Assert
		assert.NoError(t, err)
		// Assert training session
		tsUpdatesFromResp, ok := updates["training_session"].(map[string]any)
		if !ok {
			t.Errorf("response does not contain 'training_session' key")
		}
		if name, ok := tsUpdatesFromResp["name"].(string); !ok || name != newName {
			t.Errorf("training_session.Name in response %v, want %s", tsUpdatesFromResp["name"], newName)
		}
		if longDesc, ok := tsUpdatesFromResp["long_description"].(string); !ok || longDesc != newLongDescription {
			t.Errorf("training_session.LongDescription in response %v, want %s", tsUpdatesFromResp["long_description"], newLongDescription)
		}
		if tags, ok := tsUpdatesFromResp["tags"].([]string); !ok || !reflect.DeepEqual(tags, newTags) {
			t.Errorf("training_session.Tags in response %v, want %v", tsUpdatesFromResp["tags"], newTags)
		}
		if name, ok := tsUpdates["name"].(string); !ok || name != newName {
			t.Errorf("training_session.Name passed to repo %v, want %s", tsUpdates["name"], newName)
		}
		if longDesc, ok := tsUpdates["long_description"].(string); !ok || longDesc != newLongDescription {
			t.Errorf("training_session.LongDescription passed to repo %v, want %s", tsUpdates["long_description"], newLongDescription)
		}
		if tags, ok := tsUpdates["tags"].([]string); !ok || !reflect.DeepEqual(tags, newTags) {
			t.Errorf("training_session.Tags passed to repo %v, want %v", tsUpdates["tags"], newTags)
		}

		// Assert product
		productUpdatesFromResp, ok := updates["product"].(map[string]any)
		if !ok {
			t.Errorf("response does not contain 'product' key")
		}
		if price, ok := productUpdatesFromResp["price"].(float32); !ok || price != newPrice {
			t.Errorf("product.Price in response %v, want %f", productUpdatesFromResp["price"], newPrice)
		}
		if price, ok := productUpdates["price"].(float32); !ok || price != newPrice {
			t.Errorf("product.Price passed to repo %v, want %f", productUpdates["price"], newPrice)
		}
	})

	t.Run("success with no updates", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().Select(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockTrainingSession, nil)
		mockTxProductRepo.EXPECT().SelectByDetailsID(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProduct, nil)

		// Act
		_, err := testService.Update(context.Background(), &trainingsession.UpdateRequest{
			ID: tsID, // no new fields
		})

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		invalidPrice := float32(-99.55)
		invalidShortDescription := "3"

		// Act
		_, err := testService.Update(context.Background(), &trainingsession.UpdateRequest{
			ID:               tsID,
			Name:             &newName,
			Price:            &invalidPrice,
			ShortDescription: &invalidShortDescription,
		})

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().Select(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.Update(context.Background(), &trainingsession.UpdateRequest{
			ID:              tsID,
			Name:            &newName,
			LongDescription: &newLongDescription,
			Tags:            newTags,
			Price:           &newPrice,
		})

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().Select(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockTrainingSession, nil)
		mockTxProductRepo.EXPECT().SelectByDetailsID(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProduct, nil)
		dbErr := errors.New("database error")
		mockTxProductRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), dbErr).AnyTimes()

		// Act
		_, err := testService.Update(context.Background(), &trainingsession.UpdateRequest{
			ID:              tsID,
			Name:            &newName,
			LongDescription: &newLongDescription,
			Tags:            newTags,
			Price:           &newPrice,
		})

		// Assert
		assert.Error(t, err)
	})
}

func TestService_AddImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	tsID := uuid.New().String()

	addReq := &imagemodel.AddRequest{
		URL:            "https://google.com",
		SecureURL:      "https://google.com",
		PublicID:       "public/id",
		MediaServiceID: uuid.New().String(),
		OwnerID:        tsID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockImageSvc.EXPECT().AddImage(gomock.Any(), addReq, gomock.Any()).Return(nil)

		// Act
		err := testService.AddImage(context.Background(), addReq)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		invalidReq := &imagemodel.AddRequest{
			URL:     "not a url",
			OwnerID: "invalid-UUID",
		}
		mockImageSvc.EXPECT().AddImage(gomock.Any(), invalidReq, gomock.Any()).Return(imageservice.ErrInvalidArgument)

		// Act
		err := testService.AddImage(context.Background(), invalidReq)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrInvalidArgument)
	})

	t.Run("training session not found", func(t *testing.T) {
		// Arrange
		mockImageSvc.EXPECT().AddImage(gomock.Any(), addReq, gomock.Any()).Return(imageservice.ErrOwnerNotFound)

		// Act
		err := testService.AddImage(context.Background(), addReq)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrOwnerNotFound)
	})

	t.Run("image limit", func(t *testing.T) {
		// Arrange
		mockImageSvc.EXPECT().AddImage(gomock.Any(), addReq, gomock.Any()).Return(imageservice.ErrImageLimitExceeded)

		// Act
		err := testService.AddImage(context.Background(), addReq)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrImageLimitExceeded)
	})
}

func TestSesrvice_DeleteImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	imgID := uuid.New().String()

	tsID := uuid.New().String()

	deleteReq := &imagemodel.DeleteRequest{
		MediaServiceID: imgID,
		OwnerID:        tsID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockImageSvc.EXPECT().DeleteImage(gomock.Any(), deleteReq, gomock.Any()).Return(nil)

		// Act
		err := testService.DeleteImage(context.Background(), deleteReq)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid rquest payload", func(t *testing.T) {
		// Arrange
		invalidReq := &imagemodel.DeleteRequest{
			MediaServiceID: "invalid-UUID",
			OwnerID:        "invalid-UUID",
		}
		mockImageSvc.EXPECT().DeleteImage(gomock.Any(), invalidReq, gomock.Any()).Return(imageservice.ErrInvalidArgument)

		// Act
		err := testService.DeleteImage(context.Background(), invalidReq)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrInvalidArgument)
	})

	t.Run("training session not found", func(t *testing.T) {
		// Arrange
		mockImageSvc.EXPECT().DeleteImage(gomock.Any(), deleteReq, gomock.Any()).Return(imageservice.ErrOwnerNotFound)

		// Act
		err := testService.DeleteImage(context.Background(), deleteReq)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrOwnerNotFound)
	})
}

func TestService_AddImageBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tsID_1 := uuid.New().String()
	tsID_2 := uuid.New().String()

	addBatchRequest := &imagemodel.AddBatchRequest{
		URL:            "some-url",
		SecureURL:      "some-secure-url",
		PublicID:       "some-public-id",
		MediaServiceID: uuid.New().String(),
		OwnerIDs:       []string{tsID_1, tsID_2},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)

		mockImageSvc.EXPECT().AddImageBatch(gomock.Any(), addBatchRequest, gomock.Any()).Return(2, nil)

		// Act
		_, err := testService.AddImageBatch(context.Background(), addBatchRequest)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)

		mockImageSvc.EXPECT().AddImageBatch(gomock.Any(), addBatchRequest, gomock.Any()).Return(0, imageservice.ErrInvalidArgument)

		// Act
		_, err := testService.AddImageBatch(context.Background(), addBatchRequest)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrInvalidArgument)
	})

	t.Run("training sessions not found", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)

		mockImageSvc.EXPECT().AddImageBatch(gomock.Any(), addBatchRequest, gomock.Any()).Return(0, imageservice.ErrOwnersNotFound)

		// Act
		_, err := testService.AddImageBatch(context.Background(), addBatchRequest)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrOwnersNotFound)
	})
}

func TestService_DeleteImageBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tsID_1 := uuid.New().String()
	tsID_2 := uuid.New().String()

	deleteBatchRequest := &imagemodel.DeleteBatchRequst{
		MediaServiceID: uuid.New().String(),
		OwnerIDs:       []string{tsID_1, tsID_2},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)

		mockImageSvc.EXPECT().DeleteImageBatch(gomock.Any(), deleteBatchRequest, gomock.Any()).Return(2, nil)

		// Act
		_, err := testService.DeleteImageBatch(context.Background(), deleteBatchRequest)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)

		mockImageSvc.EXPECT().DeleteImageBatch(gomock.Any(), deleteBatchRequest, gomock.Any()).Return(0, imageservice.ErrInvalidArgument)

		// Act
		_, err := testService.DeleteImageBatch(context.Background(), deleteBatchRequest)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrInvalidArgument)
	})

	t.Run("training sessions not found", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)

		mockImageSvc.EXPECT().DeleteImageBatch(gomock.Any(), deleteBatchRequest, gomock.Any()).Return(0, imageservice.ErrOwnersNotFound)

		// Act
		_, err := testService.DeleteImageBatch(context.Background(), deleteBatchRequest)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrOwnersNotFound)
	})
}

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().GetWithUnpublished(gomock.Any(), tsID).Return(&trainingsession.TrainingSession{}, nil)
		mockTxTrainingSessionRepo.EXPECT().SetInStock(gomock.Any(), tsID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), tsID, false).Return(int64(1), nil)
		mockTxTrainingSessionRepo.EXPECT().Delete(gomock.Any(), tsID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().DeleteByDetailsID(gomock.Any(), tsID).Return(int64(1), nil)

		// Act
		err := testService.Delete(context.Background(), tsID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		err := testService.Delete(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().GetWithUnpublished(gomock.Any(), tsID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		err := testService.Delete(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("product not found", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().GetWithUnpublished(gomock.Any(), tsID).Return(&trainingsession.TrainingSession{}, nil)
		mockTxTrainingSessionRepo.EXPECT().SetInStock(gomock.Any(), tsID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), tsID, false).Return(int64(0), nil)

		// Act
		err := testService.Delete(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().GetWithUnpublished(gomock.Any(), tsID).Return(&trainingsession.TrainingSession{}, nil)
		mockTxTrainingSessionRepo.EXPECT().SetInStock(gomock.Any(), tsID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), tsID, false).Return(int64(1), nil)
		dbErr := errors.New("database error")
		mockTxTrainingSessionRepo.EXPECT().Delete(gomock.Any(), tsID).Return(int64(0), dbErr)

		// Act
		assert.Error(t, err)
	})
}

func TestService_DeletePermanent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().DeletePermanent(gomock.Any(), tsID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().DeletePermanentByDetailsID(gomock.Any(), tsID).Return(int64(1), nil)

		// Act
		err := testService.DeletePermanent(context.Background(), tsID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		err := testService.DeletePermanent(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().DeletePermanent(gomock.Any(), tsID).Return(int64(0), nil)

		// Act
		err := testService.DeletePermanent(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("product not found", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().DeletePermanent(gomock.Any(), tsID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().DeletePermanentByDetailsID(gomock.Any(), tsID).Return(int64(0), nil)

		// Act
		err := testService.DeletePermanent(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxTrainingSessionRepo.EXPECT().DeletePermanent(gomock.Any(), tsID).Return(int64(1), dbErr)

		// Act
		err := testService.DeletePermanent(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Restore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockTrainingSessionRepo, mockProductRepo, mockImageSvc)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().Restore(gomock.Any(), tsID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().RestoreByDetailsID(gomock.Any(), tsID).Return(int64(1), nil)

		// Act
		err := testService.Restore(context.Background(), tsID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		err := testService.Restore(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().Restore(gomock.Any(), tsID).Return(int64(0), nil)

		// Act
		err := testService.Restore(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("product not found", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxTrainingSessionRepo.EXPECT().Restore(gomock.Any(), tsID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().RestoreByDetailsID(gomock.Any(), tsID).Return(int64(0), nil)

		// Act
		err := testService.Restore(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxTrainingSessionRepo := trainingsessionmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockTrainingSessionRepo.EXPECT().DB().Return(db).AnyTimes()
		mockTrainingSessionRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxTrainingSessionRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxTrainingSessionRepo.EXPECT().Restore(gomock.Any(), tsID).Return(int64(1), dbErr)

		// Act
		err := testService.Restore(context.Background(), tsID)

		// Assert
		assert.Error(t, err)
	})
}
