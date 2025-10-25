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

package server

import (
	"context"

	"github.com/mikhail5545/product-service-go/internal/models"
	"github.com/mikhail5545/product-service-go/internal/services/product"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	productpb "github.com/mikhail5545/proto-go/proto/product/v0"
)

type Server struct {
	productpb.UnimplementedProductServiceServer
	service *product.Service
}

func New(s *product.Service) *Server {
	return &Server{service: s}
}

func (s *Server) Get(ctx context.Context, req *productpb.GetRequest) (*productpb.GetResponse, error) {
	product, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return &productpb.GetResponse{Product: types.ProductToProtobuf(product)}, nil
}

func (s *Server) List(ctx context.Context, req *productpb.ListRequest) (*productpb.ListResponse, error) {
	products, total, err := s.service.List(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbProducts []*productpb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, types.ProductToProtobuf(&product))
	}

	return &productpb.ListResponse{Products: pbProducts, Total: total}, nil
}

func (s *Server) ListByType(ctx context.Context, req *productpb.ListByTypeRequest) (*productpb.ListByTypeResponse, error) {
	products, total, err := s.service.ListByType(ctx, req.GetProductType(), int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}
	var pbProducts []*productpb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, types.ProductToProtobuf(&product))
	}

	return &productpb.ListByTypeResponse{Products: pbProducts, Total: total}, nil
}

func (s *Server) Create(ctx context.Context, req *productpb.CreateRequest) (*productpb.CreateResponse, error) {
	product := &models.Product{
		Name:             req.GetName(),
		Description:      req.GetDescription(),
		Price:            req.GetPrice(),
		Amount:           int(req.GetAmount()),
		ShippingRequired: req.GetShippingRequired(),
	}

	product, err := s.service.Create(ctx, product)
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return &productpb.CreateResponse{Product: types.ProductToProtobuf(product)}, nil
}

func (s *Server) Update(ctx context.Context, req *productpb.UpdateRequest) (*productpb.UpdateResponse, error) {
	product := &models.Product{
		Name:             req.GetName(),
		Description:      req.GetDescription(),
		Price:            req.GetPrice(),
		Amount:           int(req.GetAmount()),
		ImageUrl:         req.GetImageUrl(),
		ShippingRequired: req.GetShippingRequired(),
	}

	updates, err := s.service.Update(ctx, product, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return types.ProductToProtobufUpdate(updates), nil
}
