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
Package image provides the client-side implementation for gRPC [imagepb.ImageServiceServerClient].
It provides all client-side methods to call server-side business-logic.
*/
package image

import (
	"context"
	"fmt"
	"log"

	imagepb "github.com/mikhail5545/proto-go/proto/product_service/image/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service provides the client-side implementation for gRPC [imagepb.ImageServiceServerClient].
// It acts as an adapter between server-side [imagepb.ImageServiceServer] and
// client-side [imagepb.ImageServiceServerClient] to communicate and transport information.
type Service interface {
	// Add calls [ImageServiceServer.Add] via client connection
	// to add an image for the owner depending on specified ownerType.
	//
	// Returns an `InvalidArgument` gRPC error if the request payload or ownerType is invalid.
	// Returns a `NotFound` gRPC error if the image not found on the owner or the owner is not found.
	Add(ctx context.Context, req *imagepb.AddRequest) (*imagepb.AddResponse, error)
	// Delete calls [ImageServiceServer.Delete] via client connection
	// to delete an image from the owner depending on specified ownerType.
	//
	// Returns an `InvalidArgument` gRPC error if the request payload or ownerType is invalid.
	// Returns a `NotFound` gRPC error if the image not found on the owner or the owner is not found.
	Delete(ctx context.Context, req *imagepb.DeleteRequest) (*imagepb.DeleteResponse, error)
	// AddBatch calls [ImageServiceServer.AddBatch] via client connection
	// to add an image for the batch of owners depending on specified ownerType.
	//
	// Returns the number of affected owners.
	// Returns an `InvalidArgument` gRPC error if the request payload or ownerType is invalid.
	// Returns a `NotFound` gRPC error if none or the owners were found.
	AddBatch(ctx context.Context, req *imagepb.AddBatchRequest) (*imagepb.AddBatchResponse, error)
	// DeleteBatch calls [ImageServiceServer.DeleteBatch] via client connection
	// to delete an image from the batch of owners depending on specified ownerType.
	//
	// Returns the number of affected owners.
	// Returns an `InvalidArgument` gRPC error if the request payload or ownerType is invalid.
	// Returns a `NotFound` gRPC error if none or the owners were found or associations were not found.
	DeleteBatch(ctx context.Context, req *imagepb.DeleteBatchRequest) (*imagepb.DeleteBatchResponse, error)
	// Close tears down connection to the client and all underlying connections.
	Close() error
}

// Client holds [grpc.ClientConn] to connect to the client and
// [imagepb.ImageServiceServerClient] client to call server-side methods.
type Client struct {
	conn   *grpc.ClientConn
	client imagepb.ImageServiceClient
}

// New creates a new [image.Server] client.
func New(ctx context.Context, addr string, opt ...grpc.CallOption) (Service, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(opt...))
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection: %w", err)
	}
	log.Printf("Connection to image service at %s established", addr)

	client := imagepb.NewImageServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Add calls [ImageServiceServer.Add] via client connection
// to add an image for the owner depending on specified ownerType.
//
// Returns an `InvalidArgument` gRPC error if the request payload or ownerType is invalid.
// Returns a `NotFound` gRPC error if the image not found on the owner or the owner is not found.
func (c *Client) Add(ctx context.Context, req *imagepb.AddRequest) (*imagepb.AddResponse, error) {
	return c.client.Add(ctx, req)
}

// Delete calls [ImageServiceServer.Delete] via client connection
// to delete an image from the owner depending on specified ownerType.
//
// Returns an `InvalidArgument` gRPC error if the request payload or ownerType is invalid.
// Returns a `NotFound` gRPC error if the image not found on the owner or the owner is not found.
func (c *Client) Delete(ctx context.Context, req *imagepb.DeleteRequest) (*imagepb.DeleteResponse, error) {
	return c.client.Delete(ctx, req)
}

// AddBatch calls [ImageServiceServer.AddBatch] via client connection
// to add an image for the batch of owners depending on specified ownerType.
//
// Returns the number of affected owners.
// Returns an `InvalidArgument` gRPC error if the request payload or ownerType is invalid.
// Returns a `NotFound` gRPC error if none or the owners were found.
func (c *Client) AddBatch(ctx context.Context, req *imagepb.AddBatchRequest) (*imagepb.AddBatchResponse, error) {
	return c.client.AddBatch(ctx, req)
}

// DeleteBatch calls [ImageServiceServer.DeleteBatch] via client connection
// to delete an image from the batch of owners depending on specified ownerType.
//
// Returns the number of affected owners.
// Returns an `InvalidArgument` gRPC error if the request payload or ownerType is invalid.
// Returns a `NotFound` gRPC error if none or the owners were found or associations were not found.
func (c *Client) DeleteBatch(ctx context.Context, req *imagepb.DeleteBatchRequest) (*imagepb.DeleteBatchResponse, error) {
	return c.client.DeleteBatch(ctx, req)
}

// Close tears down connection to the client and all underlying connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
