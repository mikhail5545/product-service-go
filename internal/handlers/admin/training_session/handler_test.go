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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	trainingsession "github.com/mikhail5545/product-service-go/internal/models/training_session"
	trainingsessionservice "github.com/mikhail5545/product-service-go/internal/services/training_session"
	trainingsessinmock "github.com/mikhail5545/product-service-go/internal/test/services/training_session_mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID := uuid.New().String()

	mockTsDetails := &trainingsession.TrainingSessionDetails{
		TrainingSession: &trainingsession.TrainingSession{
			ID:               tsID,
			Name:             "Training session name",
			ShortDescription: "Training session short description",
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
		c.SetParamValues(tsID)

		mockService.EXPECT().Get(gomock.Any(), tsID).Return(mockTsDetails, nil)

		// Act
		err := handler.Get(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"training_session_details": mockTsDetails}
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
		c.SetParamValues(tsID)

		mockService.EXPECT().Get(gomock.Any(), tsID).Return(nil, trainingsessionservice.ErrNotFound)

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
		assert.Contains(t, rec.Body.String(), "Invalid training session ID")
	})
}

func TestHandler_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID := uuid.New().String()

	mockTsDetails := &trainingsession.TrainingSessionDetails{
		TrainingSession: &trainingsession.TrainingSession{
			ID:               tsID,
			Name:             "Training session name",
			ShortDescription: "Training session short description",
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
		c.SetParamValues(tsID)

		mockService.EXPECT().GetWithDeleted(gomock.Any(), tsID).Return(mockTsDetails, nil)

		// Act
		err := handler.GetWithDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"training_session_details": mockTsDetails}
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
		c.SetParamValues(tsID)

		mockService.EXPECT().GetWithDeleted(gomock.Any(), tsID).Return(nil, trainingsessionservice.ErrNotFound)

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
		assert.Contains(t, rec.Body.String(), "Invalid training session ID")
	})
}

func TestHandler_GetWithUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID := uuid.New().String()

	mockTsDetails := &trainingsession.TrainingSessionDetails{
		TrainingSession: &trainingsession.TrainingSession{
			ID:               tsID,
			Name:             "Training session name",
			ShortDescription: "Training session short description",
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
		c.SetParamValues(tsID)

		mockService.EXPECT().GetWithUnpublished(gomock.Any(), tsID).Return(mockTsDetails, nil)

		// Act
		err := handler.GetWithUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"training_session_details": mockTsDetails}
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
		c.SetParamValues(tsID)

		mockService.EXPECT().GetWithUnpublished(gomock.Any(), tsID).Return(nil, trainingsessionservice.ErrNotFound)

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
		assert.Contains(t, rec.Body.String(), "Invalid training session ID")
	})
}

func TestHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID_1 := uuid.New().String()
	tsID_2 := uuid.New().String()

	mockTsDetails := []trainingsession.TrainingSessionDetails{
		{
			TrainingSession: &trainingsession.TrainingSession{
				ID:               tsID_1,
				Name:             "Training session name 1",
				ShortDescription: "Training session short description 1",
			},
			Price:     33.33,
			ProductID: uuid.New().String(),
		},
		{
			TrainingSession: &trainingsession.TrainingSession{
				ID:               tsID_2,
				Name:             "Training session name 2",
				ShortDescription: "Training session short description 2",
			},
			Price:     32.22,
			ProductID: uuid.New().String(),
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().List(gomock.Any(), 2, 0).Return(mockTsDetails, int64(2), nil)

		// Act
		err := handler.List(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"training_session_details": mockTsDetails, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().List(gomock.Any(), 2, 0).Return(nil, int64(0), trainingsessionservice.ErrNotFound)

		// Act
		err := handler.List(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid pagination params", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.List(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_ListDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID_1 := uuid.New().String()
	tsID_2 := uuid.New().String()

	mockTsDetails := []trainingsession.TrainingSessionDetails{
		{
			TrainingSession: &trainingsession.TrainingSession{
				ID:               tsID_1,
				Name:             "Training session name 1",
				ShortDescription: "Training session short description 1",
			},
			Price:     33.33,
			ProductID: uuid.New().String(),
		},
		{
			TrainingSession: &trainingsession.TrainingSession{
				ID:               tsID_2,
				Name:             "Training session name 2",
				ShortDescription: "Training session short description 2",
			},
			Price:     32.22,
			ProductID: uuid.New().String(),
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().ListDeleted(gomock.Any(), 2, 0).Return(mockTsDetails, int64(2), nil)

		// Act
		err := handler.ListDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"training_session_details": mockTsDetails, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().ListDeleted(gomock.Any(), 2, 0).Return(nil, int64(0), trainingsessionservice.ErrNotFound)

		// Act
		err := handler.ListDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid pagination params", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.ListDeleted(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_ListUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID_1 := uuid.New().String()
	tsID_2 := uuid.New().String()

	mockTsDetails := []trainingsession.TrainingSessionDetails{
		{
			TrainingSession: &trainingsession.TrainingSession{
				ID:               tsID_1,
				Name:             "Training session name 1",
				ShortDescription: "Training session short description 1",
			},
			Price:     33.33,
			ProductID: uuid.New().String(),
		},
		{
			TrainingSession: &trainingsession.TrainingSession{
				ID:               tsID_2,
				Name:             "Training session name 2",
				ShortDescription: "Training session short description 2",
			},
			Price:     32.22,
			ProductID: uuid.New().String(),
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().ListUnpublished(gomock.Any(), 2, 0).Return(mockTsDetails, int64(2), nil)

		// Act
		err := handler.ListUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"training_session_details": mockTsDetails, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().ListUnpublished(gomock.Any(), 2, 0).Return(nil, int64(0), trainingsessionservice.ErrNotFound)

		// Act
		err := handler.ListUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid pagination params", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.ListUnpublished(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID := uuid.New().String()
	productID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		createReq := &trainingsession.CreateRequest{
			Name:             "Training session name",
			ShortDescription: "Training session description",
			Price:            33.33,
			DurationMinutes:  30,
			Format:           "online",
		}
		reqJSON, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		createResp := &trainingsession.CreateResponse{ID: tsID, ProductID: productID}
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
		createReq := &trainingsession.CreateRequest{
			Name:             "Training session name",
			ShortDescription: "Training session description",
			Price:            33.33,
			DurationMinutes:  30,
			Format:           "online",
		}
		reqJSON, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().Create(gomock.Any(), createReq).Return(nil, trainingsessionservice.ErrNotFound)

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
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(tsID)

		mockService.EXPECT().Publish(gomock.Any(), tsID).Return(nil)

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
		c.SetParamValues(tsID)

		mockService.EXPECT().Publish(gomock.Any(), tsID).Return(trainingsessionservice.ErrNotFound)

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

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(tsID)

		mockService.EXPECT().Unpublish(gomock.Any(), tsID).Return(nil)

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
		c.SetParamValues(tsID)

		mockService.EXPECT().Unpublish(gomock.Any(), tsID).Return(trainingsessionservice.ErrNotFound)

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

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		newName := "New training session name"
		newLongDescription := "New training session long description"
		updateReq := trainingsession.UpdateRequest{
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
		c.SetParamValues(tsID)

		updateReq.ID = tsID
		updates := map[string]any{"name": newName, "long_description": newLongDescription}
		mockService.EXPECT().Update(gomock.Any(), &updateReq).Return(updates, nil)

		// Act
		err := handler.Update(c)

		// Assert
		expectedResp := map[string]any{"updates": updates}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		newName := "New training session name"
		newLongDescription := "New training session long description"
		updateReq := trainingsession.UpdateRequest{
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
		newName := "New training session name"
		newLongDescription := "New training session long description"
		updateReq := trainingsession.UpdateRequest{
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
		c.SetParamValues(tsID)

		updateReq.ID = tsID
		mockService.EXPECT().Update(gomock.Any(), &updateReq).Return(nil, trainingsessionservice.ErrNotFound)

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

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(tsID)

		mockService.EXPECT().Delete(gomock.Any(), tsID).Return(nil)

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
		c.SetParamValues(tsID)

		mockService.EXPECT().Delete(gomock.Any(), tsID).Return(trainingsessionservice.ErrNotFound)

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

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(tsID)

		mockService.EXPECT().DeletePermanent(gomock.Any(), tsID).Return(nil)

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
		c.SetParamValues(tsID)

		mockService.EXPECT().DeletePermanent(gomock.Any(), tsID).Return(trainingsessionservice.ErrNotFound)

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

	mockService := trainingsessinmock.NewMockService(ctrl)
	handler := New(mockService)

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(tsID)

		mockService.EXPECT().Restore(gomock.Any(), tsID).Return(nil)

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
		c.SetParamValues(tsID)

		mockService.EXPECT().Restore(gomock.Any(), tsID).Return(trainingsessionservice.ErrNotFound)

		// Act
		err := handler.Restore(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
