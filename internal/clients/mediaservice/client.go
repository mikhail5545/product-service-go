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

package mediaservice

import (
	"context"
	"fmt"
	"log"

	muxpb "github.com/mikhail5545/proto-go/proto/media_service/mux_upload/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is a gRPC client for mux service.
type Client struct {
	conn   *grpc.ClientConn
	client muxpb.MuxUploadServiceClient
}

// NewClient creates a new media service client.
func NewClient(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	log.Printf("gRPC connection to mux service at %s established", addr)

	client := muxpb.NewMuxUploadServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection to the media service.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
