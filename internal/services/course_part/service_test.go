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
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/mikhail5545/product-service-go/internal/models/course"
	coursepart "github.com/mikhail5545/product-service-go/internal/models/course_part"
	coursemock "github.com/mikhail5545/product-service-go/internal/test/database/course_mock"
	coursepartmock "github.com/mikhail5545/product-service-go/internal/test/database/course_part_mock"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	mockPart := &coursepart.CoursePart{
		ID:               partID,
		Name:             "Course part name",
		ShortDescription: "Course part short description",
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockPartRepo.EXPECT().Get(gomock.Any(), partID).Return(mockPart, nil)

		// Act
		part, err := testService.Get(context.Background(), partID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(part, mockPart) {
			t.Errorf("Get() got = %v, want %v", part, mockPart)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-ID"

		// Act
		_, err := testService.Get(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockPartRepo.EXPECT().Get(gomock.Any(), partID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.Get(context.Background(), partID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockPartRepo.EXPECT().Get(gomock.Any(), partID).Return(nil, dbErr)

		// Act
		_, err := testService.Get(context.Background(), partID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetWithDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	mockPart := &coursepart.CoursePart{
		ID:               partID,
		Name:             "Course part name",
		ShortDescription: "Course part short description",
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockPartRepo.EXPECT().GetWithDeleted(gomock.Any(), partID).Return(mockPart, nil)

		// Act
		part, err := testService.GetWithDeleted(context.Background(), partID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(part, mockPart) {
			t.Errorf("Get() got = %v, want %v", part, mockPart)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-ID"

		// Act
		_, err := testService.GetWithDeleted(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockPartRepo.EXPECT().GetWithDeleted(gomock.Any(), partID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), partID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockPartRepo.EXPECT().GetWithDeleted(gomock.Any(), partID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithDeleted(context.Background(), partID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetWithUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	mockPart := &coursepart.CoursePart{
		ID:               partID,
		Name:             "Course part name",
		ShortDescription: "Course part short description",
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(mockPart, nil)

		// Act
		part, err := testService.GetWithUnpublished(context.Background(), partID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(part, mockPart) {
			t.Errorf("GetWithUnpublished() got = %v, want %v", part, mockPart)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-ID"

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), partID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithUnpublished(context.Background(), partID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetReduced(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	mockPart := &coursepart.CoursePart{
		ID:               partID,
		Name:             "Course part name",
		ShortDescription: "Course part short description",
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockPartRepo.EXPECT().Get(gomock.Any(), partID).Return(mockPart, nil)

		// Act
		part, err := testService.GetReduced(context.Background(), partID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(part, mockPart) {
			t.Errorf("GetReduced() got = %v, want %v", part, mockPart)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-ID"

		// Act
		_, err := testService.GetReduced(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockPartRepo.EXPECT().Get(gomock.Any(), partID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetReduced(context.Background(), partID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockPartRepo.EXPECT().Get(gomock.Any(), partID).Return(nil, dbErr)

		// Act
		_, err := testService.GetReduced(context.Background(), partID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetWithDeletedReduced(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	mockPart := &coursepart.CoursePart{
		ID:               partID,
		Name:             "Course part name",
		ShortDescription: "Course part short description",
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockPartRepo.EXPECT().GetWithDeleted(gomock.Any(), partID).Return(mockPart, nil)

		// Act
		part, err := testService.GetWithDeletedReduced(context.Background(), partID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(part, mockPart) {
			t.Errorf("GetWithDeletedReduced() got = %v, want %v", part, mockPart)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-ID"

		// Act
		_, err := testService.GetWithDeletedReduced(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockPartRepo.EXPECT().GetWithDeleted(gomock.Any(), partID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithDeletedReduced(context.Background(), partID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockPartRepo.EXPECT().GetWithDeleted(gomock.Any(), partID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithDeletedReduced(context.Background(), partID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_GetWithUnpublishedReduced(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	partID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	mockPart := &coursepart.CoursePart{
		ID:               partID,
		Name:             "Course part name",
		ShortDescription: "Course part short description",
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(mockPart, nil)

		// Act
		part, err := testService.GetWithUnpublishedReduced(context.Background(), partID)

		// Assert
		assert.NoError(t, err)
		if !reflect.DeepEqual(part, mockPart) {
			t.Errorf("GetWithUnpublishedReduced() got = %v, want %v", part, mockPart)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-ID"

		// Act
		_, err := testService.GetWithUnpublishedReduced(context.Background(), invalidID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.GetWithUnpublishedReduced(context.Background(), partID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		dbErr := errors.New("database error")
		mockPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(nil, dbErr)

		// Act
		_, err := testService.GetWithUnpublishedReduced(context.Background(), partID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	part1ID := "part-1-ID"
	part2ID := "part-2-ID"
	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	mockParts := []coursepart.CoursePart{
		{
			ID:               part1ID,
			CourseID:         courseID,
			Name:             "Course part 1 name",
			ShortDescription: "Course part 1 short description",
		},
		{
			ID:               part2ID,
			CourseID:         courseID,
			Name:             "Course part 2 name",
			ShortDescription: "Course part 2 short description",
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockPartRepo.EXPECT().List(gomock.Any(), courseID, limit, offset).Return(mockParts, nil)
		mockPartRepo.EXPECT().Count(gomock.Any(), courseID).Return(int64(2), nil)

		// Act
		parts, total, err := testService.List(context.Background(), courseID, limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, len(mockParts), len(parts))
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"
		limit, offset := 2, 0

		// Act
		_, _, err := testService.List(context.Background(), invalidID, limit, offset)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockPartRepo.EXPECT().List(gomock.Any(), courseID, limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.List(context.Background(), courseID, limit, offset)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_ListDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	part1ID := "part-1-ID"
	part2ID := "part-2-ID"
	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	mockParts := []coursepart.CoursePart{
		{
			ID:               part1ID,
			CourseID:         courseID,
			Name:             "Course part 1 name",
			ShortDescription: "Course part 1 short description",
		},
		{
			ID:               part2ID,
			CourseID:         courseID,
			Name:             "Course part 2 name",
			ShortDescription: "Course part 2 short description",
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockPartRepo.EXPECT().ListDeleted(gomock.Any(), courseID, limit, offset).Return(mockParts, nil)
		mockPartRepo.EXPECT().CountDeleted(gomock.Any(), courseID).Return(int64(2), nil)

		// Act
		parts, total, err := testService.ListDeleted(context.Background(), courseID, limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, len(mockParts), len(parts))
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"
		limit, offset := 2, 0

		// Act
		_, _, err := testService.ListDeleted(context.Background(), invalidID, limit, offset)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockPartRepo.EXPECT().ListDeleted(gomock.Any(), courseID, limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListDeleted(context.Background(), courseID, limit, offset)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_ListUnpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	part1ID := "part-1-ID"
	part2ID := "part-2-ID"
	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	mockParts := []coursepart.CoursePart{
		{
			ID:               part1ID,
			CourseID:         courseID,
			Name:             "Course part 1 name",
			ShortDescription: "Course part 1 short description",
		},
		{
			ID:               part2ID,
			CourseID:         courseID,
			Name:             "Course part 2 name",
			ShortDescription: "Course part 2 short description",
		},
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		mockPartRepo.EXPECT().ListUnpublished(gomock.Any(), courseID, limit, offset).Return(mockParts, nil)
		mockPartRepo.EXPECT().CountUnpublished(gomock.Any(), courseID).Return(int64(2), nil)

		// Act
		parts, total, err := testService.ListUnpublished(context.Background(), courseID, limit, offset)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, len(mockParts), len(parts))
	})

	t.Run("invalid UUID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-UUID"
		limit, offset := 2, 0

		// Act
		_, _, err := testService.ListUnpublished(context.Background(), invalidID, limit, offset)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		limit, offset := 2, 0
		dbErr := errors.New("database error")
		mockPartRepo.EXPECT().ListUnpublished(gomock.Any(), courseID, limit, offset).Return(nil, dbErr)

		// Act
		_, _, err := testService.ListUnpublished(context.Background(), courseID, limit, offset)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"

	createReq := coursepart.CreateRequest{
		CourseID:         courseID,
		Name:             "Course part name",
		ShortDescription: "Course part short description",
		Number:           3,
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
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)

		mockTxCourseRepo.EXPECT().Select(gomock.Any(), courseID, "id").Return(&course.Course{ID: courseID}, nil)
		mockTxPartRepo.EXPECT().CountQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), nil)

		var createdPart *coursepart.CoursePart
		mockTxPartRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, p *coursepart.CoursePart) {
				createdPart = p
			})

		// Act
		resp, err := testService.Create(context.Background(), &createReq)

		// Assert
		assert.NoError(t, err)
		if _, err := uuid.Parse(createdPart.ID); err != nil {
			t.Errorf("expected coursePart.ID to be a valid UUID, got %s", createdPart.ID)
		}
		assert.Equal(t, createReq.Name, createdPart.Name)
		assert.Equal(t, createReq.ShortDescription, createdPart.ShortDescription)
		assert.Equal(t, createReq.Number, createdPart.Number)
		assert.Equal(t, createdPart.ID, resp.ID)
		assert.Equal(t, createdPart.CourseID, resp.CourseID)

		if _, err := uuid.Parse(createdPart.CourseID); err != nil {
			t.Errorf("expected coursePart.CourseID to be a valid UUID, got %s", createdPart.CourseID)
		}
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)

		// Act
		_, err := testService.Create(context.Background(), &coursepart.CreateRequest{
			CourseID:         courseID,
			Name:             "1name", // Invalid name
			ShortDescription: "Short description",
			Number:           -2, // Invalid number,
		})

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("course not found", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)

		mockTxCourseRepo.EXPECT().Select(gomock.Any(), courseID, "id").Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.Create(context.Background(), &createReq)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("part with this number already exists", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)

		mockTxCourseRepo.EXPECT().Select(gomock.Any(), courseID, "id").Return(&course.Course{ID: courseID}, nil)
		mockTxPartRepo.EXPECT().CountQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)

		// Act
		_, err := testService.Create(context.Background(), &createReq)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)

		mockTxCourseRepo.EXPECT().Select(gomock.Any(), courseID, "id").Return(&course.Course{ID: courseID}, nil)
		mockTxPartRepo.EXPECT().CountQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), nil)
		dbErr := errors.New("database error")
		mockTxPartRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(dbErr)

		// Act
		_, err := testService.Create(context.Background(), &createReq)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"
	partID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"

	mockCourse := &course.Course{
		ID:      courseID,
		InStock: true,
	}

	mockPart := &coursepart.CoursePart{
		ID:               partID,
		CourseID:         courseID,
		Name:             "Course part name",
		ShortDescription: "Course part short description",
		Number:           2,
		Published:        false,
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
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)

		mockTxPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(mockPart, nil)
		mockTxCourseRepo.EXPECT().GetReduced(gomock.Any(), courseID).Return(mockCourse, nil)
		mockTxPartRepo.EXPECT().SetPublished(gomock.Any(), partID, true).Return(int64(1), nil)

		// Act
		err := testService.Publish(context.Background(), partID)

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

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)

		mockTxPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		err := testService.Publish(context.Background(), partID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)

		mockTxPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(mockPart, nil)
		mockTxCourseRepo.EXPECT().GetReduced(gomock.Any(), courseID).Return(mockCourse, nil)
		dbErr := errors.New("database error")
		mockTxPartRepo.EXPECT().SetPublished(gomock.Any(), partID, true).Return(int64(0), dbErr)

		// Act
		err := testService.Publish(context.Background(), partID)

		// Assert
		assert.Error(t, err)
	})

	t.Run("course not published", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)
		mockTxCourseRepo := coursemock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)
		mockCourseRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxCourseRepo)

		mockTxPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(mockPart, nil)
		mockCourse.InStock = false
		mockTxCourseRepo.EXPECT().GetReduced(gomock.Any(), courseID).Return(mockCourse, nil)

		// Act
		err := testService.Publish(context.Background(), partID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})
}

func TestService_Unpublished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	partID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"

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
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().SetPublished(gomock.Any(), partID, false).Return(int64(1), nil)

		// Act
		err := testService.Unpublish(context.Background(), partID)

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

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().SetPublished(gomock.Any(), partID, false).Return(int64(0), nil)

		// Act
		err := testService.Unpublish(context.Background(), partID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		dbErr := errors.New("database error")
		mockTxPartRepo.EXPECT().SetPublished(gomock.Any(), partID, false).Return(int64(0), dbErr)

		// Act
		err := testService.Unpublish(context.Background(), partID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	courseID := "d17081f3-4a56-4d00-b63e-f942537a702f"
	partID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"

	mockCoursePart := &coursepart.CoursePart{
		ID:               partID,
		Name:             "Old course part name",
		ShortDescription: "Old course part short description",
		Number:           4,
		CourseID:         courseID,
	}

	newName := "New course part name"
	newLongDescription := "New course part long description"
	newTags := []string{"New", "part", "tags"}
	newNumber := 2

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
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().Get(gomock.Any(), partID).Return(mockCoursePart, nil)
		mockTxPartRepo.EXPECT().CountQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), nil)

		var partUpdates map[string]any
		mockTxPartRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(_ context.Context, _ *coursepart.CoursePart, u map[string]any) {
				partUpdates = u
			})

		// Act
		updates, err := testService.Update(context.Background(), &coursepart.UpdateRequest{
			ID:              partID,
			CourseID:        courseID,
			Name:            &newName,
			LongDescription: &newLongDescription,
			Number:          &newNumber,
			Tags:            newTags,
		})

		// Assert
		assert.NoError(t, err)

		if name, ok := updates["name"].(string); !ok || name != newName {
			t.Errorf("coursePart.Name in response = %v, want %s", updates["name"], newName)
		}
		if longDescription, ok := updates["long_description"].(string); !ok || longDescription != newLongDescription {
			t.Errorf("coursePart.LongDescription in response = %v, want %s", updates["long_description"], newLongDescription)
		}
		if number, ok := updates["number"].(int); !ok || number != newNumber {
			t.Errorf("coursePart.Number in response = %v, want %d", updates["number"], newNumber)
		}
		if name, ok := partUpdates["name"].(string); !ok || name != newName {
			t.Errorf("coursePart.Name passed to repo = %v, want %s", partUpdates["name"], newName)
		}
		if longDescription, ok := partUpdates["long_description"].(string); !ok || longDescription != newLongDescription {
			t.Errorf("coursePart.LongDescription passed to repo = %v, want %s", partUpdates["long_description"], newLongDescription)
		}
		if number, ok := partUpdates["number"].(int); !ok || number != newNumber {
			t.Errorf("coursePart.Number passed to repo = %v, want %d", partUpdates["number"], newNumber)
		}
		if tags, ok := partUpdates["tags"].([]string); !ok || !reflect.DeepEqual(tags, newTags) {
			t.Errorf("coursePart.Number passed to repo = %v, want %v", partUpdates["tags"], newTags)
		}
	})

	t.Run("invalid request payload", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		invalidName := "1invalidname"
		invalidNumber := -2

		// Act
		_, err := testService.Update(context.Background(), &coursepart.UpdateRequest{
			ID:       partID,
			CourseID: courseID,
			Name:     &invalidName,
			Number:   &invalidNumber,
		})

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().Get(gomock.Any(), partID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		_, err := testService.Update(context.Background(), &coursepart.UpdateRequest{
			ID:              partID,
			CourseID:        courseID,
			Name:            &newName,
			LongDescription: &newLongDescription,
			Number:          &newNumber,
			Tags:            newTags,
		})

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("part with this number already exists", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().Get(gomock.Any(), partID).Return(mockCoursePart, nil)
		mockTxPartRepo.EXPECT().CountQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)

		// Act
		_, err := testService.Update(context.Background(), &coursepart.UpdateRequest{
			ID:              partID,
			CourseID:        courseID,
			Name:            &newName,
			LongDescription: &newLongDescription,
			Number:          &newNumber,
			Tags:            newTags,
		})

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidArgument)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().Get(gomock.Any(), partID).Return(mockCoursePart, nil)
		mockTxPartRepo.EXPECT().CountQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), nil)
		dbErr := errors.New("database error")
		mockTxPartRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), dbErr)

		// Act
		_, err := testService.Update(context.Background(), &coursepart.UpdateRequest{
			ID:              partID,
			CourseID:        courseID,
			Name:            &newName,
			LongDescription: &newLongDescription,
			Number:          &newNumber,
			Tags:            newTags,
		})

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	partID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"

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
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(&coursepart.CoursePart{}, nil)
		mockTxPartRepo.EXPECT().SetPublished(gomock.Any(), partID, false).Return(int64(1), nil)
		mockTxPartRepo.EXPECT().Delete(gomock.Any(), partID).Return(int64(1), nil)

		// Act
		err := testService.Delete(context.Background(), partID)

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
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(nil, gorm.ErrRecordNotFound)

		// Act
		err := testService.Delete(context.Background(), partID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(&coursepart.CoursePart{}, nil)
		mockTxPartRepo.EXPECT().SetPublished(gomock.Any(), partID, false).Return(int64(1), nil)
		dbErr := errors.New("database error")
		mockTxPartRepo.EXPECT().Delete(gomock.Any(), partID).Return(int64(0), dbErr)

		// Act
		err := testService.Delete(context.Background(), partID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_DeletePermanent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	partID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"

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
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().DeletePermanent(gomock.Any(), partID).Return(int64(1), nil)

		// Act
		err := testService.DeletePermanent(context.Background(), partID)

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
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().DeletePermanent(gomock.Any(), partID).Return(int64(0), nil)

		// Act
		err := testService.DeletePermanent(context.Background(), partID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		dbErr := errors.New("database error")
		mockTxPartRepo.EXPECT().DeletePermanent(gomock.Any(), partID).Return(int64(0), dbErr)

		// Act
		err := testService.DeletePermanent(context.Background(), partID)

		// Assert
		assert.Error(t, err)
	})
}

func TestService_Restore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPartRepo := coursepartmock.NewMockRepository(ctrl)
	mockCourseRepo := coursemock.NewMockRepository(ctrl)

	testService := New(mockPartRepo, mockCourseRepo)

	partID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"

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
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().Restore(gomock.Any(), partID).Return(int64(1), nil)

		// Act
		err := testService.Restore(context.Background(), partID)

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
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		mockTxPartRepo.EXPECT().Restore(gomock.Any(), partID).Return(int64(0), nil)

		// Act
		err := testService.Restore(context.Background(), partID)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		// Arrange
		mockTxPartRepo := coursepartmock.NewMockRepository(ctrl)

		mockPartRepo.EXPECT().DB().Return(db).AnyTimes()
		mockPartRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxPartRepo)

		dbErr := errors.New("database error")
		mockTxPartRepo.EXPECT().Restore(gomock.Any(), partID).Return(int64(0), dbErr)

		// Act
		err := testService.Restore(context.Background(), partID)

		// Assert
		assert.Error(t, err)
	})
}
