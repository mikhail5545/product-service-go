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

package course

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	course "github.com/mikhail5545/product-service-go/internal/models/course"
	courseservice "github.com/mikhail5545/product-service-go/internal/services/course"
	coursemock "github.com/mikhail5545/product-service-go/internal/test/services/course_mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID := uuid.New().String()

	mockCourseDetails := &course.CourseDetails{
		Course: course.Course{
			ID:               courseID,
			Name:             "Course name",
			ShortDescription: "Short course description",
		},
		Price:     44.44,
		ProductID: uuid.New().String(),
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(courseID)

		mockService.EXPECT().Get(gomock.Any(), courseID).Return(mockCourseDetails, nil)

		// Act
		err := handler.Get(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"course_details": mockCourseDetails}
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
		c.SetParamValues(courseID)

		serviceErr := &courseservice.Error{
			Msg:  "Course not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Get(gomock.Any(), courseID).Return(nil, serviceErr)

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
		assert.Contains(t, rec.Body.String(), "Invalid course ID")
	})
}

func TestHandler_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID := uuid.New().String()

	mockCourseDetails := &course.CourseDetails{
		Course: course.Course{
			ID:               courseID,
			Name:             "Course name",
			ShortDescription: "Short course description",
		},
		Price:     44.44,
		ProductID: uuid.New().String(),
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(courseID)

		mockService.EXPECT().GetWithDeleted(gomock.Any(), courseID).Return(mockCourseDetails, nil)

		// Act
		err := handler.GetWithDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"course_details": mockCourseDetails}
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
		c.SetParamValues(courseID)

		serviceErr := &courseservice.Error{
			Msg:  "Course not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().GetWithDeleted(gomock.Any(), courseID).Return(nil, serviceErr)

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
		assert.Contains(t, rec.Body.String(), "Invalid course ID")
	})
}

func TestHandler_GetWithUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID := uuid.New().String()

	mockCourseDetails := &course.CourseDetails{
		Course: course.Course{
			ID:               courseID,
			Name:             "Course name",
			ShortDescription: "Short course description",
		},
		Price:     44.44,
		ProductID: uuid.New().String(),
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(courseID)

		mockService.EXPECT().GetWithUnpublished(gomock.Any(), courseID).Return(mockCourseDetails, nil)

		// Act
		err := handler.GetWithUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"course_details": mockCourseDetails}
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
		c.SetParamValues(courseID)

		serviceErr := &courseservice.Error{
			Msg:  "Course not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), courseID).Return(nil, serviceErr)

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
		assert.Contains(t, rec.Body.String(), "Invalid course ID")
	})
}

func TestHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID_1 := uuid.New().String()
	courseID_2 := uuid.New().String()

	mockCourseDetails := []course.CourseDetails{
		{
			Course: course.Course{
				ID:               courseID_1,
				Name:             "Course 1 name",
				ShortDescription: "Course 1 short description",
			},
			Price:     55.55,
			ProductID: uuid.New().String(),
		},
		{
			Course: course.Course{
				ID:               courseID_2,
				Name:             "Course 2 name",
				ShortDescription: "Course 2 short description",
			},
			Price:     52.22,
			ProductID: uuid.New().String(),
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().List(gomock.Any(), 2, 0).Return(mockCourseDetails, int64(2), nil)

		// Act
		err := handler.List(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"course_details": mockCourseDetails, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		serviceErr := &courseservice.Error{
			Msg:  "Course not found",
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
func TestHandler_List_InvalidParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?limit=abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.List(c)
	assert.Error(t, err)
	e.HTTPErrorHandler(err, c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_ListDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID_1 := uuid.New().String()
	courseID_2 := uuid.New().String()

	mockCourseDetails := []course.CourseDetails{
		{
			Course: course.Course{
				ID:               courseID_1,
				Name:             "Course 1 name",
				ShortDescription: "Course 1 short description",
			},
			Price:     55.55,
			ProductID: uuid.New().String(),
		},
		{
			Course: course.Course{
				ID:               courseID_2,
				Name:             "Course 2 name",
				ShortDescription: "Course 2 short description",
			},
			Price:     52.22,
			ProductID: uuid.New().String(),
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().ListDeleted(gomock.Any(), 2, 0).Return(mockCourseDetails, int64(2), nil)

		// Act
		err := handler.ListDeleted(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"course_details": mockCourseDetails, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		serviceErr := &courseservice.Error{
			Msg:  "Course not found",
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
func TestHandler_ListDeleted_InvalidParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?limit=abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.ListDeleted(c)
	assert.Error(t, err)
	e.HTTPErrorHandler(err, c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_ListUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID_1 := uuid.New().String()
	courseID_2 := uuid.New().String()

	mockCourseDetails := []course.CourseDetails{
		{
			Course: course.Course{
				ID:               courseID_1,
				Name:             "Course 1 name",
				ShortDescription: "Course 1 short description",
			},
			Price:     55.55,
			ProductID: uuid.New().String(),
		},
		{
			Course: course.Course{
				ID:               courseID_2,
				Name:             "Course 2 name",
				ShortDescription: "Course 2 short description",
			},
			Price:     52.22,
			ProductID: uuid.New().String(),
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().ListUnpublished(gomock.Any(), 2, 0).Return(mockCourseDetails, int64(2), nil)

		// Act
		err := handler.ListUnpublished(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResp := map[string]any{"course_details": mockCourseDetails, "total": 2}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), rec.Body.String())
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		serviceErr := &courseservice.Error{
			Msg:  "Course not found",
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
func TestHandler_ListUnpublished_InvalidParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?limit=abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.ListUnpublished(c)
	assert.Error(t, err)
	e.HTTPErrorHandler(err, c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID := uuid.New().String()
	productID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		createReq := &course.CreateRequest{
			Name:             "Course name",
			ShortDescription: "Course short description",
			Topic:            "Course topic",
			Price:            33.33,
		}
		reqJSON, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		createResp := &course.CreateResponse{ID: courseID, ProductID: productID}
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
		createReq := course.CreateRequest{
			Name:             "Course name",
			ShortDescription: "Course short description",
			Topic:            "Course topic",
			Price:            33.33,
		}
		reqJSON, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		serviceErr := &courseservice.Error{
			Msg:  "Failed to create course",
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
		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, c)
		}
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(courseID)

		mockService.EXPECT().Publish(gomock.Any(), courseID).Return(nil)

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
		c.SetParamValues(courseID)

		serviceErr := &courseservice.Error{
			Msg:  "Course not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Publish(gomock.Any(), courseID).Return(serviceErr)

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

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(courseID)

		mockService.EXPECT().Unpublish(gomock.Any(), courseID).Return(nil)

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
		c.SetParamValues(courseID)

		serviceErr := &courseservice.Error{
			Msg:  "Course not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Unpublish(gomock.Any(), courseID).Return(serviceErr)

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

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		newName := "New course name"
		newLongDescription := "New course long description"
		updateReq := course.UpdateRequest{
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
		c.SetParamValues(courseID)

		updateReq.ID = courseID
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
		newName := "New course name"
		newLongDescription := "New course long description"
		updateReq := course.UpdateRequest{
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
		newName := "New course name"
		newLongDescription := "New course long description"
		updateReq := course.UpdateRequest{
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
		c.SetParamValues(courseID)

		serviceErr := &courseservice.Error{
			Msg:  "Failed to update course",
			Code: http.StatusInternalServerError,
		}
		updateReq.ID = courseID
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

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(courseID)

		mockService.EXPECT().Delete(gomock.Any(), courseID).Return(nil)

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
		c.SetParamValues(courseID)

		serviceErr := &courseservice.Error{
			Msg:  "Course not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Delete(gomock.Any(), courseID).Return(serviceErr)

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

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(courseID)

		mockService.EXPECT().DeletePermanent(gomock.Any(), courseID).Return(nil)

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
		c.SetParamValues(courseID)

		serviceErr := &courseservice.Error{
			Msg:  "Course not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().DeletePermanent(gomock.Any(), courseID).Return(serviceErr)

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

	mockService := coursemock.NewMockService(ctrl)
	handler := New(mockService)

	courseID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(":id")
		c.SetParamValues(courseID)

		mockService.EXPECT().Restore(gomock.Any(), courseID).Return(nil)

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
		c.SetParamValues(courseID)

		serviceErr := &courseservice.Error{
			Msg:  "Course not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Restore(gomock.Any(), courseID).Return(serviceErr)

		// Act
		err := handler.Restore(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
