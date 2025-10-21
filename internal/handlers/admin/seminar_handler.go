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

package admin

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mikhail5545/product-service-go/internal/models"
	"github.com/mikhail5545/product-service-go/internal/services"
)

type SeminarHandler struct {
	seminarService *services.SeminarService
}

func NewSeminarService(seminarService *services.SeminarService) *SeminarHandler {
	return &SeminarHandler{seminarService: seminarService}
}

func (h *SeminarHandler) GetSeminar(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid seminar id"})
	}
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

func (h *SeminarHandler) CreateSeminar(c echo.Context) error {
	var req *models.AddSeminarRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	seminar := &models.Seminar{
		Name:            req.Name,
		Description:     req.Description,
		Date:            req.Date,
		EndingDate:      req.EndingDate,
		Place:           req.Place,
		LatePaymentDate: req.LatePaymentDate,
		Details:         req.Details,
		ReservationProduct: &models.Product{
			Name:        *req.ReservationProduct.Name,
			Description: *req.ReservationProduct.Description,
			Price:       req.ReservationProduct.Price,
		},
		EarlyProduct: &models.Product{
			Name:        *req.EarlyProduct.Name,
			Description: *req.EarlyProduct.Description,
			Price:       req.EarlyProduct.Price,
		},
		LateProduct: &models.Product{
			Name:        *req.LateProduct.Name,
			Description: *req.LateProduct.Description,
			Price:       req.LateProduct.Price,
		},
		EarlySurchargeProduct: &models.Product{
			Name:        *req.EarlySurchargeProduct.Name,
			Description: *req.EarlySurchargeProduct.Description,
			Price:       req.EarlySurchargeProduct.Price,
		},
		LateSurchargeProduct: &models.Product{
			Name:        *req.LateSurchargeProduct.Name,
			Description: *req.LateSurchargeProduct.Description,
			Price:       req.LateSurchargeProduct.Price,
		},
	}
	_, err := h.seminarService.CreateSeminar(c.Request().Context(), seminar)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (h *SeminarHandler) UpdateSeminar(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid seminar id"})
	}
	var req *models.UpdateSeminarRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	seminar := &models.Seminar{
		Name:            *req.Name,
		Description:     *req.Description,
		Date:            *req.Date,
		EndingDate:      *req.EndingDate,
		Place:           *req.Place,
		LatePaymentDate: *req.LatePaymentDate,
		Details:         *req.Details,
		ReservationProduct: &models.Product{
			Name:        *req.ReservationProduct.Name,
			Description: *req.ReservationProduct.Description,
			Price:       req.ReservationProduct.Price,
		},
		EarlyProduct: &models.Product{
			Name:        *req.EarlyProduct.Name,
			Description: *req.EarlyProduct.Description,
			Price:       req.EarlyProduct.Price,
		},
		LateProduct: &models.Product{
			Name:        *req.LateProduct.Name,
			Description: *req.LateProduct.Description,
			Price:       req.LateProduct.Price,
		},
		EarlySurchargeProduct: &models.Product{
			Name:        *req.EarlySurchargeProduct.Name,
			Description: *req.EarlySurchargeProduct.Description,
			Price:       req.EarlySurchargeProduct.Price,
		},
		LateSurchargeProduct: &models.Product{
			Name:        *req.LateSurchargeProduct.Name,
			Description: *req.LateSurchargeProduct.Description,
			Price:       req.LateSurchargeProduct.Price,
		},
	}
	_, err := h.seminarService.UpdateSeminar(c.Request().Context(), seminar, id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}

func (h *SeminarHandler) DeleteSeminar(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid seminar id"})
	}
	err := h.seminarService.DeleteSeminar(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
