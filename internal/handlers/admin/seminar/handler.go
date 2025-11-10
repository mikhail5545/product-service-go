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

package seminar

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mikhail5545/product-service-go/internal/models/seminar"
	seminarservice "github.com/mikhail5545/product-service-go/internal/services/seminar"
	"github.com/mikhail5545/product-service-go/internal/util/request"
)

// Handler holds [seminarservice.Service] instance to perform service-layer logic.
type Handler struct {
	service seminarservice.Service
}

// New creates a new Handler instance.
func New(s seminarservice.Service) *Handler {
	return &Handler{service: s}
}

// ServeError is a helper function to return error response with status code as `code` and message `msg`.
//
//	h.ServeError(http.StatusBadRequest, "Invalid request payload.")
func (h *Handler) ServeError(c echo.Context, code int, msg string) error {
	return c.JSON(code, map[string]string{"error": msg})
}

// HandleServiceError handles seminar service errors and populates
// error response based on error type.
func (h *Handler) HandleServiceError(c echo.Context, err error) error {
	if errors.Is(err, seminarservice.ErrNotFound) || errors.Is(err, seminarservice.ErrImageNotFoundOnOwner) || errors.Is(err, seminarservice.ErrProductsNotFound) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	} else if errors.Is(err, seminarservice.ErrInvalidArgument) || errors.Is(err, seminarservice.ErrImageLimitExceeded) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Internal server error"})
}

func (h *Handler) Get(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid seminar ID")
	if err != nil {
		return err
	}
	details, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"seminar_details": details})
}

func (h *Handler) GetWithDeleted(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid seminar ID")
	if err != nil {
		return err
	}
	details, err := h.service.GetWithDeleted(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"seminar_details": details})
}

func (h *Handler) GetWithUnpublished(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid seminar ID")
	if err != nil {
		return err
	}
	details, err := h.service.GetWithUnpublished(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"seminar_details": details})
}

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
		"seminar_details": details,
		"total":           total,
	})
}

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
		"seminar_details": details,
		"total":           total,
	})
}

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
		"seminar_details": details,
		"total":           total,
	})
}

func (h *Handler) Create(c echo.Context) error {
	req := new(seminar.CreateRequest)
	if err := request.BindAndValidateJSON(c, req); err != nil {
		return err
	}
	resp, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusCreated, map[string]any{"response": resp})
}

func (h *Handler) Update(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid seminar ID")
	if err != nil {
		return err
	}
	req := new(seminar.UpdateRequest)
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

func (h *Handler) Publish(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid seminar ID")
	if err != nil {
		return err
	}
	if err := h.service.Publish(c.Request().Context(), id); err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) Unpublish(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid seminar ID")
	if err != nil {
		return err
	}
	if err := h.service.Unpublish(c.Request().Context(), id); err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid seminar ID")
	if err != nil {
		return err
	}
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeletePermanent(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid seminar ID")
	if err != nil {
		return err
	}
	if err := h.service.DeletePermanent(c.Request().Context(), id); err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Restore(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid seminar ID")
	if err != nil {
		return err
	}
	if err := h.service.Restore(c.Request().Context(), id); err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}
