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

package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mikhail5545/product-service-go/internal/database"
	"github.com/mikhail5545/product-service-go/internal/models"
	"gorm.io/gorm"
)

type CourseService struct {
	CourseRepo  database.CourseRepository
	ProductRepo database.ProductRepository
}

type CourseServiceError struct {
	Msg  string
	Err  error
	Code int
}

func (e *CourseServiceError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func (e *CourseServiceError) Unwrap() error {
	return e.Err
}

func (e *CourseServiceError) GetCode() int {
	return e.Code
}

func NewCourseService(cr database.CourseRepository, pr database.ProductRepository) *CourseService {
	return &CourseService{
		CourseRepo:  cr,
		ProductRepo: pr,
	}
}

func (s *CourseService) GetCourse(ctx context.Context, id string) (*models.Course, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &CourseServiceError{
			Msg:  "Invalid course ID",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	course, err := s.CourseRepo.Find(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &CourseServiceError{
				Msg:  "Course not found",
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return nil, &CourseServiceError{
			Msg:  "Failed to get course",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return course, nil
}

func (s *CourseService) GetCourses(ctx context.Context, limit, offset int) ([]models.Course, int64, error) {
	courses, err := s.CourseRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, &CourseServiceError{
			Msg:  "Failed to get courses",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	total, err := s.CourseRepo.Count(ctx)
	if err != nil {
		return nil, 0, &CourseServiceError{
			Msg:  "Failed to get courses count",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return courses, total, nil
}

func (s *CourseService) CreateCourse(ctx context.Context, course *models.Course) (*models.Course, error) {
	err := s.CourseRepo.DB().Transaction(func(tx *gorm.DB) error {
		txCourseRepo := s.CourseRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		course.ID = uuid.New().String()
		course.CreatedAt = time.Now()
		course.UpdatedAt = time.Now()
		course.Product.ID = uuid.New().String()
		course.Product.CreatedAt = time.Now()
		course.Product.UpdatedAt = time.Now()
		course.Product.ProductType = "course"
		course.Product.Amount = 0
		course.Product.ShippingRequired = false

		if err := txProductRepo.Create(ctx, course.Product); err != nil {
			return &CourseServiceError{
				Msg:  "Failed to create underlying product for course",
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}
		course.ProductID = course.Product.ID

		if err := txCourseRepo.Create(ctx, course); err != nil {
			return &CourseServiceError{
				Msg:  "Failed to create course",
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (s *CourseService) UpdateCourse(ctx context.Context, course *models.Course, id string) (*models.Course, error) {
	var courseToUpdate *models.Course
	err := s.CourseRepo.DB().Transaction(func(tx *gorm.DB) error {
		txCourseRepo := s.CourseRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		var findErr error
		courseToUpdate, findErr = txCourseRepo.Find(ctx, id)
		if findErr != nil {
			if errors.Is(findErr, gorm.ErrRecordNotFound) {
				return &CourseServiceError{
					Msg:  "Course not found",
					Err:  findErr,
					Code: http.StatusNotFound,
				}
			}
			return &CourseServiceError{
				Msg:  "Failed to get course",
				Err:  findErr,
				Code: http.StatusInternalServerError,
			}
		}

		var courseUpdated bool
		var courseProductUpdated bool
		if course.Name != "" && course.Name != courseToUpdate.Name {
			courseToUpdate.Name = course.Name
			courseToUpdate.Product.Name = fmt.Sprintf("Приобрести доступ на %s (%d  дней)", courseToUpdate.Name, courseToUpdate.AccessDuration)
			courseUpdated = true
			courseProductUpdated = true
		}
		if course.Description != "" && course.Description != courseToUpdate.Description {
			courseToUpdate.Description = course.Description
			courseToUpdate.Product.Description = course.Description
			courseUpdated = true
			courseProductUpdated = true
		}
		if course.AccessDuration >= 0 && course.AccessDuration != courseToUpdate.AccessDuration {
			courseToUpdate.AccessDuration = course.AccessDuration
			courseUpdated = true
		}
		if course.Product != nil && course.Product.Price >= 0 && course.Product.Price != courseToUpdate.Product.Price {
			courseToUpdate.Product.Price = course.Product.Price
			courseProductUpdated = true
		}
		if course.Topic != "" && course.Topic != courseToUpdate.Topic {
			courseToUpdate.Topic = course.Topic
			courseUpdated = true
		}

		if courseProductUpdated {
			if err := txProductRepo.Update(ctx, courseToUpdate.Product); err != nil {
				return &CourseServiceError{
					Msg:  "Failed to update underlying course product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		if courseUpdated {
			if err := txCourseRepo.Update(ctx, courseToUpdate); err != nil {
				return &CourseServiceError{
					Msg:  "Failed to update course",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return courseToUpdate, nil
}
