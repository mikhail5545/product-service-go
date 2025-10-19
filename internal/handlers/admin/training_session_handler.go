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

package admin

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"vitainmove.com/product-service-go/internal/models"
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
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid training session id"})
	}
	ts, err := h.tsService.GetTrainingSession(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ts)
}

func (h *TrainingSessionHandler) GetTrainingSessions(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < -1 {
		limit = 10
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < -1 {
		offset = 10
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

func (h *TrainingSessionHandler) CreateTrainingSession(c echo.Context) error {
	var req *models.AddTrainingSessionRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ts := &models.TrainingSession{
		DurationMinutes: req.DurationMinutes,
		Format:          req.Format,
		Product: &models.Product{
			Name:        req.Product.Name,
			Description: req.Product.Description,
			Price:       req.Product.Price,
		},
	}

	_, err := h.tsService.CreateTrainingSession(c.Request().Context(), ts)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (h *TrainingSessionHandler) UpdateTrainingSession(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid training session id"})
	}
	var req *models.EditTrainingSessionRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ts := &models.TrainingSession{
		DurationMinutes: req.DurationMinutes,
		Format:          req.Format,
		Product: &models.Product{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
		},
	}

	_, err := h.tsService.UpdateTrainingSession(c.Request().Context(), ts, id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}

func (h *TrainingSessionHandler) DeleteTrainingSession(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid training session id"})
	}
	err := h.tsService.DeleteTrainingSession(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
