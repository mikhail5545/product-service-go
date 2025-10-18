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

package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"vitainmove.com/product-service-go/internal/services"
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
