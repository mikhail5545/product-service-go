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

package physicalgood

import (
	"context"
	"log"
	"net"
	"net/http"
	"testing"

	"github.com/google/uuid"
	physicalgoodmodel "github.com/mikhail5545/product-service-go/internal/models/physical_good"
	physicalgoodservice "github.com/mikhail5545/product-service-go/internal/services/physical_good"
	physicalgoodmock "github.com/mikhail5545/product-service-go/internal/test/services/physical_good_mock"
	physicalgoodpb "github.com/mikhail5545/proto-go/proto/physical_good/v0"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func setupTestServer(t *testing.T) (physicalgoodpb.PhysicalGoodServiceClient, *physicalgoodmock.MockService, func()) {
	t.Helper()

	// 1. Create mock controller and mock service
	ctrl := gomock.NewController(t)
	mockService := physicalgoodmock.NewMockService(ctrl)

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

	client := physicalgoodpb.NewPhysicalGoodServiceClient(conn)

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

	goodID := uuid.New().String()
	productID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &physicalgoodmodel.PhysicalGoodDetails{
			PhysicalGood: physicalgoodmodel.PhysicalGood{
				ID:   goodID,
				Name: "Physical good name",
			},
			Price:     99.99,
			ProductID: productID,
		}

		mockService.EXPECT().Get(gomock.Any(), goodID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.Get(context.Background(), &physicalgoodpb.GetRequest{Id: goodID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, goodID, res.GetPhysicalGoodDetails().GetPhysicalGood().GetId())
		assert.Equal(t, productID, res.GetPhysicalGoodDetails().GetProductId())
		assert.Equal(t, float32(99.99), res.GetPhysicalGoodDetails().GetPrice())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Get(gomock.Any(), goodID).Return(nil, serviceErr)

		// Act
		res, err := client.Get(context.Background(), &physicalgoodpb.GetRequest{Id: goodID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Physical good not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Invalid physical good ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().Get(gomock.Any(), invalidID).Return(nil, serviceErr)

		// Act
		res, err := client.Get(context.Background(), &physicalgoodpb.GetRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid physical good ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Database error",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().Get(gomock.Any(), goodID).Return(nil, serviceErr)

		// Act
		res, err := client.Get(context.Background(), &physicalgoodpb.GetRequest{Id: goodID})

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

	goodID := uuid.New().String()
	productID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &physicalgoodmodel.PhysicalGoodDetails{
			PhysicalGood: physicalgoodmodel.PhysicalGood{
				ID:   goodID,
				Name: "Physical good name",
			},
			Price:     99.99,
			ProductID: productID,
		}

		mockService.EXPECT().GetWithDeleted(gomock.Any(), goodID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &physicalgoodpb.GetWithDeletedRequest{Id: goodID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, goodID, res.GetPhysicalGoodDetails().GetPhysicalGood().GetId())
		assert.Equal(t, productID, res.GetPhysicalGoodDetails().GetProductId())
		assert.Equal(t, float32(99.99), res.GetPhysicalGoodDetails().GetPrice())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().GetWithDeleted(gomock.Any(), goodID).Return(nil, serviceErr)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &physicalgoodpb.GetWithDeletedRequest{Id: goodID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Physical good not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Invalid physical good ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().GetWithDeleted(gomock.Any(), invalidID).Return(nil, serviceErr)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &physicalgoodpb.GetWithDeletedRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid physical good ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Database error",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().GetWithDeleted(gomock.Any(), goodID).Return(nil, serviceErr)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &physicalgoodpb.GetWithDeletedRequest{Id: goodID})

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

	goodID := uuid.New().String()
	productID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &physicalgoodmodel.PhysicalGoodDetails{
			PhysicalGood: physicalgoodmodel.PhysicalGood{
				ID:   goodID,
				Name: "Physical good name",
			},
			Price:     99.99,
			ProductID: productID,
		}

		mockService.EXPECT().GetWithUnpublished(gomock.Any(), goodID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &physicalgoodpb.GetWithUnpublishedRequest{Id: goodID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, goodID, res.GetPhysicalGoodDetails().GetPhysicalGood().GetId())
		assert.Equal(t, productID, res.GetPhysicalGoodDetails().GetProductId())
		assert.Equal(t, float32(99.99), res.GetPhysicalGoodDetails().GetPrice())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), goodID).Return(nil, serviceErr)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &physicalgoodpb.GetWithUnpublishedRequest{Id: goodID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Physical good not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Invalid physical good ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), invalidID).Return(nil, serviceErr)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &physicalgoodpb.GetWithUnpublishedRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid physical good ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Database error",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), goodID).Return(nil, serviceErr)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &physicalgoodpb.GetWithUnpublishedRequest{Id: goodID})

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

	goodID_1 := uuid.New().String()
	productID_1 := uuid.New().String()
	goodID_2 := uuid.New().String()
	productID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := []physicalgoodmodel.PhysicalGoodDetails{
			{
				PhysicalGood: physicalgoodmodel.PhysicalGood{
					ID:   goodID_1,
					Name: "Physical good 1 name",
				},
				Price:     99.99,
				ProductID: productID_1,
			},
			{
				PhysicalGood: physicalgoodmodel.PhysicalGood{
					ID:   goodID_2,
					Name: "Physical good 2 name",
				},
				Price:     199.99,
				ProductID: productID_2,
			},
		}

		mockService.EXPECT().List(gomock.Any(), limit, offset).Return(expectedDetails, int64(2), nil).Times(1)

		// Act
		res, err := client.List(context.Background(), &physicalgoodpb.ListRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetPhysicalGoodDetails()[0].GetPhysicalGood().GetName(), expectedDetails[0].PhysicalGood.Name)
		assert.Equal(t, res.GetPhysicalGoodDetails()[0].GetPhysicalGood().GetId(), expectedDetails[0].PhysicalGood.ID)
		assert.Equal(t, res.GetPhysicalGoodDetails()[0].GetPrice(), expectedDetails[0].Price)
		assert.Equal(t, res.GetPhysicalGoodDetails()[0].GetProductId(), expectedDetails[0].ProductID)
		assert.Equal(t, res.GetPhysicalGoodDetails()[1].GetPhysicalGood().GetName(), expectedDetails[1].PhysicalGood.Name)
		assert.Equal(t, res.GetPhysicalGoodDetails()[1].GetPhysicalGood().GetId(), expectedDetails[1].PhysicalGood.ID)
		assert.Equal(t, res.GetPhysicalGoodDetails()[1].GetPrice(), expectedDetails[1].Price)
		assert.Equal(t, res.GetPhysicalGoodDetails()[1].GetProductId(), expectedDetails[1].ProductID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Failed to get physical goods",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().List(gomock.Any(), limit, offset).Return(nil, int64(0), serviceErr).Times(1)

		// Act
		res, err := client.List(context.Background(), &physicalgoodpb.ListRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get physical goods")
	})
}

func TestServer_ListDeleted(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	goodID_1 := uuid.New().String()
	productID_1 := uuid.New().String()
	goodID_2 := uuid.New().String()
	productID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := []physicalgoodmodel.PhysicalGoodDetails{
			{
				PhysicalGood: physicalgoodmodel.PhysicalGood{
					ID:   goodID_1,
					Name: "Physical good 1 name",
				},
				Price:     99.99,
				ProductID: productID_1,
			},
			{
				PhysicalGood: physicalgoodmodel.PhysicalGood{
					ID:   goodID_2,
					Name: "Physical good 2 name",
				},
				Price:     199.99,
				ProductID: productID_2,
			},
		}

		mockService.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(expectedDetails, int64(2), nil).Times(1)

		// Act
		res, err := client.ListDeleted(context.Background(), &physicalgoodpb.ListDeletedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetPhysicalGoodDetails()[0].GetPhysicalGood().GetName(), expectedDetails[0].PhysicalGood.Name)
		assert.Equal(t, res.GetPhysicalGoodDetails()[0].GetPhysicalGood().GetId(), expectedDetails[0].PhysicalGood.ID)
		assert.Equal(t, res.GetPhysicalGoodDetails()[0].GetPrice(), expectedDetails[0].Price)
		assert.Equal(t, res.GetPhysicalGoodDetails()[0].GetProductId(), expectedDetails[0].ProductID)
		assert.Equal(t, res.GetPhysicalGoodDetails()[1].GetPhysicalGood().GetName(), expectedDetails[1].PhysicalGood.Name)
		assert.Equal(t, res.GetPhysicalGoodDetails()[1].GetPhysicalGood().GetId(), expectedDetails[1].PhysicalGood.ID)
		assert.Equal(t, res.GetPhysicalGoodDetails()[1].GetPrice(), expectedDetails[1].Price)
		assert.Equal(t, res.GetPhysicalGoodDetails()[1].GetProductId(), expectedDetails[1].ProductID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Failed to get physical goods",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(nil, int64(0), serviceErr).Times(1)

		// Act
		res, err := client.ListDeleted(context.Background(), &physicalgoodpb.ListDeletedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get physical goods")
	})
}

func TestServer_ListUnpublished(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	goodID_1 := uuid.New().String()
	productID_1 := uuid.New().String()
	goodID_2 := uuid.New().String()
	productID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := []physicalgoodmodel.PhysicalGoodDetails{
			{
				PhysicalGood: physicalgoodmodel.PhysicalGood{
					ID:   goodID_1,
					Name: "Physical good 1 name",
				},
				Price:     99.99,
				ProductID: productID_1,
			},
			{
				PhysicalGood: physicalgoodmodel.PhysicalGood{
					ID:   goodID_2,
					Name: "Physical good 2 name",
				},
				Price:     199.99,
				ProductID: productID_2,
			},
		}

		mockService.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(expectedDetails, int64(2), nil).Times(1)

		// Act
		res, err := client.ListUnpublished(context.Background(), &physicalgoodpb.ListUnpublishedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetPhysicalGoodDetails()[0].GetPhysicalGood().GetName(), expectedDetails[0].PhysicalGood.Name)
		assert.Equal(t, res.GetPhysicalGoodDetails()[0].GetPhysicalGood().GetId(), expectedDetails[0].PhysicalGood.ID)
		assert.Equal(t, res.GetPhysicalGoodDetails()[0].GetPrice(), expectedDetails[0].Price)
		assert.Equal(t, res.GetPhysicalGoodDetails()[0].GetProductId(), expectedDetails[0].ProductID)
		assert.Equal(t, res.GetPhysicalGoodDetails()[1].GetPhysicalGood().GetName(), expectedDetails[1].PhysicalGood.Name)
		assert.Equal(t, res.GetPhysicalGoodDetails()[1].GetPhysicalGood().GetId(), expectedDetails[1].PhysicalGood.ID)
		assert.Equal(t, res.GetPhysicalGoodDetails()[1].GetPrice(), expectedDetails[1].Price)
		assert.Equal(t, res.GetPhysicalGoodDetails()[1].GetProductId(), expectedDetails[1].ProductID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Failed to get physical goods",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(nil, int64(0), serviceErr).Times(1)

		// Act
		res, err := client.ListUnpublished(context.Background(), &physicalgoodpb.ListUnpublishedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get physical goods")
	})
}

func TestServer_Create(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	goodID := uuid.New().String()
	productID := uuid.New().String()
	createReq := physicalgoodmodel.CreateRequest{
		Name:             "Physical good name",
		ShortDescription: "Physical good short description",
		Price:            99.99,
		Amount:           33,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Create(gomock.Any(), &createReq).Return(&physicalgoodmodel.CreateResponse{ID: goodID, ProductID: productID}, nil)

		// Act
		res, err := client.Create(context.Background(), &physicalgoodpb.CreateRequest{
			Name:             createReq.Name,
			ShortDescription: createReq.ShortDescription,
			Amount:           int32(createReq.Amount),
			Price:            createReq.Price,
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), goodID)
		assert.Equal(t, res.GetProductId(), productID)
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Failed to create physical good",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().Create(gomock.Any(), &createReq).Return(nil, serviceErr)

		// Act
		res, err := client.Create(context.Background(), &physicalgoodpb.CreateRequest{
			Name:             createReq.Name,
			ShortDescription: createReq.ShortDescription,
			Amount:           int32(createReq.Amount),
			Price:            createReq.Price,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to create physical good")
	})
}

func TestServer_Publish(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	goodID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Publish(gomock.Any(), goodID).Return(nil)

		// Act
		res, err := client.Publish(context.Background(), &physicalgoodpb.PublishRequest{Id: goodID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), goodID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Invalid physical good ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().Publish(gomock.Any(), invalidID).Return(serviceErr)

		// Act
		res, err := client.Publish(context.Background(), &physicalgoodpb.PublishRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid physical good ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Publish(gomock.Any(), goodID).Return(serviceErr)

		// Act
		res, err := client.Publish(context.Background(), &physicalgoodpb.PublishRequest{Id: goodID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Physical good not found")
	})
}

func TestServer_Unpublish(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	goodID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Unpublish(gomock.Any(), goodID).Return(nil)

		// Act
		res, err := client.Unpublish(context.Background(), &physicalgoodpb.UnpublishRequest{Id: goodID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), goodID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Invalid physical good ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().Unpublish(gomock.Any(), invalidID).Return(serviceErr)

		// Act
		res, err := client.Unpublish(context.Background(), &physicalgoodpb.UnpublishRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid physical good ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Unpublish(gomock.Any(), goodID).Return(serviceErr)

		// Act
		res, err := client.Unpublish(context.Background(), &physicalgoodpb.UnpublishRequest{Id: goodID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Physical good not found")
	})
}

func TestServer_Update(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	goodID := uuid.New().String()
	newName := "New course name"
	newShortDescription := "New course short description"
	newPrice := float32(99.99)

	t.Run("success", func(t *testing.T) {
		// Arrange
		updates := make(map[string]any)
		goodUpdates := map[string]any{"name": newName, "short_description": newShortDescription}
		productUpdates := map[string]any{"price": newPrice}
		updates["physical_good"] = goodUpdates
		updates["product"] = productUpdates
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(updates, nil).Times(1)

		// Act
		res, err := client.Update(context.Background(), &physicalgoodpb.UpdateRequest{
			Id:               goodID,
			Name:             &newName,
			ShortDescription: &newShortDescription,
			Price:            &newPrice,
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.Updated.Paths[0], "updateresponse.name")
		assert.Equal(t, res.Updated.Paths[1], "updateresponse.short_description")
		assert.Equal(t, res.GetName(), goodUpdates["name"])
		assert.Equal(t, res.GetShortDescription(), goodUpdates["short_description"])
		assert.Equal(t, res.GetPrice(), productUpdates["price"])
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Invalid request payload",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, serviceErr).Times(1)

		// Act
		res, err := client.Update(context.Background(), &physicalgoodpb.UpdateRequest{
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
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, serviceErr)

		// Act
		res, err := client.Update(context.Background(), &physicalgoodpb.UpdateRequest{
			Id:               goodID,
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
		assert.Contains(t, st.Message(), "Physical good not found")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Database error",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, serviceErr)

		// Act
		res, err := client.Update(context.Background(), &physicalgoodpb.UpdateRequest{
			Id:               goodID,
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

	goodID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Delete(gomock.Any(), goodID).Return(nil)

		// Act
		res, err := client.Delete(context.Background(), &physicalgoodpb.DeleteRequest{Id: goodID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), goodID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Invalid physical good ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().Delete(gomock.Any(), invalidID).Return(serviceErr)

		// Act
		res, err := client.Delete(context.Background(), &physicalgoodpb.DeleteRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid physical good ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Delete(gomock.Any(), goodID).Return(serviceErr)

		// Act
		res, err := client.Delete(context.Background(), &physicalgoodpb.DeleteRequest{Id: goodID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Physical good not found")
	})
}

func TestServer_DeletePermanent(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	goodID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().DeletePermanent(gomock.Any(), goodID).Return(nil)

		// Act
		res, err := client.DeletePermanent(context.Background(), &physicalgoodpb.DeletePermanentRequest{Id: goodID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), goodID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Invalid physical good ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().DeletePermanent(gomock.Any(), invalidID).Return(serviceErr)

		// Act
		res, err := client.DeletePermanent(context.Background(), &physicalgoodpb.DeletePermanentRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid physical good ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().DeletePermanent(gomock.Any(), goodID).Return(serviceErr)

		// Act
		res, err := client.DeletePermanent(context.Background(), &physicalgoodpb.DeletePermanentRequest{Id: goodID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Physical good not found")
	})
}

func TestServer_Restore(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	goodID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Restore(gomock.Any(), goodID).Return(nil)

		// Act
		res, err := client.Restore(context.Background(), &physicalgoodpb.RestoreRequest{Id: goodID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), goodID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Invalid physical good ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().Restore(gomock.Any(), invalidID).Return(serviceErr)

		// Act
		res, err := client.Restore(context.Background(), &physicalgoodpb.RestoreRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid physical good ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &physicalgoodservice.Error{
			Msg:  "Physical good not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Restore(gomock.Any(), goodID).Return(serviceErr)

		// Act
		res, err := client.Restore(context.Background(), &physicalgoodpb.RestoreRequest{Id: goodID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Physical good not found")
	})
}
