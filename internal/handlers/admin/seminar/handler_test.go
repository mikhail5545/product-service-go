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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mikhail5545/product-service-go/internal/models/product"
	"github.com/mikhail5545/product-service-go/internal/models/seminar"
	seminarservice "github.com/mikhail5545/product-service-go/internal/services/seminar"
	seminarmock "github.com/mikhail5545/product-service-go/internal/test/services/seminar_mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

	seminarID := "c6248da5-a2eb-4abd-be56-a19715104c00"
	rproductID := "866561c2-a65a-4159-a5d8-a0ae5401e0c1"
	eproductID := "7baa5ff9-a864-4144-b42c-8ce6bd56ac25"
	lproductID := "38fcb2f8-d377-4b08-9eb9-8de9a89d4528"
	esproductID := "0cb3a9a5-9dd0-4ca9-b528-275071e3eb98"
	lsproductID := "14212b87-ca38-41d5-bba2-2a273fe60977"

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

	mockDetails := &seminar.SeminarDetails{
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

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		mockService.EXPECT().Get(gomock.Any(), seminarID).Return(mockDetails, nil)

		// Act
		err := handler.Get(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"seminar_details": mockDetails}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		serviceErr := &seminarservice.Error{
			Msg:  "Course not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Get(gomock.Any(), seminarID).Return(nil, serviceErr)

		// Act
		err := handler.Get(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues("invalid-uuid")

		// Act
		err := handler.Get(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid seminar ID")
	})
}

func TestHandler_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

	seminarID := "c6248da5-a2eb-4abd-be56-a19715104c00"
	rproductID := "866561c2-a65a-4159-a5d8-a0ae5401e0c1"
	eproductID := "7baa5ff9-a864-4144-b42c-8ce6bd56ac25"
	lproductID := "38fcb2f8-d377-4b08-9eb9-8de9a89d4528"
	esproductID := "0cb3a9a5-9dd0-4ca9-b528-275071e3eb98"
	lsproductID := "14212b87-ca38-41d5-bba2-2a273fe60977"

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

	mockDetails := &seminar.SeminarDetails{
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

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		mockService.EXPECT().GetWithDeleted(gomock.Any(), seminarID).Return(mockDetails, nil)

		// Act
		err := handler.GetWithDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"seminar_details": mockDetails}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		serviceErr := &seminarservice.Error{
			Msg:  "Course not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().GetWithDeleted(gomock.Any(), seminarID).Return(nil, serviceErr)

		// Act
		err := handler.GetWithDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues("invalid-uuid")

		// Act
		err := handler.GetWithDeleted(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid seminar ID")
	})
}

func TestHandler_GetWithUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

	seminarID := "c6248da5-a2eb-4abd-be56-a19715104c00"
	rproductID := "866561c2-a65a-4159-a5d8-a0ae5401e0c1"
	eproductID := "7baa5ff9-a864-4144-b42c-8ce6bd56ac25"
	lproductID := "38fcb2f8-d377-4b08-9eb9-8de9a89d4528"
	esproductID := "0cb3a9a5-9dd0-4ca9-b528-275071e3eb98"
	lsproductID := "14212b87-ca38-41d5-bba2-2a273fe60977"

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

	mockDetails := &seminar.SeminarDetails{
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

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		mockService.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(mockDetails, nil)

		// Act
		err := handler.GetWithUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"seminar_details": mockDetails}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		serviceErr := &seminarservice.Error{
			Msg:  "Course not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(nil, serviceErr)

		// Act
		err := handler.GetWithUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues("invalid-uuid")

		// Act
		err := handler.GetWithUnpublished(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid seminar ID")
	})
}

func TestHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

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

	mockDetails_1 := &seminar.SeminarDetails{
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
	mockDetails_2 := &seminar.SeminarDetails{
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
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().List(gomock.Any(), 2, 0).Return([]seminar.SeminarDetails{*mockDetails_1, *mockDetails_2}, int64(2), nil)

		// Act
		err := handler.List(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"seminar_details": []seminar.SeminarDetails{*mockDetails_1, *mockDetails_2}, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("invalid pagination params", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?a=2", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().List(gomock.Any(), 10, 0).Return([]seminar.SeminarDetails{*mockDetails_1, *mockDetails_2}, int64(2), nil)

		// Act
		err := handler.List(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"seminar_details": []seminar.SeminarDetails{*mockDetails_1, *mockDetails_2}, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		serviceErr := &seminarservice.Error{
			Msg:  "Failed to get seminars",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().List(gomock.Any(), 2, 0).Return(nil, int64(0), serviceErr)

		// Act
		err := handler.List(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestHandler_ListDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

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

	mockDetails_1 := &seminar.SeminarDetails{
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
	mockDetails_2 := &seminar.SeminarDetails{
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
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().ListDeleted(gomock.Any(), 2, 0).Return([]seminar.SeminarDetails{*mockDetails_1, *mockDetails_2}, int64(2), nil)

		// Act
		err := handler.ListDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"seminar_details": []seminar.SeminarDetails{*mockDetails_1, *mockDetails_2}, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("invalid pagination params", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?a=2", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().ListDeleted(gomock.Any(), 10, 0).Return([]seminar.SeminarDetails{*mockDetails_1, *mockDetails_2}, int64(2), nil)

		// Act
		err := handler.ListDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"seminar_details": []seminar.SeminarDetails{*mockDetails_1, *mockDetails_2}, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		serviceErr := &seminarservice.Error{
			Msg:  "Failed to get seminars",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().ListDeleted(gomock.Any(), 2, 0).Return(nil, int64(0), serviceErr)

		// Act
		err := handler.ListDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestHandler_ListUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

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

	mockDetails_1 := &seminar.SeminarDetails{
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
	mockDetails_2 := &seminar.SeminarDetails{
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
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().ListUnpublished(gomock.Any(), 2, 0).Return([]seminar.SeminarDetails{*mockDetails_1, *mockDetails_2}, int64(2), nil)

		// Act
		err := handler.ListUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"seminar_details": []seminar.SeminarDetails{*mockDetails_1, *mockDetails_2}, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("invalid pagination params", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?a=2", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().ListUnpublished(gomock.Any(), 10, 0).Return([]seminar.SeminarDetails{*mockDetails_1, *mockDetails_2}, int64(2), nil)

		// Act
		err := handler.ListUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"seminar_details": []seminar.SeminarDetails{*mockDetails_1, *mockDetails_2}, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		serviceErr := &seminarservice.Error{
			Msg:  "Failed to get seminars",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().ListUnpublished(gomock.Any(), 2, 0).Return(nil, int64(0), serviceErr)

		// Act
		err := handler.ListUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

	seminarID := uuid.New().String()
	rproductID := uuid.New().String()
	eproductID := uuid.New().String()
	lproductID := uuid.New().String()
	esproductID := uuid.New().String()
	lsproductID := uuid.New().String()

	createReq := seminar.CreateRequest{
		Name:                "Seminar name",
		ShortDescription:    "Seminar short description",
		ReservationPrice:    11.11,
		EarlyPrice:          22.22,
		LatePrice:           33.33,
		EarlySurchargePrice: 44.44,
		LateSurchargePrice:  55.55,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		reqJSON, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		createResp := &seminar.CreateResponse{
			ID:                      seminarID,
			ReservationProductID:    rproductID,
			EarlyProductID:          eproductID,
			LateProductID:           lproductID,
			EarlySurchargeProductID: esproductID,
			LateSurchargeProductID:  lsproductID,
		}
		mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(createResp, nil)

		// Act
		err := handler.Create(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		expectedResp := map[string]any{"response": createResp}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		reqJSON, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		serviceErr := &seminarservice.Error{
			Msg:  "Failed to create seminar",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, serviceErr)

		// Act
		err := handler.Create(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("invalid request JSON payload", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "bad json}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.Create(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		mockService.EXPECT().Publish(gomock.Any(), seminarID).Return(nil)

		// Act
		err := handler.Publish(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		serviceErr := &seminarservice.Error{
			Msg:  "Seminar not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Publish(gomock.Any(), seminarID).Return(serviceErr)

		// Act
		err := handler.Publish(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues("invalid-id")

		// Act
		err := handler.Publish(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_Unpublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		mockService.EXPECT().Unpublish(gomock.Any(), seminarID).Return(nil)

		// Act
		err := handler.Unpublish(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		serviceErr := &seminarservice.Error{
			Msg:  "Seminar not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Unpublish(gomock.Any(), seminarID).Return(serviceErr)

		// Act
		err := handler.Unpublish(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues("invalid-id")

		// Act
		err := handler.Unpublish(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

	seminarID := uuid.New().String()

	newName := "New seminar name"
	newDescription := "New seminar description"

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		updateReq := &seminar.UpdateRequest{
			Name:             &newName,
			ShortDescription: &newDescription,
		}
		jsonReq, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonReq))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		updates := map[string]any{"name": newName, "short_description": newDescription}
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(updates, nil)

		// Act
		err := handler.Update(c)

		// Assert
		expectedResp := map[string]any{"updates": updates}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		updateReq := seminar.UpdateRequest{
			Name:             &newName,
			ShortDescription: &newDescription,
		}
		jsonReq, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonReq))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		updateReq.ID = seminarID
		serviceErr := &seminarservice.Error{
			Msg:  "Seminar not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Update(gomock.Any(), &updateReq).Return(nil, serviceErr)

		// Act
		err := handler.Update(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid request JSON payload", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "bad json}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		// Act
		err := handler.Update(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "new name"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues("invalid-id")

		// Act
		err := handler.Update(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		mockService.EXPECT().Delete(gomock.Any(), seminarID).Return(nil)

		// Act
		err := handler.Delete(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		serviceErr := &seminarservice.Error{
			Msg:  "Seminar not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Delete(gomock.Any(), seminarID).Return(serviceErr)

		// Act
		err := handler.Delete(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues("invalid-id")

		// Act
		err := handler.Delete(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_DeletePermanent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		mockService.EXPECT().DeletePermanent(gomock.Any(), seminarID).Return(nil)

		// Act
		err := handler.DeletePermanent(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		serviceErr := &seminarservice.Error{
			Msg:  "Seminar not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().DeletePermanent(gomock.Any(), seminarID).Return(serviceErr)

		// Act
		err := handler.DeletePermanent(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues("invalid-id")

		// Act
		err := handler.DeletePermanent(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_Restore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := seminarmock.NewMockService(ctrl)
	handler := New(mockService)

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		mockService.EXPECT().Restore(gomock.Any(), seminarID).Return(nil)

		// Act
		err := handler.Restore(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(seminarID)

		serviceErr := &seminarservice.Error{
			Msg:  "Seminar not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Restore(gomock.Any(), seminarID).Return(serviceErr)

		// Act
		err := handler.Restore(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues("invalid-id")

		// Act
		err := handler.Restore(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
