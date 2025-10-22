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

package admin

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mikhail5545/product-service-go/internal/models"
	"github.com/mikhail5545/product-service-go/internal/services"
)

type CourseHandler struct {
	courseService *services.CourseService
}

func NewCourseHandler(courseService *services.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

func (h *CourseHandler) GetCourse(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid course ID"})
	}
	course, err := h.courseService.GetCourse(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) GetCourses(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset <= -1 {
		offset = 0
	}
	courses, total, err := h.courseService.GetCourses(c.Request().Context(), limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{
		"courses": courses,
		"total":   total,
	})
}

func (h *CourseHandler) CreateCourse(c echo.Context) error {
	var req *models.AddCourseRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	course := &models.Course{
		Name:           req.Name,
		Description:    req.Description,
		Topic:          req.Topic,
		AccessDuration: req.AccessDuration,
		Product: &models.Product{
			Price: req.Price,
		},
	}

	_, err := h.courseService.CreateCourse(c.Request().Context(), course)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (h *CourseHandler) UpdateCourse(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid course ID"})
	}
	var req *models.EditCourseRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	course := &models.Course{
		Name:           *req.Name,
		Description:    *req.Description,
		Topic:          *req.Topic,
		AccessDuration: *req.AccessDuration,
		Product: &models.Product{
			Price: req.Product.Price,
		},
	}
	_, err := h.courseService.UpdateCourse(c.Request().Context(), course, id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}

func (h *CourseHandler) GetCoursePart(c echo.Context) error {
	id := c.Param(":part_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid course ID"})
	}
	part, err := h.courseService.GetCoursePart(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, part)
}
