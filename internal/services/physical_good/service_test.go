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
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
	physicalgood "github.com/mikhail5545/product-service-go/internal/models/physical_good"
	"github.com/mikhail5545/product-service-go/internal/models/product"
	physicalgoodmock "github.com/mikhail5545/product-service-go/internal/test/database/physical_good_mock"
	productmock "github.com/mikhail5545/product-service-go/internal/test/database/product_mock"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
		PhysicalGood: *mockPhysicalGood,
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
		if err != nil {
			t.Errorf("Get() error = %v, wantErr %v", err, nil)
			return
		}
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
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Get() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockPhysicalGoodRepo.EXPECT().Get(gomock.Any(), physicalGoodID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.Get(context.Background(), physicalGoodID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusNotFound {
			t.Errorf("Get() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().Get(gomock.Any(), physicalGoodID).Return(nil, dbErr)

		// Act
		_, err := testService.Get(context.Background(), physicalGoodID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("Get() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}

func TestService_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
		PhysicalGood: *mockPhysicalGood,
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
		if err != nil {
			t.Errorf("GetWithDeleted() error = %v, wantErr %v", err, nil)
			return
		}
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
		if err == nil {
			t.Errorf("GetWithDeleted() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("GetWithDeleted() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("GetWithDeleted() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockPhysicalGoodRepo.EXPECT().GetWithDeleted(gomock.Any(), physicalGoodID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), physicalGoodID)

		// Assert
		if err == nil {
			t.Errorf("GetWithDeleted() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("GetWithDeleted() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusNotFound {
			t.Errorf("GetWithDeleted() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().GetWithDeleted(gomock.Any(), physicalGoodID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), physicalGoodID)

		// Assert
		if err == nil {
			t.Errorf("GetWithDeleted() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("GetWithDeleted() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("GetWithDeleted() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}

func TestService_GetWithUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
		PhysicalGood: *mockPhysicalGood,
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
		if err != nil {
			t.Errorf("GetWithUnpublished() error = %v, wantErr %v", err, nil)
			return
		}
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
		if err == nil {
			t.Errorf("GetWithUnpublished() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("GetWithUnpublished() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("GetWithUnpublished() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockPhysicalGoodRepo.EXPECT().GetWithUnpublished(gomock.Any(), physicalGoodID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), physicalGoodID)

		// Assert
		if err == nil {
			t.Errorf("GetWithUnpublished() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("GetWithUnpublished() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusNotFound {
			t.Errorf("GetWithUnpublished() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().GetWithUnpublished(gomock.Any(), physicalGoodID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), physicalGoodID)

		// Assert
		if err == nil {
			t.Errorf("GetWithUnpublished() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("GetWithUnpublished() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("GetWithUnpublished() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}

func TestService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
			PhysicalGood: mockPhysicalGoods[0],
			Price:        mockProducts[0].Price,
			ProductID:    mockProducts[0].ID,
		},
		{
			PhysicalGood: mockPhysicalGoods[1],
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
		if err != nil {
			t.Errorf("List() error = %v, wantErr %v", err, nil)
			return
		}
		if total != 2 {
			t.Errorf("List() got total %d, want %d", total, 2)
		}
		if len(details) != len(expectedDetails) {
			t.Errorf("List() got %v, want %v", details, expectedDetails)
		}
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
		if err == nil {
			t.Errorf("List() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("List() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("List() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})

	t.Run("db error on list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().List(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.List(context.Background(), limit, offset)

		// Assert
		if err == nil {
			t.Errorf("List() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("List() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("List() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}

func TestService_ListDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
			PhysicalGood: mockPhysicalGoods[0],
			Price:        mockProducts[0].Price,
			ProductID:    mockProducts[0].ID,
		},
		{
			PhysicalGood: mockPhysicalGoods[1],
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
		if err != nil {
			t.Errorf("ListDeleted() error = %v, wantErr %v", err, nil)
			return
		}
		if total != 2 {
			t.Errorf("ListDeleted() got total %d, want %d", total, 2)
		}
		if len(details) != len(expectedDetails) {
			t.Errorf("ListDeleted() got %v, want %v", details, expectedDetails)
		}
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
		if err == nil {
			t.Errorf("ListDeleted() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("ListDeleted() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("ListDeleted() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})

	t.Run("db error on list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		if err == nil {
			t.Errorf("ListDeleted() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("ListDeleted() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("ListDeleted() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}

func TestService_ListUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
			PhysicalGood: mockPhysicalGoods[0],
			Price:        mockProducts[0].Price,
			ProductID:    mockProducts[0].ID,
		},
		{
			PhysicalGood: mockPhysicalGoods[1],
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
		if err != nil {
			t.Errorf("ListUnpublished() error = %v, wantErr %v", err, nil)
			return
		}
		if total != 2 {
			t.Errorf("ListUnpublished() got total %d, want %d", total, 2)
		}
		if len(details) != len(expectedDetails) {
			t.Errorf("ListUnpublished() got %v, want %v", details, expectedDetails)
		}
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
		if err == nil {
			t.Errorf("ListUnpublished() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("ListUnpublished() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("ListUnpublished() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})

	t.Run("db error on list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockPhysicalGoodRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		if err == nil {
			t.Errorf("ListUnpublished() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("ListUnpublished() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("ListUnpublished() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}

func TestSesrvice_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
		if err != nil {
			t.Errorf("Create() error = %v, wantErr %v", err, nil)
			return
		}
		if _, err := uuid.Parse(createdPhysicalGood.ID); err != nil {
			t.Errorf("Expected physicalGood.ID to be a valid UUID, got %s", createdPhysicalGood.ID)
		}
		if createdPhysicalGood.Name != createReq.Name {
			t.Errorf("physicalGood.Name = %s, want %s", createdPhysicalGood.Name, createReq.Name)
		}
		if createdPhysicalGood.ShortDescription != createReq.ShortDescription {
			t.Errorf("physicalGood.ShortDescription = %s, want %s", createdPhysicalGood.ShortDescription, createReq.ShortDescription)
		}
		if createdPhysicalGood.Amount != createReq.Amount {
			t.Errorf("physicalGood.Amount = %d, want %d", createdPhysicalGood.Amount, createReq.Amount)
		}
		if createdPhysicalGood.ShippingRequired != createReq.ShippingRequired {
			t.Errorf("physicalGood.Amount = %v, want %v", createdPhysicalGood.ShippingRequired, createReq.ShippingRequired)
		}
		if createdPhysicalGood.InStock {
			t.Error("expected physical good to be unpublished")
		}
		if _, err := uuid.Parse(createdProduct.ID); err != nil {
			t.Errorf("Expected product.ID to be a valid UUID, got %s", createdProduct.ID)
		}
		if createdProduct.Price != createReq.Price {
			t.Errorf("product.Price = %f, want %f", createdProduct.Price, createReq.Price)
		}
		if createdProduct.DetailsID != createdPhysicalGood.ID {
			t.Errorf("product.DetailsID = %s, want %s", createdProduct.DetailsID, createdPhysicalGood.ID)
		}
		if createdProduct.DetailsType != "physical_good" {
			t.Errorf("product.DetailsType = %s, want %s", createdProduct.DetailsType, "physical_good")
		}
		if createdProduct.InStock {
			t.Error("expected product to be unpublished")
		}
		if createdProduct.ID != resp.ProductID {
			t.Errorf("response ProductID = %s, want %s", resp.ID, createdProduct.ID)
		}
		if createdPhysicalGood.ID != resp.ID {
			t.Errorf("response ID = %s, want %s", resp.ID, createdPhysicalGood.ID)
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
		_, err := testService.Create(context.Background(), &physicalgood.CreateRequest{
			Name:             "3invalidname",
			ShortDescription: "Short description",
			Amount:           -44,
			Price:            55.3,
			ShippingRequired: false,
		})

		// Assert
		if err == nil {
			t.Errorf("Create() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Create() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Create() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxPhysicalGoodRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)

		// Act
		_, err := testService.Create(context.Background(), &createReq)

		// Assert
		if err == nil {
			t.Errorf("Create() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Create() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("Create() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}

func TestService_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
		if err != nil {
			t.Errorf("Publish() error = %v, wantErr %v", err, nil)
			return
		}
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
		if err == nil {
			t.Errorf("Publish() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Publish() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Publish() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockPhysicalGoodRepo.EXPECT().SetInStock(gomock.Any(), goodID, true).Return(int64(0), nil).AnyTimes()

		// Act
		err := testService.Publish(context.Background(), goodID)

		// Assert
		if err == nil {
			t.Errorf("Publish() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Publish() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusNotFound {
			t.Errorf("Publish() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
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
		if err == nil {
			t.Errorf("Publish() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Publish() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("Publish() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}

func TestService_Unpublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
		if err != nil {
			t.Errorf("Unpublish() error = %v, wantErr %v", err, nil)
			return
		}
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
		if err == nil {
			t.Errorf("Unpublish() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Unpublish() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Unpublish() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
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
		if err == nil {
			t.Errorf("Unpublish() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Unpublish() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusNotFound {
			t.Errorf("Unpublish() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
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
		if err == nil {
			t.Errorf("Publish() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Unpublish() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("Unpublish() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
		if err != nil {
			t.Errorf("Update() error = %v, wantErr %v", err, nil)
			return
		}

		goodUpdatesFromResp, ok := updates["physical_good"].(map[string]any)
		if !ok {
			t.Fatalf("response does not contain 'physical_good' updates")
		}
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
		if !ok {
			t.Fatalf("response does not contain 'product' updates")
		}
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
		if err == nil {
			t.Errorf("Update() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Update() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Update() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockPhysicalGoodRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPhysicalGoodRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPhysicalGoodRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxPhysicalGoodRepo.EXPECT().Get(gomock.Any(), goodID).Return(nil, gorm.ErrRecordNotFound)

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
		if err == nil {
			t.Errorf("Update() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Update() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusNotFound {
			t.Errorf("Update() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
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
		if err == nil {
			t.Errorf("Update() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Update() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("Update() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
		if err != nil {
			t.Errorf("Delete() error = %v, wantErr %v", err, nil)
			return
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		err := testService.Delete(context.Background(), invalidID)

		// Assert
		if err == nil {
			t.Errorf("Delete() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Delete() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Delete() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
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
		if err == nil {
			t.Errorf("Delete() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Delete() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusNotFound {
			t.Errorf("Delete() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
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
		if err == nil {
			t.Errorf("Delete() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Delete() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("Delete() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}

func TestService_DeletePermanent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
		if err != nil {
			t.Errorf("DeletePermanent() error = %v, wantErr %v", err, nil)
			return
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		err := testService.DeletePermanent(context.Background(), invalidID)

		// Assert
		if err == nil {
			t.Errorf("DeletePermanent() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("DeletePermanent() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("DeletePermanent() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
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
		if err == nil {
			t.Errorf("DeletePermanent() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("DeletePermanent() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusNotFound {
			t.Errorf("DeletePermanent() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
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
		if err == nil {
			t.Errorf("DeletePermanent() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("DeletePermanent() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("DeletePermanent() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}

func TestService_Restore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPhysicalGoodRepo := physicalgoodmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockPhysicalGoodRepo, mockProductRepo)

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
		if err != nil {
			t.Errorf("Restore() error = %v, wantErr %v", err, nil)
			return
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		err := testService.Restore(context.Background(), invalidID)

		// Assert
		if err == nil {
			t.Errorf("Restore() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Restore() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Restore() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
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
		if err == nil {
			t.Errorf("Restore() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Restore() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusNotFound {
			t.Errorf("Restore() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
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
		if err == nil {
			t.Errorf("Restore() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Restore() expected a custom error type, got %v", err)
			return
		}
		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("Restore() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}
