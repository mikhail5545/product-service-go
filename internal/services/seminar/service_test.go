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
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
	"github.com/mikhail5545/product-service-go/internal/models/product"
	"github.com/mikhail5545/product-service-go/internal/models/seminar"
	productmock "github.com/mikhail5545/product-service-go/internal/test/database/product_mock"
	seminarmock "github.com/mikhail5545/product-service-go/internal/test/database/seminar_mock"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	seminarID := "c6248da5-a2eb-4abd-be56-a19715104c00"
	rproductID := "866561c2-a65a-4159-a5d8-a0ae5401e0c1"
	eproductID := "7baa5ff9-a864-4144-b42c-8ce6bd56ac25"
	lproductID := "38fcb2f8-d377-4b08-9eb9-8de9a89d4528"
	esproductID := "0cb3a9a5-9dd0-4ca9-b528-275071e3eb98"
	lsproductID := "14212b87-ca38-41d5-bba2-2a273fe60977"

	layout := "2006-Jan-02"

	beforeNow, _ := time.Parse(layout, "2024-Aug-03")
	afterNow, _ := time.Parse(layout, "2099-Dec-03")

	mockSeminar := &seminar.Seminar{
		ID:                      seminarID,
		Name:                    "Seminar name",
		ShortDescription:        "Seminar short description",
		ReservationProductID:    &rproductID,
		EarlyProductID:          &eproductID,
		LateProductID:           &lproductID,
		EarlySurchargeProductID: &esproductID,
		LateSurchargeProductID:  &lsproductID,
	}

	mockProducts := []product.Product{
		{
			ID:          rproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       34.44,
		},
		{
			ID:          eproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       44.44,
		},
		{
			ID:          lproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       366.44,
		},
		{
			ID:          esproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       3466.44,
		},
		{
			ID:          lsproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       346.44,
		},
	}

	t.Run("success with late_payment_date in future", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = afterNow
		mockSeminarRepo.EXPECT().Get(gomock.Any(), seminarID).Return(mockSeminar, nil)
		mockProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)

		// Act
		expectedDetails := &seminar.SeminarDetails{
			Seminar:                        mockSeminar,
			ReservationPrice:               mockProducts[0].Price,
			EarlyPrice:                     mockProducts[1].Price,
			LatePrice:                      mockProducts[2].Price,
			EarlySurchargePrice:            mockProducts[3].Price,
			LateSurchargePrice:             mockProducts[4].Price,
			CurrentPrice:                   mockProducts[1].Price,
			CurrentPriceProductID:          eproductID,
			CurrentSurchargePrice:          mockProducts[3].Price,
			CurrentSurchargePriceProductID: esproductID,
		}

		// Act
		details, err := testService.Get(context.Background(), seminarID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("Get() got %v, want %v", details, expectedDetails)
		}
	})

	t.Run("success with late_payment_date in past", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = beforeNow
		mockSeminarRepo.EXPECT().Get(gomock.Any(), seminarID).Return(mockSeminar, nil)
		mockProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)

		// Act
		expectedDetails := &seminar.SeminarDetails{
			Seminar:                        mockSeminar,
			ReservationPrice:               mockProducts[0].Price,
			EarlyPrice:                     mockProducts[1].Price,
			LatePrice:                      mockProducts[2].Price,
			EarlySurchargePrice:            mockProducts[3].Price,
			LateSurchargePrice:             mockProducts[4].Price,
			CurrentPrice:                   mockProducts[2].Price,
			CurrentPriceProductID:          lproductID,
			CurrentSurchargePrice:          mockProducts[4].Price,
			CurrentSurchargePriceProductID: lsproductID,
		}

		// Act
		details, err := testService.Get(context.Background(), seminarID)

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
		mockSeminar.LatePaymentDate = afterNow
		mockSeminarRepo.EXPECT().Get(gomock.Any(), seminarID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.Get(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = afterNow
		dbErr := errors.New("database error")
		mockSeminarRepo.EXPECT().Get(gomock.Any(), seminarID).Return(nil, dbErr)

		// Act
		_, err := testService.Get(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("seminar missing product id", func(t *testing.T) {
		// Arrange
		seminarWithMissingID := &seminar.Seminar{
			ID:                      seminarID,
			Name:                    "Seminar name",
			ShortDescription:        "Seminar short description",
			ReservationProductID:    &rproductID,
			EarlyProductID:          &eproductID,
			LateProductID:           nil, // Missing ID
			EarlySurchargeProductID: &esproductID,
			LateSurchargeProductID:  &lsproductID,
		}
		mockSeminarRepo.EXPECT().Get(gomock.Any(), seminarID).Return(seminarWithMissingID, nil)

		// Act
		_, err := testService.Get(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrIncompleteData)
	})

	t.Run("product repo returns error", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = afterNow
		mockSeminarRepo.EXPECT().Get(gomock.Any(), seminarID).Return(mockSeminar, nil)
		dbErr := errors.New("product db error")
		mockProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)

		// Act
		_, err := testService.Get(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("product repo returns incomplete products", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = afterNow
		mockSeminarRepo.EXPECT().Get(gomock.Any(), seminarID).Return(mockSeminar, nil)
		// Return only 4 products instead of 5
		incompleteProducts := mockProducts[:4]
		mockProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(incompleteProducts, nil)

		// Act
		_, err := testService.Get(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrProductsNotFound)
	})
}

func TestService_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	seminarID := "c6248da5-a2eb-4abd-be56-a19715104c00"
	rproductID := "866561c2-a65a-4159-a5d8-a0ae5401e0c1"
	eproductID := "7baa5ff9-a864-4144-b42c-8ce6bd56ac25"
	lproductID := "38fcb2f8-d377-4b08-9eb9-8de9a89d4528"
	esproductID := "0cb3a9a5-9dd0-4ca9-b528-275071e3eb98"
	lsproductID := "14212b87-ca38-41d5-bba2-2a273fe60977"

	layout := "2006-Jan-02"

	beforeNow, _ := time.Parse(layout, "2024-Aug-03")
	afterNow, _ := time.Parse(layout, "2099-Dec-03")

	mockSeminar := &seminar.Seminar{
		ID:                      seminarID,
		Name:                    "Seminar name",
		ShortDescription:        "Seminar short description",
		ReservationProductID:    &rproductID,
		EarlyProductID:          &eproductID,
		LateProductID:           &lproductID,
		EarlySurchargeProductID: &esproductID,
		LateSurchargeProductID:  &lsproductID,
	}

	mockProducts := []product.Product{
		{
			ID:          rproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       34.44,
		},
		{
			ID:          eproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       44.44,
		},
		{
			ID:          lproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       366.44,
		},
		{
			ID:          esproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       3466.44,
		},
		{
			ID:          lsproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       346.44,
		},
	}

	t.Run("success with late_payment_date in future", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = afterNow
		mockSeminarRepo.EXPECT().GetWithDeleted(gomock.Any(), seminarID).Return(mockSeminar, nil)
		mockProductRepo.EXPECT().SelectWithDeletedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)

		// Act
		expectedDetails := &seminar.SeminarDetails{
			Seminar:                        mockSeminar,
			ReservationPrice:               mockProducts[0].Price,
			EarlyPrice:                     mockProducts[1].Price,
			LatePrice:                      mockProducts[2].Price,
			EarlySurchargePrice:            mockProducts[3].Price,
			LateSurchargePrice:             mockProducts[4].Price,
			CurrentPrice:                   mockProducts[1].Price,
			CurrentPriceProductID:          eproductID,
			CurrentSurchargePrice:          mockProducts[3].Price,
			CurrentSurchargePriceProductID: esproductID,
		}

		// Act
		details, err := testService.GetWithDeleted(context.Background(), seminarID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("GetWithDeleted() got %v, want %v", details, expectedDetails)
		}
	})

	t.Run("success with late_payment_date in past", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = beforeNow
		mockSeminarRepo.EXPECT().GetWithDeleted(gomock.Any(), seminarID).Return(mockSeminar, nil)
		mockProductRepo.EXPECT().SelectWithDeletedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)

		// Act
		expectedDetails := &seminar.SeminarDetails{
			Seminar:                        mockSeminar,
			ReservationPrice:               mockProducts[0].Price,
			EarlyPrice:                     mockProducts[1].Price,
			LatePrice:                      mockProducts[2].Price,
			EarlySurchargePrice:            mockProducts[3].Price,
			LateSurchargePrice:             mockProducts[4].Price,
			CurrentPrice:                   mockProducts[2].Price,
			CurrentPriceProductID:          lproductID,
			CurrentSurchargePrice:          mockProducts[4].Price,
			CurrentSurchargePriceProductID: lsproductID,
		}

		// Act
		details, err := testService.GetWithDeleted(context.Background(), seminarID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("GeGetWithDeletedt() got %v, want %v", details, expectedDetails)
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
		mockSeminar.LatePaymentDate = afterNow
		mockSeminarRepo.EXPECT().GetWithDeleted(gomock.Any(), seminarID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = afterNow
		dbErr := errors.New("database error")
		mockSeminarRepo.EXPECT().GetWithDeleted(gomock.Any(), seminarID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("seminar missing product id", func(t *testing.T) {
		// Arrange
		seminarWithMissingID := &seminar.Seminar{
			ID:                      seminarID,
			Name:                    "Seminar name",
			ShortDescription:        "Seminar short description",
			ReservationProductID:    &rproductID,
			EarlyProductID:          &eproductID,
			LateProductID:           nil, // Missing ID
			EarlySurchargeProductID: &esproductID,
			LateSurchargeProductID:  &lsproductID,
		}
		mockSeminarRepo.EXPECT().GetWithDeleted(gomock.Any(), seminarID).Return(seminarWithMissingID, nil)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrIncompleteData)
	})

	t.Run("product repo returns error", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = afterNow
		mockSeminarRepo.EXPECT().GetWithDeleted(gomock.Any(), seminarID).Return(mockSeminar, nil)
		dbErr := errors.New("product db error")
		mockProductRepo.EXPECT().SelectWithDeletedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("product repo returns incomplete products", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = afterNow
		mockSeminarRepo.EXPECT().GetWithDeleted(gomock.Any(), seminarID).Return(mockSeminar, nil)
		// Return only 4 products instead of 5
		incompleteProducts := mockProducts[:4]
		mockProductRepo.EXPECT().SelectWithDeletedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(incompleteProducts, nil)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrProductsNotFound)
	})
}

func TestService_GetWithUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	seminarID := "c6248da5-a2eb-4abd-be56-a19715104c00"
	rproductID := "866561c2-a65a-4159-a5d8-a0ae5401e0c1"
	eproductID := "7baa5ff9-a864-4144-b42c-8ce6bd56ac25"
	lproductID := "38fcb2f8-d377-4b08-9eb9-8de9a89d4528"
	esproductID := "0cb3a9a5-9dd0-4ca9-b528-275071e3eb98"
	lsproductID := "14212b87-ca38-41d5-bba2-2a273fe60977"

	layout := "2006-Jan-02"

	beforeNow, _ := time.Parse(layout, "2024-Aug-03")
	afterNow, _ := time.Parse(layout, "2099-Dec-03")

	mockSeminar := &seminar.Seminar{
		ID:                      seminarID,
		Name:                    "Seminar name",
		ShortDescription:        "Seminar short description",
		ReservationProductID:    &rproductID,
		EarlyProductID:          &eproductID,
		LateProductID:           &lproductID,
		EarlySurchargeProductID: &esproductID,
		LateSurchargeProductID:  &lsproductID,
	}

	mockProducts := []product.Product{
		{
			ID:          rproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       34.44,
		},
		{
			ID:          eproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       44.44,
		},
		{
			ID:          lproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       366.44,
		},
		{
			ID:          esproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       3466.44,
		},
		{
			ID:          lsproductID,
			DetailsID:   seminarID,
			DetailsType: "seminar",
			Price:       346.44,
		},
	}

	t.Run("success with late_payment_date in future", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = afterNow
		mockSeminarRepo.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(mockSeminar, nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)

		// Act
		expectedDetails := &seminar.SeminarDetails{
			Seminar:                        mockSeminar,
			ReservationPrice:               mockProducts[0].Price,
			EarlyPrice:                     mockProducts[1].Price,
			LatePrice:                      mockProducts[2].Price,
			EarlySurchargePrice:            mockProducts[3].Price,
			LateSurchargePrice:             mockProducts[4].Price,
			CurrentPrice:                   mockProducts[1].Price,
			CurrentPriceProductID:          eproductID,
			CurrentSurchargePrice:          mockProducts[3].Price,
			CurrentSurchargePriceProductID: esproductID,
		}

		// Act
		details, err := testService.GetWithUnpublished(context.Background(), seminarID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("GetWithUnpublished() got %v, want %v", details, expectedDetails)
		}
	})

	t.Run("success with late_payment_date in past", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = beforeNow
		mockSeminarRepo.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(mockSeminar, nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)

		// Act
		expectedDetails := &seminar.SeminarDetails{
			Seminar:                        mockSeminar,
			ReservationPrice:               mockProducts[0].Price,
			EarlyPrice:                     mockProducts[1].Price,
			LatePrice:                      mockProducts[2].Price,
			EarlySurchargePrice:            mockProducts[3].Price,
			LateSurchargePrice:             mockProducts[4].Price,
			CurrentPrice:                   mockProducts[2].Price,
			CurrentPriceProductID:          lproductID,
			CurrentSurchargePrice:          mockProducts[4].Price,
			CurrentSurchargePriceProductID: lsproductID,
		}

		// Act
		details, err := testService.GetWithUnpublished(context.Background(), seminarID)

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
		mockSeminar.LatePaymentDate = afterNow
		mockSeminarRepo.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = afterNow
		dbErr := errors.New("database error")
		mockSeminarRepo.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("seminar missing product id", func(t *testing.T) {
		// Arrange
		seminarWithMissingID := &seminar.Seminar{
			ID:                      seminarID,
			Name:                    "Seminar name",
			ShortDescription:        "Seminar short description",
			ReservationProductID:    &rproductID,
			EarlyProductID:          &eproductID,
			LateProductID:           nil, // Missing ID
			EarlySurchargeProductID: &esproductID,
			LateSurchargeProductID:  &lsproductID,
		}
		mockSeminarRepo.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(seminarWithMissingID, nil)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrIncompleteData)
	})

	t.Run("product repo returns error", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = afterNow
		mockSeminarRepo.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(mockSeminar, nil)
		dbErr := errors.New("product db error")
		mockProductRepo.EXPECT().SelectWithUnpublishedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("product repo returns incomplete products", func(t *testing.T) {
		// Arrange
		mockSeminar.LatePaymentDate = afterNow
		mockSeminarRepo.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(mockSeminar, nil)
		// Return only 4 products instead of 5
		incompleteProducts := mockProducts[:4]
		mockProductRepo.EXPECT().SelectWithUnpublishedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(incompleteProducts, nil)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrProductsNotFound)
	})
}

func TestService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	seminarID_1 := uuid.New().String()
	rproductID_1 := uuid.New().String()
	eproductID_1 := uuid.New().String()
	lproductID_1 := uuid.New().String()
	esproductID_1 := uuid.New().String()
	lsproductID_1 := uuid.New().String()

	seminarID_2 := uuid.New().String()
	rproductID_2 := uuid.New().String()
	eproductID_2 := uuid.New().String()
	lproductID_2 := uuid.New().String()
	esproductID_2 := uuid.New().String()
	lsproductID_2 := uuid.New().String()

	layout := "2006-Jan-02"

	beforeNow, _ := time.Parse(layout, "2024-Aug-03")
	afterNow, _ := time.Parse(layout, "2099-Dec-03")

	mockSeminars := []seminar.Seminar{
		{
			ID:                      seminarID_1,
			Name:                    "Seminar 1 name",
			ShortDescription:        "Seminar 1 short description",
			ReservationProductID:    &rproductID_1,
			EarlyProductID:          &eproductID_1,
			LateProductID:           &lproductID_1,
			EarlySurchargeProductID: &esproductID_1,
			LateSurchargeProductID:  &lsproductID_1,
			LatePaymentDate:         beforeNow,
		},
		{
			ID:                      seminarID_2,
			Name:                    "Seminar 2 name",
			ShortDescription:        "Seminar 2 short description",
			ReservationProductID:    &rproductID_2,
			EarlyProductID:          &eproductID_2,
			LateProductID:           &lproductID_2,
			EarlySurchargeProductID: &esproductID_2,
			LateSurchargeProductID:  &lsproductID_2,
			LatePaymentDate:         afterNow,
		},
	}

	mockProducts := []product.Product{
		{
			ID:          rproductID_1,
			Price:       11.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          eproductID_1,
			Price:       12.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          lproductID_1,
			Price:       13.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          esproductID_1,
			Price:       14.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          lsproductID_1,
			Price:       15.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          rproductID_2,
			Price:       16.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
		{
			ID:          eproductID_2,
			Price:       17.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
		{
			ID:          lproductID_2,
			Price:       18.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
		{
			ID:          esproductID_2,
			Price:       19.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
		{
			ID:          lsproductID_2,
			Price:       20.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
	}

	expectedDetails_1 := &seminar.SeminarDetails{
		Seminar:                        &mockSeminars[0],
		ReservationPrice:               mockProducts[0].Price,
		EarlyPrice:                     mockProducts[1].Price,
		LatePrice:                      mockProducts[2].Price,
		EarlySurchargePrice:            mockProducts[3].Price,
		LateSurchargePrice:             mockProducts[4].Price,
		CurrentPrice:                   mockProducts[2].Price,
		CurrentPriceProductID:          lproductID_1,
		CurrentSurchargePrice:          mockProducts[4].Price,
		CurrentSurchargePriceProductID: lsproductID_1,
	}
	expectedDetails_2 := &seminar.SeminarDetails{
		Seminar:                        &mockSeminars[1],
		ReservationPrice:               mockProducts[5].Price,
		EarlyPrice:                     mockProducts[6].Price,
		LatePrice:                      mockProducts[7].Price,
		EarlySurchargePrice:            mockProducts[8].Price,
		LateSurchargePrice:             mockProducts[9].Price,
		CurrentPrice:                   mockProducts[6].Price,
		CurrentPriceProductID:          eproductID_2,
		CurrentSurchargePrice:          mockProducts[8].Price,
		CurrentSurchargePriceProductID: esproductID_2,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockSeminarRepo.EXPECT().List(gomock.Any(), limit, offset).Return(mockSeminars, nil)
		mockProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)
		mockSeminarRepo.EXPECT().Count(gomock.Any()).Return(int64(2), nil)

		// Act
		details, total, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(details[0].Seminar, expectedDetails_1.Seminar) {
			t.Errorf("List() got %v, want %v", details[0].Seminar, expectedDetails_1.Seminar)
		}
		if !reflect.DeepEqual(details[1].Seminar, expectedDetails_2.Seminar) {
			t.Errorf("List() got %v, want %v", details[1].Seminar, expectedDetails_2.Seminar)
		}
		if details[0].CurrentPrice != expectedDetails_1.CurrentPrice {
			t.Errorf("List() got %f, want %f", details[0].CurrentPrice, expectedDetails_1.CurrentPrice)
		}
		if details[1].CurrentPrice != expectedDetails_2.CurrentPrice {
			t.Errorf("List() got %f, want %f", details[1].CurrentPrice, expectedDetails_2.CurrentPrice)
		}
		if total != 2 {
			t.Errorf("List() got total %d, want %d", total, 2)
		}
	})

	t.Run("db error", func(t *testing.T) {
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockSeminarRepo.EXPECT().List(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})

	t.Run("db error on count", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockSeminarRepo.EXPECT().List(gomock.Any(), limit, offset).Return(mockSeminars, nil)
		mockProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)
		dbErr := errors.New("db count error")
		mockSeminarRepo.EXPECT().Count(gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, _, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})

	t.Run("success with one seminar having a missing product id", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		seminarWithMissingID := mockSeminars[0]
		seminarWithMissingID.LateProductID = nil // one ID is missing

		seminarsWithOneBad := []seminar.Seminar{
			seminarWithMissingID,
			mockSeminars[1],
		}

		validProducts := mockProducts[5:]

		mockSeminarRepo.EXPECT().List(gomock.Any(), limit, offset).Return(seminarsWithOneBad, nil)
		mockProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(validProducts, nil)
		mockSeminarRepo.EXPECT().Count(gomock.Any()).Return(int64(2), nil)

		// Act
		details, total, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		if len(details) != 1 {
			t.Fatalf("List() expected 1 seminar detail, but got %d", len(details))
		}
		if details[0].Seminar.ID != seminarID_2 {
			t.Errorf("List() returned wrong seminar, got ID %s, want %s", details[0].Seminar.ID, seminarID_2)
		}
		if total != 2 {
			t.Errorf("List() got total %d, want %d", total, 2)
		}
	})

	t.Run("success with incomplete products from repo", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0

		// Products for the first seminar are missing from the response
		incompleteProducts := mockProducts[5:]

		mockSeminarRepo.EXPECT().List(gomock.Any(), limit, offset).Return(mockSeminars, nil)
		mockProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(incompleteProducts, nil)
		mockSeminarRepo.EXPECT().Count(gomock.Any()).Return(int64(2), nil)

		// Act
		details, total, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		if len(details) != 1 {
			t.Fatalf("List() expected 1 seminar detail, but got %d", len(details))
		}
		if details[0].Seminar.ID != seminarID_2 {
			t.Errorf("List() returned wrong seminar, got ID %s, want %s", details[0].Seminar.ID, seminarID_2)
		}
		if total != 2 {
			t.Errorf("List() got total %d, want %d", total, 2)
		}
	})

	t.Run("success empty list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockSeminarRepo.EXPECT().List(gomock.Any(), limit, offset).Return([]seminar.Seminar{}, nil)
		mockProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]product.Product{}, nil)
		mockSeminarRepo.EXPECT().Count(gomock.Any()).Return(int64(0), nil)

		// Act
		details, total, err := testService.List(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, len(details), 0)
		assert.Equal(t, int64(0), total)
	})
}

func TestService_ListDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	seminarID_1 := uuid.New().String()
	rproductID_1 := uuid.New().String()
	eproductID_1 := uuid.New().String()
	lproductID_1 := uuid.New().String()
	esproductID_1 := uuid.New().String()
	lsproductID_1 := uuid.New().String()

	seminarID_2 := uuid.New().String()
	rproductID_2 := uuid.New().String()
	eproductID_2 := uuid.New().String()
	lproductID_2 := uuid.New().String()
	esproductID_2 := uuid.New().String()
	lsproductID_2 := uuid.New().String()

	layout := "2006-Jan-02"

	beforeNow, _ := time.Parse(layout, "2024-Aug-03")
	afterNow, _ := time.Parse(layout, "2099-Dec-03")

	mockSeminars := []seminar.Seminar{
		{
			ID:                      seminarID_1,
			Name:                    "Seminar 1 name",
			ShortDescription:        "Seminar 1 short description",
			ReservationProductID:    &rproductID_1,
			EarlyProductID:          &eproductID_1,
			LateProductID:           &lproductID_1,
			EarlySurchargeProductID: &esproductID_1,
			LateSurchargeProductID:  &lsproductID_1,
			LatePaymentDate:         beforeNow,
		},
		{
			ID:                      seminarID_2,
			Name:                    "Seminar 2 name",
			ShortDescription:        "Seminar 2 short description",
			ReservationProductID:    &rproductID_2,
			EarlyProductID:          &eproductID_2,
			LateProductID:           &lproductID_2,
			EarlySurchargeProductID: &esproductID_2,
			LateSurchargeProductID:  &lsproductID_2,
			LatePaymentDate:         afterNow,
		},
	}

	mockProducts := []product.Product{
		{
			ID:          rproductID_1,
			Price:       11.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          eproductID_1,
			Price:       12.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          lproductID_1,
			Price:       13.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          esproductID_1,
			Price:       14.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          lsproductID_1,
			Price:       15.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          rproductID_2,
			Price:       16.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
		{
			ID:          eproductID_2,
			Price:       17.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
		{
			ID:          lproductID_2,
			Price:       18.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
		{
			ID:          esproductID_2,
			Price:       19.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
		{
			ID:          lsproductID_2,
			Price:       20.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
	}

	expectedDetails_1 := &seminar.SeminarDetails{
		Seminar:                        &mockSeminars[0],
		ReservationPrice:               mockProducts[0].Price,
		EarlyPrice:                     mockProducts[1].Price,
		LatePrice:                      mockProducts[2].Price,
		EarlySurchargePrice:            mockProducts[3].Price,
		LateSurchargePrice:             mockProducts[4].Price,
		CurrentPrice:                   mockProducts[2].Price,
		CurrentPriceProductID:          lproductID_1,
		CurrentSurchargePrice:          mockProducts[4].Price,
		CurrentSurchargePriceProductID: lsproductID_1,
	}
	expectedDetails_2 := &seminar.SeminarDetails{
		Seminar:                        &mockSeminars[1],
		ReservationPrice:               mockProducts[5].Price,
		EarlyPrice:                     mockProducts[6].Price,
		LatePrice:                      mockProducts[7].Price,
		EarlySurchargePrice:            mockProducts[8].Price,
		LateSurchargePrice:             mockProducts[9].Price,
		CurrentPrice:                   mockProducts[6].Price,
		CurrentPriceProductID:          eproductID_2,
		CurrentSurchargePrice:          mockProducts[8].Price,
		CurrentSurchargePriceProductID: esproductID_2,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockSeminarRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(mockSeminars, nil)
		mockProductRepo.EXPECT().SelectWithDeletedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)
		mockSeminarRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(2), nil)

		// Act
		details, total, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(details[0].Seminar, expectedDetails_1.Seminar) {
			t.Errorf("ListDeleted() got %v, want %v", details[0].Seminar, expectedDetails_1.Seminar)
		}
		if !reflect.DeepEqual(details[1].Seminar, expectedDetails_2.Seminar) {
			t.Errorf("ListDeleted() got %v, want %v", details[1].Seminar, expectedDetails_2.Seminar)
		}
		if details[0].CurrentPrice != expectedDetails_1.CurrentPrice {
			t.Errorf("ListDeleted() got %f, want %f", details[0].CurrentPrice, expectedDetails_1.CurrentPrice)
		}
		if details[1].CurrentPrice != expectedDetails_2.CurrentPrice {
			t.Errorf("ListDeleted() got %f, want %f", details[1].CurrentPrice, expectedDetails_2.CurrentPrice)
		}
		if total != 2 {
			t.Errorf("ListDeleted() got total %d, want %d", total, 2)
		}
	})

	t.Run("db error", func(t *testing.T) {
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockSeminarRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListDeleted(context.Background(), limit, offset)
		// Assert
		assert.Error(t, err)
	})

	t.Run("db error on count", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockSeminarRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(mockSeminars, nil)
		mockProductRepo.EXPECT().SelectWithDeletedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)
		dbErr := errors.New("db count error")
		mockSeminarRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, _, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})

	t.Run("success with one seminar having a missing product id", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		seminarWithMissingID := mockSeminars[0]
		seminarWithMissingID.LateProductID = nil // one ID is missing

		seminarsWithOneBad := []seminar.Seminar{
			seminarWithMissingID,
			mockSeminars[1],
		}

		validProducts := mockProducts[5:]

		mockSeminarRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(seminarsWithOneBad, nil)
		mockProductRepo.EXPECT().SelectWithDeletedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(validProducts, nil)
		mockSeminarRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(2), nil)

		// Act
		details, total, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, details, 1)
		assert.Equal(t, details[0].Seminar.ID, seminarID_2)
		assert.Equal(t, int64(2), total)
	})

	t.Run("success with incomplete products from repo", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0

		// Products for the first seminar are missing from the response
		incompleteProducts := mockProducts[5:]

		mockSeminarRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(mockSeminars, nil)
		mockProductRepo.EXPECT().SelectWithDeletedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(incompleteProducts, nil)
		mockSeminarRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(2), nil)

		// Act
		details, total, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, details, 1)
		assert.Equal(t, details[0].Seminar.ID, seminarID_2)
		assert.Equal(t, int64(2), total)
	})

	t.Run("success empty list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockSeminarRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return([]seminar.Seminar{}, nil)
		mockProductRepo.EXPECT().SelectWithDeletedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]product.Product{}, nil)
		mockSeminarRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(0), nil)

		// Act
		details, total, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, details, 0)
		assert.Equal(t, int64(0), total)
	})
}

func TestService_ListUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	seminarID_1 := uuid.New().String()
	rproductID_1 := uuid.New().String()
	eproductID_1 := uuid.New().String()
	lproductID_1 := uuid.New().String()
	esproductID_1 := uuid.New().String()
	lsproductID_1 := uuid.New().String()

	seminarID_2 := uuid.New().String()
	rproductID_2 := uuid.New().String()
	eproductID_2 := uuid.New().String()
	lproductID_2 := uuid.New().String()
	esproductID_2 := uuid.New().String()
	lsproductID_2 := uuid.New().String()

	layout := "2006-Jan-02"

	beforeNow, _ := time.Parse(layout, "2024-Aug-03")
	afterNow, _ := time.Parse(layout, "2099-Dec-03")

	mockSeminars := []seminar.Seminar{
		{
			ID:                      seminarID_1,
			Name:                    "Seminar 1 name",
			ShortDescription:        "Seminar 1 short description",
			ReservationProductID:    &rproductID_1,
			EarlyProductID:          &eproductID_1,
			LateProductID:           &lproductID_1,
			EarlySurchargeProductID: &esproductID_1,
			LateSurchargeProductID:  &lsproductID_1,
			LatePaymentDate:         beforeNow,
		},
		{
			ID:                      seminarID_2,
			Name:                    "Seminar 2 name",
			ShortDescription:        "Seminar 2 short description",
			ReservationProductID:    &rproductID_2,
			EarlyProductID:          &eproductID_2,
			LateProductID:           &lproductID_2,
			EarlySurchargeProductID: &esproductID_2,
			LateSurchargeProductID:  &lsproductID_2,
			LatePaymentDate:         afterNow,
		},
	}

	mockProducts := []product.Product{
		{
			ID:          rproductID_1,
			Price:       11.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          eproductID_1,
			Price:       12.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          lproductID_1,
			Price:       13.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          esproductID_1,
			Price:       14.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          lsproductID_1,
			Price:       15.11,
			DetailsID:   seminarID_1,
			DetailsType: "seminar",
		},
		{
			ID:          rproductID_2,
			Price:       16.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
		{
			ID:          eproductID_2,
			Price:       17.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
		{
			ID:          lproductID_2,
			Price:       18.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
		{
			ID:          esproductID_2,
			Price:       19.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
		{
			ID:          lsproductID_2,
			Price:       20.11,
			DetailsID:   seminarID_2,
			DetailsType: "seminar",
		},
	}

	expectedDetails_1 := &seminar.SeminarDetails{
		Seminar:                        &mockSeminars[0],
		ReservationPrice:               mockProducts[0].Price,
		EarlyPrice:                     mockProducts[1].Price,
		LatePrice:                      mockProducts[2].Price,
		EarlySurchargePrice:            mockProducts[3].Price,
		LateSurchargePrice:             mockProducts[4].Price,
		CurrentPrice:                   mockProducts[2].Price,
		CurrentPriceProductID:          lproductID_1,
		CurrentSurchargePrice:          mockProducts[4].Price,
		CurrentSurchargePriceProductID: lsproductID_1,
	}
	expectedDetails_2 := &seminar.SeminarDetails{
		Seminar:                        &mockSeminars[1],
		ReservationPrice:               mockProducts[5].Price,
		EarlyPrice:                     mockProducts[6].Price,
		LatePrice:                      mockProducts[7].Price,
		EarlySurchargePrice:            mockProducts[8].Price,
		LateSurchargePrice:             mockProducts[9].Price,
		CurrentPrice:                   mockProducts[6].Price,
		CurrentPriceProductID:          eproductID_2,
		CurrentSurchargePrice:          mockProducts[8].Price,
		CurrentSurchargePriceProductID: esproductID_2,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockSeminarRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(mockSeminars, nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)
		mockSeminarRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(2), nil)

		// Act
		details, total, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(details[0].Seminar, expectedDetails_1.Seminar) {
			t.Errorf("ListUnpublished() got %v, want %v", details[0].Seminar, expectedDetails_1.Seminar)
		}
		if !reflect.DeepEqual(details[1].Seminar, expectedDetails_2.Seminar) {
			t.Errorf("ListUnpublished() got %v, want %v", details[1].Seminar, expectedDetails_2.Seminar)
		}
		if details[0].CurrentPrice != expectedDetails_1.CurrentPrice {
			t.Errorf("ListUnpublished() got %f, want %f", details[0].CurrentPrice, expectedDetails_1.CurrentPrice)
		}
		if details[1].CurrentPrice != expectedDetails_2.CurrentPrice {
			t.Errorf("ListUnpublished() got %f, want %f", details[1].CurrentPrice, expectedDetails_2.CurrentPrice)
		}
		if total != 2 {
			t.Errorf("ListUnpublished() got total %d, want %d", total, 2)
		}
	})

	t.Run("db error", func(t *testing.T) {
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockSeminarRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListUnpublished(context.Background(), limit, offset)
		// Assert
		assert.Error(t, err)
	})

	t.Run("db error on count", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockSeminarRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(mockSeminars, nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil)
		dbErr := errors.New("db count error")
		mockSeminarRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, _, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.Error(t, err)
	})

	t.Run("success with one seminar having a missing product id", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		seminarWithMissingID := mockSeminars[0]
		seminarWithMissingID.LateProductID = nil // one ID is missing

		seminarsWithOneBad := []seminar.Seminar{
			seminarWithMissingID,
			mockSeminars[1],
		}

		validProducts := mockProducts[5:]

		mockSeminarRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(seminarsWithOneBad, nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(validProducts, nil)
		mockSeminarRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(2), nil)

		// Act
		details, total, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, details, 1)
		assert.Equal(t, details[0].Seminar.ID, seminarID_2)
		assert.Equal(t, int64(2), total)
	})

	t.Run("success with incomplete products from repo", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0

		// Products for the first seminar are missing from the response
		incompleteProducts := mockProducts[5:]

		mockSeminarRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(mockSeminars, nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(incompleteProducts, nil)
		mockSeminarRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(2), nil)

		// Act
		details, total, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, details, 1)
		assert.Equal(t, details[0].Seminar.ID, seminarID_2)
		assert.Equal(t, int64(2), total)
	})

	t.Run("success empty list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockSeminarRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return([]seminar.Seminar{}, nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return([]product.Product{}, nil)
		mockSeminarRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(0), nil)

		// Act
		details, total, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, details, 0)
		assert.Equal(t, int64(0), total)
	})
}

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	layout := "2006-Jan-02"

	date, _ := time.Parse(layout, "2033-Dec-05")
	endingDate, _ := time.Parse(layout, "2033-Dec-07")
	latePaymentDate, _ := time.Parse(layout, "2033-Nov-03")

	createReq := &seminar.CreateRequest{
		Name:                "Seminar name",
		ShortDescription:    "Seminar short description",
		ReservationPrice:    11.11,
		EarlyPrice:          12.22,
		LatePrice:           13.33,
		EarlySurchargePrice: 14.44,
		LateSurchargePrice:  15.55,
		Date:                date,
		EndingDate:          endingDate,
		LatePaymentDate:     latePaymentDate,
		Place:               "Seminar place",
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		var createdSeminar *seminar.Seminar
		mockTxSeminarRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, s *seminar.Seminar) {
				createdSeminar = s
			}).Return(nil)

		var createdProducts []*product.Product
		mockTxProductRepo.EXPECT().CreateBatch(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, products ...*product.Product) {
				createdProducts = products
			}).Return(nil).AnyTimes()

		// Act
		resp, err := testService.Create(context.Background(), createReq)

		// Assert
		assert.NoError(t, err)

		// Assert Seminar
		if _, err := uuid.Parse(createdSeminar.ID); err != nil {
			t.Errorf("Expected seminar.ID to be a valid UUID, got %s", createdSeminar.ID)
		}
		assert.Equal(t, createReq.Name, createdSeminar.Name)
		assert.False(t, createdSeminar.InStock)

		// Assert Products
		assert.Len(t, createdProducts, 5)

		productPriceMap := map[float32]bool{
			createReq.ReservationPrice:    false,
			createReq.EarlyPrice:          false,
			createReq.LatePrice:           false,
			createReq.EarlySurchargePrice: false,
			createReq.LateSurchargePrice:  false,
		}

		for _, p := range createdProducts {
			if _, err := uuid.Parse(p.ID); err != nil {
				t.Errorf("Expected product.ID to be a valid UUID, got %s", p.ID)
			}
			assert.Equal(t, createdSeminar.ID, p.DetailsID)
			assert.Equal(t, "seminar", p.DetailsType)
			assert.False(t, p.InStock)
			if _, ok := productPriceMap[p.Price]; ok {
				productPriceMap[p.Price] = true
			}
		}

		for price, found := range productPriceMap {
			if !found {
				t.Errorf("product with price %f was not created", price)
			}
		}

		// Assert Response
		assert.Equal(t, createdSeminar.ID, resp.ID)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(seminarmock.NewMockRepository(ctrl))
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(productmock.NewMockRepository(ctrl))

		invalidReq := &seminar.CreateRequest{Name: "a"} // Invalid name length

		// Act
		_, err := testService.Create(context.Background(), invalidReq)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxProductRepo.EXPECT().CreateBatch(gomock.Any(), gomock.Any()).Return(nil)
		dbErr := errors.New("database error")
		mockTxSeminarRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr).AnyTimes()

		// Act
		_, err := testService.Create(context.Background(), createReq)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().SetInStock(gomock.Any(), seminarID, true).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), seminarID, true).Return(int64(5), nil)

		// Act
		err := testService.Publish(context.Background(), seminarID)

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

	t.Run("seminar not found", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().SetInStock(gomock.Any(), seminarID, true).Return(int64(0), nil)

		// Act
		err := testService.Publish(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("not all products are found", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().SetInStock(gomock.Any(), seminarID, true).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), seminarID, true).Return(int64(3), nil)

		// Act
		err := testService.Publish(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxSeminarRepo.EXPECT().SetInStock(gomock.Any(), seminarID, true).Return(int64(0), dbErr)

		// Act
		err := testService.Publish(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Unpublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().SetInStock(gomock.Any(), seminarID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), seminarID, false).Return(int64(5), nil)

		// Act
		err := testService.Unpublish(context.Background(), seminarID)

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

	t.Run("seminar not found", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().SetInStock(gomock.Any(), seminarID, false).Return(int64(0), nil)

		// Act
		err := testService.Unpublish(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("not all products are found", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().SetInStock(gomock.Any(), seminarID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), seminarID, false).Return(int64(3), nil)

		// Act
		err := testService.Unpublish(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxSeminarRepo.EXPECT().SetInStock(gomock.Any(), seminarID, false).Return(int64(0), dbErr)

		// Act
		err := testService.Unpublish(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	seminarID := uuid.New().String()
	rproductID := uuid.New().String()
	eproductID := uuid.New().String()
	lproductID := uuid.New().String()
	esproductID := uuid.New().String()
	lsproductID := uuid.New().String()

	layout := "2006-Jan-02"

	date, _ := time.Parse(layout, "2025-Dec-03")
	endingDate, _ := time.Parse(layout, "2025-Dec-06")
	latePaymentDate, _ := time.Parse(layout, "2025-Nov-20")

	mockSeminar := &seminar.Seminar{
		ID:                      seminarID,
		Name:                    "Seminar old name",
		ShortDescription:        "Seminar old short description",
		Date:                    date,
		EndingDate:              endingDate,
		LatePaymentDate:         latePaymentDate,
		ReservationProductID:    &rproductID,
		EarlyProductID:          &eproductID,
		LateProductID:           &lproductID,
		EarlySurchargeProductID: &esproductID,
		LateSurchargeProductID:  &lsproductID,
	}

	mockProducts := []product.Product{
		{
			ID:          rproductID,
			Price:       11.11,
			DetailsID:   seminarID,
			DetailsType: "seminar",
		},
		{
			ID:          eproductID,
			Price:       12.22,
			DetailsID:   seminarID,
			DetailsType: "seminar",
		},
		{
			ID:          lproductID,
			Price:       13.33,
			DetailsID:   seminarID,
			DetailsType: "seminar",
		},
		{
			ID:          esproductID,
			Price:       14.44,
			DetailsID:   seminarID,
			DetailsType: "seminar",
		},
		{
			ID:          lsproductID,
			Price:       15.55,
			DetailsID:   seminarID,
			DetailsType: "seminar",
		},
	}

	newName := "New seminar name"
	newLongDescription := "New seminar long description"
	newReservationPrice := float32(44.44)
	newLatePaymentDate, _ := time.Parse(layout, "2025-Nov-12")
	newTags := []string{"new", "seminar", "tags"}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().Get(gomock.Any(), seminarID).Return(mockSeminar, nil).AnyTimes()
		mockTxProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil).AnyTimes()

		var seminarUpdates map[string]any
		mockTxSeminarRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, _ *seminar.Seminar, u map[string]any) {
				seminarUpdates = u
			}).Return(int64(1), nil).AnyTimes()

		var productUpdates map[string]any
		mockTxProductRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, _ *product.Product, u map[string]any) {
				productUpdates = u
			}).Return(int64(1), nil).AnyTimes()

		// Act
		updates, err := testService.Update(context.Background(), &seminar.UpdateRequest{
			ID:               seminarID,
			Name:             &newName,
			LongDescription:  &newLongDescription,
			ReservationPrice: &newReservationPrice,
			LatePaymentDate:  &newLatePaymentDate,
			Tags:             newTags,
		})

		// Assert
		assert.NoError(t, err)
		// Assert seminar
		seminarUpdatesFromResp, ok := updates["seminar"].(map[string]any)
		if !ok {
			t.Error("response does not contain 'seminar' key")
		}
		if name, ok := seminarUpdatesFromResp["name"].(string); !ok || name != newName {
			t.Errorf("seminar.Name in response %v, want %s", seminarUpdatesFromResp["name"], newName)
		}
		if longDesc, ok := seminarUpdatesFromResp["long_description"].(string); !ok || longDesc != newLongDescription {
			t.Errorf("seminar.LongDescription in response %v, want %s", seminarUpdatesFromResp["long_description"], newLongDescription)
		}
		if latePaymentDate, ok := seminarUpdatesFromResp["late_payment_date"].(time.Time); !ok || latePaymentDate != newLatePaymentDate {
			t.Errorf("seminar.LatePaymentDate in response %v, want %v", seminarUpdatesFromResp["late_payment_date"], newLatePaymentDate)
		}
		if tags, ok := seminarUpdatesFromResp["tags"].([]string); !ok || !reflect.DeepEqual(tags, newTags) {
			t.Errorf("seminar.Tags in response %v, want %v", seminarUpdatesFromResp["tags"], newTags)
		}
		if name, ok := seminarUpdates["name"].(string); !ok || name != newName {
			t.Errorf("seminar.Name passed to repo %v, want %s", seminarUpdates["name"], newName)
		}
		if longDesc, ok := seminarUpdates["long_description"].(string); !ok || longDesc != newLongDescription {
			t.Errorf("seminar.LongDescription passed to repo %v, want %s", seminarUpdates["long_description"], newLongDescription)
		}
		if latePaymentDate, ok := seminarUpdates["late_payment_date"].(time.Time); !ok || latePaymentDate != newLatePaymentDate {
			t.Errorf("seminar.LatePaymentDate passed to repo %v, want %v", seminarUpdates["late_payment_date"], newLatePaymentDate)
		}
		if tags, ok := seminarUpdates["tags"].([]string); !ok || !reflect.DeepEqual(tags, newTags) {
			t.Errorf("seminar.Tags passed to repo %v, want %v", seminarUpdates["tags"], newTags)
		}

		// Assert product
		productUpdatesFromResp, ok := updates["reservation_product"].(map[string]any)
		if !ok {
			t.Error("response does not contain 'reservation_product' key")
		}
		if price, ok := productUpdatesFromResp["price"].(float32); !ok || price != newReservationPrice {
			t.Errorf("product.Price in response %v, want %f", productUpdatesFromResp["price"], newReservationPrice)
		}
		if price, ok := productUpdates["price"].(float32); !ok || price != newReservationPrice {
			t.Errorf("product.Price passed to repo %v, want %f", productUpdates["price"], newReservationPrice)
		}
	})

	t.Run("success with multiple updated products", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().Get(gomock.Any(), seminarID).Return(mockSeminar, nil).AnyTimes()
		mockTxProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil).AnyTimes()

		mockTxSeminarRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil).AnyTimes()

		allProductUpdates := make(map[string]any)
		mockTxProductRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, p *product.Product, u map[string]any) {
				allProductUpdates[p.ID] = u
			}).Return(int64(1), nil).AnyTimes()

		newLatePrice := float32(23.55)
		newLateSurchargePrice := float32(99.99)

		// Act
		updates, err := testService.Update(context.Background(), &seminar.UpdateRequest{
			ID:                 seminarID,
			ReservationPrice:   &newReservationPrice,
			LatePrice:          &newLatePrice,
			LateSurchargePrice: &newLateSurchargePrice,
		})

		// Assert
		assert.NoError(t, err)
		rproductUpdates, ok := allProductUpdates[rproductID].(map[string]any)
		if !ok {
			t.Error("reservation product updates was not passed to the repo")
		}
		if price, ok := rproductUpdates["price"].(float32); !ok || price != newReservationPrice {
			t.Errorf("reservation_product.Price passed to repo %v, want %f", rproductUpdates["price"], newReservationPrice)
		}
		lproductUpdates, ok := allProductUpdates[lproductID].(map[string]any)
		if !ok {
			t.Error("late product updates was not passed to the repo")
		}
		if price, ok := lproductUpdates["price"].(float32); !ok || price != newLatePrice {
			t.Errorf("late_product.Price passed to repo %v, want %f", lproductUpdates["price"], newLatePrice)
		}
		lsproductUpdates, ok := allProductUpdates[lsproductID].(map[string]any)
		if !ok {
			t.Error("late surcharge product updates was not passed to the repo")
		}
		if price, ok := lsproductUpdates["price"].(float32); !ok || price != newLateSurchargePrice {
			t.Errorf("late_surcharge_product.Price passed to repo %v, want %f", lsproductUpdates["price"], newLateSurchargePrice)
		}
		rproductUpdatesFromResp, ok := updates["reservation_product"].(map[string]any)
		if !ok {
			t.Error("response does not have 'reservation_product' key")
		}
		if price, ok := rproductUpdatesFromResp["price"].(float32); !ok || price != newReservationPrice {
			t.Errorf("reservation_product.Price from response %v, want %f", rproductUpdatesFromResp["price"], newReservationPrice)
		}
		lproductUpdatesFromResp, ok := updates["late_product"].(map[string]any)
		if !ok {
			t.Error("response does not have 'late_product' key")
		}
		if price, ok := lproductUpdatesFromResp["price"].(float32); !ok || price != newLatePrice {
			t.Errorf("late_product.Price from response %v, want %f", lproductUpdatesFromResp["price"], newLatePrice)
		}
		lsproductUpdatesFromResp, ok := updates["late_surcharge_product"].(map[string]any)
		if !ok {
			t.Error("response does not have 'late_surcharge_product' key")
		}
		if price, ok := lsproductUpdatesFromResp["price"].(float32); !ok || price != newLateSurchargePrice {
			t.Errorf("late_surcharge_product.Price from response %v, want %f", lsproductUpdatesFromResp["price"], newLateSurchargePrice)
		}
	})

	t.Run("success with no updates", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().Get(gomock.Any(), seminarID).Return(mockSeminar, nil).AnyTimes()
		mockTxProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil).AnyTimes()

		// Act
		_, err := testService.Update(context.Background(), &seminar.UpdateRequest{
			ID: seminarID, // no new fields
		})

		// Assert
		assert.NoError(t, err)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		invalidDate, _ := time.Parse(layout, "2003-Nov-01") // Date in the past
		invalidName := "1invalidseminarname"

		// Act
		_, err := testService.Update(context.Background(), &seminar.UpdateRequest{
			ID:   seminarID,
			Date: &invalidDate,
			Name: &invalidName,
		})

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().Get(gomock.Any(), seminarID).Return(nil, gorm.ErrRecordNotFound).AnyTimes()

		// Act
		_, err := testService.Update(context.Background(), &seminar.UpdateRequest{
			ID: seminarID,
		})

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().Get(gomock.Any(), seminarID).Return(mockSeminar, nil).AnyTimes()
		mockTxProductRepo.EXPECT().SelectByIDs(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockProducts, nil).AnyTimes()
		dbErr := errors.New("database error")
		mockTxSeminarRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), dbErr).AnyTimes()

		// Act
		_, err := testService.Update(context.Background(), &seminar.UpdateRequest{
			ID:   seminarID,
			Name: &newName,
		})

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(&seminar.Seminar{}, nil)
		mockTxSeminarRepo.EXPECT().SetInStock(gomock.Any(), seminarID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), seminarID, false).Return(int64(5), nil)
		mockTxSeminarRepo.EXPECT().Delete(gomock.Any(), seminarID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().DeleteByDetailsID(gomock.Any(), seminarID).Return(int64(5), nil)

		// Act
		err := testService.Delete(context.Background(), seminarID)

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
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		err := testService.Delete(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("cannot unpublish all products", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(&seminar.Seminar{}, nil)
		mockTxSeminarRepo.EXPECT().SetInStock(gomock.Any(), seminarID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), seminarID, false).Return(int64(3), nil)

		// Act
		err := testService.Delete(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(&seminar.Seminar{}, nil)
		mockTxSeminarRepo.EXPECT().SetInStock(gomock.Any(), seminarID, false).Return(int64(1), nil)
		dbErr := errors.New("database error")
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), seminarID, false).Return(int64(0), dbErr)

		// Act
		err := testService.Delete(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_DeletePermanent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().DeletePermanent(gomock.Any(), seminarID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().DeletePermanentByDetailsID(gomock.Any(), seminarID).Return(int64(5), nil)

		// Act
		err := testService.DeletePermanent(context.Background(), seminarID)

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
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().DeletePermanent(gomock.Any(), seminarID).Return(int64(0), nil)

		// Act
		err := testService.DeletePermanent(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("cannot delete all products", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().DeletePermanent(gomock.Any(), seminarID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().DeletePermanentByDetailsID(gomock.Any(), seminarID).Return(int64(3), nil)

		// Act
		err := testService.DeletePermanent(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxSeminarRepo.EXPECT().DeletePermanent(gomock.Any(), seminarID).Return(int64(0), dbErr)

		// Act
		err := testService.DeletePermanent(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Restore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeminarRepo := seminarmock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)

	testService := New(mockSeminarRepo, mockProductRepo)

	// Use an in-memory SQLite DB for testing transactions.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// This prevents GORM from starting a real DB transaction,
		// allowing the mock repositories to work as expected.
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().Restore(gomock.Any(), seminarID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().RestoreByDetailsID(gomock.Any(), seminarID).Return(int64(5), nil)

		// Act
		err := testService.Restore(context.Background(), seminarID)

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
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().Restore(gomock.Any(), seminarID).Return(int64(0), nil)

		// Act
		err := testService.Restore(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("cannot restore all products", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxSeminarRepo.EXPECT().Restore(gomock.Any(), seminarID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().RestoreByDetailsID(gomock.Any(), seminarID).Return(int64(3), nil)

		// Act
		err := testService.Restore(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxSeminarRepo := seminarmock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockSeminarRepo.EXPECT().DB().Return(db).AnyTimes()
		mockSeminarRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxSeminarRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxSeminarRepo.EXPECT().Restore(gomock.Any(), seminarID).Return(int64(0), dbErr)

		// Act
		err := testService.Restore(context.Background(), seminarID)

		// Assert
		assert.Error(t, err)
	})
}
