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
	product := &models.Product{
		Name:             req.GetName(),
		Description:      req.GetDescription(),
		Price:            req.GetPrice(),
		Amount:           int(req.GetAmount()),
		ShippingRequired: req.GetShippingRequired(),
	}

	product, err := s.productService.CreateProduct(ctx, product)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &productpb.CreateProductResponse{Product: utils.ConvertToProtobufProduct(product)}, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error) {
	product := &models.Product{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Amount:      int(req.GetAmount()),
	}

	product, err := s.productService.UpdateProduct(ctx, product, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &productpb.UpdateProductResponse{Product: utils.ConvertToProtobufProduct(product)}, nil
}
