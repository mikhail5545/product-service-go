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
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mikhail5545/product-service-go/internal/models"
	"github.com/mikhail5545/product-service-go/internal/services/seminar"
	"github.com/mikhail5545/product-service-go/internal/util/request"
)

type Handler struct {
	service *seminar.Service
}

func New(s *seminar.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Get(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid seminar ID")
	seminar, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, seminar)
}

func (h *Handler) List(c echo.Context) error {
	limit, offset := request.GetPaginationParams(c, 10, 0)
	seminars, total, err := h.service.List(c.Request().Context(), limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{
		"seminars": seminars,
		"total":    total,
	})
}

func (h *Handler) Create(c echo.Context) error {
	var req *models.AddSeminarRequest
	if err := request.BindAndValidateJSON(c, &req); err != nil {
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
	_, err := h.service.Create(c.Request().Context(), seminar)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Update(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid seminar ID")
	var req *models.UpdateSeminarRequest
	if err := request.BindAndValidateJSON(c, &req); err != nil {
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
	_, err = h.service.Update(c.Request().Context(), seminar, id)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid seminar ID")
	err = h.service.Delete(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
