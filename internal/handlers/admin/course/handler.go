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

// Package course provides http-handling logic for course routes. It acts as
// an adapter between http server logic and service-layer business logic [seminarserice.Service].
package course

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	coursemodel "github.com/mikhail5545/product-service-go/internal/models/course"
	courseservice "github.com/mikhail5545/product-service-go/internal/services/course"
	"github.com/mikhail5545/product-service-go/internal/util/request"
)

// Handler holds the course service to handle HTTP requests.
type Handler struct {
	service courseservice.Service
}

// New creates a new course handler.
func New(s courseservice.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) ServeError(c echo.Context, code int, msg string) error {
	return c.JSON(code, map[string]string{"error": msg})
}

func (h *Handler) HandleServiceError(c echo.Context, err error) error {
	var se *courseservice.Error
	if errors.As(err, &se) {
		return c.JSON(se.GetCode(), map[string]any{"error": se.Msg})
	}
	return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Internal server error"})
}

// Get handles the retrieval of a single published course by its ID.
// @Summary Get a course by ID
// @Description Retrieves details for a specific course.
// @Success 200 {object} map[string]any{course_details=course.CourseDetails}
func (h *Handler) Get(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course ID")
	if err != nil {
		return err
	}
	details, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"course_details": details})
}

// GetWithDeleted handles the retrieval of a course by its ID, including soft-deleted ones.
// @Summary Get a course by ID (including deleted)
// @Description Retrieves details for a specific course, even if it has been soft-deleted.
// @Success 200 {object} map[string]any{course_details=course.CourseDetails}
func (h *Handler) GetWithDeleted(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course ID")
	if err != nil {
		return err
	}
	details, err := h.service.GetWithDeleted(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"course_details": details})
}

// GetWithUnpublished handles the retrieval of a course by its ID, including unpublished ones.
// @Summary Get a course by ID (including unpublished)
// @Description Retrieves details for a specific course, even if it is not published.
// @Success 200 {object} map[string]any{course_details=course.CourseDetails}
func (h *Handler) GetWithUnpublished(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course ID")
	if err != nil {
		return err
	}
	details, err := h.service.GetWithUnpublished(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"course_details": details})
}

// List handles the retrieval of a paginated list of published courses.
// @Summary List published courses
// @Description Retrieves a paginated list of courses that are currently published.
// @Success 200 {object} map[string]any{course_details=[]course.CourseDetails, total=int64}
func (h *Handler) List(c echo.Context) error {
	limit, offset, err := request.GetPaginationParams(c, 10, 0)
	if err != nil {
		return err
	}
	details, total, err := h.service.List(c.Request().Context(), limit, offset)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"course_details": details,
		"total":          total,
	})
}

// ListDeleted handles the retrieval of a paginated list of soft-deleted courses.
// @Summary List soft-deleted courses
// @Description Retrieves a paginated list of courses that have been soft-deleted.
// @Success 200 {object} map[string]any{course_details=[]course.CourseDetails, total=int64}
func (h *Handler) ListDeleted(c echo.Context) error {
	limit, offset, err := request.GetPaginationParams(c, 10, 0)
	if err != nil {
		return err
	}
	details, total, err := h.service.ListDeleted(c.Request().Context(), limit, offset)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"course_details": details,
		"total":          total,
	})
}

// ListUnpublished handles the retrieval of a paginated list of unpublished courses.
// @Summary List unpublished courses
// @Description Retrieves a paginated list of courses that are not currently published.
// @Success 200 {object} map[string]any{course_details=[]course.CourseDetails, total=int64}
func (h *Handler) ListUnpublished(c echo.Context) error {
	limit, offset, err := request.GetPaginationParams(c, 10, 0)
	if err != nil {
		return err
	}
	details, total, err := h.service.ListUnpublished(c.Request().Context(), limit, offset)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"course_details": details,
		"total":          total,
	})
}

// Create handles the creation of a new course and its associated product.
// @Summary Create a new course
// @Description Creates a new course with the provided details. The course is created in an unpublished state.
// @Success 201 "Created" {object} map[string]any{response=course.CreateResponse}
func (h *Handler) Create(c echo.Context) error {
	req := new(coursemodel.CreateRequest)
	if err := request.BindAndValidateJSON(c, req); err != nil {
		return err
	}
	resp, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusCreated, map[string]any{"response": resp})
}

// Update handles the partial update of an existing course and its product.
// @Summary Update a course
// @Description Updates a course's details. Only the provided fields will be updated.
// @Success 202 {object} map[string]any{updates=course.UpdateResponse}
func (h *Handler) Update(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course ID")
	if err != nil {
		return err
	}
	req := new(coursemodel.UpdateRequest)
	if err := request.BindAndValidateJSON(c, req); err != nil {
		return err
	}
	req.ID = id
	updates, err := h.service.Update(c.Request().Context(), req)
	if err != nil {
		return h.HandleServiceError(c, err)
	}

	return c.JSON(http.StatusAccepted, map[string]any{"updates": updates})
}

// Delete handles the soft-deletion of a course.
// @Summary Soft-delete a course
// @Description Soft-deletes a course and its associated product. The course is also unpublished.
// @Success 204 "No Content"
func (h *Handler) Delete(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course ID")
	if err != nil {
		return err
	}
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// DeletePermanent handles the permanent deletion of a course.
// @Summary Permanently delete a course
// @Description Permanently deletes a course and its associated product from the database. This action is irreversible.
// @Success 204 "No Content"
func (h *Handler) DeletePermanent(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course ID")
	if err != nil {
		return err
	}
	if err := h.service.DeletePermanent(c.Request().Context(), id); err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// Restore handles the restoration of a soft-deleted course.
// @Summary Restore a soft-deleted course
// @Description Restores a soft-deleted course and its associated product. The course will be in an unpublished state after restoration.
// @Success 202 "Accepted"
func (h *Handler) Restore(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course ID")
	if err != nil {
		return err
	}
	if err := h.service.Restore(c.Request().Context(), id); err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}

// Publish handles the publishing of a course.
// @Summary Publish a course
// @Description Publishes a course and its associated product, making them available.
// @Success 202 "Accepted"
func (h *Handler) Publish(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course ID")
	if err != nil {
		return err
	}
	if err := h.service.Publish(c.Request().Context(), id); err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}

// Unpublish handles the unpublishing of a course.
// @Summary Unpublish a course
// @Description Unpublishes a course, its product, and all its parts.
// @Success 202 "Accepted"
func (h *Handler) Unpublish(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course ID")
	if err != nil {
		return err
	}
	if err := h.service.Unpublish(c.Request().Context(), id); err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}
