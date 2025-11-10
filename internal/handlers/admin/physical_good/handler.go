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

package physicalgood

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	physicalgood "github.com/mikhail5545/product-service-go/internal/models/physical_good"
	physicalgoodservice "github.com/mikhail5545/product-service-go/internal/services/physical_good"
	"github.com/mikhail5545/product-service-go/internal/util/request"
)

type Handler struct {
	service physicalgoodservice.Service
}

func New(s physicalgoodservice.Service) *Handler {
	return &Handler{service: s}
}

// ServeError is a helper function to return error response with status code as `code` and message `msg`.
//
//	h.ServeError(http.StatusBadRequest, "Invalid request payload.")
func (h *Handler) ServeError(c echo.Context, code int, msg string) error {
	return c.JSON(code, map[string]string{"error": msg})
}

// HandleServiceError handles physical good service errors and populates
// error response based on error type.
func (h *Handler) HandleServiceError(c echo.Context, err error) error {
	if errors.Is(err, physicalgoodservice.ErrNotFound) || errors.Is(err, physicalgoodservice.ErrImageNotFoundOnOwner) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	} else if errors.Is(err, physicalgoodservice.ErrInvalidArgument) || errors.Is(err, physicalgoodservice.ErrImageLimitExceeded) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Internal server error"})
}

func (h *Handler) Get(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid physical good ID")
	if err != nil {
		return err
	}
	details, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"physical_good_details": details})
}

func (h *Handler) GetWithDeleted(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid physical good ID")
	if err != nil {
		return err
	}
	details, err := h.service.GetWithDeleted(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"physical_good_details": details})
}

func (h *Handler) GetWithUnpublished(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid physical good ID")
	if err != nil {
		return err
	}
	details, err := h.service.GetWithUnpublished(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"physical_good_details": details})
}

// List handles the retrieval of a paginated list of published physical goods.
// @Summary List published physical goods
// @Description Retrieves a paginated list of physical goods that are currently published.
// @Success 200 {object} map[string]any{physical_good_details=[]physicalgood.PhysicalGoodDetails, total=int64}
func (h *Handler) List(c echo.Context) error {
	limit, offset, err := request.GetPaginationParams(c, 10, 0)
	if err != nil {
		return err
	}
	details, total, err := h.service.List(c.Request().Context(), limit, offset)
	if err != nil {
		h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"physical_good_details": details,
		"total":                 total,
	})
}

// ListDeleted handles the retrieval of a paginated list of soft-deleted physical goods.
// @Summary List soft-deleted physical goods
// @Description Retrieves a paginated list of physical goods that have been soft-deleted.
// @Success 200 {object} map[string]any{physical_good_details=[]physicalgood.PhysicalGoodDetails, total=int64}
func (h *Handler) ListDeleted(c echo.Context) error {
	limit, offset, err := request.GetPaginationParams(c, 10, 0)
	if err != nil {
		return err
	}
	details, total, err := h.service.ListDeleted(c.Request().Context(), limit, offset)
	if err != nil {
		h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"physical_good_details": details,
		"total":                 total,
	})
}

// ListUnpublished handles the retrieval of a paginated list of unpublished physical goods.
// @Summary List unpublished physical goods
// @Description Retrieves a paginated list of physical goods that are not currently published.
// @Success 200 {object} map[string]any{physical_good_details=[]physicalgood.PhysicalGoodDetails, total=int64}
func (h *Handler) ListUnpublished(c echo.Context) error {
	limit, offset, err := request.GetPaginationParams(c, 10, 0)
	if err != nil {
		return err
	}
	details, total, err := h.service.ListUnpublished(c.Request().Context(), limit, offset)
	if err != nil {
		h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"physical_good_details": details,
		"total":                 total,
	})
}

func (h *Handler) Create(c echo.Context) error {
	var req *physicalgood.CreateRequest
	if err := c.Bind(&req); err != nil {
		return h.ServeError(c, http.StatusBadRequest, "Invalid request JSON payload")
	}
	resp, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusCreated, map[string]any{"response": resp})
}

func (h *Handler) Publish(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid physical good ID")
	if err != nil {
		return err
	}
	err = h.service.Publish(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) Unpublish(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid physical good ID")
	if err != nil {
		return err
	}
	err = h.service.Unpublish(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) Update(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid physical good ID")
	if err != nil {
		return err
	}
	var req *physicalgood.UpdateRequest
	if err := c.Bind(&req); err != nil {
		return h.ServeError(c, http.StatusBadRequest, "Invalid request JSON payload")
	}
	req.ID = id
	updates, err := h.service.Update(c.Request().Context(), req)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusAccepted, map[string]any{"updates": updates})
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid physical good ID")
	if err != nil {
		return err
	}
	err = h.service.Delete(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeletePermanent(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid physical good ID")
	if err != nil {
		return err
	}
	err = h.service.DeletePermanent(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Restore(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid physical good ID")
	if err != nil {
		return err
	}
	err = h.service.Restore(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.NoContent(http.StatusAccepted)
}
