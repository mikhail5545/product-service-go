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
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/mikhail5545/product-service-go/internal/models/course"
	"github.com/mikhail5545/product-service-go/internal/models/product"
	coursemock "github.com/mikhail5545/product-service-go/internal/test/database/course_mock"
	coursepartmock "github.com/mikhail5545/product-service-go/internal/test/database/course_part_mock"
	productmock "github.com/mikhail5545/product-service-go/internal/test/database/product_mock"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	expectedCourse := &course.Course{
		ID:               courseID,
		Name:             "Test Course",
		ShortDescription: "A short description",
	}

	expectedProduct := &product.Product{
		ID:        "product-uuid",
		Price:     99.99,
		DetailsID: courseID,
	}

	expectedDetails := &course.CourseDetails{
		Course:    *expectedCourse,
		Price:     expectedProduct.Price,
		ProductID: expectedProduct.ID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().Get(context.Background(), courseID).Return(expectedCourse, nil)
		mockProductRepo.EXPECT().GetByDetailsID(context.Background(), courseID).Return(expectedProduct, nil)

		// Act
		details, err := testService.Get(context.Background(), courseID)

		// Assert
		if err != nil {
			t.Errorf("Get() error = %v, wantErr %v", err, nil)
			return
		}
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("Get() got = %v, want %v", details, expectedDetails)
		}
	})

	t.Run("course not found", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().Get(context.Background(), courseID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.Get(context.Background(), courseID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Get() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("course product not found", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().Get(context.Background(), courseID).Return(expectedCourse, nil)
		mockProductRepo.EXPECT().GetByDetailsID(context.Background(), courseID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.Get(context.Background(), courseID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Get() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"

		// Act
		_, err := testService.Get(context.Background(), invalidID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error for invalid UUID, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Get() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})
}

func TestService_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	expectedCourse := &course.Course{
		ID:               courseID,
		Name:             "Test Course",
		ShortDescription: "A short description",
	}

	expectedProduct := &product.Product{
		ID:        "product-uuid",
		Price:     99.99,
		DetailsID: courseID,
	}

	expectedDetails := &course.CourseDetails{
		Course:    *expectedCourse,
		Price:     expectedProduct.Price,
		ProductID: expectedProduct.ID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().GetWithDeleted(gomock.Any(), courseID).Return(expectedCourse, nil)
		mockProductRepo.EXPECT().GetWithDeletedByDetailsID(gomock.Any(), courseID).Return(expectedProduct, nil)

		// Act
		details, err := testService.GetWithDeleted(context.Background(), courseID)

		// Assert
		if err != nil {
			t.Errorf("Get() error = %v, wantErr %v", err, nil)
			return
		}
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("Get() got = %v, want %v", details, expectedDetails)
		}
	})

	t.Run("course not found", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().GetWithDeleted(gomock.Any(), courseID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), courseID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Get() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("course product not found", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().GetWithDeleted(gomock.Any(), courseID).Return(expectedCourse, nil)
		mockProductRepo.EXPECT().GetWithDeletedByDetailsID(gomock.Any(), courseID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), courseID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Get() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"

		// Act
		_, err := testService.GetWithDeleted(context.Background(), invalidID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error for invalid UUID, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Get() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})
}

func TestService_GetWithUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	expectedCourse := &course.Course{
		ID:               courseID,
		Name:             "Test Course",
		ShortDescription: "A short description",
	}

	expectedProduct := &product.Product{
		ID:        "product-uuid",
		Price:     99.99,
		DetailsID: courseID,
	}

	expectedDetails := &course.CourseDetails{
		Course:    *expectedCourse,
		Price:     expectedProduct.Price,
		ProductID: expectedProduct.ID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().GetWithUnpublished(gomock.Any(), courseID).Return(expectedCourse, nil)
		mockProductRepo.EXPECT().GetWithUnpublishedByDetailsID(gomock.Any(), courseID).Return(expectedProduct, nil)

		// Act
		details, err := testService.GetWithUnpublished(context.Background(), courseID)

		// Assert
		if err != nil {
			t.Errorf("Get() error = %v, wantErr %v", err, nil)
			return
		}
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("Get() got = %v, want %v", details, expectedDetails)
		}
	})

	t.Run("course not found", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().GetWithUnpublished(gomock.Any(), courseID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), courseID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Get() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("course product not found", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().GetWithUnpublished(gomock.Any(), courseID).Return(expectedCourse, nil)
		mockProductRepo.EXPECT().GetWithUnpublishedByDetailsID(gomock.Any(), courseID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), courseID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Get() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), invalidID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error for invalid UUID, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Get() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})
}

func TestService_GetReduced(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	expectedCourse := &course.Course{
		ID:               courseID,
		Name:             "Test Course",
		ShortDescription: "A short description",
	}

	expectedProduct := &product.Product{
		ID:        "product-uuid",
		Price:     99.99,
		DetailsID: courseID,
	}

	expectedDetails := &course.CourseDetails{
		Course:    *expectedCourse,
		Price:     expectedProduct.Price,
		ProductID: expectedProduct.ID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().GetReduced(context.Background(), courseID).Return(expectedCourse, nil)
		mockProductRepo.EXPECT().GetByDetailsID(context.Background(), courseID).Return(expectedProduct, nil)

		// Act
		details, err := testService.GetReduced(context.Background(), courseID)

		// Assert
		if err != nil {
			t.Errorf("Get() error = %v, wantErr %v", err, nil)
			return
		}
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("Get() got = %v, want %v", details, expectedDetails)
		}
	})

	t.Run("course not found", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().GetReduced(context.Background(), courseID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetReduced(context.Background(), courseID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Get() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("course product not found", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().GetReduced(context.Background(), courseID).Return(expectedCourse, nil)
		mockProductRepo.EXPECT().GetByDetailsID(context.Background(), courseID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetReduced(context.Background(), courseID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Get() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"

		// Act
		_, err := testService.GetReduced(context.Background(), invalidID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error for invalid UUID, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Get() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})
}

func TestService_GetReducedWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	expectedCourse := &course.Course{
		ID:               courseID,
		Name:             "Test Course",
		ShortDescription: "A short description",
	}

	expectedProduct := &product.Product{
		ID:        "product-uuid",
		Price:     99.99,
		DetailsID: courseID,
	}

	expectedDetails := &course.CourseDetails{
		Course:    *expectedCourse,
		Price:     expectedProduct.Price,
		ProductID: expectedProduct.ID,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().GetReducedWithDeleted(context.Background(), courseID).Return(expectedCourse, nil)
		mockProductRepo.EXPECT().GetWithDeletedByDetailsID(context.Background(), courseID).Return(expectedProduct, nil)

		// Act
		details, err := testService.GetReducedWithDeleted(context.Background(), courseID)

		// Assert
		if err != nil {
			t.Errorf("Get() error = %v, wantErr %v", err, nil)
			return
		}
		if !reflect.DeepEqual(details, expectedDetails) {
			t.Errorf("Get() got = %v, want %v", details, expectedDetails)
		}
	})

	t.Run("course not found", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().GetReducedWithDeleted(context.Background(), courseID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetReducedWithDeleted(context.Background(), courseID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Get() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("course product not found", func(t *testing.T) {
		// Arrange
		mockCourseRepo.EXPECT().GetReducedWithDeleted(context.Background(), courseID).Return(expectedCourse, nil)
		mockProductRepo.EXPECT().GetWithDeletedByDetailsID(context.Background(), courseID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetReducedWithDeleted(context.Background(), courseID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Get() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"

		// Act
		_, err := testService.GetReducedWithDeleted(context.Background(), invalidID)

		// Assert
		if err == nil {
			t.Errorf("Get() expected an error for invalid UUID, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Get() expected a custom error type, got %T", err)
			return
		}
		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Get() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})
}

func TestService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	course1ID := "d17081f3-4a56-4d00-b63e-f942537a702f"
	course2ID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"

	mockCourses := []course.Course{
		{
			ID:               course1ID,
			Name:             "Course name 1",
			ShortDescription: "Short course description 1",
		},
		{
			ID:               course2ID,
			Name:             "Course name 2",
			ShortDescription: "Short course description 2",
		},
	}

	mockProducts := []product.Product{
		{
			ID:        "prod-1",
			Price:     99.99,
			DetailsID: course1ID,
		},
		{
			ID:        "prod-2",
			Price:     199.99,
			DetailsID: course2ID,
		},
	}

	expectedDetails := []course.CourseDetails{
		{
			Course:    mockCourses[0],
			Price:     mockProducts[0].Price,
			ProductID: mockProducts[0].ID,
		},
		{
			Course:    mockCourses[1],
			Price:     mockProducts[1].Price,
			ProductID: mockProducts[1].ID,
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockCourseRepo.EXPECT().List(gomock.Any(), limit, offset).Return(mockCourses, nil)
		mockCourseRepo.EXPECT().Count(gomock.Any()).Return(int64(2), nil)
		mockProductRepo.EXPECT().SelectByDetailsIDs(gomock.Any(), []string{course1ID, course2ID}, "id", "price", "details_id").Return(mockProducts, nil)

		// Act
		courses, total, err := testService.List(context.Background(), limit, offset)

		// Assert
		if err != nil {
			t.Errorf("List() error = %v, wantErr %v", err, nil)
			return
		}
		if total != 2 {
			t.Errorf("List() got = %v, want %v", total, 2)
		}
		// The order from SelectByDetailsIDs is not guaranteed, so we need to compare them in a way that ignores order.
		if len(courses) != len(expectedDetails) {
			t.Errorf("List() got %d items, want %d", len(courses), len(expectedDetails))
		}
	})

	t.Run("db error on count", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockCourseRepo.EXPECT().List(gomock.Any(), limit, offset).Return(mockCourses, nil)
		mockCourseRepo.EXPECT().Count(gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, _, err := testService.List(context.Background(), limit, offset)

		// Assert
		if err == nil {
			t.Errorf("List() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) || customErr.Code != http.StatusInternalServerError {
			t.Errorf("List() expected an internal server error, got %v", err)
		}
		if errors.As(err, &customErr) && customErr.Msg != "Failed to get courses count" {
			t.Errorf("List() expected a 'Failed to get courses count' error message, got %s", customErr.Msg)
		}
	})

	t.Run("db error on list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockCourseRepo.EXPECT().List(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.List(context.Background(), limit, offset)

		// Assert
		if err == nil {
			t.Errorf("List() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) || customErr.Code != http.StatusInternalServerError {
			t.Errorf("List() expected an internal server error, got %v", err)
		}
	})
}

func TestService_ListDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	course1ID := "d17081f3-4a56-4d00-b63e-f942537a702f"
	course2ID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"

	mockCourses := []course.Course{
		{
			ID:               course1ID,
			Name:             "Course name 1",
			ShortDescription: "Short course description 1",
		},
		{
			ID:               course2ID,
			Name:             "Course name 2",
			ShortDescription: "Short course description 2",
		},
	}

	mockProducts := []product.Product{
		{
			ID:        "prod-1",
			Price:     99.99,
			DetailsID: course1ID,
		},
		{
			ID:        "prod-2",
			Price:     199.99,
			DetailsID: course2ID,
		},
	}

	expectedDetails := []course.CourseDetails{
		{
			Course:    mockCourses[0],
			Price:     mockProducts[0].Price,
			ProductID: mockProducts[0].ID,
		},
		{
			Course:    mockCourses[1],
			Price:     mockProducts[1].Price,
			ProductID: mockProducts[1].ID,
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockCourseRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(mockCourses, nil)
		mockCourseRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(2), nil)
		mockProductRepo.EXPECT().SelectWithDeletedByDetailsIDs(gomock.Any(), []string{course1ID, course2ID}, "id", "price", "details_id").Return(mockProducts, nil)

		// Act
		courses, total, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		if err != nil {
			t.Errorf("ListDeleted() error = %v, wantErr %v", err, nil)
			return
		}
		if total != 2 {
			t.Errorf("ListDeleted() got = %v, want %v", total, 2)
		}
		if len(courses) != len(expectedDetails) {
			t.Errorf("ListDeleted() got %d items, want %d", len(courses), len(expectedDetails))
		}
	})

	t.Run("db error on count", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockCourseRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(mockCourses, nil)
		mockCourseRepo.EXPECT().CountDeleted(gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, _, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		if err == nil {
			t.Errorf("ListDeleted() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) || customErr.Code != http.StatusInternalServerError {
			t.Errorf("ListDeleted() expected an internal server error, got %v", err)
		}
		if errors.As(err, &customErr) && customErr.Msg != "Failed to get courses count" {
			t.Errorf("ListDeleted() expected a 'Failed to get courses count' error message, got %s", customErr.Msg)
		}
	})

	t.Run("db error on list", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockCourseRepo.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListDeleted(context.Background(), limit, offset)

		// Assert
		if err == nil {
			t.Errorf("ListDeleted() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) || customErr.Code != http.StatusInternalServerError {
			t.Errorf("ListDeleted() expected an internal server error, got %v", err)
		}
	})
}

func TestService_ListUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	course1ID := "d17081f3-4a56-4d00-b63e-f942537a702f"
	course2ID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"

	mockCourses := []course.Course{
		{ID: course1ID, Name: "Unpublished Course 1"},
		{ID: course2ID, Name: "Unpublished Course 2"},
	}

	mockProducts := []product.Product{
		{ID: "prod-1", Price: 99.99, DetailsID: course1ID},
		{ID: "prod-2", Price: 199.99, DetailsID: course2ID},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockCourseRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(mockCourses, nil)
		mockCourseRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(2), nil)
		mockProductRepo.EXPECT().SelectWithUnpublishedByDetailsIDs(gomock.Any(), []string{course1ID, course2ID}, "id", "price", "details_id").Return(mockProducts, nil)

		// Act
		courses, total, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		if err != nil {
			t.Errorf("ListUnpublished() error = %v, wantErr %v", err, nil)
			return
		}
		if total != 2 {
			t.Errorf("ListUnpublished() got total = %v, want %v", total, 2)
		}
		if len(courses) != 2 {
			t.Errorf("ListUnpublished() got %d items, want %d", len(courses), 2)
		}
	})

	t.Run("db error on count", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockCourseRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(mockCourses, nil)
		mockCourseRepo.EXPECT().CountUnpublished(gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, _, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		if err == nil {
			t.Errorf("ListUnpublished() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) || customErr.Code != http.StatusInternalServerError {
			t.Errorf("ListUnpublished() expected an internal server error, got %v", err)
		}
		if errors.As(err, &customErr) && customErr.Msg != "Failed to get courses count" {
			t.Errorf("ListUnpublished() expected a 'Failed to get courses count' error message, got %s", customErr.Msg)
		}
	})

	t.Run("db error on list unpublished", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockCourseRepo.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListUnpublished(context.Background(), limit, offset)

		// Assert
		if err == nil {
			t.Errorf("ListUnpublished() expected an error, but got nil")
			return
		}
		var customErr *Error
		if !errors.As(err, &customErr) || customErr.Code != http.StatusInternalServerError {
			t.Errorf("ListUnpublished() expected an internal server error, got %v", err)
		}
	})
}

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	createReq := &course.CreateRequest{
		Name:             "Course name",
		ShortDescription: "Course short description",
		Topic:            "Course topic",
		Price:            99.99,
		AccessDuration:   30,
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
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		var createdCourse *course.Course
		mockTxCourseRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, c *course.Course) {
				createdCourse = c
			}).Return(nil)

		var createdProduct *product.Product
		mockTxProductRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, p *product.Product) {
				createdProduct = p
			}).Return(nil)

		// Act
		resp, err := testService.Create(context.Background(), createReq)
		if err != nil {
			t.Fatalf("Create() error = %v, wantErr %v", err, nil)
		}

		// Assert
		if _, err := uuid.Parse(createdCourse.ID); err != nil {
			t.Errorf("expected course.ID to be a valid UUID, got %s", createdCourse.ID)
		}
		if createdCourse.Name != createReq.Name {
			t.Errorf("course.Name = %s, want %s", createdCourse.Name, createReq.Name)
		}

		if _, err := uuid.Parse(createdProduct.ID); err != nil {
			t.Errorf("expected product.ID to be a valid UUID, got %s", createdProduct.ID)
		}
		if createdProduct.DetailsID != createdCourse.ID {
			t.Errorf("product.DetailsID = %s, want %s", createdProduct.DetailsID, createdCourse.ID)
		}
		if createdProduct.Price != createReq.Price {
			t.Errorf("product.Price = %f, want %f", createdProduct.Price, createReq.Price)
		}

		if resp.ID != createdCourse.ID {
			t.Errorf("response ID = %s, want %s", resp.ID, createdCourse.ID)
		}
		if resp.ProductID != createdProduct.ID {
			t.Errorf("response ProductID = %s, want %s", resp.ProductID, createdProduct.ID)
		}
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		// Act
		// Invalid price and empty topic
		_, err = testService.Create(context.Background(), &course.CreateRequest{Name: "Name", ShortDescription: "ShortDescription", Price: -2.3, Topic: ""})

		// Assert
		if err == nil {
			t.Errorf("Create() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Create() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Create() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		dbErr := errors.New("database error")
		mockTxCourseRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)

		// Act
		_, err = testService.Create(context.Background(), createReq)

		// Assert
		if err == nil {
			t.Errorf("Create() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Create() expected a custom error type, got %T", err)
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

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

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
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxCourseRepo.EXPECT().SetInStock(gomock.Any(), courseID, true).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), courseID, true).Return(int64(1), nil)

		// Act
		err = testService.Publish(context.Background(), courseID)

		// Assert
		if err != nil {
			t.Fatalf("Publish() error = %v, wantErr %v", err, nil)
		}
	})

	t.Run("invalid course UUID", func(t *testing.T) {
		// Act
		err := testService.Publish(context.Background(), "Invalid-UUID")

		// Assert
		if err == nil {
			t.Errorf("Publish() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Publish() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Publish() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})
}

func TestService_Unpublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

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
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxCourseRepo.EXPECT().SetInStock(gomock.Any(), courseID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), courseID, false).Return(int64(1), nil)
		mockTxPartRepo.EXPECT().SetPublishedByCourseID(gomock.Any(), courseID, false).Return(int64(1), nil)

		// Act
		err = testService.Unpublish(context.Background(), courseID)

		// Assert
		if err != nil {
			t.Fatalf("Unpublish() error = %v, wantErr %v", err, nil)
		}
	})

	t.Run("invalid course UUID", func(t *testing.T) {
		// Act
		err := testService.Unpublish(context.Background(), "Invalid-UUID")

		// Assert
		if err == nil {
			t.Errorf("Unpublish() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Unpublish() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Unpublish() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})
}

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	newName := "New course name"
	newShortDescription := "New course description"
	newPrice := float32(192.33)
	newTags := []string{"course", "tags", "new"}

	mockCourse := &course.Course{
		ID:               courseID,
		Name:             "Old course name",
		ShortDescription: "Old course description",
	}

	mockProduct := &product.Product{
		ID:          "product-ID",
		Price:       33.4,
		DetailsID:   courseID,
		DetailsType: "course",
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
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxCourseRepo.EXPECT().Get(gomock.Any(), courseID).Return(mockCourse, nil)
		mockTxProductRepo.EXPECT().GetByDetailsID(gomock.Any(), courseID).Return(mockProduct, nil)

		var courseUpdates map[string]any
		mockTxCourseRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, _ *course.Course, u map[string]any) {
				courseUpdates = u
			})

		var productUpdates map[string]any
		mockTxProductRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, _ *product.Product, u map[string]any) {
				productUpdates = u
			})

		// Act
		updates, err := testService.Update(context.Background(), &course.UpdateRequest{
			ID:               courseID,
			Name:             &newName,
			ShortDescription: &newShortDescription,
			Tags:             newTags,
			Price:            &newPrice,
		})

		// Assert
		if err != nil {
			t.Fatalf("Update() error = %v, wantErr %v", err, nil)
		}

		courseUpdatesFromResp, ok := updates["course"].(map[string]any)
		if !ok {
			t.Fatalf("response does not contain 'course' updates")
		}

		if name, ok := courseUpdatesFromResp["name"].(string); !ok || name != newName {
			t.Errorf("course.Name in response = %v, want %s", courseUpdatesFromResp["name"], newName)
		}
		if shortDescription, ok := courseUpdatesFromResp["short_description"].(string); !ok || shortDescription != newShortDescription {
			t.Errorf("course.ShortDescription in response = %v, want %s", courseUpdatesFromResp["short_description"], newShortDescription)
		}

		productUpdatesFromResp, ok := updates["product"].(map[string]any)
		if !ok {
			t.Fatalf("response does not contain 'product' updates")
		}

		if price, ok := productUpdatesFromResp["price"].(float32); !ok || price != newPrice {
			t.Errorf("product.Price in response = %v, want %f", productUpdatesFromResp["price"], newPrice)
		}

		// Check what was passed to the mock repo update functions
		if name, ok := courseUpdates["name"].(string); !ok || name != newName {
			t.Errorf("course.Name passed to repo = %s, want %s", name, newName)
		}
		if shortDescription, ok := courseUpdates["short_description"].(string); !ok || shortDescription != newShortDescription {
			t.Errorf("course.ShortDescription passed to repo = %s, want %s", shortDescription, newShortDescription)
		}
		if tags, ok := courseUpdates["tags"].([]string); !ok || !reflect.DeepEqual(tags, newTags) {
			t.Errorf("course.Tags passed to repo = %v, want %v", tags, newTags)
		}
		if price, ok := productUpdates["price"].(float32); !ok || price != newPrice {
			t.Errorf("product.Price passed to repo = %f, want %f", price, newPrice)
		}
	})

	t.Run("invalid request payload", func(t *testing.T) {
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		invalidName := "1somename" // must start with a letter
		_, err := testService.Update(context.Background(), &course.UpdateRequest{
			ID:   courseID,
			Name: &invalidName,
		})
		if err == nil {
			t.Errorf("Update() expected an error, but got nil")
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Update() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Update() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})

	t.Run("course not found", func(t *testing.T) {
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxCourseRepo.EXPECT().Get(gomock.Any(), courseID).Return(nil, gorm.ErrRecordNotFound)

		_, err := testService.Update(context.Background(), &course.UpdateRequest{
			ID:   courseID,
			Name: &newName,
		})
		if err == nil {
			t.Errorf("Update() expected an error, but got nil")
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Update() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Update() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)

		mockTxCourseRepo.EXPECT().Get(gomock.Any(), courseID).Return(mockCourse, nil)
		mockTxProductRepo.EXPECT().GetByDetailsID(gomock.Any(), courseID).Return(mockProduct, nil)

		dbErr := errors.New("database error")
		mockTxCourseRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), dbErr)

		_, err := testService.Update(context.Background(), &course.UpdateRequest{
			ID:   courseID,
			Name: &newName,
		})
		if err == nil {
			t.Errorf("Update() expected an error, but got nil")
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Update() expected a custom error type, got %T", err)
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

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

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
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxCourseRepo.EXPECT().SetInStock(gomock.Any(), courseID, false).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().SetInStockByDetailsID(gomock.Any(), courseID, false).Return(int64(1), nil)
		mockTxPartRepo.EXPECT().SetPublishedByCourseID(gomock.Any(), courseID, false).Return(int64(1), nil)

		mockTxCourseRepo.EXPECT().Delete(gomock.Any(), courseID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().DeleteByDetailsID(gomock.Any(), courseID).Return(int64(1), nil)
		mockTxPartRepo.EXPECT().DeleteByCourseID(gomock.Any(), courseID).Return(int64(1), nil)

		err = testService.Delete(context.Background(), courseID)

		if err != nil {
			t.Fatalf("Delete() error = %v, wantErr %v", err, nil)
		}
	})

	t.Run("invalid course UUID", func(t *testing.T) {
		err := testService.Delete(context.Background(), "Invalid-UUID")

		if err == nil {
			t.Errorf("Delete() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Delete() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Delete() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})

	t.Run("course not found", func(t *testing.T) {
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxCourseRepo.EXPECT().SetInStock(gomock.Any(), courseID, false).Return(int64(0), nil)

		err = testService.Delete(context.Background(), courseID)

		if err == nil {
			t.Errorf("Delete() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Delete() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Delete() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		dbErr := errors.New("database error")
		mockTxCourseRepo.EXPECT().SetInStock(gomock.Any(), courseID, false).Return(int64(0), dbErr)

		err = testService.Delete(context.Background(), courseID)

		if err == nil {
			t.Errorf("Delete() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Delete() expected a custom error type, got %T", err)
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

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

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
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxCourseRepo.EXPECT().DeletePermanent(gomock.Any(), courseID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().DeletePermanentByDetailsID(gomock.Any(), courseID).Return(int64(1), nil)
		mockTxPartRepo.EXPECT().DeletePermanentByCourseID(gomock.Any(), courseID).Return(int64(1), nil)

		err := testService.DeletePermanent(context.Background(), courseID)
		if err != nil {
			t.Fatalf("DeletePermanent() error = %v, wantErr %v", err, nil)
		}
	})

	t.Run("invalid course UUID", func(t *testing.T) {
		err := testService.DeletePermanent(context.Background(), "Invalid-UUID")

		if err == nil {
			t.Errorf("DeletePermanent() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("DeletePermanent() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusBadRequest {
			t.Errorf("DeletePermanent() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxCourseRepo.EXPECT().DeletePermanent(gomock.Any(), courseID).Return(int64(0), nil)

		err := testService.DeletePermanent(context.Background(), courseID)
		if err == nil {
			t.Errorf("DeletePermanent() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("DeletePermanent() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("DeletePermanent() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("db error", func(t *testing.T) {
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		dbErr := errors.New("database error")
		mockTxCourseRepo.EXPECT().DeletePermanent(gomock.Any(), courseID).Return(int64(0), dbErr)

		err := testService.DeletePermanent(context.Background(), courseID)
		if err == nil {
			t.Errorf("DeletePermanent() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("DeletePermanent() expected a custom error type, got %T", err)
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

	mockCourseRepo := coursemock.NewMockRepository(ctrl)
	mockProductRepo := productmock.NewMockRepository(ctrl)
	mockPartRepo := coursepartmock.NewMockRepository(ctrl)

	testService := &Service{
		CourseRepo:  mockCourseRepo,
		ProductRepo: mockProductRepo,
		PartRepo:    mockPartRepo,
	}

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

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
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxCourseRepo.EXPECT().Restore(gomock.Any(), courseID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().RestoreByDetailsID(gomock.Any(), courseID).Return(int64(1), nil)
		mockTxPartRepo.EXPECT().RestoreByCourseID(gomock.Any(), courseID).Return(int64(1), nil)

		err := testService.Restore(context.Background(), courseID)
		if err != nil {
			t.Fatalf("Restore() error = %v, wantErr %v", err, nil)
		}
	})

	t.Run("invalid course UUID", func(t *testing.T) {
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxCourseRepo.EXPECT().Restore(gomock.Any(), courseID).Return(int64(1), nil)
		mockTxProductRepo.EXPECT().RestoreByDetailsID(gomock.Any(), courseID).Return(int64(1), nil)
		mockTxPartRepo.EXPECT().RestoreByCourseID(gomock.Any(), courseID).Return(int64(1), nil)

		err := testService.Restore(context.Background(), "Invalid-UUID")
		if err == nil {
			t.Errorf("Restore() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Restore() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusBadRequest {
			t.Errorf("Restore() expected status code %d, got %d", http.StatusBadRequest, customErr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxCourseRepo.EXPECT().Restore(gomock.Any(), courseID).Return(int64(0), nil)

		err := testService.Restore(context.Background(), courseID)
		if err == nil {
			t.Errorf("Restore() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Restore() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusNotFound {
			t.Errorf("Restore() expected status code %d, got %d", http.StatusNotFound, customErr.Code)
		}
	})

	t.Run("database error", func(t *testing.T) {
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)
		mockTxProductRepo := productmock.NewMockRepository(ctrl)
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockCourseRepo.EXPECT().DB().Return(db).AnyTimes()
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)
		mockProductRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxProductRepo)
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		dbErr := errors.New("database error")
		mockTxCourseRepo.EXPECT().Restore(gomock.Any(), courseID).Return(int64(0), dbErr)

		err := testService.Restore(context.Background(), courseID)
		if err == nil {
			t.Errorf("Restore() expected an error, but got nil")
			return
		}

		var customErr *Error
		if !errors.As(err, &customErr) {
			t.Errorf("Restore() expected a custom error type, got %T", err)
			return
		}

		if customErr.Code != http.StatusInternalServerError {
			t.Errorf("Restore() expected status code %d, got %d", http.StatusInternalServerError, customErr.Code)
		}
	})
}
