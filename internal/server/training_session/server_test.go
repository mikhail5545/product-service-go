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

package trainingsession

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"github.com/google/uuid"
	trainingsession "github.com/mikhail5545/product-service-go/internal/models/training_session"
	trainingsessionmodel "github.com/mikhail5545/product-service-go/internal/models/training_session"
	trainingsessionservice "github.com/mikhail5545/product-service-go/internal/services/training_session"
	trainingsessionmock "github.com/mikhail5545/product-service-go/internal/test/services/training_session_mock"
	trainingsessionpb "github.com/mikhail5545/proto-go/proto/product_service/training_session/v0"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func setupTestServer(t *testing.T) (trainingsessionpb.TrainingSessionServiceClient, *trainingsessionmock.MockService, func()) {
	t.Helper()

	// 1. Create mock controller and mock service
	ctrl := gomock.NewController(t)
	mockService := trainingsessionmock.NewMockService(ctrl)

	// 2. Create an in-memory listener
	lis := bufconn.Listen(1024 * 1024)

	// 3. Create and register the gRPC server
	s := grpc.NewServer()
	Register(s, mockService)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// 4. Create a client connection to the in-memory server
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.NewClient("passthrough:///", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	client := trainingsessionpb.NewTrainingSessionServiceClient(conn)

	// Teardown function
	cleanup := func() {
		ctrl.Finish()
		conn.Close()
		s.Stop()
	}

	return client, mockService, cleanup
}

func TestServer_Get(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID := uuid.New().String()
	productID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &trainingsessionmodel.TrainingSessionDetails{
			TrainingSession: &trainingsession.TrainingSession{
				ID:   tsID,
				Name: "Training session name",
			},
			Price:     99.99,
			ProductID: productID,
		}

		mockService.EXPECT().Get(gomock.Any(), tsID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.Get(context.Background(), &trainingsessionpb.GetRequest{Id: tsID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, tsID, res.GetTrainingSessionDetails().GetTrainingSession().GetId())
		assert.Equal(t, productID, res.GetTrainingSessionDetails().GetProductId())
		assert.Equal(t, float32(99.99), res.GetTrainingSessionDetails().GetPrice())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Get(gomock.Any(), tsID).Return(nil, trainingsessionservice.ErrNotFound)

		// Act
		res, err := client.Get(context.Background(), &trainingsessionpb.GetRequest{Id: tsID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Training session not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Get(gomock.Any(), invalidID).Return(nil, trainingsessionservice.ErrInvalidArgument)

		// Act
		res, err := client.Get(context.Background(), &trainingsessionpb.GetRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid training session ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().Get(gomock.Any(), tsID).Return(nil, svcErr)

		// Act
		res, err := client.Get(context.Background(), &trainingsessionpb.GetRequest{Id: tsID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Database error")
	})
}

func TestServer_GetWithDeleted(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID := uuid.New().String()
	productID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &trainingsessionmodel.TrainingSessionDetails{
			TrainingSession: &trainingsession.TrainingSession{
				ID:   tsID,
				Name: "Training session name",
			},
			Price:     99.99,
			ProductID: productID,
		}

		mockService.EXPECT().GetWithDeleted(gomock.Any(), tsID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &trainingsessionpb.GetWithDeletedRequest{Id: tsID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, tsID, res.GetTrainingSessionDetails().GetTrainingSession().GetId())
		assert.Equal(t, productID, res.GetTrainingSessionDetails().GetProductId())
		assert.Equal(t, float32(99.99), res.GetTrainingSessionDetails().GetPrice())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().GetWithDeleted(gomock.Any(), tsID).Return(nil, trainingsessionservice.ErrNotFound)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &trainingsessionpb.GetWithDeletedRequest{Id: tsID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Training session not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().GetWithDeleted(gomock.Any(), invalidID).Return(nil, trainingsessionservice.ErrInvalidArgument)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &trainingsessionpb.GetWithDeletedRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid training session ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().GetWithDeleted(gomock.Any(), tsID).Return(nil, svcErr)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &trainingsessionpb.GetWithDeletedRequest{Id: tsID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Database error")
	})
}

func TestServer_GetWithUnpublished(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID := uuid.New().String()
	productID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &trainingsessionmodel.TrainingSessionDetails{
			TrainingSession: &trainingsession.TrainingSession{
				ID:   tsID,
				Name: "Training session name",
			},
			Price:     99.99,
			ProductID: productID,
		}

		mockService.EXPECT().GetWithUnpublished(gomock.Any(), tsID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &trainingsessionpb.GetWithUnpublishedRequest{Id: tsID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, tsID, res.GetTrainingSessionDetails().GetTrainingSession().GetId())
		assert.Equal(t, productID, res.GetTrainingSessionDetails().GetProductId())
		assert.Equal(t, float32(99.99), res.GetTrainingSessionDetails().GetPrice())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), tsID).Return(nil, trainingsessionservice.ErrNotFound)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &trainingsessionpb.GetWithUnpublishedRequest{Id: tsID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Training session not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), invalidID).Return(nil, trainingsessionservice.ErrInvalidArgument)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &trainingsessionpb.GetWithUnpublishedRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid training session ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), tsID).Return(nil, svcErr)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &trainingsessionpb.GetWithUnpublishedRequest{Id: tsID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Database error")
	})
}

func TestServer_List(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID_1 := uuid.New().String()
	productID_1 := uuid.New().String()
	tsID_2 := uuid.New().String()
	productID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := []trainingsessionmodel.TrainingSessionDetails{
			{
				TrainingSession: &trainingsessionmodel.TrainingSession{
					ID:   tsID_1,
					Name: "Training session 1 name",
				},
				Price:     99.99,
				ProductID: productID_1,
			},
			{
				TrainingSession: &trainingsessionmodel.TrainingSession{
					ID:   tsID_2,
					Name: "Training session 2 name",
				},
				Price:     199.99,
				ProductID: productID_2,
			},
		}

		mockService.EXPECT().List(gomock.Any(), limit, offset).Return(expectedDetails, int64(2), nil).Times(1)

		// Act
		res, err := client.List(context.Background(), &trainingsessionpb.ListRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetTrainingSessionsDetails()[0].GetTrainingSession().GetName(), expectedDetails[0].TrainingSession.Name)
		assert.Equal(t, res.GetTrainingSessionsDetails()[0].GetTrainingSession().GetId(), expectedDetails[0].TrainingSession.ID)
		assert.Equal(t, res.GetTrainingSessionsDetails()[0].GetPrice(), expectedDetails[0].Price)
		assert.Equal(t, res.GetTrainingSessionsDetails()[0].GetProductId(), expectedDetails[0].ProductID)
		assert.Equal(t, res.GetTrainingSessionsDetails()[1].GetTrainingSession().GetName(), expectedDetails[1].TrainingSession.Name)
		assert.Equal(t, res.GetTrainingSessionsDetails()[1].GetTrainingSession().GetId(), expectedDetails[1].TrainingSession.ID)
		assert.Equal(t, res.GetTrainingSessionsDetails()[1].GetPrice(), expectedDetails[1].Price)
		assert.Equal(t, res.GetTrainingSessionsDetails()[1].GetProductId(), expectedDetails[1].ProductID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().List(gomock.Any(), limit, offset).Return(nil, int64(0), svcErr).Times(1)

		// Act
		res, err := client.List(context.Background(), &trainingsessionpb.ListRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get training sessions")
	})
}

func TestServer_ListDeleted(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID_1 := uuid.New().String()
	productID_1 := uuid.New().String()
	tsID_2 := uuid.New().String()
	productID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := []trainingsessionmodel.TrainingSessionDetails{
			{
				TrainingSession: &trainingsessionmodel.TrainingSession{
					ID:   tsID_1,
					Name: "Training session 1 name",
				},
				Price:     99.99,
				ProductID: productID_1,
			},
			{
				TrainingSession: &trainingsessionmodel.TrainingSession{
					ID:   tsID_2,
					Name: "Training session 2 name",
				},
				Price:     199.99,
				ProductID: productID_2,
			},
		}

		mockService.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(expectedDetails, int64(2), nil).Times(1)

		// Act
		res, err := client.ListDeleted(context.Background(), &trainingsessionpb.ListDeletedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetTrainingSessionsDetails()[0].GetTrainingSession().GetName(), expectedDetails[0].TrainingSession.Name)
		assert.Equal(t, res.GetTrainingSessionsDetails()[0].GetTrainingSession().GetId(), expectedDetails[0].TrainingSession.ID)
		assert.Equal(t, res.GetTrainingSessionsDetails()[0].GetPrice(), expectedDetails[0].Price)
		assert.Equal(t, res.GetTrainingSessionsDetails()[0].GetProductId(), expectedDetails[0].ProductID)
		assert.Equal(t, res.GetTrainingSessionsDetails()[1].GetTrainingSession().GetName(), expectedDetails[1].TrainingSession.Name)
		assert.Equal(t, res.GetTrainingSessionsDetails()[1].GetTrainingSession().GetId(), expectedDetails[1].TrainingSession.ID)
		assert.Equal(t, res.GetTrainingSessionsDetails()[1].GetPrice(), expectedDetails[1].Price)
		assert.Equal(t, res.GetTrainingSessionsDetails()[1].GetProductId(), expectedDetails[1].ProductID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(nil, int64(0), svcErr).Times(1)

		// Act
		res, err := client.ListDeleted(context.Background(), &trainingsessionpb.ListDeletedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get training sessions")
	})
}

func TestServer_ListUnpublished(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID_1 := uuid.New().String()
	productID_1 := uuid.New().String()
	tsID_2 := uuid.New().String()
	productID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := []trainingsessionmodel.TrainingSessionDetails{
			{
				TrainingSession: &trainingsessionmodel.TrainingSession{
					ID:   tsID_1,
					Name: "Training session 1 name",
				},
				Price:     99.99,
				ProductID: productID_1,
			},
			{
				TrainingSession: &trainingsessionmodel.TrainingSession{
					ID:   tsID_2,
					Name: "Training session 2 name",
				},
				Price:     199.99,
				ProductID: productID_2,
			},
		}

		mockService.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(expectedDetails, int64(2), nil).Times(1)

		// Act
		res, err := client.ListUnpublished(context.Background(), &trainingsessionpb.ListUnpublishedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetTrainingSessionsDetails()[0].GetTrainingSession().GetName(), expectedDetails[0].TrainingSession.Name)
		assert.Equal(t, res.GetTrainingSessionsDetails()[0].GetTrainingSession().GetId(), expectedDetails[0].TrainingSession.ID)
		assert.Equal(t, res.GetTrainingSessionsDetails()[0].GetPrice(), expectedDetails[0].Price)
		assert.Equal(t, res.GetTrainingSessionsDetails()[0].GetProductId(), expectedDetails[0].ProductID)
		assert.Equal(t, res.GetTrainingSessionsDetails()[1].GetTrainingSession().GetName(), expectedDetails[1].TrainingSession.Name)
		assert.Equal(t, res.GetTrainingSessionsDetails()[1].GetTrainingSession().GetId(), expectedDetails[1].TrainingSession.ID)
		assert.Equal(t, res.GetTrainingSessionsDetails()[1].GetPrice(), expectedDetails[1].Price)
		assert.Equal(t, res.GetTrainingSessionsDetails()[1].GetProductId(), expectedDetails[1].ProductID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(nil, int64(0), svcErr).Times(1)

		// Act
		res, err := client.ListUnpublished(context.Background(), &trainingsessionpb.ListUnpublishedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get training sessions")
	})
}

func TestServer_Create(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID := uuid.New().String()
	productID := uuid.New().String()
	createReq := trainingsessionmodel.CreateRequest{
		Name:             "training session name",
		ShortDescription: "training session short description",
		Price:            99.99,
		DurationMinutes:  30,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Create(gomock.Any(), &createReq).Return(&trainingsessionmodel.CreateResponse{ID: tsID, ProductID: productID}, nil)

		// Act
		res, err := client.Create(context.Background(), &trainingsessionpb.CreateRequest{
			Name:             createReq.Name,
			ShortDescription: createReq.ShortDescription,
			DurationMinutes:  int32(createReq.DurationMinutes),
			Price:            createReq.Price,
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), tsID)
		assert.Equal(t, res.GetProductId(), productID)
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().Create(gomock.Any(), &createReq).Return(nil, svcErr)

		// Act
		res, err := client.Create(context.Background(), &trainingsessionpb.CreateRequest{
			Name:             createReq.Name,
			ShortDescription: createReq.ShortDescription,
			DurationMinutes:  int32(createReq.DurationMinutes),
			Price:            createReq.Price,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to create training session")
	})
}

func TestServer_Publish(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Publish(gomock.Any(), tsID).Return(nil)

		// Act
		res, err := client.Publish(context.Background(), &trainingsessionpb.PublishRequest{Id: tsID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), tsID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Publish(gomock.Any(), invalidID).Return(trainingsessionservice.ErrInvalidArgument)

		// Act
		res, err := client.Publish(context.Background(), &trainingsessionpb.PublishRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid training session ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Publish(gomock.Any(), tsID).Return(trainingsessionservice.ErrNotFound)

		// Act
		res, err := client.Publish(context.Background(), &trainingsessionpb.PublishRequest{Id: tsID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Training session not found")
	})
}

func TestServer_Unpublish(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Unpublish(gomock.Any(), tsID).Return(nil)

		// Act
		res, err := client.Unpublish(context.Background(), &trainingsessionpb.UnpublishRequest{Id: tsID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), tsID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Unpublish(gomock.Any(), invalidID).Return(trainingsessionservice.ErrInvalidArgument)

		// Act
		res, err := client.Unpublish(context.Background(), &trainingsessionpb.UnpublishRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid training session ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Unpublish(gomock.Any(), tsID).Return(trainingsessionservice.ErrNotFound)

		// Act
		res, err := client.Unpublish(context.Background(), &trainingsessionpb.UnpublishRequest{Id: tsID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Training session not found")
	})
}

func TestServer_Update(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID := uuid.New().String()
	newName := "New training session name"
	newShortDescription := "New training session short description"
	newPrice := float32(99.99)

	t.Run("success", func(t *testing.T) {
		// Arrange
		updates := make(map[string]any)
		tsUpdates := map[string]any{"name": newName, "short_description": newShortDescription}
		productUpdates := map[string]any{"price": newPrice}
		updates["training_session"] = tsUpdates
		updates["product"] = productUpdates
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(updates, nil).Times(1)

		// Act
		res, err := client.Update(context.Background(), &trainingsessionpb.UpdateRequest{
			Id:               tsID,
			Name:             &newName,
			ShortDescription: &newShortDescription,
			Price:            &newPrice,
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.Updated.Paths[0], "updateresponse.name")
		assert.Equal(t, res.Updated.Paths[1], "updateresponse.short_description")
		assert.Equal(t, res.GetName(), tsUpdates["name"])
		assert.Equal(t, res.GetShortDescription(), tsUpdates["short_description"])
		assert.Equal(t, res.GetPrice(), productUpdates["price"])
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, trainingsessionservice.ErrInvalidArgument).Times(1)

		// Act
		res, err := client.Update(context.Background(), &trainingsessionpb.UpdateRequest{
			Id:               invalidID,
			Name:             &newName,
			ShortDescription: &newShortDescription,
			Price:            &newPrice,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid request payload")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, trainingsessionservice.ErrNotFound)

		// Act
		res, err := client.Update(context.Background(), &trainingsessionpb.UpdateRequest{
			Id:               tsID,
			Name:             &newName,
			ShortDescription: &newShortDescription,
			Price:            &newPrice,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Training session not found")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, svcErr)

		// Act
		res, err := client.Update(context.Background(), &trainingsessionpb.UpdateRequest{
			Id:               tsID,
			Name:             &newName,
			ShortDescription: &newShortDescription,
			Price:            &newPrice,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Database error")
	})
}

func TestServer_Delete(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Delete(gomock.Any(), tsID).Return(nil)

		// Act
		res, err := client.Delete(context.Background(), &trainingsessionpb.DeleteRequest{Id: tsID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), tsID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Delete(gomock.Any(), invalidID).Return(trainingsessionservice.ErrInvalidArgument)

		// Act
		res, err := client.Delete(context.Background(), &trainingsessionpb.DeleteRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid training session ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Delete(gomock.Any(), tsID).Return(trainingsessionservice.ErrNotFound)

		// Act
		res, err := client.Delete(context.Background(), &trainingsessionpb.DeleteRequest{Id: tsID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Training session not found")
	})
}

func TestServer_DeletePermanent(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().DeletePermanent(gomock.Any(), tsID).Return(nil)

		// Act
		res, err := client.DeletePermanent(context.Background(), &trainingsessionpb.DeletePermanentRequest{Id: tsID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), tsID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().DeletePermanent(gomock.Any(), invalidID).Return(trainingsessionservice.ErrInvalidArgument)

		// Act
		res, err := client.DeletePermanent(context.Background(), &trainingsessionpb.DeletePermanentRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid training session ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().DeletePermanent(gomock.Any(), tsID).Return(trainingsessionservice.ErrNotFound)

		// Act
		res, err := client.DeletePermanent(context.Background(), &trainingsessionpb.DeletePermanentRequest{Id: tsID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Training session not found")
	})
}

func TestServer_Restore(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	tsID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Restore(gomock.Any(), tsID).Return(nil)

		// Act
		res, err := client.Restore(context.Background(), &trainingsessionpb.RestoreRequest{Id: tsID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), tsID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Restore(gomock.Any(), invalidID).Return(trainingsessionservice.ErrInvalidArgument)

		// Act
		res, err := client.Restore(context.Background(), &trainingsessionpb.RestoreRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid training session ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Restore(gomock.Any(), tsID).Return(trainingsessionservice.ErrNotFound)

		// Act
		res, err := client.Restore(context.Background(), &trainingsessionpb.RestoreRequest{Id: tsID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Training session not found")
	})
}
