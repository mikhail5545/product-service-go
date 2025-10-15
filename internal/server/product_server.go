package server

import (
	"vitainmove.com/product-service-go/internal/services"
)

type ProductServer struct {
	productpb.UnimplementedProductServiceServer
	productService *services.ProductService
}

func NewProductServer(productService *services.ProductService) *ProductServer {
	return &ProductServer{productService: productService}
}
