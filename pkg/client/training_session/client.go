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

	trainingsessionpb "github.com/mikhail5545/proto-go/proto/product_service/training_session/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service provides the client-side implementation for gRPC [trainingsessionpb.TrainingSessionServiceClient].
// It acts as an adapter between client-side [trainingsessionpb.TrainingSessionServiceServer] and
// client-side [trainingsessionpb.TrainingSessionServiceClient] to communicate and transport information.
type Service interface {
	// Get calls [TrainingSessionServiceServer.Get] method via client connection
	// to retrieve a training session by their ID.
	// It returns the full training session object with its associated product details.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Get(ctx context.Context, req *trainingsessionpb.GetRequest) (*trainingsessionpb.GetResponse, error)
	// GetWithDeleted calls [TrainingSessionServiceServer.GetWithDeleted] method via client connection
	// to retrieve a training session by their ID, including soft-deleted ones.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	GetWithDeleted(ctx context.Context, req *trainingsessionpb.GetWithDeletedRequest) (*trainingsessionpb.GetWithDeletedResponse, error)
	// GetWithUnpublished calls [TrainingSessionServiceServer.GetWithUnpublished] method via client connection
	// to retrieve a training session by their ID, including unpublished ones.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	GetWithUnpublished(ctx context.Context, req *trainingsessionpb.GetWithUnpublishedRequest) (*trainingsessionpb.GetWithUnpublishedResponse, error)
	// List calls [TrainingSessionServiceServer.List] method via client connection
	// to retrieve a paginated list of all training sessions.
	// The response contains a list of training sessions and the total count.
	List(ctx context.Context, req *trainingsessionpb.ListRequest) (*trainingsessionpb.ListResponse, error)
	// ListDeleted calls [TrainingSessionServiceServer.ListDeleted] method via client connection
	// to retrieve a paginated list of all soft-deleted training sessions.
	ListDeleted(ctx context.Context, req *trainingsessionpb.ListDeletedRequest) (*trainingsessionpb.ListDeletedResponse, error)
	// ListUnpublished calls [TrainingSessionServiceServer.ListUnpublished] method via client connection
	// to retrieve a paginated list of all unpublished training sessions.
	ListUnpublished(ctx context.Context, req *trainingsessionpb.ListUnpublishedRequest) (*trainingsessionpb.ListUnpublishedResponse, error)
	// Publish calls [TrainingSessionServiceServer.Publish] method via client connection
	// to make a training session and its associated product available in the catalog.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Publish(ctx context.Context, req *trainingsessionpb.PublishRequest) (*trainingsessionpb.PublishResponse, error)
	// Unpublish calls [TrainingSessionServiceServer.Unpublish] method via client connection
	// to archive a training session and its associated product from the catalog.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Unpublish(ctx context.Context, req *trainingsessionpb.UnpublishRequest) (*trainingsessionpb.UnpublishResponse, error)
	// Create calls [TrainingSessionServiceServer.Create] method via client connection
	// to create a new training session and its associated product.
	// Both are created in an unpublished state.
	//
	// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
	// It returns the ID of the newly created training session and its associated product.
	Create(ctx context.Context, req *trainingsessionpb.CreateRequest) (*trainingsessionpb.CreateResponse, error)
	// Update calls [TrainingSessionServiceServer.Update] method via client connection
	// to update training session fields that have been actually changed. All request fields
	// except ID are optional, so service will update training session only if at least one field
	// has been updated.
	// It populates only updated fields in the response.
	//
	// Returns a `NotFound` gRPC error if the record is not found.
	// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
	Update(ctx context.Context, req *trainingsessionpb.UpdateRequest) (*trainingsessionpb.UpdateResponse, error)
	// Delete calls [TrainingSessionServiceServer.Delete] method via client connection
	// to perform a soft-delete on a training session and its associated product.
	// It also unpublishes them, requiring manual re-publishing after restoration.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Delete(ctx context.Context, req *trainingsessionpb.DeleteRequest) (*trainingsessionpb.DeleteResponse, error)
	// DeletePermanent calls [TrainingSessionServiceServer.DeletePermanent] method via client connection
	// to permanently delete a training session and its associated product from the database.
	// This action is irreversible.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	DeletePermanent(ctx context.Context, req *trainingsessionpb.DeletePermanentRequest) (*trainingsessionpb.DeletePermanentResponse, error)
	// Restore calls [TrainingSessionServiceServer.Restore] method via client connection
	// to restore a soft-deleted training session and its associated product.
	// The restored records are not automatically published and must be published manually.
	//
	// Returns a `NotFound` gRPC error if any of the records are not found.
	// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
	Restore(ctx context.Context, req *trainingsessionpb.RestoreRequest) (*trainingsessionpb.RestoreResponse, error)
	// AddImage calls [TrainingSessionServiceServer.AddImage] method via client connection
	// to add a new image to a training session. It's called by media-service-go upon successful image upload.
	// It validates the request, checks the image limit and appends the new information.
	//
	// Returns `InvalidArgument` gRPC error if the request payload is invalid/image limit is exceeded.
	// Returns `NotFound` gRPC error if the record is not found.
	AddImage(ctx context.Context, req *trainingsessionpb.AddImageRequest) (*trainingsessionpb.AddImageResponse, error)
	// DeleteImage calls [TrainingSessionServiceServer.DeleteImage] method via client connection
	// to delete an image from a training session. It's called by media-service-go upon successful image deletion.
	// The function validates the request and removes the image information from the training session.
	// This action is irreversable.
	//
	// Returns `InvalidArgument` gRPC error if the request payload is invalid.
	// Returns `NotFound` gRPC error if any of records is not found.
	DeleteImage(ctx context.Context, req *trainingsessionpb.DeleteImageRequest) (*trainingsessionpb.DeleteImageResponse, error)
	// AddImageBatch calls [TrainingSessionServiceServer.AddImageBatch] method via client connection
	// to add an image for a batch of training sessions. It's called by media-service-go
	// upon successful image uplaod.
	//
	// Returns the number of affected training sessions.
	// Returns `InvalidArgument` gRPC error if the request payload is invalid.
	// Returns `NotFound` gRPC error none of the training sessions were found.
	AddImageBatch(ctx context.Context, req *trainingsessionpb.AddImageBatchRequest) (*trainingsessionpb.AddImageBatchResponse, error)
	// DeleteImageBatch calls [TrainingSessionServiceServer.DeleteImageBatch] method via client connection
	// to delete an image from a batch of training sessions. It's called by media-service-go
	// upon successful image deletion.
	//
	// Returns the number of affected training sessions.
	// Returns `InvalidArgument` gRPC error if the request payload is invalid.
	// Returns `NotFound` gRPC error none of the training sessions were found or the image was not found.
	DeleteImageBatch(ctx context.Context, req *trainingsessionpb.DeleteImageBatchRequest) (*trainingsessionpb.DeleteImageBatchResponse, error)

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
	log.Printf("Connection to training session service at %s established", addr)

	client := trainingsessionpb.NewTrainingSessionServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Get calls [TrainingSessionServiceServer.Get] method via client connection
// to retrieve a training session by their ID.
// It returns the full training session object with its associated product details.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Get(ctx context.Context, req *trainingsessionpb.GetRequest) (*trainingsessionpb.GetResponse, error) {
	return c.client.Get(ctx, req)
}

// GetWithDeleted calls [TrainingSessionServiceServer.GetWithDeleted] method via client connection
// to retrieve a training session by their ID, including soft-deleted ones.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) GetWithDeleted(ctx context.Context, req *trainingsessionpb.GetWithDeletedRequest) (*trainingsessionpb.GetWithDeletedResponse, error) {
	return c.client.GetWithDeleted(ctx, req)
}

// GetWithUnpublished calls [TrainingSessionServiceServer.GetWithUnpublished] method via client connection
// to retrieve a training session by their ID, including unpublished ones.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) GetWithUnpublished(ctx context.Context, req *trainingsessionpb.GetWithUnpublishedRequest) (*trainingsessionpb.GetWithUnpublishedResponse, error) {
	return c.client.GetWithUnpublished(ctx, req)
}

// List calls [TrainingSessionServiceServer.List] method via client connection
// to retrieve a paginated list of all training sessions.
// The response contains a list of training sessions and the total count.
func (c *Client) List(ctx context.Context, req *trainingsessionpb.ListRequest) (*trainingsessionpb.ListResponse, error) {
	return c.client.List(ctx, req)
}

// ListDeleted calls [TrainingSessionServiceServer.ListDeleted] method via client connection
// to retrieve a paginated list of all soft-deleted training sessions.
func (c *Client) ListDeleted(ctx context.Context, req *trainingsessionpb.ListDeletedRequest) (*trainingsessionpb.ListDeletedResponse, error) {
	return c.client.ListDeleted(ctx, req)
}

// ListUnpublished calls [TrainingSessionServiceServer.ListUnpublished] method via client connection
// to retrieve a paginated list of all unpublished training sessions.
func (c *Client) ListUnpublished(ctx context.Context, req *trainingsessionpb.ListUnpublishedRequest) (*trainingsessionpb.ListUnpublishedResponse, error) {
	return c.client.ListUnpublished(ctx, req)
}

// Publish calls [TrainingSessionServiceServer.Publish] method via client connection
// to make a training session and its associated product available in the catalog.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Publish(ctx context.Context, req *trainingsessionpb.PublishRequest) (*trainingsessionpb.PublishResponse, error) {
	return c.client.Publish(ctx, req)
}

// Unpublish calls [TrainingSessionServiceServer.Unpublish] method via client connection
// to archive a training session and its associated product from the catalog.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Unpublish(ctx context.Context, req *trainingsessionpb.UnpublishRequest) (*trainingsessionpb.UnpublishResponse, error) {
	return c.client.Unpublish(ctx, req)
}

// Create calls [TrainingSessionServiceServer.Create] method via client connection
// to create a new training session and its associated product.
// Both are created in an unpublished state.
//
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
// It returns the ID of the newly created training session and its associated product.
func (c *Client) Create(ctx context.Context, req *trainingsessionpb.CreateRequest) (*trainingsessionpb.CreateResponse, error) {
	return c.client.Create(ctx, req)
}

// Update calls [TrainingSessionServiceServer.Update] method via client connection
// to update training session fields that have been actually changed. All request fields
// except ID are optional, so service will update training session only if at least one field
// has been updated.
// It populates only updated fields in the response.
//
// Returns a `NotFound` gRPC error if the record is not found.
// Returns an `InvalidArgument` gRPC error if the request payload is invalid.
func (c *Client) Update(ctx context.Context, req *trainingsessionpb.UpdateRequest) (*trainingsessionpb.UpdateResponse, error) {
	return c.client.Update(ctx, req)
}

// Delete calls [TrainingSessionServiceServer.Delete] method via client connection
// to perform a soft-delete on a training session and its associated product.
// It also unpublishes them, requiring manual re-publishing after restoration.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Delete(ctx context.Context, req *trainingsessionpb.DeleteRequest) (*trainingsessionpb.DeleteResponse, error) {
	return c.client.Delete(ctx, req)
}

// DeletePermanent calls [TrainingSessionServiceServer.DeletePermanent] method via client connection
// to permanently delete a training session and its associated product from the database.
// This action is irreversible.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) DeletePermanent(ctx context.Context, req *trainingsessionpb.DeletePermanentRequest) (*trainingsessionpb.DeletePermanentResponse, error) {
	return c.client.DeletePermanent(ctx, req)
}

// Restore calls [TrainingSessionServiceServer.Restore] method via client connection
// to restore a soft-deleted training session and its associated product.
// The restored records are not automatically published and must be published manually.
//
// Returns a `NotFound` gRPC error if any of the records are not found.
// Returns an `InvalidArgument` gRPC error if the provided ID is not a valid UUID.
func (c *Client) Restore(ctx context.Context, req *trainingsessionpb.RestoreRequest) (*trainingsessionpb.RestoreResponse, error) {
	return c.client.Restore(ctx, req)
}

// AddImage calls [TrainingSessionServiceServer.AddImage] method via client connection
// to add a new image to a training session. It's called by media-service-go upon successful image upload.
// It validates the request, checks the image limit and appends the new information.
//
// Returns `InvalidArgument` gRPC error if the request payload is invalid/image limit is exceeded.
// Returns `NotFound` gRPC error if the record is not found.
func (c *Client) AddImage(ctx context.Context, req *trainingsessionpb.AddImageRequest) (*trainingsessionpb.AddImageResponse, error) {
	return c.client.AddImage(ctx, req)
}

// DeleteImage calls [TrainingSessionServiceServer.DeleteImage] method via client connection
// to delete an image from a training session. It's called by media-service-go upon successful image deletion.
// The function validates the request and removes the image information from the training session.
// This action is irreversable.
//
// Returns `InvalidArgument` gRPC error if the request payload is invalid.
// Returns `NotFound` gRPC error if any of records is not found.
func (c *Client) DeleteImage(ctx context.Context, req *trainingsessionpb.DeleteImageRequest) (*trainingsessionpb.DeleteImageResponse, error) {
	return c.client.DeleteImage(ctx, req)
}

// AddImageBatch calls [TrainingSessionServiceServer.AddImageBatch] method via client connection
// to add an image for a batch of training sessions. It's called by media-service-go
// upon successful image uplaod.
//
// Returns the number of affected training sessions.
// Returns `InvalidArgument` gRPC error if the request payload is invalid.
// Returns `NotFound` gRPC error none of the training sessions were found.
func (c *Client) AddImageBatch(ctx context.Context, req *trainingsessionpb.AddImageBatchRequest) (*trainingsessionpb.AddImageBatchResponse, error) {
	return c.client.AddImageBatch(ctx, req)
}

// DeleteImageBatch calls [TrainingSessionServiceServer.DeleteImageBatch] method via client connection
// to delete an image from a batch of training sessions. It's called by media-service-go
// upon successful image deletion.
//
// Returns the number of affected training sessions.
// Returns `InvalidArgument` gRPC error if the request payload is invalid.
// Returns `NotFound` gRPC error none of the training sessions were found or the image was not found.
func (c *Client) DeleteImageBatch(ctx context.Context, req *trainingsessionpb.DeleteImageBatchRequest) (*trainingsessionpb.DeleteImageBatchResponse, error) {
	return c.client.DeleteImageBatch(ctx, req)
}

// Close tears down connection to the client and all underlying connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
