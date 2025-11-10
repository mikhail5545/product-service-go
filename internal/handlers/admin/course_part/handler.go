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
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	coursepartmodel "github.com/mikhail5545/product-service-go/internal/models/course_part"
	coursepart "github.com/mikhail5545/product-service-go/internal/services/course_part"
	"github.com/mikhail5545/product-service-go/internal/util/request"
)

type Handler struct {
	service coursepart.Service
}

func New(s coursepart.Service) *Handler {
	return &Handler{service: s}
}

// ServeError is a helper function to return error response with status code as `code` and message `msg`.
//
//	h.ServeError(http.StatusBadRequest, "Invalid request payload.")
func (h *Handler) ServeError(c echo.Context, code int, msg string) error {
	return c.JSON(code, map[string]string{"error": msg})
}

// HandleServiceError handles course service errors and populates
// error response based on error type.
func (h *Handler) HandleServiceError(c echo.Context, err error) error {
	if errors.Is(err, coursepart.ErrNotFound) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	} else if errors.Is(err, coursepart.ErrInvalidArgument) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Internal server error"})
}

// Get handles the retrieval of a single published course_part by its ID.
// @Summary Get a course_part by ID
// @Description Retrieves details for a specific course_part.
// @Tags admin-course-parts
// @Param id path string true "Course Part ID"
// @Success 200 {object} map[string]any{course_part=coursepartmodel.CoursePart}
// @Failure 400 {object} map[string]string{error=string} "Invalid course part ID"
// @Failure 404 {object} map[string]string{error=string} "Course part not found"
// @Router /admin/course-parts/{id} [get]
func (h *Handler) Get(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course part ID")
	if err != nil {
		return err
	}
	part, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"course_part": part})
}

// GetWithDeleted handles the retrieval of a course_part by its ID, including soft-deleted ones.
// @Summary Get a course_part by ID (including deleted)
// @Description Retrieves details for a specific course_part, even if it has been soft-deleted.
// @Tags admin-course-parts
// @Param id path string true "Course Part ID"
// @Success 200 {object} map[string]any{course_part=coursepartmodel.CoursePart}
// @Failure 400 {object} map[string]string{error=string} "Invalid course part ID"
// @Failure 404 {object} map[string]string{error=string} "Course part not found"
// @Router /admin/course-parts/deleted/{id} [get]
func (h *Handler) GetWithDeleted(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course part ID")
	if err != nil {
		return err
	}
	part, err := h.service.GetWithDeleted(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"course_part": part})
}

// GetWithUnpublished handles the retrieval of a course_part by its ID, including unpublished ones.
// @Summary Get a course_part by ID (including unpublished)
// @Description Retrieves details for a specific course_part, even if it is not published.
// @Tags admin-course-parts
// @Param id path string true "Course Part ID"
// @Success 200 {object} map[string]any{course_part=coursepartmodel.CoursePart}
// @Failure 400 {object} map[string]string{error=string} "Invalid course part ID"
// @Failure 404 {object} map[string]string{error=string} "Course part not found"
// @Router /admin/course-parts/unpublished/{id} [get]
func (h *Handler) GetWithUnpublished(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course part ID")
	if err != nil {
		return err
	}
	part, err := h.service.GetWithUnpublished(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"course_part": part})
}

// List handles the retrieval of a paginated list of published course_parts.
// @Summary List published course_parts
// @Description Retrieves a paginated list of course_parts that are currently published.
// @Tags admin-course-parts
// @Param cid path string true "Course ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]any{course_parts=[]coursepartmodel.CoursePart, total=int64}
// @Failure 400 {object} map[string]string{error=string} "Invalid course ID"
// @Router /admin/courses/{cid}/parts [get]
func (h *Handler) List(c echo.Context) error {
	cid, err := request.GetIDParam(c, ":cid", "Invalid course ID")
	if err != nil {
		return err
	}
	limit, offset, err := request.GetPaginationParams(c, 10, 0)
	if err != nil {
		return err
	}
	parts, total, err := h.service.List(c.Request().Context(), cid, limit, offset)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"course_parts": parts,
		"total":        total,
	})
}

// ListDeleted handles the retrieval of a paginated list of soft-deleted course_parts.
// @Summary List soft-deleted course_parts
// @Description Retrieves a paginated list of course_parts that have been soft-deleted.
// @Tags admin-course-parts
// @Param cid path string true "Course ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]any{course_parts=[]coursepartmodel.CoursePart, total=int64}
// @Failure 400 {object} map[string]string{error=string} "Invalid course ID"
// @Router /admin/courses/{cid}/parts/deleted [get]
func (h *Handler) ListDeleted(c echo.Context) error {
	cid, err := request.GetIDParam(c, ":cid", "Invalid course ID")
	if err != nil {
		return err
	}
	limit, offset, err := request.GetPaginationParams(c, 10, 0)
	if err != nil {
		return err
	}
	parts, total, err := h.service.ListDeleted(c.Request().Context(), cid, limit, offset)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"course_parts": parts,
		"total":        total,
	})
}

// ListUnpublished handles the retrieval of a paginated list of unpublished course_parts.
// @Summary List unpublished course_parts
// @Description Retrieves a paginated list of course_parts that are not currently published.
// @Tags admin-course-parts
// @Param cid path string true "Course ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]any{course_parts=[]coursepartmodel.CoursePart, total=int64}
// @Failure 400 {object} map[string]string{error=string} "Invalid course ID"
// @Router /admin/courses/{cid}/parts/unpublished [get]
func (h *Handler) ListUnpublished(c echo.Context) error {
	cid, err := request.GetIDParam(c, ":cid", "Invalid course ID")
	if err != nil {
		return err
	}
	limit, offset, err := request.GetPaginationParams(c, 10, 0)
	if err != nil {
		return err
	}
	parts, total, err := h.service.ListUnpublished(c.Request().Context(), cid, limit, offset)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"course_parts": parts,
		"total":        total,
	})
}

// Create handles the creation of a new course_part and its associated product.
// @Summary Create a new course_part
// @Description Creates a new course_part for a given course. The course_part is created in an unpublished state.
// @Tags admin-course-parts
// @Accept json
// @Param cid path string true "Course ID"
// @Param course_part body coursepartmodel.CreateRequest true "Course Part Create Request"
// @Success 201 "Created" {object} map[string]any{response=coursepartmodel.CreateResponse}
// @Failure 400 {object} map[string]string{error=string} "Invalid request payload or ID"
// @Failure 404 {object} map[string]string{error=string} "Course not found"
// @Router /admin/courses/{cid}/parts [post]
func (h *Handler) Create(c echo.Context) error {
	cid, err := request.GetIDParam(c, ":cid", "Invalid course ID")
	if err != nil {
		return err
	}
	var req *coursepartmodel.CreateRequest
	err = c.Bind(&req)
	if err != nil {
		return h.ServeError(c, http.StatusBadRequest, "Invalid request JSON payload")
	}
	req.CourseID = cid
	resp, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusCreated, map[string]any{"response": resp})
}

// Publish handles the publishing of a course_part.
// @Summary Publish a course_part
// @Description Publishes a course_part, making it available. Fails if the parent course is not published.
// @Tags admin-course-parts
// @Param id path string true "Course Part ID"
// @Success 202 "Accepted"
// @Failure 400 {object} map[string]string{error=string} "Invalid course part ID or parent course is unpublished"
// @Failure 404 {object} map[string]string{error=string} "Course part not found"
// @Failure 500 {object} map[string]string{error=string} "Internal server error"
// @Router /admin/course-parts/publish/{id} [post]
func (h *Handler) Publish(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course part ID")
	if err != nil {
		return err
	}
	err = h.service.Publish(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}

// Unpublish handles the unpublishing of a course_part.
// @Summary Unpublish a course_part
// @Description Unpublishes a course_part.
// @Tags admin-course-parts
// @Param id path string true "Course Part ID"
// @Success 202 "Accepted"
// @Failure 400 {object} map[string]string{error=string} "Invalid course part ID"
// @Failure 404 {object} map[string]string{error=string} "Course part not found"
// @Router /admin/course-parts/unpublish/{id} [post]
func (h *Handler) Unpublish(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course part ID")
	if err != nil {
		return err
	}
	err = h.service.Unpublish(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}

// AddVideo handles the association of a course_part with
// the MuxVideo.
// @Summary Add MuxVideo to the course_part
// @Description Associates a MuxVideo with an existing course_part by providing the MuxVideoID.
// @Tags admin-course-parts
// @Accept json
// @Param id path string true "Course Part ID"
// @Param video_request body coursepartmodel.AddVideoRequest true "Add Video Request"
// @Success 200 "OK" {object} map[string]any{response=coursepartmodel.AddVideoResponse}
// @Failure 400 {object} map[string]string{error=string} "Invalid request payload or ID"
// @Router /admin/course-parts/video/{id} [post]
func (h *Handler) AddVideo(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course part ID")
	if err != nil {
		return err
	}
	var req *coursepartmodel.AddVideoRequest
	err = c.Bind(&req)
	if err != nil {
		return h.ServeError(c, http.StatusBadRequest, "Invalid request JSON payload")
	}
	req.ID = id
	updates, err := h.service.AddVideo(c.Request().Context(), req)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"updates": updates})
}

// Update handles the partial update of an existing course_part and its product.
// @Summary Update a course_part
// @Description Updates a course_part's details. Only the provided fields will be updated.
// @Tags admin-course-parts
// @Accept json
// @Param id path string true "Course Part ID"
// @Param course_part body coursepartmodel.UpdateRequest true "Course Part Update Request"
// @Success 202 {object} map[string]any{updates=coursepartmodel.UpdateResponse}
// @Failure 400 {object} map[string]string{error=string} "Invalid request payload or ID"
// @Router /admin/course-parts/{id} [patch]
func (h *Handler) Update(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course part ID")
	if err != nil {
		return err
	}
	var req *coursepartmodel.UpdateRequest
	err = c.Bind(&req)
	if err != nil {
		return h.ServeError(c, http.StatusBadRequest, "Invalid request JSON payload")
	}
	req.ID = id
	updates, err := h.service.Update(c.Request().Context(), req)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusAccepted, map[string]any{"updates": updates})
}

// Delete handles the soft-deletion of a course_part.
// @Summary Soft-delete a course_part
// @Description Soft-deletes a course_part. The course_part is also unpublished.
// @Tags admin-course-parts
// @Param id path string true "Course Part ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string{error=string} "Invalid course part ID"
// @Failure 404 {object} map[string]string{error=string} "Course part not found"
// @Router /admin/course-parts/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course part ID")
	if err != nil {
		return err
	}
	err = h.service.Delete(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// DeletePermanent handles the permanent deletion of a course_part.
// @Summary Permanently delete a course_part
// @Description Permanently deletes a course_part from the database. This action is irreversible.
// @Tags admin-course-parts
// @Param id path string true "Course Part ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string{error=string} "Invalid course part ID"
// @Failure 404 {object} map[string]string{error=string} "Course part not found"
// @Router /admin/course-parts/permanent/{id} [delete]
func (h *Handler) DeletePermanent(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course part ID")
	if err != nil {
		return err
	}
	err = h.service.DeletePermanent(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// Restore handles the restoration of a soft-deleted course_part.
// @Summary Restore a soft-deleted course_part
// @Description Restores a soft-deleted course_part. The course_part will be in an unpublished state after restoration.
// @Tags admin-course-parts
// @Param id path string true "Course Part ID"
// @Success 202 "Accepted"
// @Failure 400 {object} map[string]string{error=string} "Invalid course part ID"
// @Failure 404 {object} map[string]string{error=string} "Course part not found"
// @Router /admin/course-parts/restore/{id} [post]
func (h *Handler) Restore(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course part ID")
	if err != nil {
		return err
	}
	err = h.service.Restore(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}
