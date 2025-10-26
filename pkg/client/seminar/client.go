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
	// It returns the full seminar object.
	// If the seminar is not found, it returns a `NotFound` gRPC error.
	Get(ctx context.Context, req *seminarpb.GetRequest) (*seminarpb.GetResponse, error)
	// List calls [SeminarServiceServer.List] method via client connection
	// to retrieve a paginated list of all seminars.
	// The response contains a list of seminars
	// and the total number of seminars in the system.
	List(ctx context.Context, req *seminarpb.ListRequest) (*seminarpb.ListResponse, error)
	// Create calls [SeminarServiceServer.Create] method via client connection
	// to create a new seminar record, typically in the process of direct seminar
	// creation. It automatically creates all underlying products and populdates they're `name` and `description`
	// fields from [models.Seminar.Name] and [models.Seminar.Description] if not provided.
	//
	// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
	// It returns newly created seminar model with all fields.
	Create(ctx context.Context, req *seminarpb.CreateRequest) (*seminarpb.CreateResponse, error)
	// Update calls [SeminarServiceServer.Update] method via client connection
	// to update seminar fields that have been acually changed. All request fields
	// except ID are optional, so service will update seminar only if at least one field
	// has been updated.
	//
	// It populates only updated fields in the response along with the `fieldmaskpb.UpdateMask` which contains
	// paths to updated fields.
	Update(ctx context.Context, req *seminarpb.UpdateRequest) (*seminarpb.UpdateResponse, error)
	// Delete calls [SeminarServiceServer.Delete] method via client connection
	// to completely delete Seminar record from the system.
	Delete(ctx context.Context, req *seminarpb.DeleteRequest) (*seminarpb.DeleteResponse, error)

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
// It returns the full seminar object.
// If the seminar is not found, it returns a `NotFound` gRPC error.
func (c *Client) Get(ctx context.Context, req *seminarpb.GetRequest) (*seminarpb.GetResponse, error) {
	return c.client.Get(ctx, req)
}

// List calls [SeminarServiceServer.List] method via client connection
// to retrieve a paginated list of all seminars.
// The response contains a list of seminars
// and the total number of seminars in the system.
func (c *Client) List(ctx context.Context, req *seminarpb.ListRequest) (*seminarpb.ListResponse, error) {
	return c.client.List(ctx, req)
}

// Create calls [SeminarServiceServer.Create] method via client connection
// to create a new seminar record, typically in the process of direct seminar
// creation. It automatically creates all underlying products and populdates they're `name` and `description`
// fields from [models.Seminar.Name] and [models.Seminar.Description] if not provided.
//
// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
// It returns newly created seminar model with all fields.
func (c *Client) Create(ctx context.Context, req *seminarpb.CreateRequest) (*seminarpb.CreateResponse, error) {
	return c.client.Create(ctx, req)
}

// Update calls [SeminarServiceServer.Update] method via client connection
// to update seminar fields that have been acually changed. All request fields
// except ID are optional, so service will update seminar only if at least one field
// has been updated.
//
// It populates only updated fields in the response along with the `fieldmaskpb.UpdateMask` which contains
// paths to updated fields.
func (c *Client) Update(ctx context.Context, req *seminarpb.UpdateRequest) (*seminarpb.UpdateResponse, error) {
	return c.client.Update(ctx, req)
}

// Delete calls [SeminarServiceServer.Delete] method via client connection
// to completely delete Seminar record from the system.
func (c *Client) Delete(ctx context.Context, req *seminarpb.DeleteRequest) (*seminarpb.DeleteResponse, error) {
	return c.client.Delete(ctx, req)
}

// Close tears down connection to the client and all underlying connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
