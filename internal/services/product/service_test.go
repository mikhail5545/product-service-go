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

// Package product provides service-layer business logic for products.
package product

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
	"github.com/mikhail5545/product-service-go/internal/models/product"
	productmock "github.com/mikhail5545/product-service-go/internal/test/database/product_mock"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockProductRepo)

	productID := uuid.New().String()
	mockProduct := &product.Product{
		ID:          productID,
		DetailsID:   uuid.New().String(),
		DetailsType: "course",
		InStock:     false,
		Price:       33.33,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockProductRepo.EXPECT().Get(gomock.Any(), productID).Return(mockProduct, nil)

		// Act
		product, err := testService.Get(context.Background(), productID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(product, mockProduct) {
			t.Errorf("Get() expected %w, got %w", mockProduct, product)
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
		mockProductRepo.EXPECT().Get(gomock.Any(), productID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.Get(context.Background(), productID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockProductRepo.EXPECT().Get(gomock.Any(), productID).Return(nil, dbErr)

		// Act
		_, err := testService.Get(context.Background(), productID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockProductRepo)

	productID := uuid.New().String()
	mockProduct := &product.Product{
		ID:          productID,
		DetailsID:   uuid.New().String(),
		DetailsType: "course",
		InStock:     false,
		Price:       33.33,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockProductRepo.EXPECT().GetWithDeleted(gomock.Any(), productID).Return(mockProduct, nil)

		// Act
		product, err := testService.GetWithDeleted(context.Background(), productID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(product, mockProduct) {
			t.Errorf("GetWithDeleted() expected %w, got %w", mockProduct, product)
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
		mockProductRepo.EXPECT().GetWithDeleted(gomock.Any(), productID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), productID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockProductRepo.EXPECT().GetWithDeleted(gomock.Any(), productID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), productID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetWithUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockProductRepo)

	productID := uuid.New().String()
	mockProduct := &product.Product{
		ID:          productID,
		DetailsID:   uuid.New().String(),
		DetailsType: "course",
		InStock:     false,
		Price:       33.33,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockProductRepo.EXPECT().GetWithUnpublished(gomock.Any(), productID).Return(mockProduct, nil)

		// Act
		product, err := testService.GetWithUnpublished(context.Background(), productID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(product, mockProduct) {
			t.Errorf("GetWithUnpublished() expected %w, got %w", mockProduct, product)
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
		mockProductRepo.EXPECT().GetWithUnpublished(gomock.Any(), productID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), productID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockProductRepo.EXPECT().GetWithUnpublished(gomock.Any(), productID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), productID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetByDetailsID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockProductRepo)

	productID := uuid.New().String()
	detailsID := uuid.New().String()

	mockProduct := &product.Product{
		ID:          productID,
		DetailsID:   detailsID,
		DetailsType: "course",
		InStock:     false,
		Price:       33.33,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockProductRepo.EXPECT().GetByDetailsID(gomock.Any(), detailsID).Return(mockProduct, nil)

		// Act
		product, err := testService.GetByDetailsID(context.Background(), detailsID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(product, mockProduct) {
			t.Errorf("GetByDetailsID() expected %w, got %w", mockProduct, product)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		_, err := testService.GetByDetailsID(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockProductRepo.EXPECT().GetByDetailsID(gomock.Any(), detailsID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetByDetailsID(context.Background(), detailsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockProductRepo.EXPECT().GetByDetailsID(gomock.Any(), detailsID).Return(nil, dbErr)

		// Act
		_, err := testService.GetByDetailsID(context.Background(), detailsID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetWithDeletedByDetailsID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockProductRepo)

	productID := uuid.New().String()
	detailsID := uuid.New().String()

	mockProduct := &product.Product{
		ID:          productID,
		DetailsID:   detailsID,
		DetailsType: "course",
		InStock:     false,
		Price:       33.33,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockProductRepo.EXPECT().GetWithDeletedByDetailsID(gomock.Any(), detailsID).Return(mockProduct, nil)

		// Act
		product, err := testService.GetWithDeletedByDetailsID(context.Background(), detailsID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(product, mockProduct) {
			t.Errorf("GetWithDeletedByDetailsID() expected %w, got %w", mockProduct, product)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		_, err := testService.GetWithDeletedByDetailsID(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockProductRepo.EXPECT().GetWithDeletedByDetailsID(gomock.Any(), detailsID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithDeletedByDetailsID(context.Background(), detailsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockProductRepo.EXPECT().GetWithDeletedByDetailsID(gomock.Any(), detailsID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithDeletedByDetailsID(context.Background(), detailsID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetWithUnpublishedByDetailsID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockProductRepo)

	productID := uuid.New().String()
	detailsID := uuid.New().String()

	mockProduct := &product.Product{
		ID:          productID,
		DetailsID:   detailsID,
		DetailsType: "course",
		InStock:     false,
		Price:       33.33,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockProductRepo.EXPECT().GetWithUnpublishedByDetailsID(gomock.Any(), detailsID).Return(mockProduct, nil)

		// Act
		product, err := testService.GetWithUnpublishedByDetailsID(context.Background(), detailsID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(product, mockProduct) {
			t.Errorf("GetWithUnpublishedByDetailsID() expected %w, got %w", mockProduct, product)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"

		// Act
		_, err := testService.GetWithUnpublishedByDetailsID(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockProductRepo.EXPECT().GetWithUnpublishedByDetailsID(gomock.Any(), detailsID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithUnpublishedByDetailsID(context.Background(), detailsID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockProductRepo.EXPECT().GetWithUnpublishedByDetailsID(gomock.Any(), detailsID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithUnpublishedByDetailsID(context.Background(), detailsID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockProductRepo)

	productID_1 := uuid.New().String()
	productID_2 := uuid.New().String()

	mockProducts := []product.Product{
		{
			ID:          productID_1,
			DetailsID:   uuid.New().String(),
			DetailsType: "course",
			InStock:     false,
			Price:       33.33,
		},
		{
			ID:          productID_2,
			DetailsID:   uuid.New().String(),
			DetailsType: "training_session",
			InStock:     false,
			Price:       32.22,
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockProductRepo.EXPECT().List(gomock.Any(), limit, offset).Return(mockProducts, nil)
		mockProductRepo.EXPECT().Count(gomock.Any()).Return(int64(2), nil)

		// Act
		products, total, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, products, 2)
	})

	t.Run("success with empty list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockProductRepo.EXPECT().List(gomock.Any(), limit, offset).Return([]product.Product{}, nil)
		mockProductRepo.EXPECT().Count(gomock.Any()).Return(int64(0), nil)

		// Act
		products, total, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, products, 0)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockProductRepo.EXPECT().List(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_ListDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockProductRepo)

	productID_1 := uuid.New().String()
	productID_2 := uuid.New().String()

	mockProducts := []product.Product{
		{
			ID:          productID_1,
			DetailsID:   uuid.New().String(),
			DetailsType: "course",
			InStock:     false,
			Price:       33.33,
		},
		{
			ID:          productID_2,
			DetailsID:   uuid.New().String(),
			DetailsType: "training_session",
			InStock:     false,
			Price:       32.22,
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockProductRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(mockProducts, nil)
		mockProductRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(2), nil)

		// Act
		products, total, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, products, 2)
	})

	t.Run("success with empty list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockProductRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return([]product.Product{}, nil)
		mockProductRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(0), nil)

		// Act
		products, total, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, products, 0)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockProductRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_ListUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockProductRepo)

	productID_1 := uuid.New().String()
	productID_2 := uuid.New().String()

	mockProducts := []product.Product{
		{
			ID:          productID_1,
			DetailsID:   uuid.New().String(),
			DetailsType: "course",
			InStock:     false,
			Price:       33.33,
		},
		{
			ID:          productID_2,
			DetailsID:   uuid.New().String(),
			DetailsType: "training_session",
			InStock:     false,
			Price:       32.22,
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockProductRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(mockProducts, nil)
		mockProductRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(2), nil)

		// Act
		products, total, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, products, 2)
	})

	t.Run("success with empty list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockProductRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return([]product.Product{}, nil)
		mockProductRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(0), nil)

		// Act
		products, total, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, products, 0)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockProductRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_ListByDetailsType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockProductRepo)

	productID_1 := uuid.New().String()
	productID_2 := uuid.New().String()
	detailsType := "course"

	mockProducts := []product.Product{
		{
			ID:          productID_1,
			DetailsID:   uuid.New().String(),
			DetailsType: detailsType,
			InStock:     false,
			Price:       33.33,
		},
		{
			ID:          productID_2,
			DetailsID:   uuid.New().String(),
			DetailsType: detailsType,
			InStock:     false,
			Price:       32.22,
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockProductRepo.EXPECT().ListByDetailsType(gomock.Any(), detailsType, limit, offset).Return(mockProducts, nil)
		mockProductRepo.EXPECT().CountByDetailsType(gomock.Any(), detailsType).Return(int64(2), nil)

		// Act
		products, total, err := testService.ListByDetailsType(context.Background(), detailsType, limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, products, 2)
	})

	t.Run("success with empty list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockProductRepo.EXPECT().ListByDetailsType(gomock.Any(), detailsType, limit, offset).Return([]product.Product{}, nil)
		mockProductRepo.EXPECT().CountByDetailsType(gomock.Any(), detailsType).Return(int64(0), nil)

		// Act
		products, total, err := testService.ListByDetailsType(context.Background(), detailsType, limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, products, 0)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockProductRepo.EXPECT().ListByDetailsType(gomock.Any(), detailsType, limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListByDetailsType(context.Background(), detailsType, limit, offset)

		// Assert
		assert.Error(t, err)
	})
}
