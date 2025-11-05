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
Package seminar provides the client-side implementation for gRPC [seminarpb.SeminarServiceClient].
It provides all client-side methods to call server-side business-logic.
*/
package seminar

import (
	"context"
	"fmt"
	"log"

	seminarpb "github.com/mikhail5545/proto-go/proto/seminar/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service provides the client-side implementation for gRPC [seminarpb.SeminarServiceClient].
// It acts as an adapter between client-side [seminarpb.SeminarServiceServer] and
// client-side [seminarpb.SeminarServiceClient] to communicate and transport information.
type Service interface {
	// Get calls [SeminarServiceServer.Get] method via client connection
	// to retrieve a seminar by their ID.
	// It returns the full seminar object with all associated product details.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Get(ctx context.Context, req *seminarpb.GetRequest) (*seminarpb.GetResponse, error)
	// GetWithDeleted calls [SeminarServiceServer.GetWithDeleted] method via client connection
	// to retrieve a seminar by their ID, including soft-deleted ones.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	GetWithDeleted(ctx context.Context, req *seminarpb.GetWithDeletedRequest) (*seminarpb.GetWithDeletedResponse, error)
	// GetWithUnpublished calls [SeminarServiceServer.GetWithUnpublished] method via client connection
	// to retrieve a seminar by their ID, including unpublished ones.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	GetWithUnpublished(ctx context.Context, req *seminarpb.GetWithUnpublishedRequest) (*seminarpb.GetWithUnpublishedResponse, error)
	// List calls [SeminarServiceServer.List] method via client connection
	// to retrieve a paginated list of all seminars.
	// The response contains a list of seminars and the total count.
	List(ctx context.Context, req *seminarpb.ListRequest) (*seminarpb.ListResponse, error)
	// ListDeleted calls [SeminarServiceServer.ListDeleted] method via client connection
	// to retrieve a paginated list of all soft-deleted seminars.
	ListDeleted(ctx context.Context, req *seminarpb.ListDeletedRequest) (*seminarpb.ListDeletedResponse, error)
	// ListUnpublished calls [SeminarServiceServer.ListUnpublished] method via client connection
	// to retrieve a paginated list of all unpublished seminars.
	ListUnpublished(ctx context.Context, req *seminarpb.ListUnpublishedRequest) (*seminarpb.ListUnpublishedResponse, error)
	// Create calls [SeminarServiceServer.Create] method via client connection
	// to create a new seminar and its five associated products.
	// All records are created in an unpublished state.
	//
	// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
	// It returns the ID of the newly created seminar and its associated product IDs.
	Create(ctx context.Context, req *seminarpb.CreateRequest) (*seminarpb.CreateResponse, error)
	// Publish calls [SeminarServiceServer.Publish] method via client connection
	// to make a seminar and all its associated products available in the catalog.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Publish(ctx context.Context, req *seminarpb.PublishRequest) (*seminarpb.PublishResponse, error)
	// Unpublish calls [SeminarServiceServer.Unpublish] method via client connection
	// to archive a seminar and all its associated products from the catalog.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Unpublish(ctx context.Context, req *seminarpb.UnpublishRequest) (*seminarpb.UnpublishResponse, error)
	// Update calls [SeminarServiceServer.Update] method via client connection
	// to update seminar fields that have been actually changed. All request fields
	// except ID are optional, so service will update seminar only if at least one field
	// has been updated.
	// It populates only updated fields in the response.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
	Update(ctx context.Context, req *seminarpb.UpdateRequest) (*seminarpb.UpdateResponse, error)
	// Delete calls [SeminarServiceServer.Delete] method via client connection
	// to perform a soft-delete on a seminar and all of its associated products.
	// It also unpublishes them, requiring manual re-publishing after restoration.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Delete(ctx context.Context, req *seminarpb.DeleteRequest) (*seminarpb.DeleteResponse, error)
	// DeletePermanent calls [SeminarServiceServer.DeletePermanent] method via client connection
	// to permanently delete a seminar and all of its associated products from the database.
	// This action is irreversible.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	DeletePermanent(ctx context.Context, req *seminarpb.DeletePermanentRequest) (*seminarpb.DeletePermanentResponse, error)
	// Restore calls [SeminarServiceServer.Restore] method via client connection
	// to restore a soft-deleted seminar and all of its associated products.
	// The restored records are not automatically published and must be published manually.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Restore(ctx context.Context, req *seminarpb.RestoreRequest) (*seminarpb.RestoreResponse, error)

	// Close tears down connection to the client and all underlying connections.
	Close() error
}

// Client holds [grpc.ClientConn] to connect to the client and
// [seminarpb.SeminarServiceClient] client to call server-side methods.
type Client struct {
	conn   *grpc.ClientConn
	client seminarpb.SeminarServiceClient
}

// New creates a new [seminar.Server] client.
func New(ctx context.Context, addr string, opt ...grpc.CallOption) (Service, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(opt...))
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection: %v", err)
	}
	log.Printf("Connection to seminar service at %s established", addr)

	client := seminarpb.NewSeminarServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Get calls [SeminarServiceServer.Get] method via client connection
// to retrieve a seminar by their ID.
// It returns the full seminar object with all associated product details.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Get(ctx context.Context, req *seminarpb.GetRequest) (*seminarpb.GetResponse, error) {
	return c.client.Get(ctx, req)
}

// GetWithDeleted calls [SeminarServiceServer.GetWithDeleted] method via client connection
// to retrieve a seminar by their ID, including soft-deleted ones.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) GetWithDeleted(ctx context.Context, req *seminarpb.GetWithDeletedRequest) (*seminarpb.GetWithDeletedResponse, error) {
	return c.client.GetWithDeleted(ctx, req)
}

// GetWithUnpublished calls [SeminarServiceServer.GetWithUnpublished] method via client connection
// to retrieve a seminar by their ID, including unpublished ones.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) GetWithUnpublished(ctx context.Context, req *seminarpb.GetWithUnpublishedRequest) (*seminarpb.GetWithUnpublishedResponse, error) {
	return c.client.GetWithUnpublished(ctx, req)
}

// List calls [SeminarServiceServer.List] method via client connection
// to retrieve a paginated list of all seminars.
// The response contains a list of seminars and the total count.
func (c *Client) List(ctx context.Context, req *seminarpb.ListRequest) (*seminarpb.ListResponse, error) {
	return c.client.List(ctx, req)
}

// ListDeleted calls [SeminarServiceServer.ListDeleted] method via client connection
// to retrieve a paginated list of all soft-deleted seminars.
func (c *Client) ListDeleted(ctx context.Context, req *seminarpb.ListDeletedRequest) (*seminarpb.ListDeletedResponse, error) {
	return c.client.ListDeleted(ctx, req)
}

// ListUnpublished calls [SeminarServiceServer.ListUnpublished] method via client connection
// to retrieve a paginated list of all unpublished seminars.
func (c *Client) ListUnpublished(ctx context.Context, req *seminarpb.ListUnpublishedRequest) (*seminarpb.ListUnpublishedResponse, error) {
	return c.client.ListUnpublished(ctx, req)
}

// Create calls [SeminarServiceServer.Create] method via client connection
// to create a new seminar and its five associated products.
// All records are created in an unpublished state.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
// It returns the ID of the newly created seminar and its associated product IDs.
func (c *Client) Create(ctx context.Context, req *seminarpb.CreateRequest) (*seminarpb.CreateResponse, error) {
	return c.client.Create(ctx, req)
}

// Publish calls [SeminarServiceServer.Publish] method via client connection
// to make a seminar and all its associated products available in the catalog.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Publish(ctx context.Context, req *seminarpb.PublishRequest) (*seminarpb.PublishResponse, error) {
	return c.client.Publish(ctx, req)
}

// Unpublish calls [SeminarServiceServer.Unpublish] method via client connection
// to archive a seminar and all its associated products from the catalog.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Unpublish(ctx context.Context, req *seminarpb.UnpublishRequest) (*seminarpb.UnpublishResponse, error) {
	return c.client.Unpublish(ctx, req)
}

// Update calls [SeminarServiceServer.Update] method via client connection
// to update seminar fields that have been actually changed. All request fields
// except ID are optional, so service will update seminar only if at least one field
// has been updated.
// It populates only updated fields in the response.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (c *Client) Update(ctx context.Context, req *seminarpb.UpdateRequest) (*seminarpb.UpdateResponse, error) {
	return c.client.Update(ctx, req)
}

// Delete calls [SeminarServiceServer.Delete] method via client connection
// to perform a soft-delete on a seminar and all of its associated products.
// It also unpublishes them, requiring manual re-publishing after restoration.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Delete(ctx context.Context, req *seminarpb.DeleteRequest) (*seminarpb.DeleteResponse, error) {
	return c.client.Delete(ctx, req)
}

// DeletePermanent calls [SeminarServiceServer.DeletePermanent] method via client connection
// to permanently delete a seminar and all of its associated products from the database.
// This action is irreversible.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) DeletePermanent(ctx context.Context, req *seminarpb.DeletePermanentRequest) (*seminarpb.DeletePermanentResponse, error) {
	return c.client.DeletePermanent(ctx, req)
}

// Restore calls [SeminarServiceServer.Restore] method via client connection
// to restore a soft-deleted seminar and all of its associated products.
// The restored records are not automatically published and must be published manually.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Restore(ctx context.Context, req *seminarpb.RestoreRequest) (*seminarpb.RestoreResponse, error) {
	return c.client.Restore(ctx, req)
}

// Close tears down connection to the client and all underlying connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
