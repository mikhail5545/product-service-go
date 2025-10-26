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

	coursepartpb "github.com/mikhail5545/proto-go/proto/course_part/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service provides the client-side implementation for gRPC [coursepartpb.CoursePartServiceClient].
// It acts as an adapter between client-side coursepartpb.CoursePartServiceServer] and
// client-side [coursepartpb.CoursePartServiceClient] to communicate and transport information.
type Service interface {
	// Get calls [CoursePartServiceServer.Get] method via client connection
	// to retrieve a course part by their ID.
	// It returns the full course part object.
	// If the course part is not found, it returns a `NotFound` gRPC error.
	Get(ctx context.Context, req *coursepartpb.GetRequest) (*coursepartpb.GetPartResponse, error)
	// Unimplemented
	AddVideo(ctx context.Context, req *coursepartpb.AddVideoRequest) (*coursepartpb.AddVideoResponse, error)

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
// It returns the full course part object.
// If the course part is not found, it returns a `NotFound` gRPC error.
func (c *Client) Get(ctx context.Context, req *coursepartpb.GetRequest) (*coursepartpb.GetPartResponse, error) {
	return c.client.Get(ctx, req)
}

// Unimplemented
func (c *Client) AddVideo(ctx context.Context, req *coursepartpb.AddVideoRequest) (*coursepartpb.AddVideoResponse, error) {
	return c.client.AddVideo(ctx, req)
}

// Close tears down connection to the client and all underlying connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
