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
Package physicalgood provides the client-side implementation for gRPC [physicalgoodpb.PhysicalGoodServiceClient].
It provides all client-side methods to call server-side business-logic.
*/
package physicalgood

import (
	"context"
	"fmt"
	"log"

	physicalgoodpb "github.com/mikhail5545/proto-go/proto/physical_good/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service provides the client-side implementation for gRPC [physicalgoodpb.PhysicalGoodServiceServer].
// It acts as an adapter between client-side [physicalgoodpb.PhysicalGoodServiceServer] and
// client-side [physicalgoodpb.PhysicalGoodServiceClient] to communicate and transport information.
type Service interface {
	// Get calls [PhysicalGoodServiceServer.Get] method via client connection
	// to retrieve a physical good by its ID.
	// It returns the full physical good object with its associated product details.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Get(ctx context.Context, req *physicalgoodpb.GetRequest) (*physicalgoodpb.GetResponse, error)
	// GetWithDeleted calls [PhysicalGoodServiceServer.GetWithDeleted] method via client connection
	// to retrieve a physical good by its ID, including soft-deleted ones.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	GetWithDeleted(ctx context.Context, req *physicalgoodpb.GetWithDeletedRequest) (*physicalgoodpb.GetWithDeletedResponse, error)
	// GetWithUnpublished calls [PhysicalGoodServiceServer.GetWithUnpublished] method via client connection
	// to retrieve a physical good by its ID, including unpublished ones.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	GetWithUnpublished(ctx context.Context, req *physicalgoodpb.GetWithUnpublishedRequest) (*physicalgoodpb.GetWithUnpublishedResponse, error)
	// List calls [PhysicalGoodServiceServer.List] method via client connection
	// to retrieve a paginated list of all physical goods.
	// The response contains a list of physical goods and the total count.
	List(ctx context.Context, req *physicalgoodpb.ListRequest) (*physicalgoodpb.ListResponse, error)
	// ListDeleted calls [PhysicalGoodServiceServer.ListDeleted] method via client connection
	// to retrieve a paginated list of all soft-deleted physical goods.
	ListDeleted(ctx context.Context, req *physicalgoodpb.ListDeletedRequest) (*physicalgoodpb.ListDeletedResponse, error)
	// ListUnpublished calls [PhysicalGoodServiceServer.ListUnpublished] method via client connection
	// to retrieve a paginated list of all unpublished physical goods.
	ListUnpublished(ctx context.Context, req *physicalgoodpb.ListUnpublishedRequest) (*physicalgoodpb.ListUnpublishedResponse, error)
	// Publish calls [PhysicalGoodServiceServer.Publish] method via client connection
	// to make a physical good and its associated product available in the catalog.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Publish(ctx context.Context, req *physicalgoodpb.PublishRequest) (*physicalgoodpb.PublishResponse, error)
	// Unpublish calls [PhysicalGoodServiceServer.Unpublish] method via client connection
	// to archive a physical good and its associated product from the catalog.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Unpublish(ctx context.Context, req *physicalgoodpb.UnpublishRequest) (*physicalgoodpb.UnpublishResponse, error)
	// Create calls [PhysicalGoodServiceServer.Create] method via client connection
	// to create a new physical good and its associated product.
	// Both are created in an unpublished state.
	//
	// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
	// It returns the ID of the newly created physical good and its associated product.
	Create(ctx context.Context, req *physicalgoodpb.CreateRequest) (*physicalgoodpb.CreateResponse, error)
	// Update calls [PhysicalGoodServiceServer.Update] method via client connection
	// to update physical good fields that have been actually changed.
	// It populates only updated fields in the response.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
	Update(ctx context.Context, req *physicalgoodpb.UpdateRequest) (*physicalgoodpb.UpdateResponse, error)
	// Delete calls [PhysicalGoodServiceServer.Delete] method via client connection
	// to perform a soft-delete on a physical good and its associated product.
	// It also unpublishes them, requiring manual re-publishing after restoration.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Delete(ctx context.Context, req *physicalgoodpb.DeleteRequest) (*physicalgoodpb.DeleteResponse, error)
	// DeletePermanent calls [PhysicalGoodServiceServer.DeletePermanent] method via client connection
	// to permanently delete a physical good and its associated product from the database.
	// This action is irreversible.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	DeletePermanent(ctx context.Context, req *physicalgoodpb.DeletePermanentRequest) (*physicalgoodpb.DeletePermanentResponse, error)
	// Restore calls [PhysicalGoodServiceServer.Restore] method via client connection
	// to restore a soft-deleted physical good and its associated product.
	// The restored records are not automatically published and must be published manually.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Restore(ctx context.Context, req *physicalgoodpb.RestoreRequest) (*physicalgoodpb.RestoreResponse, error)
	// AddImage calls [PhysicalGoodServiceServer.AddImage] method via client connection
	// to add a new image to a physical good. It's called by media-service-go upon successful image upload.
	// It validates the request, checks the image limit and appends the new information.
	//
	// Returns `InvalidArgument` gRPC error if the request payload is invalid/image limit is exceeded.
	// Returns `NotFound` gRPC error if the record is not found.
	AddImage(ctx context.Context, req *physicalgoodpb.AddImageRequest) (*physicalgoodpb.AddImageResponse, error)
	// DeleteImage calls [PhysicalGoodServiceServer.DeleteImage] method via client connection
	// to delete an image from a physical good. It's called by media-service-go upon successful image deletion.
	// The function validates the request and removes the image information from the physical good.
	// This action is irreversable.
	//
	// Returns `InvalidArgument` gRPC error if the request payload is invalid.
	// Returns `NotFound` gRPC error if any of records is not found.
	DeleteImage(ctx context.Context, req *physicalgoodpb.DeleteImageRequest) (*physicalgoodpb.DeleteImageResponse, error)
	// Close tears down connection to the client and all underlying connections.
	Close() error
}

// Client holds [grpc.ClientConn] to connect to the client and
// [physicalgoodpb.PhysicalGoodServiceClient] client to call server-side methods.
type Client struct {
	conn   *grpc.ClientConn
	client physicalgoodpb.PhysicalGoodServiceClient
}

// New creates a new physical good client.
func New(ctx context.Context, addr string, opt ...grpc.CallOption) (Service, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(opt...))
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection: %v", err)
	}
	log.Printf("Connection to physical good service at %s established", addr)

	client := physicalgoodpb.NewPhysicalGoodServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Get calls [PhysicalGoodServiceServer.Get] method via client connection
// to retrieve a physical good by its ID.
// It returns the full physical good object with its associated product details.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Get(ctx context.Context, req *physicalgoodpb.GetRequest) (*physicalgoodpb.GetResponse, error) {
	return c.client.Get(ctx, req)
}

// GetWithDeleted calls [PhysicalGoodServiceServer.GetWithDeleted] method via client connection
// to retrieve a physical good by its ID, including soft-deleted ones.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) GetWithDeleted(ctx context.Context, req *physicalgoodpb.GetWithDeletedRequest) (*physicalgoodpb.GetWithDeletedResponse, error) {
	return c.client.GetWithDeleted(ctx, req)
}

// GetWithUnpublished calls [PhysicalGoodServiceServer.GetWithUnpublished] method via client connection
// to retrieve a physical good by its ID, including unpublished ones.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) GetWithUnpublished(ctx context.Context, req *physicalgoodpb.GetWithUnpublishedRequest) (*physicalgoodpb.GetWithUnpublishedResponse, error) {
	return c.client.GetWithUnpublished(ctx, req)
}

// List calls [PhysicalGoodServiceServer.List] method via client connection
// to retrieve a paginated list of all physical goods.
// The response contains a list of physical goods and the total count.
func (c *Client) List(ctx context.Context, req *physicalgoodpb.ListRequest) (*physicalgoodpb.ListResponse, error) {
	return c.client.List(ctx, req)
}

// ListDeleted calls [PhysicalGoodServiceServer.ListDeleted] method via client connection
// to retrieve a paginated list of all soft-deleted physical goods.
func (c *Client) ListDeleted(ctx context.Context, req *physicalgoodpb.ListDeletedRequest) (*physicalgoodpb.ListDeletedResponse, error) {
	return c.client.ListDeleted(ctx, req)
}

// ListUnpublished calls [PhysicalGoodServiceServer.ListUnpublished] method via client connection
// to retrieve a paginated list of all unpublished physical goods.
func (c *Client) ListUnpublished(ctx context.Context, req *physicalgoodpb.ListUnpublishedRequest) (*physicalgoodpb.ListUnpublishedResponse, error) {
	return c.client.ListUnpublished(ctx, req)
}

// Publish calls [PhysicalGoodServiceServer.Publish] method via client connection
// to make a physical good and its associated product available in the catalog.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Publish(ctx context.Context, req *physicalgoodpb.PublishRequest) (*physicalgoodpb.PublishResponse, error) {
	return c.client.Publish(ctx, req)
}

// Unpublish calls [PhysicalGoodServiceServer.Unpublish] method via client connection
// to archive a physical good and its associated product from the catalog.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Unpublish(ctx context.Context, req *physicalgoodpb.UnpublishRequest) (*physicalgoodpb.UnpublishResponse, error) {
	return c.client.Unpublish(ctx, req)
}

// Create calls [PhysicalGoodServiceServer.Create] method via client connection
// to create a new physical good and its associated product.
// Both are created in an unpublished state.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
// It returns the ID of the newly created physical good and its associated product.
func (c *Client) Create(ctx context.Context, req *physicalgoodpb.CreateRequest) (*physicalgoodpb.CreateResponse, error) {
	return c.client.Create(ctx, req)
}

// Update calls [PhysicalGoodServiceServer.Update] method via client connection
// to update physical good fields that have been actually changed.
// It populates only updated fields in the response.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (c *Client) Update(ctx context.Context, req *physicalgoodpb.UpdateRequest) (*physicalgoodpb.UpdateResponse, error) {
	return c.client.Update(ctx, req)
}

// Delete calls [PhysicalGoodServiceServer.Delete] method via client connection
// to perform a soft-delete on a physical good and its associated product.
// It also unpublishes them, requiring manual re-publishing after restoration.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Delete(ctx context.Context, req *physicalgoodpb.DeleteRequest) (*physicalgoodpb.DeleteResponse, error) {
	return c.client.Delete(ctx, req)
}

// DeletePermanent calls [PhysicalGoodServiceServer.DeletePermanent] method via client connection
// to permanently delete a physical good and its associated product from the database.
// This action is irreversible.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) DeletePermanent(ctx context.Context, req *physicalgoodpb.DeletePermanentRequest) (*physicalgoodpb.DeletePermanentResponse, error) {
	return c.client.DeletePermanent(ctx, req)
}

// Restore calls [PhysicalGoodServiceServer.Restore] method via client connection
// to restore a soft-deleted physical good and its associated product.
// The restored records are not automatically published and must be published manually.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Restore(ctx context.Context, req *physicalgoodpb.RestoreRequest) (*physicalgoodpb.RestoreResponse, error) {
	return c.client.Restore(ctx, req)
}

// AddImage calls [PhysicalGoodServiceServer.AddImage] method via client connection
// to add a new image to a physical good. It's called by media-service-go upon successful image upload.
// It validates the request, checks the image limit and appends the new information.
//
// Returns `InvalidArgument` gRPC error if the request payload is invalid/image limit is exceeded.
// Returns `NotFound` gRPC error if the record is not found.
func (c *Client) AddImage(ctx context.Context, req *physicalgoodpb.AddImageRequest) (*physicalgoodpb.AddImageResponse, error) {
	return c.client.AddImage(ctx, req)
}

// DeleteImage calls [PhysicalGoodServiceServer.DeleteImage] method via client connection
// to delete an image from a physical good. It's called by media-service-go upon successful image deletion.
// The function validates the request and removes the image information from the physical good.
// This action is irreversable.
//
// Returns `InvalidArgument` gRPC error if the request payload is invalid.
// Returns `NotFound` gRPC error if any of records is not found.
func (c *Client) DeleteImage(ctx context.Context, req *physicalgoodpb.DeleteImageRequest) (*physicalgoodpb.DeleteImageResponse, error) {
	return c.client.DeleteImage(ctx, req)
}

// Close tears down connection to the client and all underlying connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
