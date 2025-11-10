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
	coursepartservice "github.com/mikhail5545/product-service-go/internal/services/course_part"
	"github.com/mikhail5545/product-service-go/internal/util/request"
)

type Handler struct {
	service coursepartservice.Service
}

func New(s coursepartservice.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) ServeError(c echo.Context, code int, msg string) error {
	return c.JSON(code, map[string]string{"error": msg})
}

func (h *Handler) HandleServiceError(c echo.Context, err error) error {
	if errors.Is(err, coursepartservice.ErrNotFound) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	} else if errors.Is(err, coursepartservice.ErrInvalidArgument) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Internal server error"})
}

func (h *Handler) Get(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid course part ID")
	if err != nil {
		return err
	}
	details, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"course_part_details": details})
}

func (h *Handler) List(c echo.Context) error {
	cid, err := request.GetIDParam(c, ":cid", "Invalid course ID")
	if err != nil {
		return err
	}
	limit, offset, err := request.GetPaginationParams(c, -1, 0)
	if err != nil {
		return err
	}
	details, total, err := h.service.List(c.Request().Context(), cid, limit, offset)
	if err != nil {
		return h.HandleServiceError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"course_part_details": details,
		"total":               total,
	})
}
