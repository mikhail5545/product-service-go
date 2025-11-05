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
Package product provides the client-side implementation for gRPC [productpb.ProductServiceClient].
It provides all client-side methods to call server-side business-logic.
*/
package product

import (
	"context"
	"fmt"
	"log"

	productpb "github.com/mikhail5545/proto-go/proto/product/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service provides the client-side implementation for gRPC [productpb.ProductServiceClient].
// It acts as an adapter between client-side [productpb.ProductServiceServer] and
// client-side [productpb.ProductSesrviceClient] to communicate and transport information.
type Service interface {
	// Get retrieves a product by their ID.
	// It returns the full product object.
	// If the product is not found, it returns a `NotFound` gRPC error.
	Get(ctx context.Context, req *productpb.GetRequest) (*productpb.GetResponse, error)
	// List retrieves a paginated list of all products.
	// The response contains a list of products
	// and the total number of products in the system.
	List(ctx context.Context, req *productpb.ListRequest) (*productpb.ListResponse, error)
	// ListByType retrieves a paginated list of all products by their `type` field.
	// The response contains a list of products that have specified `type`
	// and the total number of products with that `type` in the system.
	ListByDetailsType(ctx context.Context, req *productpb.ListByDetailsTypeRequest) (*productpb.ListByDetailsTypeResponse, error)

	// Close tears down connection to the client and all underlying connections.
	Close() error
}

// Client [holds grpc.ClientConn] to connect to the client and
// [productpb.ProductServiceClient] client to call server-side methods.
type Client struct {
	conn   *grpc.ClientConn
	client productpb.ProductServiceClient
}

// New creates a new Product service server client.
func New(ctx context.Context, addr string, opt ...grpc.CallOption) (Service, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(opt...))
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection: %v", err)
	}
	log.Printf("Connection to product service at %s established", addr)

	client := productpb.NewProductServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Get calls [ProductServiceServer.Get] method via client connection
// to retrieve a product by their ID.
// It returns the full product object.
// If the product is not found, it returns a `NotFound` gRPC error.
func (c *Client) Get(ctx context.Context, req *productpb.GetRequest) (*productpb.GetResponse, error) {
	return c.client.Get(ctx, req)
}

// List calls [ProductServiceServer.List] method via client connection
// to retrieve a paginated list of all products.
// The response contains a list of products
// and the total number of products in the system.
func (c *Client) List(ctx context.Context, req *productpb.ListRequest) (*productpb.ListResponse, error) {
	return c.client.List(ctx, req)
}

// ListByDetailsType calls [ProductServiceServer.ListByType] method via client connection
// to retrieve a paginated list of all products by their `DetailsType` field.
// The response contains a list of products that have specified `type`
// and the total number of products with that `type` in the system.
func (c *Client) ListByDetailsType(ctx context.Context, req *productpb.ListByDetailsTypeRequest) (*productpb.ListByDetailsTypeResponse, error) {
	return c.client.ListByDetailsType(ctx, req)
}

// Close tears down connection to the client and all underlying connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
