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
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	mediaclient "github.com/mikhail5545/product-service-go/internal/clients/mediaservice"
	"github.com/mikhail5545/product-service-go/internal/database/course"
	"github.com/mikhail5545/product-service-go/internal/database/product"
	"github.com/mikhail5545/product-service-go/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	CourseRepo  course.Repository
	ProductRepo product.Repository
	MediaClient *mediaclient.Client
}

type Error struct {
	Msg  string
	Err  error
	Code int
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) GetCode() int {
	return e.Code
}

func New(
	cr course.Repository,
	pr product.Repository,
	mc *mediaclient.Client,
) *Service {
	return &Service{
		CourseRepo:  cr,
		ProductRepo: pr,
		MediaClient: mc,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*models.Course, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{
			Msg:  "Invalid course ID",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	course, err := s.CourseRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{
				Msg:  "Course not found",
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return nil, &Error{
			Msg:  "Failed to get course",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return course, nil
}

func (s *Service) GetReduced(ctx context.Context, id string) (*models.Course, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{
			Msg:  "Invalid course ID",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	course, err := s.CourseRepo.GetReduced(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{
				Msg:  "Course not found",
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return nil, &Error{
			Msg:  "Failed to get course",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return course, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]models.Course, int64, error) {
	courses, err := s.CourseRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to get courses",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	total, err := s.CourseRepo.Count(ctx)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to get courses count",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return courses, total, nil
}

func (s *Service) Create(ctx context.Context, course *models.Course) (*models.Course, error) {
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
			return &Error{
				Msg:  "Failed to create underlying product for course",
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}
		course.ProductID = course.Product.ID

		if err := txCourseRepo.Create(ctx, course); err != nil {
			return &Error{
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

// Update performs a partial update of a course information.
// The request should contain the course's ID and the fields to be updated.
// The response includes only the fields that were actually changed in map representation.
// It returns updates of course model itself and updates of underlying product model.
func (s *Service) Update(ctx context.Context, req *models.Course, id string) (map[string]any, map[string]any, error) {
	updates := make(map[string]any)
	productUpdates := make(map[string]any)
	err := s.CourseRepo.DB().Transaction(func(tx *gorm.DB) error {
		txCourseRepo := s.CourseRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		course, err := txCourseRepo.Get(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{
					Msg:  "Course not found",
					Err:  err,
					Code: http.StatusNotFound,
				}
			}
			return &Error{
				Msg:  "Failed to get course",
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}

		if req.Name != "" && req.Name != course.Name {
			updates["name"] = req.Name
			productUpdates["name"] = fmt.Sprintf("Приобрести доступ на %s (%d  дней)", course.Name, course.AccessDuration)
		}
		if req.Description != "" && req.Description != course.Description {
			updates["description"] = req.Description
			productUpdates["description"] = req.Description
		}
		if req.AccessDuration >= 0 && req.AccessDuration != course.AccessDuration {
			updates["access_duration"] = req.AccessDuration
		}
		if req.Product != nil && req.Product.Price >= 0 && req.Product.Price != course.Product.Price {
			productUpdates["price"] = req.Product.Price
		}
		if req.Topic != "" && req.Topic != course.Topic {
			updates["topic"] = req.Topic
		}

		if len(productUpdates) > 0 {
			if _, err := txProductRepo.Update(ctx, course.Product, productUpdates); err != nil {
				return &Error{
					Msg:  "Failed to update underlying course product",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		if len(updates) > 0 {
			if _, err := txCourseRepo.Update(ctx, course, updates); err != nil {
				return &Error{
					Msg:  "Failed to update course",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}
	return updates, productUpdates, nil
}
