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

package product

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mikhail5545/product-service-go/internal/models"
	"github.com/mikhail5545/product-service-go/internal/services/product"
	"github.com/mikhail5545/product-service-go/internal/util/request"
)

type Handler struct {
	service *product.Service
}

func New(s *product.Service) *Handler {
	return &Handler{service: s}
}

// AddProductPayload defines the expected structure of the incoming JSON from the frontend.
type AddProductPayload struct {
	ProductFormState models.AddProductRequest `json:"productFormState"`
}

func (h *Handler) Get(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid product ID")
	if err != nil {
		return err
	}
	product, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, product)
}

func (h *Handler) List(c echo.Context) error {
	limit, offset := request.GetPaginationParams(c, 10, 0)
	products, total, err := h.service.List(c.Request().Context(), limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{
		"products": products,
		"total":    total,
	})
}

func (h *Handler) Create(c echo.Context) error {
	var req *AddProductPayload
	if err := request.BindAndValidateJSON(c, &req); err != nil {
		return err
	}
	product := &models.Product{
		Name:             req.ProductFormState.Name,
		Description:      req.ProductFormState.Description,
		Price:            req.ProductFormState.Price,
		Amount:           req.ProductFormState.Amount,
		ShippingRequired: req.ProductFormState.ShippingRequired,
	}
	_, err := h.service.Create(c.Request().Context(), product)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Update(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid product ID")
	if err != nil {
		return err
	}
	var req *models.EditProductRequest
	if err := request.BindAndValidateJSON(c, &req); err != nil {
		return err
	}
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Amount:      req.Amount,
	}
	_, err = h.service.Update(c.Request().Context(), product, id)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := request.GetIDParam(c, ":id", "Invalid product ID")
	if err != nil {
		return err
	}
	err = h.service.Delete(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
