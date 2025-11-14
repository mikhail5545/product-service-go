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
Package video provides the client-side implementation for gRPC [videopb.VideoServiceServerClient].
It provides all client-side methods to call server-side business-logic.
*/
package video

import (
	"context"
	"fmt"
	"log"

	videopb "github.com/mikhail5545/proto-go/proto/product_service/video/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service provides the client-side implementation for gRPC [videopb.VideoServiceClient].
// It acts as an adapter between server-side [videopb.VideoServiceServer] and
// client-side [videopb.VideoServiceClient] to communicate and transport information.
// See more details about [underlying protobuf services].
//
// [underlying protobuf services]: https://github.com/mikhail5545/proto-go
type Service interface {
	// Add calls [VideoServiceServer.Add] via client connection
	// to associate a video with a single owner for specified owner type.
	// If there was another video, associated with this owner, it will be replaced with the new one. It also
	// should be deassociated in the corresponding service separately. This function handles only local owner-video relations.
	// It first validates that the video exists in the media service.
	//
	// Returns an `InvalidArgument` gRPC error if the request payload is invalid or video already associated with this owner.
	// Returns a `NotFound` gRPC error if any of the video is not found or the owner is not found.
	Add(ctx context.Context, req *videopb.AddRequest) (*videopb.AddResponse, error)
	// Remove calls [VideoServiceServer.Remove] via client connection
	// to disassociate a video from a single owner for specified owner type.
	// This function handles only local owner-video relations.
	// Owner should be also deassociated from the video in the corresponding service.
	//
	// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
	// Returns a `NotFound` gRPC error if any of the video is not found or the owner is not found.
	Remove(ctx context.Context, req *videopb.RemoveRequest) (*videopb.RemoveResponse, error)
	// GetOwner calls [VideoServiceServer.GetOwner] via client connection
	// to retrieve a single owner information including unpublished ones.
	// Returns minimal necessary owner information. If more owner information is needed,
	// specific owner gRPC service's Get method should be called.
	//
	// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
	// Returns a `NotFound` gRPC error if owner is not found.
	GetOwner(ctx context.Context, req *videopb.GetOwnerRequest) (*videopb.GetOwnerResponse, error)
	// Close tears down connection to the client and all underlying connections.
	Close() error
}

// Client holds [grpc.ClientConn] to connect to the client and
// [videopb.VideoServiceClient] client to call server-side methods.
type Client struct {
	conn   *grpc.ClientConn
	client videopb.VideoServiceClient
}

// New creates a new [video.Server] client.
func New(ctx context.Context, addr string, opt ...grpc.CallOption) (Service, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(opt...))
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection: %w", err)
	}
	log.Printf("Connection to image service at %s established", addr)

	client := videopb.NewVideoServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Add calls [VideoServiceServer.Add] via client connection
// to associate a video with a single owner for specified owner type.
// If there was another video, associated with this owner, it will be replaced with the new one. It also
// should be deassociated in the corresponding service separately. This function handles only local owner-video relations.
// It first validates that the video exists in the media service.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid or video already associated with this owner.
// Returns a `NotFound` gRPC error if any of the video is not found or the owner is not found.
func (c *Client) Add(ctx context.Context, req *videopb.AddRequest) (*videopb.AddResponse, error) {
	return c.client.Add(ctx, req)
}

// Remove calls [VideoServiceServer.Remove] via client connection
// to disassociate a video from a single owner for specified owner type.
// This function handles only local owner-video relations.
// Owner should be also deassociated from the video in the corresponding service.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
// Returns a `NotFound` gRPC error if any of the video is not found or the owner is not found.
func (c *Client) Remove(ctx context.Context, req *videopb.RemoveRequest) (*videopb.RemoveResponse, error) {
	return c.client.Remove(ctx, req)
}

// GetOwner calls [VideoServiceServer.GetOwner] via client connection
// to retrieve a single owner information including unpublished ones.
// Returns minimal necessary owner information. If more owner information is needed,
// specific owner gRPC service's Get method should be called.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
// Returns a `NotFound` gRPC error if owner is not found.
func (c *Client) GetOwner(ctx context.Context, req *videopb.GetOwnerRequest) (*videopb.GetOwnerResponse, error) {
	return c.client.GetOwner(ctx, req)
}

// Close tears down connection to the client and all underlying connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
