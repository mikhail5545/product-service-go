package server

import (
	"context"

	"vitainmove.com/product-service-go/internal/models"
	"vitainmove.com/product-service-go/internal/services"
	"vitainmove.com/product-service-go/internal/utils"
	productpb "vitainmove.com/product-service-go/proto/product/v0"
)

type ProductServer struct {
	productpb.UnimplementedProductServiceServer
	productService *services.ProductService
}

func NewProductServer(productService *services.ProductService) *ProductServer {
	return &ProductServer{productService: productService}
}

func (s *ProductServer) ListProducts(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
	products, total, err := s.productService.GetProducts(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, toGRPCError(err)
	}
	var pbProducts []*productpb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, utils.ConvertToProtobufProduct(&product))
	}

	return &productpb.ListProductsResponse{Products: pbProducts, Total: total}, nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.GetProductResponse, error) {
	product, err := s.productService.GetProduct(ctx, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &productpb.GetProductResponse{Product: utils.ConvertToProtobufProduct(product)}, nil
}

func (s *ProductServer) ListProductsByType(ctx context.Context, req *productpb.ListProductsByTypeRequest) (*productpb.ListProductsByTypeResponse, error) {
	products, total, err := s.productService.GetProductByType(ctx, req.GetProductType(), int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, toGRPCError(err)
	}
	var pbProducts []*productpb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, utils.ConvertToProtobufProduct(&product))
	}

	return &productpb.ListProductsByTypeResponse{Products: pbProducts, Total: total}, nil
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error) {
	createReq := &models.AddProductRequest{
		Name:             req.GetName(),
		Description:      req.GetDescription(),
		Price:            req.GetPrice(),
		Amount:           int(req.GetAmount()),
		ShippingRequired: req.GetShippingRequired(),
	}

	product, err := s.productService.CreateProduct(ctx, createReq)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &productpb.CreateProductResponse{Product: utils.ConvertToProtobufProduct(product)}, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error) {
	updateReq := &models.EditProductRequest{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Amount:      int(req.GetAmount()),
	}

	product, err := s.productService.UpdateProduct(ctx, updateReq, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &productpb.UpdateProductResponse{Product: utils.ConvertToProtobufProduct(product)}, nil
}
