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

/*
Package product provides the implementation of the gRPC
[productpb.ProductServiceServer] interface and provides
various operations for Product models.
*/
package product

import (
	"context"

	"github.com/mikhail5545/product-service-go/internal/models"
	"github.com/mikhail5545/product-service-go/internal/services/product"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	productpb "github.com/mikhail5545/proto-go/proto/product/v0"
)

// Server implements the gRPC [productpb.ProductServiceServer] interface and provides
// operations for Product models. It acts as an adapter between the gRPC transport layer
// and the server-layer buusiness logic of microservice, defined in the [product.Service].
//
// For more information about underlying gRPC server, see [github.com/mikhail5545/proto-go].
type Server struct {
	productpb.UnimplementedProductServiceServer
	// service provides service-layer business logic for Product operations.
	service *product.Service
}

// New creates a new Server instance.
func New(s *product.Service) *Server {
	return &Server{service: s}
}

// Get retrieves a product by their ID.
// It returns the full product object.
// If the product is not found, it returns a `NotFound` gRPC error.
func (s *Server) Get(ctx context.Context, req *productpb.GetRequest) (*productpb.GetResponse, error) {
	product, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCError(err)
	}

	return &productpb.GetResponse{Product: types.ProductToProtobuf(product)}, nil
}

// List retrieves a paginated list of all products.
// The response contains a list of products
// and the total number of products in the system.
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

// ListByType retrieves a paginated list of all products by their `type` field.
// The response contains a list of products that have specified `type`
// and the total number of products with that `type` in the system.
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

// Create creates a new product record, typically in the process of direct product
// creation. To create underlying product for other data types ([models.Course], [models.Seminar]),
// use specified Create methods of this data types.
//
// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
// It returns newly created Product model with all fields.
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

// Update updates product fields that have been acually changed. All request fields
// except ID are optional, so service will update product only if at least one field
// has been updated.
//
// It populates only updated fields in the response along with the `fieldmaskpb.UpdateMask` which contains
// paths to updated fields.
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
