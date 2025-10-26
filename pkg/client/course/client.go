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
	// If the course is not found, it returns a `NotFound` gRPC error.
	Get(ctx context.Context, req *coursepb.GetRequest) (*coursepb.GetResponse, error)
	// GetReduced calls [CourseServiceServer.GetReduced] method via client connection
	// to retrieve a course by their ID.
	// It returns the reduced course object (not all fields are presented, especially it does not provide
	// list of [models.CoursePart] for this course).
	// If the course is not found, it returns a `NotFound` gRPC error.
	GetReduced(ctx context.Context, req *coursepb.GetReducedRequest) (*coursepb.GetReducedResponse, error)
	// List calls [CourseServiceServer.List] method via client connection
	// to retrieve a paginated list of all courses.
	// The response contains a list of courses
	// and the total number of courses in the system.
	List(ctx context.Context, req *coursepb.ListRequest) (*coursepb.ListResponse, error)
	// Create calls [CourseServiceServer.Create] method via client connection
	// to create a new course record, typically in the process of direct course
	// creation. It automatically creates all underlying products and populdates they're `name` and `description`
	// fields from [models.Course.Name] and [models.Course.Description] if not provided.
	//
	// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
	// It returns newly created course model with all fields.
	Create(ctx context.Context, req *coursepb.CreateRequest) (*coursepb.CreateResponse, error)
	// Update calls [CourseServiceServer.Update] method via client connection
	// to update course fields that have been acually changed. All request fields
	// except ID are optional, so service will update course only if at least one field
	// has been updated.
	//
	// It populates only updated fields in the response along with the `fieldmaskpb.UpdateMask` which contains
	// paths to updated fields.
	Update(ctx context.Context, req *coursepb.UpdateRequest) (*coursepb.UpdateResponse, error)

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
// If the course is not found, it returns a `NotFound` gRPC error.
func (c *Client) Get(ctx context.Context, req *coursepb.GetRequest) (*coursepb.GetResponse, error) {
	return c.client.Get(ctx, req)
}

// GetReduced calls [CourseServiceServer.GetReduced] method via client connection
// to retrieve a course by their ID.
// It returns the reduced course object (not all fields are presented, especially it does not provide
// list of [models.CoursePart] for this course).
// If the course is not found, it returns a `NotFound` gRPC error.
func (c *Client) GetReduced(ctx context.Context, req *coursepb.GetReducedRequest) (*coursepb.GetReducedResponse, error) {
	return c.client.GetReduced(ctx, req)
}

// List calls [CourseServiceServer.List] method via client connection
// to retrieve a paginated list of all courses.
// The response contains a list of courses
// and the total number of courses in the system.
func (c *Client) List(ctx context.Context, req *coursepb.ListRequest) (*coursepb.ListResponse, error) {
	return c.client.List(ctx, req)
}

// Create calls [CourseServiceServer.Create] method via client connection
// to create a new course record, typically in the process of direct course
// creation. It automatically creates all underlying products and populdates they're `name` and `description`
// fields from [models.Course.Name] and [models.Course.Description] if not provided.
//
// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
// It returns newly created course model with all fields.
func (c *Client) Create(ctx context.Context, req *coursepb.CreateRequest) (*coursepb.CreateResponse, error) {
	return c.client.Create(ctx, req)
}

// Update calls [CourseServiceServer.Update] method via client connection
// to update course fields that have been acually changed. All request fields
// except ID are optional, so service will update course only if at least one field
// has been updated.
//
// It populates only updated fields in the response along with the `fieldmaskpb.UpdateMask` which contains
// paths to updated fields.
func (c *Client) Update(ctx context.Context, req *coursepb.UpdateRequest) (*coursepb.UpdateResponse, error) {
	return c.client.Updae(ctx, req)
}

// Close tears down connection to the client and all underlying connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
