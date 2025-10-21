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

package public

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mikhail5545/product-service-go/internal/services"
)

type TrainingSessionHandler struct {
	tsService *services.TrainingSessionService
}

func NewTrainingSessionHandler(tsService *services.TrainingSessionService) *TrainingSessionHandler {
	return &TrainingSessionHandler{tsService: tsService}
}

func (h *TrainingSessionHandler) GetTrainingSession(c echo.Context) error {
	id := c.Param("id")
	ts, err := h.tsService.GetTrainingSession(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ts)
}

func (h *TrainingSessionHandler) GetTrainingSessions(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	sessions, total, err := h.tsService.GetTrainingSessions(c.Request().Context(), limit, offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"training_sessions": sessions,
		"total":             total,
	})
}
