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

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// AddProductPayload defines the expected structure of the incoming JSON from the frontend.
type AddProductPayload struct {
	ProductFormState models.AddProductRequest `json:"productFormState"`
}

func (h *ProductHandler) GetProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}
	product, err := h.productService.GetProduct(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < -1 {
		limit = 10
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	products, total, err := h.productService.GetProducts(c.Request().Context(), limit, offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"products": products,
		"total":    total,
	})
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req *AddProductPayload
	if err := c.Bind(&req); err != nil {
		return err
	}

	product := &models.Product{
		Name:             req.ProductFormState.Name,
		Description:      req.ProductFormState.Description,
		Price:            req.ProductFormState.Price,
		Amount:           req.ProductFormState.Amount,
		ShippingRequired: req.ProductFormState.ShippingRequired,
	}

	_, err := h.productService.CreateProduct(c.Request().Context(), product)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}
	var req *models.EditProductRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Amount:      req.Amount,
	}

	_, err := h.productService.UpdateProduct(c.Request().Context(), product, id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}
	err := h.productService.DeleteProduct(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
