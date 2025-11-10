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

package physicalgood

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	physicalgood "github.com/mikhail5545/product-service-go/internal/models/physical_good"
	"github.com/mikhail5545/product-service-go/internal/models/product"
	imageservice "github.com/mikhail5545/product-service-go/internal/services/image_manager"
	physicalgoodmock "github.com/mikhail5545/product-service-go/internal/test/database/physical_good_mock"
	productmock "github.com/mikhail5545/product-service-go/internal/test/database/product_mock"

	imageservicemock "github.com/mikhail5545/product-service-go/internal/test/services/image_mock"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	physicalGoodID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	mockPhysicalGood := &physicalgood.PhysicalGood{
		ID:               physicalGoodID,
		Name:             "Physical good name",
		ShortDescription: "Physical good short description",
		InStock:          true,
		ShippingRequired: true,
	}

	mockProduct := &product.Product{
		ID:          "prod-id",
		InStock:     true,
		DetailsID:   physicalGoodID,
		Price:       35.55,
		DetailsType: "physical_good",
	}

	expectedDetails := &physicalgood.PhysicalGoodDetails{
		PhysicalGood: mockPhysicalGood,
		Price:        mockProduct.Price,
		ProductID:    mockProduct.ID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockPhysicalGoodRepo.EXPECT().Get(gomock.Any(), physicalGoodID).Return(mockPhysicalGood, nil)
		mockProductRepo.EXPECT().SelectByDetailsID(gomock.Any(), physicalGoodID, gomock.Any()).Return(mockProduct, nil)

		// Act
		details, err := testService.Get(context.Background(), physicalGoodID)

		// Assert
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
		mockPhysicalGoodRepo.EXPECT().Get(gomock.Any(), physicalGoodID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.Get(context.Background(), physicalGoodID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().Get(gomock.Any(), physicalGoodID).Return(nil, dbErr)

		// Act
		_, err := testService.Get(context.Background(), physicalGoodID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	physicalGoodID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	mockPhysicalGood := &physicalgood.PhysicalGood{
		ID:               physicalGoodID,
		Name:             "Physical good name",
		ShortDescription: "Physical good short description",
		InStock:          true,
		ShippingRequired: true,
	}

	mockProduct := &product.Product{
		ID:          "prod-id",
		InStock:     true,
		DetailsID:   physicalGoodID,
		Price:       35.55,
		DetailsType: "physical_good",
	}

	expectedDetails := &physicalgood.PhysicalGoodDetails{
		PhysicalGood: mockPhysicalGood,
		Price:        mockProduct.Price,
		ProductID:    mockProduct.ID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockPhysicalGoodRepo.EXPECT().GetWithDeleted(gomock.Any(), physicalGoodID).Return(mockPhysicalGood, nil)
		mockProductRepo.EXPECT().SelectWithDeletedByDetailsID(gomock.Any(), physicalGoodID, gomock.Any()).Return(mockProduct, nil)

		// Act
		details, err := testService.GetWithDeleted(context.Background(), physicalGoodID)

		// Assert
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
		mockPhysicalGoodRepo.EXPECT().GetWithDeleted(gomock.Any(), physicalGoodID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), physicalGoodID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().GetWithDeleted(gomock.Any(), physicalGoodID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), physicalGoodID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetWithUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	physicalGoodID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	mockPhysicalGood := &physicalgood.PhysicalGood{
		ID:               physicalGoodID,
		Name:             "Physical good name",
		ShortDescription: "Physical good short description",
		InStock:          true,
		ShippingRequired: true,
	}

	mockProduct := &product.Product{
		ID:          "prod-id",
		InStock:     true,
		DetailsID:   physicalGoodID,
		Price:       35.55,
		DetailsType: "physical_good",
	}

	expectedDetails := &physicalgood.PhysicalGoodDetails{
		PhysicalGood: mockPhysicalGood,
		Price:        mockProduct.Price,
		ProductID:    mockProduct.ID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockPhysicalGoodRepo.EXPECT().GetWithUnpublished(gomock.Any(), physicalGoodID).Return(mockPhysicalGood, nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByDetailsID(gomock.Any(), physicalGoodID, gomock.Any()).Return(mockProduct, nil)

		// Act
		details, err := testService.GetWithUnpublished(context.Background(), physicalGoodID)

		// Assert
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
		mockPhysicalGoodRepo.EXPECT().GetWithUnpublished(gomock.Any(), physicalGoodID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), physicalGoodID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().GetWithUnpublished(gomock.Any(), physicalGoodID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), physicalGoodID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	phg1ID := "0d9828df-c57b-4629-9729-8c9641598e17"
	phg2ID := "a33845f2-1c3c-4397-9380-7ecdb1d8c853"

	mockPhysicalGoods := []physicalgood.PhysicalGood{
		{
			ID:               phg1ID,
			Name:             "First physical good name",
			ShortDescription: "First physical good description",
		},
		{
			ID:               phg2ID,
			Name:             "Second physical good name",
			ShortDescription: "Second physical good description",
		},
	}

	mockProducts := []product.Product{
		{
			ID:          "prod-1-ID",
			Price:       34.24,
			DetailsID:   phg1ID,
			DetailsType: "physical_good",
		},
		{
			ID:          "prod-2-ID",
			Price:       3443.254,
			DetailsID:   phg2ID,
			DetailsType: "physical_good",
		},
	}

	expectedDetails := []physicalgood.PhysicalGoodDetails{
		{
			PhysicalGood: &mockPhysicalGoods[0],
			Price:        mockProducts[0].Price,
			ProductID:    mockProducts[0].ID,
		},
		{
			PhysicalGood: &mockPhysicalGoods[1],
			Price:        mockProducts[1].Price,
			ProductID:    mockProducts[1].ID,
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockPhysicalGoodRepo.EXPECT().List(gomock.Any(), limit, offset).Return(mockPhysicalGoods, nil)
		mockPhysicalGoodRepo.EXPECT().Count(gomock.Any()).Return(int64(2), nil)
		mockProductRepo.EXPECT().SelectByDetailsIDs(gomock.Any(), []string{phg1ID, phg2ID}, gomock.Any()).Return(mockProducts, nil)

		// Act
		details, total, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, len(expectedDetails), len(details))
	})

	t.Run("db error on count", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockPhysicalGoodRepo.EXPECT().List(gomock.Any(), limit, offset).Return(mockPhysicalGoods, nil)

		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().Count(gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, _, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})

	t.Run("db error on list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().List(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_ListDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	phg1ID := "0d9828df-c57b-4629-9729-8c9641598e17"
	phg2ID := "a33845f2-1c3c-4397-9380-7ecdb1d8c853"

	mockPhysicalGoods := []physicalgood.PhysicalGood{
		{
			ID:               phg1ID,
			Name:             "First physical good name",
			ShortDescription: "First physical good description",
		},
		{
			ID:               phg2ID,
			Name:             "Second physical good name",
			ShortDescription: "Second physical good description",
		},
	}

	mockProducts := []product.Product{
		{
			ID:          "prod-1-ID",
			Price:       34.24,
			DetailsID:   phg1ID,
			DetailsType: "physical_good",
		},
		{
			ID:          "prod-2-ID",
			Price:       3443.254,
			DetailsID:   phg2ID,
			DetailsType: "physical_good",
		},
	}

	expectedDetails := []physicalgood.PhysicalGoodDetails{
		{
			PhysicalGood: &mockPhysicalGoods[0],
			Price:        mockProducts[0].Price,
			ProductID:    mockProducts[0].ID,
		},
		{
			PhysicalGood: &mockPhysicalGoods[1],
			Price:        mockProducts[1].Price,
			ProductID:    mockProducts[1].ID,
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockPhysicalGoodRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(mockPhysicalGoods, nil)
		mockPhysicalGoodRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(2), nil)
		mockProductRepo.EXPECT().SelectWithDeletedByDetailsIDs(gomock.Any(), []string{phg1ID, phg2ID}, gomock.Any()).Return(mockProducts, nil)

		// Act
		details, total, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, len(expectedDetails), len(details))
	})

	t.Run("db error on count", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockPhysicalGoodRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(mockPhysicalGoods, nil)

		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, _, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})

	t.Run("db error on list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_ListUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	phg1ID := "0d9828df-c57b-4629-9729-8c9641598e17"
	phg2ID := "a33845f2-1c3c-4397-9380-7ecdb1d8c853"

	mockPhysicalGoods := []physicalgood.PhysicalGood{
		{
			ID:               phg1ID,
			Name:             "First physical good name",
			ShortDescription: "First physical good description",
		},
		{
			ID:               phg2ID,
			Name:             "Second physical good name",
			ShortDescription: "Second physical good description",
		},
	}

	mockProducts := []product.Product{
		{
			ID:          "prod-1-ID",
			Price:       34.24,
			DetailsID:   phg1ID,
			DetailsType: "physical_good",
		},
		{
			ID:          "prod-2-ID",
			Price:       3443.254,
			DetailsID:   phg2ID,
			DetailsType: "physical_good",
		},
	}

	expectedDetails := []physicalgood.PhysicalGoodDetails{
		{
			PhysicalGood: &mockPhysicalGoods[0],
			Price:        mockProducts[0].Price,
			ProductID:    mockProducts[0].ID,
		},
		{
			PhysicalGood: &mockPhysicalGoods[1],
			Price:        mockProducts[1].Price,
			ProductID:    mockProducts[1].ID,
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockPhysicalGoodRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(mockPhysicalGoods, nil)
		mockPhysicalGoodRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(2), nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByDetailsIDs(gomock.Any(), []string{phg1ID, phg2ID}, gomock.Any()).Return(mockProducts, nil)

		// Act
		details, total, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, len(expectedDetails), len(details))
	})

	t.Run("db error on count", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockPhysicalGoodRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(mockPhysicalGoods, nil)

		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, _, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})

	t.Run("db error on list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})
}

func TestSesrvice_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	createReq := physicalgood.CreateRequest{
		Name:             "Physical good name",
		ShortDescription: "Physical good short description",
		Price:            43.22,
		Amount:           2,
		ShippingRequired: false,
	}

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		var createdPhysicalGood *physicalgood.PhysicalGood
		mockTxPhysicalGoodRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, g *physicalgood.PhysicalGood) {
				createdPhysicalGood = g
			})

		var createdProduct *product.Product
		mockTxProductRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, p *product.Product) {
				createdProduct = p
			})

		// Act
		resp, err := testService.Create(context.Background(), &createReq)

		// Assert
		assert.NoError(t, err)
		if _, err := uuid.Parse(createdPhysicalGood.ID); err != nil {
			t.Errorf("Expected physicalGood.ID to be a valid UUID, got %s", createdPhysicalGood.ID)
		}
		assert.Equal(t, createReq.Name, createdPhysicalGood.Name)
		assert.Equal(t, createReq.ShortDescription, createdPhysicalGood.ShortDescription)
		assert.Equal(t, createReq.Amount, createdPhysicalGood.Amount)
		assert.Equal(t, createReq.Amount, createdPhysicalGood.Amount)
		assert.Equal(t, createReq.ShippingRequired, createdPhysicalGood.ShippingRequired)
		assert.False(t, createdPhysicalGood.InStock)

		if _, err := uuid.Parse(createdProduct.ID); err != nil {
			t.Errorf("Expected product.ID to be a valid UUID, got %s", createdProduct.ID)
		}
		assert.Equal(t, createReq.Price, createdProduct.Price)
		assert.Equal(t, createdPhysicalGood.ID, createdProduct.DetailsID)
		assert.Equal(t, "physical_good", createdProduct.DetailsType)
		assert.False(t, createdProduct.InStock)
		assert.Equal(t, createdProduct.ID, resp.ProductID)
		assert.Equal(t, createdPhysicalGood.ID, resp.ID)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		// Act
		_, err := testService.Create(context.Background(), &physicalgood.CreateRequest{
			Name:             "3invalidname",
			ShortDescription: "Short description",
			Amount:           -44,
			Price:            55.3,
			ShippingRequired: false,
		})

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxPhysicalGoodRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr).AnyTimes()

		// Act
		_, err := testService.Create(context.Background(), &createReq)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	goodID := "0d9828df-c57b-4629-9729-8c9641598e17"

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().SetInStock(gomock.Any(), goodID, true).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), goodID, true).Return(int64(1), nil)

		// Act
		err := testService.Publish(context.Background(), goodID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		ivnalidID := "invalid-UUID"

		// Act
		err := testService.Publish(context.Background(), ivnalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().SetInStock(gomock.Any(), goodID, true).Return(int64(0), nil).AnyTimes()

		// Act
		err := testService.Publish(context.Background(), goodID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxPhysicalGoodRepo.EXPECT().SetInStock(gomock.Any(), goodID, true).Return(int64(0), dbErr)

		// Act
		err := testService.Publish(context.Background(), goodID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Unpublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	goodID := "0d9828df-c57b-4629-9729-8c9641598e17"

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().SetInStock(gomock.Any(), goodID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), goodID, false).Return(int64(1), nil)

		// Act
		err := testService.Unpublish(context.Background(), goodID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		ivnalidID := "invalid-UUID"

		// Act
		err := testService.Unpublish(context.Background(), ivnalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().SetInStock(gomock.Any(), goodID, false).Return(int64(0), nil)

		// Act
		err := testService.Unpublish(context.Background(), goodID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxPhysicalGoodRepo.EXPECT().SetInStock(gomock.Any(), goodID, false).Return(int64(0), dbErr)

		// Act
		err := testService.Unpublish(context.Background(), goodID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	goodID := "0d9828df-c57b-4629-9729-8c9641598e17"

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	mockPhysicalGood := &physicalgood.PhysicalGood{
		ID:               goodID,
		Name:             "Old physical good name",
		ShortDescription: "Old physical good short description",
		Amount:           4,
		ShippingRequired: false,
	}

	mockProduct := &product.Product{
		ID:          "product-ID",
		DetailsID:   goodID,
		Price:       34.22,
		DetailsType: "physical_good",
	}

	newName := "New physical good name"
	newAmount := 66
	newPrice := float32(88.34)
	newLongDescription := "Long description"
	newTags := []string{"new", "tags", "physicalgood"}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mockPhysicalGood, nil).AnyTimes()
		mockTxProductRepo.EXPECT().SelectByDetailsID(gomock.Any(), goodID, gomock.Any()).Return(mockProduct, nil).AnyTimes()

		var goodUpdates map[string]any
		mockTxPhysicalGoodRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, _ *physicalgood.PhysicalGood, u map[string]any) {
				goodUpdates = u
			}).Return(int64(1), nil).AnyTimes()

		var productUpdates map[string]any
		mockTxProductRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, _ *product.Product, u map[string]any) {
				productUpdates = u
			}).Return(int64(1), nil).AnyTimes()

		// Act
		updates, err := testService.Update(context.Background(), &physicalgood.UpdateRequest{
			ID:              goodID,
			Name:            &newName,
			LongDescription: &newLongDescription,
			Tags:            newTags,
			Price:           &newPrice,
			Amount:          &newAmount,
		})

		// Assert
		assert.NoError(t, err)

		goodUpdatesFromResp, ok := updates["physical_good"].(map[string]any)
		assert.True(t, ok)

		if name, ok := goodUpdatesFromResp["name"].(string); !ok || name != newName {
			t.Errorf("physicalGood.Name in response = %v, want %s", goodUpdatesFromResp["name"], newName)
		}
		if longDesc, ok := goodUpdatesFromResp["long_description"].(string); !ok || longDesc != newLongDescription {
			t.Errorf("physicalGood.LongDescription in response = %v, want %s", goodUpdatesFromResp["long_description"], newLongDescription)
		}
		if amount, ok := goodUpdatesFromResp["amount"].(int); !ok || amount != newAmount {
			t.Errorf("physicalGood.Amount in response = %v, want %d", goodUpdatesFromResp["amount"], newAmount)
		}
		if tags, ok := goodUpdatesFromResp["tags"].([]string); !ok || !reflect.DeepEqual(tags, newTags) {
			t.Errorf("physicalGood.Tags in response = %v, want %v", goodUpdatesFromResp["tags"], newTags)
		}

		productUpdatesFromResp, ok := updates["product"].(map[string]any)
		assert.True(t, ok)
		if price, ok := productUpdatesFromResp["price"].(float32); !ok || price != newPrice {
			t.Errorf("product.Price in response = %v, want %f", goodUpdatesFromResp["price"], newPrice)
		}

		if name, ok := goodUpdates["name"].(string); !ok || name != newName {
			t.Errorf("physicalGood.Name passed to repo = %s, want %s", name, newName)
		}
		if longDesc, ok := goodUpdates["long_description"].(string); !ok || longDesc != newLongDescription {
			t.Errorf("physicalGood.LongDescription passed to repo = %s, want %s", longDesc, newLongDescription)
		}
		if amount, ok := goodUpdates["amount"].(int); !ok || amount != newAmount {
			t.Errorf("physicalGood.Amount passed to repo = %d, want %d", amount, newAmount)
		}
		if tags, ok := goodUpdates["tags"].([]string); !ok || !reflect.DeepEqual(tags, newTags) {
			t.Errorf("physicalGood.Tags passed to repo = %v, want %v", tags, newTags)
		}
		if price, ok := productUpdates["price"].(float32); !ok || price != newPrice {
			t.Errorf("product.Price passed to repo = %f, want %f", price, newPrice)
		}
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		// Act
		invalidName := "1invalidname"
		invalidAmount := -21
		_, err := testService.Update(context.Background(), &physicalgood.UpdateRequest{
			ID:              goodID,
			Name:            &invalidName,
			LongDescription: &newLongDescription,
			Tags:            newTags,
			Price:           &newPrice,
			Amount:          &invalidAmount,
		})

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().Get(gomock.Any(), goodID).Return(nil, gorm.ErrRecordNotFound).AnyTimes()

		// Act
		_, err := testService.Update(context.Background(), &physicalgood.UpdateRequest{
			ID:              goodID,
			Name:            &newName,
			LongDescription: &newLongDescription,
			Tags:            newTags,
			Price:           &newPrice,
			Amount:          &newAmount,
		})

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().Get(gomock.Any(), goodID).Return(mockPhysicalGood, nil)
		mockTxProductRepo.EXPECT().SelectByDetailsID(gomock.Any(), goodID, gomock.Any()).Return(mockProduct, nil)
		dbErr := errors.New("database error")
		mockTxPhysicalGoodRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, err := testService.Update(context.Background(), &physicalgood.UpdateRequest{
			ID:              goodID,
			Name:            &newName,
			LongDescription: &newLongDescription,
			Tags:            newTags,
			Price:           &newPrice,
			Amount:          &newAmount,
		})

		// Assert
		assert.Error(t, err)
	})
}

func TestService_AddImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	goodID := uuid.New().String()

	addReq := &imagemodel.AddRequest{
		URL:            "https://google.com",
		SecureURL:      "https://google.com",
		PublicID:       "public/id",
		MediaServiceID: uuid.New().String(),
		OwnerID:        goodID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockImgeSvc.EXPECT().AddImage(gomock.Any(), addReq, gomock.Any()).Return(nil)

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
		mockImgeSvc.EXPECT().AddImage(gomock.Any(), invalidReq, gomock.Any()).Return(imageservice.ErrInvalidArgument)

		// Act
		err := testService.AddImage(context.Background(), invalidReq)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrInvalidArgument)
	})

	t.Run("physical good not found", func(t *testing.T) {
		// Arrange
		mockImgeSvc.EXPECT().AddImage(gomock.Any(), addReq, gomock.Any()).Return(imageservice.ErrOwnerNotFound)

		// Act
		err := testService.AddImage(context.Background(), addReq)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrOwnerNotFound)
	})

	t.Run("image limit", func(t *testing.T) {
		// Arrange
		mockImgeSvc.EXPECT().AddImage(gomock.Any(), addReq, gomock.Any()).Return(imageservice.ErrImageLimitExceeded)

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

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	imgID := uuid.New().String()
	goodID := uuid.New().String()

	deleteReq := &imagemodel.DeleteRequest{
		MediaServiceID: imgID,
		OwnerID:        goodID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockImgeSvc.EXPECT().DeleteImage(gomock.Any(), deleteReq, gomock.Any()).Return(nil)

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
		mockImgeSvc.EXPECT().DeleteImage(gomock.Any(), invalidReq, gomock.Any()).Return(imageservice.ErrInvalidArgument)

		// Act
		err := testService.DeleteImage(context.Background(), invalidReq)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrInvalidArgument)
	})

	t.Run("physical good not found", func(t *testing.T) {
		// Arrange
		mockImgeSvc.EXPECT().DeleteImage(gomock.Any(), deleteReq, gomock.Any()).Return(imageservice.ErrOwnerNotFound)

		// Act
		err := testService.DeleteImage(context.Background(), deleteReq)

		// Assert
		assert.Error(t, err)
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrOwnerNotFound)
	})

	t.Run("associations not found", func(t *testing.T) {
		// Arrange
		mockImgeSvc.EXPECT().DeleteImage(gomock.Any(), deleteReq, gomock.Any()).Return(imageservice.ErrAssociationsNotFound)

		// Act
		err := testService.DeleteImage(context.Background(), deleteReq)

		// Assert
		assert.Error(t, err)
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrAssociationsNotFound)
	})
}

func TestService_AddImageBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImageSvc)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	goodID_1 := uuid.New().String()
	goodID_2 := uuid.New().String()

	addBatchRequest := &imagemodel.AddBatchRequest{
		URL:            "some-url",
		SecureURL:      "some-secure-url",
		PublicID:       "some-public-id",
		MediaServiceID: uuid.New().String(),
		OwnerIDs:       []string{goodID_1, goodID_2},
	}

	t.Run("success", func(t *testing.T) {
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)

		mockImageSvc.EXPECT().AddImageBatch(gomock.Any(), addBatchRequest, gomock.Any()).Return(2, nil)

		// Act
		_, err := testService.AddImageBatch(context.Background(), addBatchRequest)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)

		mockImageSvc.EXPECT().AddImageBatch(gomock.Any(), addBatchRequest, gomock.Any()).Return(0, imageservice.ErrInvalidArgument)

		// Act
		_, err := testService.AddImageBatch(context.Background(), addBatchRequest)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrInvalidArgument)
	})

	t.Run("physical goods not found", func(t *testing.T) {
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)

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

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImageSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImageSvc)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	goodID_1 := uuid.New().String()
	goodID_2 := uuid.New().String()

	deleteBatchRequest := &imagemodel.DeleteBatchRequst{
		MediaServiceID: uuid.New().String(),
		OwnerIDs:       []string{goodID_1, goodID_2},
	}

	t.Run("success", func(t *testing.T) {
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)

		mockImageSvc.EXPECT().DeleteImageBatch(gomock.Any(), deleteBatchRequest, gomock.Any()).Return(2, nil)

		// Act
		_, err := testService.DeleteImageBatch(context.Background(), deleteBatchRequest)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)

		mockImageSvc.EXPECT().DeleteImageBatch(gomock.Any(), deleteBatchRequest, gomock.Any()).Return(0, imageservice.ErrInvalidArgument)

		// Act
		_, err := testService.DeleteImageBatch(context.Background(), deleteBatchRequest)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrInvalidArgument)
	})

	t.Run("physical goods not found", func(t *testing.T) {
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)

		mockImageSvc.EXPECT().DeleteImageBatch(gomock.Any(), deleteBatchRequest, gomock.Any()).Return(0, imageservice.ErrOwnersNotFound)

		// Act
		_, err := testService.DeleteImageBatch(context.Background(), deleteBatchRequest)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrOwnersNotFound)
	})

	t.Run("associations not found", func(t *testing.T) {
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)

		mockImageSvc.EXPECT().DeleteImageBatch(gomock.Any(), deleteBatchRequest, gomock.Any()).Return(0, imageservice.ErrAssociationsNotFound)

		// Act
		_, err := testService.DeleteImageBatch(context.Background(), deleteBatchRequest)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, imageservice.ErrAssociationsNotFound)
	})
}

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	goodID := "0d9828df-c57b-4629-9729-8c9641598e17"

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().GetWithUnpublished(gomock.Any(), goodID).Return(&physicalgood.PhysicalGood{}, nil)
		mockTxPhysicalGoodRepo.EXPECT().SetInStock(gomock.Any(), goodID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), goodID, false).Return(int64(1), nil)
		mockTxPhysicalGoodRepo.EXPECT().Delete(gomock.Any(), goodID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().DeleteByDetailsID(gomock.Any(), goodID).Return(int64(1), nil)

		// Act
		err := testService.Delete(context.Background(), goodID)

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
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().GetWithUnpublished(gomock.Any(), goodID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		err := testService.Delete(context.Background(), goodID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().GetWithUnpublished(gomock.Any(), goodID).Return(&physicalgood.PhysicalGood{}, nil)
		mockTxPhysicalGoodRepo.EXPECT().SetInStock(gomock.Any(), goodID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), goodID, false).Return(int64(1), nil)
		dbErr := errors.New("database error")
		mockTxPhysicalGoodRepo.EXPECT().Delete(gomock.Any(), goodID).Return(int64(0), dbErr)

		// Act
		err := testService.Delete(context.Background(), goodID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_DeletePermanent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	goodID := "0d9828df-c57b-4629-9729-8c9641598e17"

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().DeletePermanent(gomock.Any(), goodID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().DeletePermanentByDetailsID(gomock.Any(), goodID).Return(int64(1), nil)

		// Act
		err := testService.DeletePermanent(context.Background(), goodID)

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
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().DeletePermanent(gomock.Any(), goodID).Return(int64(0), nil)

		// Act
		err := testService.DeletePermanent(context.Background(), goodID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxPhysicalGoodRepo.EXPECT().DeletePermanent(gomock.Any(), goodID).Return(int64(0), dbErr)

		// Act
		err := testService.DeletePermanent(context.Background(), goodID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Restore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockImgeSvc := imageservicemock.NewMockService(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo, mockImgeSvc)

	goodID := "0d9828df-c57b-4629-9729-8c9641598e17"

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().Restore(gomock.Any(), goodID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().RestoreByDetailsID(gomock.Any(), goodID).Return(int64(1), nil)

		// Act
		err := testService.Restore(context.Background(), goodID)

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
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().Restore(gomock.Any(), goodID).Return(int64(0), nil)

		// Act
		err := testService.Restore(context.Background(), goodID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxPhysicalGoodRepo.EXPECT().Restore(gomock.Any(), goodID).Return(int64(0), dbErr)

		// Act
		err := testService.Restore(context.Background(), goodID)

		// Assert
		assert.Error(t, err)
	})
}
