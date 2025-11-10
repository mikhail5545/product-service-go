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

package coursepart

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	coursepart "github.com/mikhail5545/product-service-go/internal/models/course_part"
	coursepartservice "github.com/mikhail5545/product-service-go/internal/services/course_part"
	coursepartmock "github.com/mikhail5545/product-service-go/internal/test/services/course_part_mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockPart := &coursepart.CoursePart{
			ID:   partID,
			Name: "Test Part",
		}

		expectedPart := &coursepart.CoursePart{
			ID:        partID,
			Name:      "Test Part",
			Published: false,
		}

		mockService.EXPECT().Get(gomock.Any(), partID).Return(mockPart, nil)

		// Act
		err := handler.Get(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponse := map[string]any{"course_part": expectedPart}

		expectedJSON, err := json.Marshal(expectedResponse)
		if err != nil {
			t.FailNow()
		}
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().Get(gomock.Any(), partID).Return(nil, coursepartservice.ErrNotFound)

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
		assert.Contains(t, rec.Body.String(), "Invalid course part ID")
	})
}

func TestHandler_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockPart := &coursepart.CoursePart{
			ID:   partID,
			Name: "Test Part",
		}

		expectedPart := &coursepart.CoursePart{
			ID:        partID,
			Name:      "Test Part",
			Published: false,
		}

		mockService.EXPECT().GetWithDeleted(gomock.Any(), partID).Return(mockPart, nil)

		// Act
		err := handler.GetWithDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponse := map[string]any{"course_part": expectedPart}

		expectedJSON, err := json.Marshal(expectedResponse)
		if err != nil {
			t.FailNow()
		}
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().GetWithDeleted(gomock.Any(), partID).Return(nil, coursepartservice.ErrNotFound)

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
		assert.Contains(t, rec.Body.String(), "Invalid course part ID")
	})
}

func TestHandler_GetWithUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockPart := &coursepart.CoursePart{
			ID:   partID,
			Name: "Test Part",
		}

		expectedPart := &coursepart.CoursePart{
			ID:        partID,
			Name:      "Test Part",
			Published: false,
		}

		mockService.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(mockPart, nil)

		// Act
		err := handler.GetWithUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponse := map[string]any{"course_part": expectedPart}

		expectedJSON, err := json.Marshal(expectedResponse)
		if err != nil {
			t.FailNow()
		}
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(nil, coursepartservice.ErrNotFound)

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
		assert.Contains(t, rec.Body.String(), "Invalid course part ID")
	})
}

func TestHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID_1 := uuid.New().String()
	partID_2 := uuid.New().String()
	courseID := uuid.New().String()

	mockParts := []coursepart.CoursePart{
		{
			ID:               partID_1,
			Name:             "Course part name 1",
			ShortDescription: "Course part name 2",
		},
		{
			ID:               partID_2,
			Name:             "Course part name 2",
			ShortDescription: "Course part name 2",
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":cid")
		c.SetParamValues(courseID)

		mockService.EXPECT().List(gomock.Any(), courseID, 2, 0).Return(mockParts, int64(2), nil)

		// Act
		err := handler.List(c)

		// Assert
		expectedResp := map[string]any{"course_parts": mockParts, "total": 2}
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedJSON, err := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":cid")
		c.SetParamValues(courseID)

		mockService.EXPECT().List(gomock.Any(), courseID, 2, 0).Return(nil, int64(0), coursepartservice.ErrNotFound)

		// Act
		err := handler.List(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("missing id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
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

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID_1 := uuid.New().String()
	partID_2 := uuid.New().String()
	courseID := uuid.New().String()

	mockParts := []coursepart.CoursePart{
		{
			ID:               partID_1,
			Name:             "Course part name 1",
			ShortDescription: "Course part name 2",
		},
		{
			ID:               partID_2,
			Name:             "Course part name 2",
			ShortDescription: "Course part name 2",
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":cid")
		c.SetParamValues(courseID)

		mockService.EXPECT().ListDeleted(gomock.Any(), courseID, 2, 0).Return(mockParts, int64(2), nil)

		// Act
		err := handler.ListDeleted(c)

		// Assert
		expectedResp := map[string]any{"course_parts": mockParts, "total": 2}
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedJSON, err := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":cid")
		c.SetParamValues(courseID)

		mockService.EXPECT().ListDeleted(gomock.Any(), courseID, 2, 0).Return(nil, int64(0), coursepartservice.ErrNotFound)

		// Act
		err := handler.ListDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("missing id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
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

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID_1 := uuid.New().String()
	partID_2 := uuid.New().String()
	courseID := uuid.New().String()

	mockParts := []coursepart.CoursePart{
		{
			ID:               partID_1,
			Name:             "Course part name 1",
			ShortDescription: "Course part name 2",
		},
		{
			ID:               partID_2,
			Name:             "Course part name 2",
			ShortDescription: "Course part name 2",
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":cid")
		c.SetParamValues(courseID)

		mockService.EXPECT().ListUnpublished(gomock.Any(), courseID, 2, 0).Return(mockParts, int64(2), nil)

		// Act
		err := handler.ListUnpublished(c)

		// Assert
		expectedResp := map[string]any{"course_parts": mockParts, "total": 2}
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedJSON, err := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":cid")
		c.SetParamValues(courseID)

		mockService.EXPECT().ListUnpublished(gomock.Any(), courseID, 2, 0).Return(nil, int64(0), coursepartservice.ErrNotFound)

		// Act
		err := handler.ListUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("missing id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
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

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	courseID := "c17081f3-4a56-4d00-b63e-f942537a702f"
	partID := "p17081f3-4a56-4d00-b63e-f942537a702f"

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		createReq := &coursepart.CreateRequest{
			Name:             "New Part",
			ShortDescription: "A new part",
			Number:           1,
		}
		jsonBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":cid")
		c.SetParamValues(courseID)

		// The handler will modify the request object, so we need to match it.
		expectedReq := *createReq
		expectedReq.CourseID = courseID

		mockResponse := &coursepart.CreateResponse{
			ID:       partID,
			CourseID: courseID,
		}

		mockService.EXPECT().Create(gomock.Any(), &expectedReq).Return(mockResponse, nil)

		// Act
		err := handler.Create(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		expectedJSON := `{"response":{"id":"p17081f3-4a56-4d00-b63e-f942537a702f","course_id":"c17081f3-4a56-4d00-b63e-f942537a702f"}}`
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	})

	t.Run("bind error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "bad json`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":cid")
		c.SetParamValues(courseID)

		// Act
		err := handler.Create(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid request JSON payload")
	})

	t.Run("validation error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		// Missing required 'Name' field
		createReq := &coursepart.CreateRequest{
			ShortDescription: "A new part",
			Number:           1,
		}
		jsonBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":cid")
		c.SetParamValues(courseID)

		mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, coursepartservice.ErrNotFound)

		// Act
		err := handler.Create(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		createReq := &coursepart.CreateRequest{
			Name:             "New Part",
			ShortDescription: "A new part",
			Number:           1,
		}
		jsonBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":cid")
		c.SetParamValues(courseID)

		mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, coursepartservice.ErrNotFound)

		// Act
		err := handler.Create(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().Delete(gomock.Any(), partID).Return(nil)

		// Act
		err := handler.Delete(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().Delete(gomock.Any(), partID).Return(coursepartservice.ErrNotFound)

		// Act
		err := handler.Delete(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// No param set

		c.SetParamNames(":id")
		c.SetParamValues("invalid-uuid")
		// Act
		err := handler.Delete(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid course part ID")
	})
}

func TestHandler_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().Publish(gomock.Any(), partID).Return(nil)

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
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().Publish(gomock.Any(), partID).Return(coursepartservice.ErrNotFound)

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
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues("invalid-uuid")
		// Act
		err := handler.Publish(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid course part ID")
	})
}

func TestHandler_Unpublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().Unpublish(gomock.Any(), partID).Return(nil)

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
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().Unpublish(gomock.Any(), partID).Return(coursepartservice.ErrNotFound)

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
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues("invalid-uuid")
		// Act
		err := handler.Unpublish(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid course part ID")
	})
}

func TestHandler_AddVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"
	videoID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		addVideoReq := &coursepart.AddVideoRequest{
			MUXVideoID: videoID,
		}
		jsonBody, _ := json.Marshal(addVideoReq)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		addVideoReq.ID = partID
		mockService.EXPECT().AddVideo(gomock.Any(), addVideoReq).Return(map[string]any{"mux_video_id": videoID}, nil)

		// Act
		err := handler.AddVideo(c)

		// Assert
		expectedResp := map[string]any{"updates": map[string]any{"mux_video_id": videoID}}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		addVideoReq := &coursepart.AddVideoRequest{
			MUXVideoID: videoID,
		}
		jsonBody, _ := json.Marshal(addVideoReq)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		addVideoReq.ID = partID
		mockService.EXPECT().AddVideo(gomock.Any(), addVideoReq).Return(nil, coursepartservice.ErrNotFound)

		// Act
		err := handler.AddVideo(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("missing id", func(t *testing.T) {
		// Arrange
		addVideoReq := &coursepart.AddVideoRequest{
			MUXVideoID: videoID,
		}
		jsonBody, _ := json.Marshal(addVideoReq)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.AddVideo(c)

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

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	t.Run("success", func(t *testing.T) {
		// Arrange
		newName := "New course part name"
		newShortDescription := "New course part short description"
		updateReq := &coursepart.UpdateRequest{
			Name:             &newName,
			ShortDescription: &newShortDescription,
		}
		jsonBody, _ := json.Marshal(updateReq)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		updates := map[string]any{"name": newName, "short_description": newShortDescription}
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(updates, nil)

		// Act
		err := handler.Update(c)

		// Assert
		expectedResp := map[string]any{"updates": updates}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		newName := "New course part name"
		newShortDescription := "New course part short description"
		updateReq := &coursepart.UpdateRequest{
			Name:             &newName,
			ShortDescription: &newShortDescription,
		}
		jsonBody, _ := json.Marshal(updateReq)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, coursepartservice.ErrNotFound)

		// Act
		err := handler.Update(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("bind error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "bad json`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		// Act
		err := handler.Update(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		// Arrange
		invalidName := "a" // too short
		updateReq := &coursepart.UpdateRequest{
			Name: &invalidName,
		}
		jsonBody, _ := json.Marshal(updateReq)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, coursepartservice.ErrNotFound)

		// Act
		err := handler.Update(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"name": "new name"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Act
		err := handler.Update(c)

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

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().DeletePermanent(gomock.Any(), partID).Return(nil)

		// Act
		err := handler.DeletePermanent(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().DeletePermanent(gomock.Any(), partID).Return(coursepartservice.ErrNotFound)

		// Act
		err := handler.DeletePermanent(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// No param set

		c.SetParamNames(":id")
		c.SetParamValues("invalid-uuid")
		// Act
		err := handler.DeletePermanent(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid course part ID")
	})
}

func TestHandler_Restore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursepartmock.NewMockService(ctrl)
	handler := New(mockService)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().Restore(gomock.Any(), partID).Return(nil)

		// Act
		err := handler.Restore(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(partID)

		mockService.EXPECT().Restore(gomock.Any(), partID).Return(coursepartservice.ErrNotFound)

		// Act
		err := handler.Restore(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// No param set

		c.SetParamNames(":id")
		c.SetParamValues("invalid-uuid")
		// Act
		err := handler.Restore(c)

		// Assert
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid course part ID")
	})
}
