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
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	trainingsession "github.com/mikhail5545/product-service-go/internal/models/training_session"
	trainingsessionservice "github.com/mikhail5545/product-service-go/internal/services/training_session"
	"github.com/mikhail5545/product-service-go/internal/util/request"
)

type Handler struct {
	tsService trainingsessionservice.Service
}

func New(tsService trainingsessionservice.Service) *Handler {
	return &Handler{tsService: tsService}
}

func (h *Handler) ServeError(c echo.Context, code int, msg string) error {
	return c.JSON(code, map[string]string{"error": msg})
}

func (h *Handler) HandleServiceError(c echo.Context, err error) error {
	var se *trainingsessionservice.Error
	if errors.As(err, &se) {
		return c.JSON(se.GetCode(), map[string]any{"error": se.Msg})
	}
	return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Internal server error"})
}

func (h *Handler) Get(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid training session ID")
	if err != nil {
		return err
	}
	details, err := h.tsService.Get(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"training_session_details": details})
}

func (h *Handler) GetWithDeleted(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid training session ID")
	if err != nil {
		return err
	}
	details, err := h.tsService.GetWithDeleted(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"training_session_details": details})
}

func (h *Handler) GetWithUnpublished(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid training session ID")
	if err != nil {
		return err
	}
	details, err := h.tsService.GetWithUnpublished(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"training_session_details": details})
}

// List handles the retrieval of a paginated list of published training sessions.
// @Summary List published training sessions
// @Description Retrieves a paginated list of training sessions that are currently published.
// @Success 200 {object} map[string]any{training_session_details=[]trainingsession.TrainingSessionDetails, total=int64}
func (h *Handler) List(c echo.Context) error {
	limit, offset, err := request.GetPaginationParams(c, 10, 0)
	if err != nil {
		return err
	}
	details, total, err := h.tsService.List(c.Request().Context(), limit, offset)
	if err != nil {
		h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"training_session_details": details,
		"total":                    total,
	})
}

// ListDeleted handles the retrieval of a paginated list of soft-deleted training sessions.
// @Summary List soft-deleted training sessions
// @Description Retrieves a paginated list of training sessions that have been soft-deleted.
// @Success 200 {object} map[string]any{training_session_details=[]trainingsession.TrainingSessionDetails, total=int64}
func (h *Handler) ListDeleted(c echo.Context) error {
	limit, offset, err := request.GetPaginationParams(c, 10, 0)
	if err != nil {
		return err
	}
	details, total, err := h.tsService.ListDeleted(c.Request().Context(), limit, offset)
	if err != nil {
		h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"training_session_details": details,
		"total":                    total,
	})
}

// ListUnpublished handles the retrieval of a paginated list of unpublished training sessions.
// @Summary List unpublished training sessions
// @Description Retrieves a paginated list of training sessions that are not currently published.
// @Success 200 {object} map[string]any{training_session_details=[]trainingsession.TrainingSessionDetails, total=int64}
func (h *Handler) ListUnpublished(c echo.Context) error {
	limit, offset, err := request.GetPaginationParams(c, 10, 0)
	if err != nil {
		return err
	}
	details, total, err := h.tsService.ListUnpublished(c.Request().Context(), limit, offset)
	if err != nil {
		h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"training_session_details": details,
		"total":                    total,
	})
}

func (h *Handler) Create(c echo.Context) error {
	var req *trainingsession.CreateRequest
	if err := c.Bind(&req); err != nil {
		return h.ServeError(c, http.StatusBadRequest, "Invalid request JSON payload")
	}
	resp, err := h.tsService.Create(c.Request().Context(), req)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusCreated, map[string]any{"response": resp})
}

func (h *Handler) Publish(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid training session ID")
	if err != nil {
		return err
	}
	err = h.tsService.Publish(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) Unpublish(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid training session ID")
	if err != nil {
		return err
	}
	err = h.tsService.Unpublish(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) Update(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid training session ID")
	if err != nil {
		return err
	}
	var req *trainingsession.UpdateRequest
	if err := c.Bind(&req); err != nil {
		return h.ServeError(c, http.StatusBadRequest, "Invalid request JSON payload")
	}
	req.ID = id
	updates, err := h.tsService.Update(c.Request().Context(), req)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusCreated, map[string]any{"updates": updates})
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid training session ID")
	if err != nil {
		return err
	}
	err = h.tsService.Delete(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeletePermanent(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid training session ID")
	if err != nil {
		return err
	}
	err = h.tsService.DeletePermanent(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Restore(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid training session ID")
	if err != nil {
		return err
	}
	err = h.tsService.Restore(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}
