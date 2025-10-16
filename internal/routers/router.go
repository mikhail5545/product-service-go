package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"vitainmove.com/product-service-go/internal/handlers"
	"vitainmove.com/product-service-go/internal/services"
)

func SetupRouter(e *echo.Echo, productService *services.ProductService) {
	e.HTTPErrorHandler = handlers.HTTPErrorHandler

	api := e.Group("/api")
	ver := api.Group("/v0")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	productHandler := handlers.NewProductHandler(productService)
	products := ver.Group("/products")
	{
		products.GET("", productHandler.GetProducts)
	}
}
