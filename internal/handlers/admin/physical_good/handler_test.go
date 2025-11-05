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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	physicalgood "github.com/mikhail5545/product-service-go/internal/models/physical_good"
	pjysicalgood "github.com/mikhail5545/product-service-go/internal/models/physical_good"
	physicalgoodservice "github.com/mikhail5545/product-service-go/internal/services/physical_good"
	physicalgoodmock "github.com/mikhail5545/product-service-go/internal/test/services/physical_good_mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID := uuid.New().String()

	mockPhysicalGoodDetails := &pjysicalgood.PhysicalGoodDetails{
		PhysicalGood: physicalgood.PhysicalGood{
			ID:               goodID,
			Name:             "Physical good name",
			ShortDescription: "Physical good short description",
		},
		Price:     33.33,
		ProductID: uuid.New().String(),
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		mockService.EXPECT().Get(gomock.Any(), goodID).Return(mockPhysicalGoodDetails, nil)

		// Act
		err := handler.Get(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"physical_good_details": mockPhysicalGoodDetails}
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
		c.SetParamValues(goodID)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Get(gomock.Any(), goodID).Return(nil, serviceErr)

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
		assert.Contains(t, rec.Body.String(), "Invalid physical good ID")
	})
}

func TestHandler_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID := uuid.New().String()

	mockPhysicalGoodDetails := &pjysicalgood.PhysicalGoodDetails{
		PhysicalGood: physicalgood.PhysicalGood{
			ID:               goodID,
			Name:             "Physical good name",
			ShortDescription: "Physical good short description",
		},
		Price:     33.33,
		ProductID: uuid.New().String(),
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		mockService.EXPECT().GetWithDeleted(gomock.Any(), goodID).Return(mockPhysicalGoodDetails, nil)

		// Act
		err := handler.GetWithDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"physical_good_details": mockPhysicalGoodDetails}
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
		c.SetParamValues(goodID)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().GetWithDeleted(gomock.Any(), goodID).Return(nil, serviceErr)

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
		assert.Contains(t, rec.Body.String(), "Invalid physical good ID")
	})
}

func TestHandler_GetWithUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID := uuid.New().String()

	mockPhysicalGoodDetails := &pjysicalgood.PhysicalGoodDetails{
		PhysicalGood: physicalgood.PhysicalGood{
			ID:               goodID,
			Name:             "Physical good name",
			ShortDescription: "Physical good short description",
		},
		Price:     33.33,
		ProductID: uuid.New().String(),
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		mockService.EXPECT().GetWithUnpublished(gomock.Any(), goodID).Return(mockPhysicalGoodDetails, nil)

		// Act
		err := handler.GetWithUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"physical_good_details": mockPhysicalGoodDetails}
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
		c.SetParamValues(goodID)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), goodID).Return(nil, serviceErr)

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
		assert.Contains(t, rec.Body.String(), "Invalid physical good ID")
	})
}

func TestHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID_1 := uuid.New().String()
	goodID_2 := uuid.New().String()

	mockPhysicalGoodDetails := []physicalgood.PhysicalGoodDetails{
		{
			PhysicalGood: pjysicalgood.PhysicalGood{
				ID:               goodID_1,
				Name:             "Physical good name 1",
				ShortDescription: "Physical good short description 1",
			},
			Price:     11.11,
			ProductID: uuid.New().String(),
		},
		{
			PhysicalGood: pjysicalgood.PhysicalGood{
				ID:               goodID_2,
				Name:             "Physical good name 2",
				ShortDescription: "Physical good short description 2",
			},
			Price:     22.22,
			ProductID: uuid.New().String(),
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().List(gomock.Any(), 2, 0).Return(mockPhysicalGoodDetails, int64(2), nil)

		// Act
		err := handler.List(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"physical_good_details": mockPhysicalGoodDetails, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().List(gomock.Any(), 2, 0).Return(nil, int64(0), serviceErr)

		// Act
		err := handler.List(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestHandler_ListUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID_1 := uuid.New().String()
	goodID_2 := uuid.New().String()

	mockPhysicalGoodDetails := []physicalgood.PhysicalGoodDetails{
		{
			PhysicalGood: pjysicalgood.PhysicalGood{
				ID:               goodID_1,
				Name:             "Physical good name 1",
				ShortDescription: "Physical good short description 1",
			},
			Price:     11.11,
			ProductID: uuid.New().String(),
		},
		{
			PhysicalGood: pjysicalgood.PhysicalGood{
				ID:               goodID_2,
				Name:             "Physical good name 2",
				ShortDescription: "Physical good short description 2",
			},
			Price:     22.22,
			ProductID: uuid.New().String(),
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().ListUnpublished(gomock.Any(), 2, 0).Return(mockPhysicalGoodDetails, int64(2), nil)

		// Act
		err := handler.ListUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"physical_good_details": mockPhysicalGoodDetails, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().ListUnpublished(gomock.Any(), 2, 0).Return(nil, int64(0), serviceErr)

		// Act
		err := handler.ListUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestHandler_ListDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID_1 := uuid.New().String()
	goodID_2 := uuid.New().String()

	mockPhysicalGoodDetails := []physicalgood.PhysicalGoodDetails{
		{
			PhysicalGood: pjysicalgood.PhysicalGood{
				ID:               goodID_1,
				Name:             "Physical good name 1",
				ShortDescription: "Physical good short description 1",
			},
			Price:     11.11,
			ProductID: uuid.New().String(),
		},
		{
			PhysicalGood: pjysicalgood.PhysicalGood{
				ID:               goodID_2,
				Name:             "Physical good name 2",
				ShortDescription: "Physical good short description 2",
			},
			Price:     22.22,
			ProductID: uuid.New().String(),
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().ListDeleted(gomock.Any(), 2, 0).Return(mockPhysicalGoodDetails, int64(2), nil)

		// Act
		err := handler.ListDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"physical_good_details": mockPhysicalGoodDetails, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().ListDeleted(gomock.Any(), 2, 0).Return(nil, int64(0), serviceErr)

		// Act
		err := handler.ListDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID := uuid.New().String()
	productID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		createReq := &physicalgood.CreateRequest{
			Name:             "Physical good name",
			ShortDescription: "Physical good short description",
			Amount:           3,
			Price:            33.33,
		}
		reqJSON, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		createResp := &physicalgood.CreateResponse{ID: goodID, ProductID: productID}
		mockService.EXPECT().Create(gomock.Any(), createReq).Return(createResp, nil)

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
		createReq := physicalgood.CreateRequest{
			Name:             "Physical good name",
			ShortDescription: "Physical good short description",
			Amount:           3,
			Price:            33.33,
		}
		reqJSON, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Failed to create physical good",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().Create(gomock.Any(), &createReq).Return(nil, serviceErr)

		// Act
		err := handler.Create(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("bind error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "bad json`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.Create(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		mockService.EXPECT().Publish(gomock.Any(), goodID).Return(nil)

		// Act
		err := handler.Publish(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
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

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Publish(gomock.Any(), goodID).Return(serviceErr)

		// Act
		err := handler.Publish(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestHandler_Unpublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		mockService.EXPECT().Unpublish(gomock.Any(), goodID).Return(nil)

		// Act
		err := handler.Unpublish(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
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

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Unpublish(gomock.Any(), goodID).Return(serviceErr)

		// Act
		err := handler.Unpublish(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		newName := "New physical good name"
		newLongDescription := "New physical good long description"
		updateReq := physicalgood.UpdateRequest{
			Name:            &newName,
			LongDescription: &newLongDescription,
		}
		reqJSON, _ := json.Marshal(updateReq)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		updateReq.ID = goodID
		updates := map[string]any{"name": newName, "long_description": newLongDescription}
		mockService.EXPECT().Update(gomock.Any(), &updateReq).Return(updates, nil)

		// Act
		err := handler.Update(c)

		// Assert
		expectedResp := map[string]any{"updates": updates}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		newName := "New physical good name"
		newLongDescription := "New physical good long description"
		updateReq := physicalgood.UpdateRequest{
			Name:            &newName,
			LongDescription: &newLongDescription,
		}
		reqJSON, _ := json.Marshal(updateReq)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
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
	t.Run("service error", func(t *testing.T) {
		// Arrange
		newName := "New physical good name"
		newLongDescription := "New physical good long description"
		updateReq := physicalgood.UpdateRequest{
			Name:            &newName,
			LongDescription: &newLongDescription,
		}
		reqJSON, _ := json.Marshal(updateReq)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Failed to update physical good",
			Code: http.StatusInternalServerError,
		}
		updateReq.ID = goodID
		mockService.EXPECT().Update(gomock.Any(), &updateReq).Return(nil, serviceErr)

		// Act
		err := handler.Update(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		mockService.EXPECT().Delete(gomock.Any(), goodID).Return(nil)

		// Act
		err := handler.Delete(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
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

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Delete(gomock.Any(), goodID).Return(serviceErr)

		// Act
		err := handler.Delete(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestHandler_DeletePermanent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		mockService.EXPECT().DeletePermanent(gomock.Any(), goodID).Return(nil)

		// Act
		err := handler.DeletePermanent(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
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

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().DeletePermanent(gomock.Any(), goodID).Return(serviceErr)

		// Act
		err := handler.DeletePermanent(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestHandler_Restore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := physicalgoodmock.NewMockService(ctrl)
	handler := New(mockService)

	goodID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		mockService.EXPECT().Restore(gomock.Any(), goodID).Return(nil)

		// Act
		err := handler.Restore(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
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

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(goodID)

		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Restore(gomock.Any(), goodID).Return(serviceErr)

		// Act
		err := handler.Restore(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
