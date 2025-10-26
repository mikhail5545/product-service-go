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
Package trainingsession provides the client-side implementation for gRPC [trainingsessionpb.TrainingSessionServiceClient].
It provides all client-side methods to call server-side business-logic.
*/
package trainingsession

import (
	"context"
	"fmt"
	"log"

	trainingsessionpb "github.com/mikhail5545/proto-go/proto/training_session/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service provides the client-side implementation for gRPC [trainingsessionpb.TrainingSessionServiceClient].
// It acts as an adapter between client-side [trainingsessionpb.TrainingSessionServiceServer] and
// client-side [trainingsessionpb.TrainingSessionServiceClient] to communicate and transport information.
type Service interface {
	// Get calls [TrainingSessionServiceServer.Get] method via client connection
	// to retrieve a training session by their ID.
	// It returns the full training session object.
	// If the training session is not found, it returns a `NotFound` gRPC error.
	Get(ctx context.Context, req *trainingsessionpb.GetRequest) (*trainingsessionpb.GetResponse, error)
	// List calls [TrainingSessionServiceServer.List] method via client connection
	// to retrieve a paginated list of all training sessions.
	// The response contains a list of full training session objects.
	// and the total number of training sessions in the system.
	List(ctx context.Context, req *trainingsessionpb.ListRequest) (*trainingsessionpb.ListResponse, error)
	// Create calls [TrainingSessionServiceServer.Create] method via client connection
	// to create a new training session record, typically in the process of direct training session
	// creation. It automatically creates an underlying product.
	//
	// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
	// It returns newly created course training session with all fields.
	Create(ctx context.Context, req *trainingsessionpb.CreateRequest) (*trainingsessionpb.CreateResponse, error)
	// Update calls [TrainingSessionServiceServer.Update] method via client connection
	// to update training session fields that have been acually changed. All request fields
	// except ID are optional, so service will update training session only if at least one field
	// has been updated.
	//
	// It populates only updated fields in the response along with the `fieldmaskpb.UpdateMask` which contains
	// paths to updated fields.
	Update(ctx context.Context, req *trainingsessionpb.UpdateRequest) (*trainingsessionpb.UpdateResponse, error)
	// Delete calls [TrainingSessionServiceServer.Delete] method via client connection
	// to completely deletes training session from the system.
	Delete(ctx context.Context, req *trainingsessionpb.DeleteRequest) (*trainingsessionpb.DeleteResponse, error)

	// Close tears down connection to the client and all underlying connections.
	Close() error
}

// Client holds [grpc.ClientConn] to connect to the client and
// [trainingsessionpb.TrainingSessionServiceClient] client to call server-side methods.
type Client struct {
	conn   *grpc.ClientConn
	client trainingsessionpb.TrainingSessionServiceClient
}

// New creates a new [trainingsession.Server] client.
func New(ctx context.Context, addr string, opt ...grpc.CallOption) (Service, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(opt...))
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection: %v", err)
	}
	log.Printf("Connection to course service at %s established", addr)

	client := trainingsessionpb.NewTrainingSessionServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Get calls [TrainingSessionServiceServer.Get] method via client connection
// to retrieve a training session by their ID.
// It returns the full training session object.
// If the training session is not found, it returns a `NotFound` gRPC error.
func (c *Client) Get(ctx context.Context, req *trainingsessionpb.GetRequest) (*trainingsessionpb.GetResponse, error) {
	return c.client.Get(ctx, req)
}

// List calls [TrainingSessionServiceServer.List] method via client connection
// to retrieve a paginated list of all training sessions.
// The response contains a list of full training session objects.
// and the total number of training sessions in the system.
func (c *Client) List(ctx context.Context, req *trainingsessionpb.ListRequest) (*trainingsessionpb.ListResponse, error) {
	return c.client.List(ctx, req)
}

// Create calls [TrainingSessionServiceServer.Create] method via client connection
// to create a new training session record, typically in the process of direct training session
// creation. It automatically creates an underlying product.
//
// If request payload not satisfies service expectations, it returns a `InvalidArgument` gRPC error.
// It returns newly created course training session with all fields.
func (c *Client) Create(ctx context.Context, req *trainingsessionpb.CreateRequest) (*trainingsessionpb.CreateResponse, error) {
	return c.client.Create(ctx, req)
}

// Update calls [TrainingSessionServiceServer.Update] method via client connection
// to update training session fields that have been acually changed. All request fields
// except ID are optional, so service will update training session only if at least one field
// has been updated.
//
// It populates only updated fields in the response along with the `fieldmaskpb.UpdateMask` which contains
// paths to updated fields.
func (c *Client) Update(ctx context.Context, req *trainingsessionpb.UpdateRequest) (*trainingsessionpb.UpdateResponse, error) {
	return c.client.Update(ctx, req)
}

// Delete calls [TrainingSessionServiceServer.Delete] method via client connection
// to completely deletes training session from the system.
func (c *Client) Delete(ctx context.Context, req *trainingsessionpb.DeleteRequest) (*trainingsessionpb.DeleteResponse, error) {
	return c.client.Delete(ctx, req)
}

// Close tears down connection to the client and all underlying connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
