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
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mikhail5545/product-service-go/internal/models"
	trainingsession "github.com/mikhail5545/product-service-go/internal/services/training_session"
	"github.com/mikhail5545/product-service-go/internal/util/request"
)

type Handler struct {
	tsService *trainingsession.Service
}

func New(tsService *trainingsession.Service) *Handler {
	return &Handler{tsService: tsService}
}

func (h *Handler) Get(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid training session ID")
	ts, err := h.tsService.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ts)
}

func (h *Handler) List(c echo.Context) error {
	limit, offset := request.GetPaginationParams(c, 10, 0)
	sessions, total, err := h.tsService.List(c.Request().Context(), limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{
		"training_sessions": sessions,
		"total":             total,
	})
}

func (h *Handler) Create(c echo.Context) error {
	var req *models.AddTrainingSessionRequest
	if err := request.BindAndValidateJSON(c, &req); err != nil {
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
	_, err := h.tsService.Create(c.Request().Context(), ts)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Update(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid training session ID")
	var req *models.EditTrainingSessionRequest
	if err := request.BindAndValidateJSON(c, &req); err != nil {
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
	_, _, err = h.tsService.Update(c.Request().Context(), ts, id)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid training session ID")
	err = h.tsService.Delete(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
