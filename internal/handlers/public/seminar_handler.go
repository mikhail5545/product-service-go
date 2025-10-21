// vitainmove.com/product-service-go
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

package public

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"vitainmove.com/product-service-go/internal/services"
)

type SeminarHandler struct {
	seminarService *services.SeminarService
}

func NewSeminarHandler(seminarService *services.SeminarService) *SeminarHandler {
	return &SeminarHandler{seminarService: seminarService}
}

func (h *SeminarHandler) GetSeminar(c echo.Context) error {
	id := c.Param("id")
	seminar, err := h.seminarService.GetSeminar(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, seminar)
}

func (h *SeminarHandler) GetSeminars(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset <= -1 {
		offset = 0
	}

	seminars, total, err := h.seminarService.GetSeminars(c.Request().Context(), limit, offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"seminars": seminars,
		"total":    total,
	})
}
