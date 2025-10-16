package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"vitainmove.com/product-service-go/internal/services"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
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
