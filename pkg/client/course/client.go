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
Package course provides the client-side implementation for gRPC [coursepb.CourseServiceClient].
It provides all client-side methods to call server-side business-logic.
*/
package course

import (
	"context"
	"fmt"
	"log"

	coursepb "github.com/mikhail5545/proto-go/proto/course/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service provides the client-side implementation for gRPC [coursepb.CourseServiceClient].
// It acts as an adapter between client-side [coursepb.CourseServiceServer] and
// client-side [coursepb.CourseServiceClient] to communicate and transport information.
type Service interface {
	// Get calls [CourseServiceServer.Get] method via client connection
	// to retrieve a course by their ID.
	// It returns the full course object.
	//
	// If the course is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	Get(ctx context.Context, req *coursepb.GetRequest) (*coursepb.GetResponse, error)
	// GetWithDeleted calls [CourseServiceServer.GetWithDeleted] method via client connection
	// to retrieve a course by their ID, including soft-deleted ones.
	//
	// If the course is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	GetWithDeleted(ctx context.Context, req *coursepb.GetWithDeletedRequest) (*coursepb.GetWithDeletedResponse, error)
	// GetWithUnpublished calls [CourseServiceServer.GetWithUnpublished] method via client connection
	// to retrieve a course by their ID, including unpublished ones.
	//
	// If the course is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	GetWithUnpublished(ctx context.Context, req *coursepb.GetWithUnpublishedRequest) (*coursepb.GetWithUnpublishedResponse, error)
	// GetReduced calls [CourseServiceServer.GetReduced] method via client connection
	// to retrieve a course by their ID.
	// It returns the reduced course object (without course parts).
	//
	// If the course is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	GetReduced(ctx context.Context, req *coursepb.GetReducedRequest) (*coursepb.GetReducedResponse, error)
	// GetReducedWithDeleted calls [CourseServiceServer.GetReducedWithDeleted] method via client connection
	// to retrieve a course by their ID, including soft-deleted ones (without course parts).
	//
	// If the course is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	GetReducedWithDeleted(ctx context.Context, req *coursepb.GetReducedWithDeletedRequest) (*coursepb.GetReducedWithDeletedResponse, error)
	// List calls [CourseServiceServer.List] method via client connection
	// to retrieve a paginated list of all courses.
	// The response contains a list of courses
	// and the total number of courses in the system.
	List(ctx context.Context, req *coursepb.ListRequest) (*coursepb.ListResponse, error)
	// ListDeleted calls [CourseServiceServer.ListDeleted] method via client connection
	// to retrieve a paginated list of all soft-deleted courses.
	ListDeleted(ctx context.Context, req *coursepb.ListDeletedRequest) (*coursepb.ListDeletedResponse, error)
	// ListUnpublished calls [CourseServiceServer.ListUnpublished] method via client connection
	// to retrieve a paginated list of all unpublished courses.
	ListUnpublished(ctx context.Context, req *coursepb.ListUnpublishedRequest) (*coursepb.ListUnpublishedResponse, error)
	// Publish calls [CourseServiceServer.Publish] method via client connection
	// to publish a course.
	//
	// If the course is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	Publish(ctx context.Context, req *coursepb.PublishRequest) (*coursepb.PublishResponse, error)
	// Unpublish calls [CourseServiceServer.Unpublish] method via client connection
	// to unpublish a course.
	//
	// If the course is not found, it returns a `NotFound` gRPC error.
	// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
	Unpublish(ctx context.Context, req *coursepb.UnpublishRequest) (*coursepb.UnpublishResponse, error)
	// Create calls [CourseServiceServer.Create] method via client connection
	// to create a new course record. It automatically creates an underlying product.
	//
	// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
	// It returns the ID of the newly created course and its associated product.
	Create(ctx context.Context, req *coursepb.CreateRequest) (*coursepb.CreateResponse, error)
	// Update calls [CourseServiceServer.Update] method via client connection
	// to update course fields that have been actually changed. All request fields
	// except ID are optional, so service will update course only if at least one field
	// has been updated.
	// It populates only updated fields in the response.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
	Update(ctx context.Context, req *coursepb.UpdateRequest) (*coursepb.UpdateResponse, error)
	// Delete calls [CourseServiceServer.Delete] method via client connection
	// to perform a soft-delete on a course, its associated product, and all its course parts.
	// It also unpublishes them, requiring manual re-publishing after restoration.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Delete(ctx context.Context, req *coursepb.DeleteRequest) (*coursepb.DeleteResponse, error)
	// DeletePermanent calls [CourseServiceServer.DeletePermanent] method via client connection
	// to permanently delete a course, its associated product, and all its course parts from the database.
	// This action is irreversible.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	DeletePermanent(ctx context.Context, req *coursepb.DeletePermanentRequest) (*coursepb.DeletePermanentResponse, error)
	// Restore calls [CourseServiceServer.Restore] method via client connection
	// to restore a soft-deleted course, its associated product, and all its course parts.
	// The restored records are not automatically published and must be published manually.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Restore(ctx context.Context, req *coursepb.RestoreRequest) (*coursepb.RestoreResponse, error)

	// Close tears down connection to the client and all underlying connections.
	Close() error
}

// Client holds [grpc.ClientConn] to connect to the client and
// [coursepb.CourseServiceClient] client to call server-side methods.
type Client struct {
	conn   *grpc.ClientConn
	client coursepb.CourseServiceClient
}

// New creates a new [course.Server] client.
func New(ctx context.Context, addr string, opt ...grpc.CallOption) (Service, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(opt...))
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection: %v", err)
	}
	log.Printf("Connection to course service at %s established", addr)

	client := coursepb.NewCourseServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Get calls [CourseServiceServer.Get] method via client connection
// to retrieve a course by their ID.
// It returns the full course object.
//
// If the course is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
func (c *Client) Get(ctx context.Context, req *coursepb.GetRequest) (*coursepb.GetResponse, error) {
	return c.client.Get(ctx, req)
}

// GetWithDeleted calls [CourseServiceServer.GetWithDeleted] method via client connection
// to retrieve a course by their ID, including soft-deleted ones.
//
// If the course is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
func (c *Client) GetWithDeleted(ctx context.Context, req *coursepb.GetWithDeletedRequest) (*coursepb.GetWithDeletedResponse, error) {
	return c.client.GetWithDeleted(ctx, req)
}

// GetWithUnpublished calls [CourseServiceServer.GetWithUnpublished] method via client connection
// to retrieve a course by their ID, including unpublished ones.
//
// If the course is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
func (c *Client) GetWithUnpublished(ctx context.Context, req *coursepb.GetWithUnpublishedRequest) (*coursepb.GetWithUnpublishedResponse, error) {
	return c.client.GetWithUnpublished(ctx, req)
}

// GetReduced calls [CourseServiceServer.GetReduced] method via client connection
// to retrieve a course by their ID.
// It returns the reduced course object (without course parts).
//
// If the course is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
func (c *Client) GetReduced(ctx context.Context, req *coursepb.GetReducedRequest) (*coursepb.GetReducedResponse, error) {
	return c.client.GetReduced(ctx, req)
}

// GetReducedWithDeleted calls [CourseServiceServer.GetReducedWithDeleted] method via client connection
// to retrieve a course by their ID, including soft-deleted ones (without course parts).
//
// If the course is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
func (c *Client) GetReducedWithDeleted(ctx context.Context, req *coursepb.GetReducedWithDeletedRequest) (*coursepb.GetReducedWithDeletedResponse, error) {
	return c.client.GetReducedWithDeleted(ctx, req)
}

// List calls [CourseServiceServer.List] method via client connection
// to retrieve a paginated list of all courses.
// The response contains a list of courses
// and the total number of courses in the system.
func (c *Client) List(ctx context.Context, req *coursepb.ListRequest) (*coursepb.ListResponse, error) {
	return c.client.List(ctx, req)
}

// ListDeleted calls [CourseServiceServer.ListDeleted] method via client connection
// to retrieve a paginated list of all soft-deleted courses.
func (c *Client) ListDeleted(ctx context.Context, req *coursepb.ListDeletedRequest) (*coursepb.ListDeletedResponse, error) {
	return c.client.ListDeleted(ctx, req)
}

// ListUnpublished calls [CourseServiceServer.ListUnpublished] method via client connection
// to retrieve a paginated list of all unpublished courses.
func (c *Client) ListUnpublished(ctx context.Context, req *coursepb.ListUnpublishedRequest) (*coursepb.ListUnpublishedResponse, error) {
	return c.client.ListUnpublished(ctx, req)
}

// Publish calls [CourseServiceServer.Publish] method via client connection
// to publish a course.
//
// If the course is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
func (c *Client) Publish(ctx context.Context, req *coursepb.PublishRequest) (*coursepb.PublishResponse, error) {
	return c.client.Publish(ctx, req)
}

// Unpublish calls [CourseServiceServer.Unpublish] method via client connection
// to unpublish a course.
//
// If the course is not found, it returns a `NotFound` gRPC error.
// If the provided ID is not a valid UUID, it returns an `InvalidArgument` gRPC error.
func (c *Client) Unpublish(ctx context.Context, req *coursepb.UnpublishRequest) (*coursepb.UnpublishResponse, error) {
	return c.client.Unpublish(ctx, req)
}

// Create calls [CourseServiceServer.Create] method via client connection
// to create a new course record. It automatically creates an underlying product.
//
// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
// It returns the ID of the newly created course and its associated product.
func (c *Client) Create(ctx context.Context, req *coursepb.CreateRequest) (*coursepb.CreateResponse, error) {
	return c.client.Create(ctx, req)
}

// Update calls [CourseServiceServer.Update] method via client connection
// to update course fields that have been actually changed. All request fields
// except ID are optional, so service will update course only if at least one field
// has been updated.
// It populates only updated fields in the response.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (c *Client) Update(ctx context.Context, req *coursepb.UpdateRequest) (*coursepb.UpdateResponse, error) {
	return c.client.Update(ctx, req)
}

// Delete calls [CourseServiceServer.Delete] method via client connection
// to perform a soft-delete on a course, its associated product, and all its course parts.
// It also unpublishes them, requiring manual re-publishing after restoration.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Delete(ctx context.Context, req *coursepb.DeleteRequest) (*coursepb.DeleteResponse, error) {
	return c.client.Delete(ctx, req)
}

// DeletePermanent calls [CourseServiceServer.DeletePermanent] method via client connection
// to permanently delete a course, its associated product, and all its course parts from the database.
// This action is irreversible.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) DeletePermanent(ctx context.Context, req *coursepb.DeletePermanentRequest) (*coursepb.DeletePermanentResponse, error) {
	return c.client.DeletePermanent(ctx, req)
}

// Restore calls [CourseServiceServer.Restore] method via client connection
// to restore a soft-deleted course, its associated product, and all its course parts.
// The restored records are not automatically published and must be published manually.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Restore(ctx context.Context, req *coursepb.RestoreRequest) (*coursepb.RestoreResponse, error) {
	return c.client.Restore(ctx, req)
}

// Close tears down connection to the client and all underlying connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
