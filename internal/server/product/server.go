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

	productservice "github.com/mikhail5545/product-service-go/internal/services/product"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
	"github.com/mikhail5545/product-service-go/internal/util/types"
	productpb "github.com/mikhail5545/proto-go/proto/product_service/product/v0"
	"google.golang.org/grpc"
)

// Server implements the gRPC [productpb.ProductServiceServer] interface and provides
// operations for Product models. It acts as an adapter between the gRPC transport layer
// and the server-layer buusiness logic of microservice, defined in the [product.Service].
//
// For more information about underlying gRPC server, see [github.com/mikhail5545/proto-go].
type Server struct {
	productpb.UnimplementedProductServiceServer
	// service provides service-layer business logic for Product operations.
	service productservice.Service
}

// New creates a new Server instance.
func New(svc productservice.Service) *Server {
	return &Server{service: svc}
}

// Register registers the course server with a gRPC server instance.
func Register(s *grpc.Server, svc productservice.Service) {
	productpb.RegisterProductServiceServer(s, New(svc))
}

// Get retrieves a product by their ID.
// It returns the full product object.
// If the product is not found, it returns a `NotFound` gRPC error.
func (s *Server) Get(ctx context.Context, req *productpb.GetRequest) (*productpb.GetResponse, error) {
	product, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	return &productpb.GetResponse{Product: types.ProductToProtobuf(product)}, nil
}

// List retrieves a paginated list of all products.
// The response contains a list of products
// and the total number of products in the system.
func (s *Server) List(ctx context.Context, req *productpb.ListRequest) (*productpb.ListResponse, error) {
	products, total, err := s.service.List(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.HandleServiceError(err)
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
func (s *Server) ListByDetailsType(ctx context.Context, req *productpb.ListByDetailsTypeRequest) (*productpb.ListByDetailsTypeResponse, error) {
	products, total, err := s.service.ListByDetailsType(ctx, req.GetDetailsType(), int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, errors.HandleServiceError(err)
	}
	var pbProducts []*productpb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, types.ProductToProtobuf(&product))
	}
	return &productpb.ListByDetailsTypeResponse{Products: pbProducts, Total: total}, nil
}
