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
Package coursepart provides the client-side implementation for gRPC [coursepartpb.CoursePartServiceClient].
It provides all client-side methods to call server-side business-logic.
*/
package coursepart

import (
	"context"
	"fmt"
	"log"

	coursepartpb "github.com/mikhail5545/proto-go/proto/product_service/course_part/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service provides the client-side implementation for gRPC [coursepartpb.CoursePartServiceClient].
// It acts as an adapter between client-side coursepartpb.CoursePartServiceServer] and
// client-side [coursepartpb.CoursePartServiceClient] to communicate and transport information.
type Service interface {
	// Get calls [CoursePartServiceServer.Get] method via client connection
	// to retrieve a course part by their ID.
	// It attemps to retrieve MUXVideo information by calling the media service.
	//
	// If the course part is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	// If call to media service fails, it returns an `Unavaliable` gRPC error.
	Get(ctx context.Context, req *coursepartpb.GetRequest) (*coursepartpb.GetResponse, error)
	// GetWithDeleted calls [CoursePartServiceServer.GetWithDeleted] method via client connection
	// to retrieve a course part by their ID, including soft-deleted ones.
	// It attemps to retrieve MUXVideo information by calling the media service.
	//
	// If the course part is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	// If call to media service fails, it returns an `Unavaliable` gRPC error.
	GetWithDeleted(ctx context.Context, req *coursepartpb.GetWithDeletedRequest) (*coursepartpb.GetWithDeletedResponse, error)
	// GetWithUnpublished calls [CoursePartServiceServer.GetWithUnpublished] method via client connection
	// to retrieve a course part by their ID, including unpublished ones.
	// It attemps to retrieve MUXVideo information by calling the media service.
	//
	// If the course part is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	// If call to media service fails, it returns an `Unavaliable` gRPC error.
	GetWithUnpublished(ctx context.Context, req *coursepartpb.GetWithUnpublishedRequest) (*coursepartpb.GetWithUnpublishedResponse, error)
	// GetReduced calls [CoursePartServiceServer.GetReduced] method via client connection
	// to retrieve a course part by their ID.
	// It does not populate MUXVideo details; the MUXVideo field in the returned course part struct will be nil.
	// This is a lighter version of the `Get` method.
	//
	// If the course is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	GetReduced(ctx context.Context, req *coursepartpb.GetReducedRequest) (*coursepartpb.GetReducedResponse, error)
	// GetReducedWithDeleted calls [CoursePartServiceServer.GetReducedWithDeleted] method via client connection
	// to retrieve a course part by their ID, including soft-deleted ones.
	// It does not populate MUXVideo details; the MUXVideo field in the returned course part struct will be nil.
	// This is a lighter version of the `GetWithDeleted` method.
	//
	// If the course is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	GetWithDeletedReduced(ctx context.Context, req *coursepartpb.GetWithDeletedReducedRequest) (*coursepartpb.GetWithDeletedReducedResponse, error)
	// GetWithUnpublishedReduced calls [CoursePartServiceServer.GetWithUnpublishedReduced] method via client connection
	// to retrieve a course part by their ID, including unpublished, but not soft-deleted ones.
	// It does not populate MUXVideo details; the MUXVideo field in the returned course part struct will be nil.
	// This is a lighter version of the `GetWithUnpublished` method.
	//
	// If the course is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	GetWithUnpublishedReduced(ctx context.Context, req *coursepartpb.GetWithUnpublishedReducedRequest) (*coursepartpb.GetWithUnpublishedReducedResponse, error)
	// List calls [CoursePartServiceServer.List] method via client connection
	// to retrieve a paginated list of all course parts.
	// It does not populate MUXVideo details; the MUXVideo field in the returned course part structs will be nil.
	// The response contains a list of course parts
	// and the total number of course parts in the system.
	List(ctx context.Context, req *coursepartpb.ListRequest) (*coursepartpb.ListResponse, error)
	// ListDeleted calls [CoursePartServiceServer.ListDeleted] method via client connection
	// to retrieve a paginated list of all soft-deleted course parts.
	// It does not populate MUXVideo details; the MUXVideo field in the returned course part structs will be nil.
	// The response contains a list of course parts
	// and the total number of soft-deleted course parts in the system.
	ListDeleted(ctx context.Context, req *coursepartpb.ListDeletedRequest) (*coursepartpb.ListDeletedResponse, error)
	// ListUnpublished calls [CoursePartServiceServer.ListUnpublished] method via client connection
	// to retrieve a paginated list of all unpublished course parts.
	// It does not populate MUXVideo details; the MUXVideo field in the returned course part structs will be nil.
	// The response contains a list of course parts
	// and the total number of unpublished course parts in the system.
	ListUnpublished(ctx context.Context, req *coursepartpb.ListUnpublishedRequest) (*coursepartpb.ListUnpublishedResponse, error)
	// Create calls [CoursePartServiceServer.Create] method via client connection
	// to create a new course part record. It automatically creates an underlying product.
	//
	// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
	// It returns the ID of the newly created course part and its associated product.
	Create(ctx context.Context, req *coursepartpb.CreateRequest) (*coursepartpb.CreateResponse, error)
	// Publish calls [CoursePartServiceServer.Publish] method via client connection
	// to publish a course part.
	// It will fail if parent course is not published.
	//
	// If the course or course part is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	Publish(ctx context.Context, req *coursepartpb.PublishRequest) (*coursepartpb.PublishResponse, error)
	// Unpublish calls [CoursePartServiceServer.Unpublish] method via client connection
	// to unpublish a course part.
	//
	// If the course or course part is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	Unpublish(ctx context.Context, req *coursepartpb.UnpublishRequest) (*coursepartpb.UnpublishResponse, error)
	// Update calls [CoursePartServiceServer.Update] method via client connection
	// to update course part fields that have been actually changed. All request fields
	// except ID are optional, so service will update course only if at least one field
	// has been updated.
	// It populates only updated fields in the response.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
	Update(ctx context.Context, req *coursepartpb.UpdateRequest) (*coursepartpb.UpdateResponse, error)
	// Delete calls [CoursePartServiceServer.Delete] method via client connection
	// to perform a soft-delete on a course part and its associated product.
	// It also unpublishes them, requiring manual re-publishing after restoration.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Delete(ctx context.Context, req *coursepartpb.DeleteRequest) (*coursepartpb.DeleteResponse, error)
	// DeletePermanent calls [CoursePartServiceServer.DeletePermanent] method via client connection
	// to permanently delete a course part and its associated product from the database.
	// This action is irreversible.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	DeletePermanent(ctx context.Context, req *coursepartpb.DeletePermanentRequest) (*coursepartpb.DeletePermanentResponse, error)
	// Restore calls [CoursePartServiceServer.Restore] method via client connection
	// to restore a soft-deleted course part and its associated product.
	// The restored records are not automatically published and must be published manually.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Restore(ctx context.Context, req *coursepartpb.RestoreRequest) (*coursepartpb.RestoreResponse, error)

	// Close tears down connection to the client and all underlying connections.
	Close() error
}

// Client holds [grpc.ClientConn] to connect to the client and
// [coursepartpb.CoursePartServiceClient] client to call server-side methods.
type Client struct {
	conn   *grpc.ClientConn
	client coursepartpb.CoursePartServiceClient
}

// New creates a new [seminar.Server] client.
func New(ctx context.Context, addr string, opt ...grpc.CallOption) (Service, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(opt...))
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection: %v", err)
	}
	log.Printf("Connection to course part service at %s established", addr)

	client := coursepartpb.NewCoursePartServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Get calls [CoursePartServiceServer.Get] method via client connection
// to retrieve a course part by their ID.
// It attemps to retrieve MUXVideo information by calling the media service.
//
// If the course part is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
// If call to media service fails, it returns an `Unavaliable` gRPC error.
func (c *Client) Get(ctx context.Context, req *coursepartpb.GetRequest) (*coursepartpb.GetResponse, error) {
	return c.client.Get(ctx, req)
}

// GetWithDeleted calls [CoursePartServiceServer.GetWithDeleted] method via client connection
// to retrieve a course part by their ID, including soft-deleted ones.
// It attemps to retrieve MUXVideo information by calling the media service.
//
// If the course part is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
// If call to media service fails, it returns an `Unavaliable` gRPC error.
func (c *Client) GetWithDeleted(ctx context.Context, req *coursepartpb.GetWithDeletedRequest) (*coursepartpb.GetWithDeletedResponse, error) {
	return c.client.GetWithDeleted(ctx, req)
}

// GetWithUnpublished calls [CoursePartServiceServer.GetWithUnpublished] method via client connection
// to retrieve a course part by their ID, including unpublished ones.
// It attemps to retrieve MUXVideo information by calling the media service.
//
// If the course part is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
// If call to media service fails, it returns an `Unavaliable` gRPC error.
func (c *Client) GetWithUnpublished(ctx context.Context, req *coursepartpb.GetWithUnpublishedRequest) (*coursepartpb.GetWithUnpublishedResponse, error) {
	return c.client.GetWithUnpublished(ctx, req)
}

// GetReduced calls [CoursePartServiceServer.GetReduced] method via client connection
// to retrieve a course part by their ID.
// It does not populate MUXVideo details; the MUXVideo field in the returned course part struct will be nil.
// This is a lighter version of the `Get` method.
//
// If the course is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
func (c *Client) GetReduced(ctx context.Context, req *coursepartpb.GetReducedRequest) (*coursepartpb.GetReducedResponse, error) {
	return c.client.GetReduced(ctx, req)
}

// GetReducedWithDeleted calls [CoursePartServiceServer.GetReducedWithDeleted] method via client connection
// to retrieve a course part by their ID, including soft-deleted ones.
// It does not populate MUXVideo details; the MUXVideo field in the returned course part struct will be nil.
// This is a lighter version of the `GetWithDeleted` method.
//
// If the course is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
func (c *Client) GetWithDeletedReduced(ctx context.Context, req *coursepartpb.GetWithDeletedReducedRequest) (*coursepartpb.GetWithDeletedReducedResponse, error) {
	return c.client.GetWithDeletedReduced(ctx, req)
}

// GetWithUnpublishedReduced calls [CoursePartServiceServer.GetReducedWithDeleted] method via client connection
// to retrieve a course part by their ID, including unpublished, but not soft-deleted ones.
// It does not populate MUXVideo details; the MUXVideo field in the returned course part struct will be nil.
// This is a lighter version of the `GetWithUnpublished` method.
//
// If the course is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
func (c *Client) GetWithUnpublishedReduced(ctx context.Context, req *coursepartpb.GetWithUnpublishedReducedRequest) (*coursepartpb.GetWithUnpublishedReducedResponse, error) {
	return c.client.GetWithUnpublishedReduced(ctx, req)
}

// List calls [CoursePartServiceServer.List] method via client connection
// to retrieve a paginated list of all course parts.
// It does not populate MUXVideo details; the MUXVideo field in the returned course part structs will be nil.
// The response contains a list of course parts
// and the total number of course parts in the system.
func (c *Client) List(ctx context.Context, req *coursepartpb.ListRequest) (*coursepartpb.ListResponse, error) {
	return c.client.List(ctx, req)
}

// ListDeleted calls [CoursePartServiceServer.ListDeleted] method via client connection
// to retrieve a paginated list of all soft-deleted course parts.
// It does not populate MUXVideo details; the MUXVideo field in the returned course part structs will be nil.
// The response contains a list of course parts
// and the total number of soft-deleted course parts in the system.
func (c *Client) ListDeleted(ctx context.Context, req *coursepartpb.ListDeletedRequest) (*coursepartpb.ListDeletedResponse, error) {
	return c.client.ListDeleted(ctx, req)
}

// ListUnpublished calls [CoursePartServiceServer.ListUnpublished] method via client connection
// to retrieve a paginated list of all unpublished course parts.
// It does not populate MUXVideo details; the MUXVideo field in the returned course part structs will be nil.
// The response contains a list of course parts
// and the total number of unpublished course parts in the system.
func (c *Client) ListUnpublished(ctx context.Context, req *coursepartpb.ListUnpublishedRequest) (*coursepartpb.ListUnpublishedResponse, error) {
	return c.client.ListUnpublished(ctx, req)
}

// Create calls [CoursePartServiceServer.Create] method via client connection
// to create a new course part record. It automatically creates an underlying product.
//
// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
// It returns the ID of the newly created course part and its associated product.
func (c *Client) Create(ctx context.Context, req *coursepartpb.CreateRequest) (*coursepartpb.CreateResponse, error) {
	return c.client.Create(ctx, req)
}

// Publish calls [CoursePartServiceServer.Publish] method via client connection
// to publish a course part.
// It will fail if parent course is not published.
//
// If the course or course part is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
func (c *Client) Publish(ctx context.Context, req *coursepartpb.PublishRequest) (*coursepartpb.PublishResponse, error) {
	return c.client.Publish(ctx, req)
}

// Unpublish calls [CoursePartServiceServer.Unpublish] method via client connection
// to unpublish a course part.
//
// If the course or course part is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
func (c *Client) Unpublish(ctx context.Context, req *coursepartpb.UnpublishRequest) (*coursepartpb.UnpublishResponse, error) {
	return c.client.Unpublish(ctx, req)
}

// Update calls [CoursePartServiceServer.Update] method via client connection
// to update course part fields that have been actually changed. All request fields
// except ID are optional, so service will update course only if at least one field
// has been updated.
// It populates only updated fields in the response.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (c *Client) Update(ctx context.Context, req *coursepartpb.UpdateRequest) (*coursepartpb.UpdateResponse, error) {
	return c.client.Update(ctx, req)
}

// Delete calls [CoursePartServiceServer.Delete] method via client connection
// to perform a soft-delete on a course part and its associated product.
// It also unpublishes them, requiring manual re-publishing after restoration.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Delete(ctx context.Context, req *coursepartpb.DeleteRequest) (*coursepartpb.DeleteResponse, error) {
	return c.client.Delete(ctx, req)
}

// DeletePermanent calls [CoursePartServiceServer.DeletePermanent] method via client connection
// to permanently delete a course part and its associated product from the database.
// This action is irreversible.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) DeletePermanent(ctx context.Context, req *coursepartpb.DeletePermanentRequest) (*coursepartpb.DeletePermanentResponse, error) {
	return c.client.DeletePermanent(ctx, req)
}

// Restore calls [CoursePartServiceServer.Restore] method via client connection
// to restore a soft-deleted course part and its associated product.
// The restored records are not automatically published and must be published manually.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Restore(ctx context.Context, req *coursepartpb.RestoreRequest) (*coursepartpb.RestoreResponse, error) {
	return c.client.Restore(ctx, req)
}

// Close tears down connection to the client and all underlying connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
