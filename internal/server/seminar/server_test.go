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

package seminar

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"github.com/google/uuid"
	seminarmodel "github.com/mikhail5545/product-service-go/internal/models/seminar"
	seminarservice "github.com/mikhail5545/product-service-go/internal/services/seminar"
	seminarmock "github.com/mikhail5545/product-service-go/internal/test/services/seminar_mock"
	seminarpb "github.com/mikhail5545/proto-go/proto/product_service/seminar/v0"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func setupTestServer(t *testing.T) (seminarpb.SeminarServiceClient, *seminarmock.MockService, func()) {
	t.Helper()

	// 1. Create mock controller and mock service
	ctrl := gomock.NewController(t)
	mockService := seminarmock.NewMockService(ctrl)

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

	client := seminarpb.NewSeminarServiceClient(conn)

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

	seminarID := uuid.New().String()
	rproductID := uuid.New().String()
	eproductID := uuid.New().String()
	lproductID := uuid.New().String()
	esproductID := uuid.New().String()
	lsproductID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &seminarmodel.SeminarDetails{
			Seminar: &seminarmodel.Seminar{
				ID:                      seminarID,
				Name:                    "Seminar name",
				ReservationProductID:    &rproductID,
				EarlyProductID:          &eproductID,
				LateProductID:           &lproductID,
				EarlySurchargeProductID: &esproductID,
				LateSurchargeProductID:  &lsproductID,
			},
			ReservationPrice:               11.11,
			EarlyPrice:                     22.22,
			LatePrice:                      33.33,
			EarlySurchargePrice:            44.44,
			LateSurchargePrice:             55.55,
			CurrentPrice:                   22.22,
			CurrentPriceProductID:          eproductID,
			CurrentSurchargePrice:          44.44,
			CurrentSurchargePriceProductID: esproductID,
		}

		mockService.EXPECT().Get(gomock.Any(), seminarID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.Get(context.Background(), &seminarpb.GetRequest{Id: seminarID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, seminarID, res.GetSeminarDetails().GetSeminar().GetId())
		assert.Equal(t, rproductID, res.GetSeminarDetails().GetSeminar().GetReservationProductId())
		assert.Equal(t, eproductID, res.GetSeminarDetails().GetSeminar().GetEarlyProductId())
		assert.Equal(t, lproductID, res.GetSeminarDetails().GetSeminar().GetLateProductId())
		assert.Equal(t, esproductID, res.GetSeminarDetails().GetSeminar().GetEarlySurchargeProductId())
		assert.Equal(t, lsproductID, res.GetSeminarDetails().GetSeminar().GetLateSurchargeProductId())
		assert.Equal(t, float32(22.22), res.GetSeminarDetails().GetEarlyPrice())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Get(gomock.Any(), seminarID).Return(nil, seminarservice.ErrNotFound)

		// Act
		res, err := client.Get(context.Background(), &seminarpb.GetRequest{Id: seminarID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Seminar not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Get(gomock.Any(), invalidID).Return(nil, seminarservice.ErrInvalidArgument)

		// Act
		res, err := client.Get(context.Background(), &seminarpb.GetRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid Seminar ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().Get(gomock.Any(), seminarID).Return(nil, svcErr)

		// Act
		res, err := client.Get(context.Background(), &seminarpb.GetRequest{Id: seminarID})

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

	seminarID := uuid.New().String()
	rproductID := uuid.New().String()
	eproductID := uuid.New().String()
	lproductID := uuid.New().String()
	esproductID := uuid.New().String()
	lsproductID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &seminarmodel.SeminarDetails{
			Seminar: &seminarmodel.Seminar{
				ID:                      seminarID,
				Name:                    "Seminar name",
				ReservationProductID:    &rproductID,
				EarlyProductID:          &eproductID,
				LateProductID:           &lproductID,
				EarlySurchargeProductID: &esproductID,
				LateSurchargeProductID:  &lsproductID,
			},
			ReservationPrice:               11.11,
			EarlyPrice:                     22.22,
			LatePrice:                      33.33,
			EarlySurchargePrice:            44.44,
			LateSurchargePrice:             55.55,
			CurrentPrice:                   22.22,
			CurrentPriceProductID:          eproductID,
			CurrentSurchargePrice:          44.44,
			CurrentSurchargePriceProductID: esproductID,
		}

		mockService.EXPECT().GetWithDeleted(gomock.Any(), seminarID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &seminarpb.GetWithDeletedRequest{Id: seminarID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, seminarID, res.GetSeminarDetails().GetSeminar().GetId())
		assert.Equal(t, rproductID, res.GetSeminarDetails().GetSeminar().GetReservationProductId())
		assert.Equal(t, eproductID, res.GetSeminarDetails().GetSeminar().GetEarlyProductId())
		assert.Equal(t, lproductID, res.GetSeminarDetails().GetSeminar().GetLateProductId())
		assert.Equal(t, esproductID, res.GetSeminarDetails().GetSeminar().GetEarlySurchargeProductId())
		assert.Equal(t, lsproductID, res.GetSeminarDetails().GetSeminar().GetLateSurchargeProductId())
		assert.Equal(t, float32(22.22), res.GetSeminarDetails().GetEarlyPrice())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().GetWithDeleted(gomock.Any(), seminarID).Return(nil, seminarservice.ErrNotFound)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &seminarpb.GetWithDeletedRequest{Id: seminarID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Seminar not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().GetWithDeleted(gomock.Any(), invalidID).Return(nil, seminarservice.ErrInvalidArgument)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &seminarpb.GetWithDeletedRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid Seminar ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().GetWithDeleted(gomock.Any(), seminarID).Return(nil, svcErr)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &seminarpb.GetWithDeletedRequest{Id: seminarID})

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

	seminarID := uuid.New().String()
	rproductID := uuid.New().String()
	eproductID := uuid.New().String()
	lproductID := uuid.New().String()
	esproductID := uuid.New().String()
	lsproductID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &seminarmodel.SeminarDetails{
			Seminar: &seminarmodel.Seminar{
				ID:                      seminarID,
				Name:                    "Seminar name",
				ReservationProductID:    &rproductID,
				EarlyProductID:          &eproductID,
				LateProductID:           &lproductID,
				EarlySurchargeProductID: &esproductID,
				LateSurchargeProductID:  &lsproductID,
			},
			ReservationPrice:               11.11,
			EarlyPrice:                     22.22,
			LatePrice:                      33.33,
			EarlySurchargePrice:            44.44,
			LateSurchargePrice:             55.55,
			CurrentPrice:                   22.22,
			CurrentPriceProductID:          eproductID,
			CurrentSurchargePrice:          44.44,
			CurrentSurchargePriceProductID: esproductID,
		}

		mockService.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &seminarpb.GetWithUnpublishedRequest{Id: seminarID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, seminarID, res.GetSeminarDetails().GetSeminar().GetId())
		assert.Equal(t, rproductID, res.GetSeminarDetails().GetSeminar().GetReservationProductId())
		assert.Equal(t, eproductID, res.GetSeminarDetails().GetSeminar().GetEarlyProductId())
		assert.Equal(t, lproductID, res.GetSeminarDetails().GetSeminar().GetLateProductId())
		assert.Equal(t, esproductID, res.GetSeminarDetails().GetSeminar().GetEarlySurchargeProductId())
		assert.Equal(t, lsproductID, res.GetSeminarDetails().GetSeminar().GetLateSurchargeProductId())
		assert.Equal(t, float32(22.22), res.GetSeminarDetails().GetEarlyPrice())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(nil, seminarservice.ErrNotFound)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &seminarpb.GetWithUnpublishedRequest{Id: seminarID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Seminar not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), invalidID).Return(nil, seminarservice.ErrInvalidArgument)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &seminarpb.GetWithUnpublishedRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid Seminar ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), seminarID).Return(nil, svcErr)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &seminarpb.GetWithUnpublishedRequest{Id: seminarID})

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

	seminarID_1 := uuid.New().String()
	rproductID_1 := uuid.New().String()
	eproductID_1 := uuid.New().String()
	lproductID_1 := uuid.New().String()
	esproductID_1 := uuid.New().String()
	lsproductID_1 := uuid.New().String()
	seminarID_2 := uuid.New().String()
	rproductID_2 := uuid.New().String()
	eproductID_2 := uuid.New().String()
	lproductID_2 := uuid.New().String()
	esproductID_2 := uuid.New().String()
	lsproductID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := []seminarmodel.SeminarDetails{
			{
				Seminar: &seminarmodel.Seminar{
					ID:                      seminarID_1,
					Name:                    "Seminar 1 name",
					ReservationProductID:    &rproductID_1,
					EarlyProductID:          &eproductID_1,
					LateProductID:           &lproductID_1,
					EarlySurchargeProductID: &esproductID_1,
					LateSurchargeProductID:  &lsproductID_1,
				},
				CurrentPrice:          99.99,
				CurrentPriceProductID: eproductID_1,
			},
			{
				Seminar: &seminarmodel.Seminar{
					ID:                      seminarID_2,
					Name:                    "Seminar 2 name",
					ReservationProductID:    &rproductID_2,
					EarlyProductID:          &eproductID_2,
					LateProductID:           &lproductID_2,
					EarlySurchargeProductID: &esproductID_2,
					LateSurchargeProductID:  &lsproductID_2,
				},
				CurrentPrice:          199.99,
				CurrentPriceProductID: eproductID_2,
			},
		}

		mockService.EXPECT().List(gomock.Any(), limit, offset).Return(expectedDetails, int64(2), nil).Times(1)

		// Act
		res, err := client.List(context.Background(), &seminarpb.ListRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetSeminarDetails()[0].GetSeminar().GetName(), expectedDetails[0].Seminar.Name)
		assert.Equal(t, res.GetSeminarDetails()[0].GetSeminar().GetId(), expectedDetails[0].Seminar.ID)
		assert.Equal(t, res.GetSeminarDetails()[0].GetCurrentPrice(), expectedDetails[0].CurrentPrice)
		assert.Equal(t, res.GetSeminarDetails()[0].GetCurrentPriceProductId(), expectedDetails[0].CurrentPriceProductID)
		assert.Equal(t, res.GetSeminarDetails()[1].GetSeminar().GetName(), expectedDetails[1].Seminar.Name)
		assert.Equal(t, res.GetSeminarDetails()[1].GetSeminar().GetId(), expectedDetails[1].Seminar.ID)
		assert.Equal(t, res.GetSeminarDetails()[1].GetCurrentPrice(), expectedDetails[1].CurrentPrice)
		assert.Equal(t, res.GetSeminarDetails()[1].GetCurrentPriceProductId(), expectedDetails[1].CurrentPriceProductID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().List(gomock.Any(), limit, offset).Return(nil, int64(0), svcErr).Times(1)

		// Act
		res, err := client.List(context.Background(), &seminarpb.ListRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get Seminars")
	})
}

func TestServer_ListDeleted(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	seminarID_1 := uuid.New().String()
	rproductID_1 := uuid.New().String()
	eproductID_1 := uuid.New().String()
	lproductID_1 := uuid.New().String()
	esproductID_1 := uuid.New().String()
	lsproductID_1 := uuid.New().String()
	seminarID_2 := uuid.New().String()
	rproductID_2 := uuid.New().String()
	eproductID_2 := uuid.New().String()
	lproductID_2 := uuid.New().String()
	esproductID_2 := uuid.New().String()
	lsproductID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := []seminarmodel.SeminarDetails{
			{
				Seminar: &seminarmodel.Seminar{
					ID:                      seminarID_1,
					Name:                    "Seminar 1 name",
					ReservationProductID:    &rproductID_1,
					EarlyProductID:          &eproductID_1,
					LateProductID:           &lproductID_1,
					EarlySurchargeProductID: &esproductID_1,
					LateSurchargeProductID:  &lsproductID_1,
				},
				CurrentPrice:          99.99,
				CurrentPriceProductID: eproductID_1,
			},
			{
				Seminar: &seminarmodel.Seminar{
					ID:                      seminarID_2,
					Name:                    "Seminar 2 name",
					ReservationProductID:    &rproductID_2,
					EarlyProductID:          &eproductID_2,
					LateProductID:           &lproductID_2,
					EarlySurchargeProductID: &esproductID_2,
					LateSurchargeProductID:  &lsproductID_2,
				},
				CurrentPrice:          199.99,
				CurrentPriceProductID: eproductID_2,
			},
		}

		mockService.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(expectedDetails, int64(2), nil).Times(1)

		// Act
		res, err := client.ListDeleted(context.Background(), &seminarpb.ListDeletedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetSeminarDetails()[0].GetSeminar().GetName(), expectedDetails[0].Seminar.Name)
		assert.Equal(t, res.GetSeminarDetails()[0].GetSeminar().GetId(), expectedDetails[0].Seminar.ID)
		assert.Equal(t, res.GetSeminarDetails()[0].GetCurrentPrice(), expectedDetails[0].CurrentPrice)
		assert.Equal(t, res.GetSeminarDetails()[0].GetCurrentPriceProductId(), expectedDetails[0].CurrentPriceProductID)
		assert.Equal(t, res.GetSeminarDetails()[1].GetSeminar().GetName(), expectedDetails[1].Seminar.Name)
		assert.Equal(t, res.GetSeminarDetails()[1].GetSeminar().GetId(), expectedDetails[1].Seminar.ID)
		assert.Equal(t, res.GetSeminarDetails()[1].GetCurrentPrice(), expectedDetails[1].CurrentPrice)
		assert.Equal(t, res.GetSeminarDetails()[1].GetCurrentPriceProductId(), expectedDetails[1].CurrentPriceProductID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(nil, int64(0), svcErr).Times(1)

		// Act
		res, err := client.ListDeleted(context.Background(), &seminarpb.ListDeletedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get Seminars")
	})
}

func TestServer_ListUnpublished(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	seminarID_1 := uuid.New().String()
	rproductID_1 := uuid.New().String()
	eproductID_1 := uuid.New().String()
	lproductID_1 := uuid.New().String()
	esproductID_1 := uuid.New().String()
	lsproductID_1 := uuid.New().String()
	seminarID_2 := uuid.New().String()
	rproductID_2 := uuid.New().String()
	eproductID_2 := uuid.New().String()
	lproductID_2 := uuid.New().String()
	esproductID_2 := uuid.New().String()
	lsproductID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := []seminarmodel.SeminarDetails{
			{
				Seminar: &seminarmodel.Seminar{
					ID:                      seminarID_1,
					Name:                    "Seminar 1 name",
					ReservationProductID:    &rproductID_1,
					EarlyProductID:          &eproductID_1,
					LateProductID:           &lproductID_1,
					EarlySurchargeProductID: &esproductID_1,
					LateSurchargeProductID:  &lsproductID_1,
				},
				CurrentPrice:          99.99,
				CurrentPriceProductID: eproductID_1,
			},
			{
				Seminar: &seminarmodel.Seminar{
					ID:                      seminarID_2,
					Name:                    "Seminar 2 name",
					ReservationProductID:    &rproductID_2,
					EarlyProductID:          &eproductID_2,
					LateProductID:           &lproductID_2,
					EarlySurchargeProductID: &esproductID_2,
					LateSurchargeProductID:  &lsproductID_2,
				},
				CurrentPrice:          199.99,
				CurrentPriceProductID: eproductID_2,
			},
		}

		mockService.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(expectedDetails, int64(2), nil).Times(1)

		// Act
		res, err := client.ListUnpublished(context.Background(), &seminarpb.ListUnpublishedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetSeminarDetails()[0].GetSeminar().GetName(), expectedDetails[0].Seminar.Name)
		assert.Equal(t, res.GetSeminarDetails()[0].GetSeminar().GetId(), expectedDetails[0].Seminar.ID)
		assert.Equal(t, res.GetSeminarDetails()[0].GetCurrentPrice(), expectedDetails[0].CurrentPrice)
		assert.Equal(t, res.GetSeminarDetails()[0].GetCurrentPriceProductId(), expectedDetails[0].CurrentPriceProductID)
		assert.Equal(t, res.GetSeminarDetails()[1].GetSeminar().GetName(), expectedDetails[1].Seminar.Name)
		assert.Equal(t, res.GetSeminarDetails()[1].GetSeminar().GetId(), expectedDetails[1].Seminar.ID)
		assert.Equal(t, res.GetSeminarDetails()[1].GetCurrentPrice(), expectedDetails[1].CurrentPrice)
		assert.Equal(t, res.GetSeminarDetails()[1].GetCurrentPriceProductId(), expectedDetails[1].CurrentPriceProductID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(nil, int64(0), svcErr).Times(1)

		// Act
		res, err := client.ListUnpublished(context.Background(), &seminarpb.ListUnpublishedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get Seminars")
	})
}

func TestServer_Create(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	seminarID := uuid.New().String()
	rproductID := uuid.New().String()
	eproductID := uuid.New().String()
	lproductID := uuid.New().String()
	esproductID := uuid.New().String()
	lsproductID := uuid.New().String()

	createReq := seminarmodel.CreateRequest{
		Name:             "seminar name",
		ShortDescription: "seminar short description",
		ReservationPrice: 99.99,
		EarlyPrice:       22.22,
		LatePrice:        33.33,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&seminarmodel.CreateResponse{
			ID:                      seminarID,
			ReservationProductID:    rproductID,
			EarlyProductID:          eproductID,
			LateProductID:           lproductID,
			EarlySurchargeProductID: esproductID,
			LateSurchargeProductID:  lsproductID,
		},
			nil,
		)

		// Act
		res, err := client.Create(context.Background(), &seminarpb.CreateRequest{
			Name:             createReq.Name,
			ShortDescription: createReq.ShortDescription,
			ReservationPrice: createReq.ReservationPrice,
			LatePrice:        createReq.LatePrice,
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, seminarID, res.GetId())
		assert.Equal(t, rproductID, res.GetReservationProductId())
		assert.Equal(t, eproductID, res.GetEarlyProductId())
		assert.Equal(t, lproductID, res.GetLateProductId())
		assert.Equal(t, esproductID, res.GetEarlySurchargeProductId())
		assert.Equal(t, lsproductID, res.GetLateSurchargeProductId())
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, svcErr)

		// Act
		res, err := client.Create(context.Background(), &seminarpb.CreateRequest{
			Name:             createReq.Name,
			ShortDescription: createReq.ShortDescription,
			ReservationPrice: createReq.ReservationPrice,
			LatePrice:        createReq.LatePrice,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to create seminar")
	})
}

func TestServer_Publish(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Publish(gomock.Any(), seminarID).Return(nil)

		// Act
		res, err := client.Publish(context.Background(), &seminarpb.PublishRequest{Id: seminarID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), seminarID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Publish(gomock.Any(), invalidID).Return(seminarservice.ErrInvalidArgument)

		// Act
		res, err := client.Publish(context.Background(), &seminarpb.PublishRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid seminar ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Publish(gomock.Any(), seminarID).Return(seminarservice.ErrNotFound)

		// Act
		res, err := client.Publish(context.Background(), &seminarpb.PublishRequest{Id: seminarID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "seminar not found")
	})
}

func TestServer_Unpublish(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Unpublish(gomock.Any(), seminarID).Return(nil)

		// Act
		res, err := client.Unpublish(context.Background(), &seminarpb.UnpublishRequest{Id: seminarID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), seminarID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Unpublish(gomock.Any(), invalidID).Return(seminarservice.ErrInvalidArgument)

		// Act
		res, err := client.Unpublish(context.Background(), &seminarpb.UnpublishRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid seminar ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Unpublish(gomock.Any(), seminarID).Return(seminarservice.ErrNotFound)

		// Act
		res, err := client.Unpublish(context.Background(), &seminarpb.UnpublishRequest{Id: seminarID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "seminar not found")
	})
}

func TestServer_Update(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	seminarID := uuid.New().String()
	newName := "New seminar name"
	newShortDescription := "New seminar short description"
	newEarlyPrice := float32(99.99)
	newEarlySurchargePrice := float32(199.99)

	t.Run("success", func(t *testing.T) {
		// Arrange
		updates := make(map[string]any)
		seminarUpdates := map[string]any{"name": newName, "short_description": newShortDescription}
		eproductUpdates := map[string]any{"price": newEarlyPrice}
		esproductUpdates := map[string]any{"price": newEarlySurchargePrice}
		updates["seminar"] = seminarUpdates
		updates["early_product"] = eproductUpdates
		updates["early_surcharge_product"] = esproductUpdates
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(updates, nil).Times(1)

		// Act
		res, err := client.Update(context.Background(), &seminarpb.UpdateRequest{
			Id:                  seminarID,
			Name:                &newName,
			ShortDescription:    &newShortDescription,
			EarlyPrice:          &newEarlyPrice,
			EarlySurchargePrice: &newEarlySurchargePrice,
		})

		// Assert
		assert.NoError(t, err)
		assert.Contains(t, res.Updated.Paths, "updateresponse.name")
		assert.Contains(t, res.Updated.Paths, "updateresponse.short_description")
		assert.Contains(t, res.Updated.Paths, "updateresponse.early_price")
		assert.Contains(t, res.Updated.Paths, "updateresponse.early_surcharge_price")
		assert.Equal(t, seminarUpdates["name"], res.GetName())
		assert.Equal(t, seminarUpdates["short_description"], res.GetShortDescription())
		assert.Equal(t, eproductUpdates["price"], res.GetEarlyPrice())
		assert.Equal(t, esproductUpdates["price"], res.GetEarlySurchargePrice())
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, seminarservice.ErrInvalidArgument).Times(1)

		// Act
		res, err := client.Update(context.Background(), &seminarpb.UpdateRequest{
			Id:                  invalidID,
			Name:                &newName,
			ShortDescription:    &newShortDescription,
			EarlyPrice:          &newEarlyPrice,
			EarlySurchargePrice: &newEarlySurchargePrice,
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
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, seminarservice.ErrNotFound)

		// Act
		res, err := client.Update(context.Background(), &seminarpb.UpdateRequest{
			Id:                  seminarID,
			Name:                &newName,
			ShortDescription:    &newShortDescription,
			EarlyPrice:          &newEarlyPrice,
			EarlySurchargePrice: &newEarlySurchargePrice,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "seminar not found")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, svcErr)

		// Act
		res, err := client.Update(context.Background(), &seminarpb.UpdateRequest{
			Id:                  seminarID,
			Name:                &newName,
			ShortDescription:    &newShortDescription,
			EarlyPrice:          &newEarlyPrice,
			EarlySurchargePrice: &newEarlySurchargePrice,
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

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Delete(gomock.Any(), seminarID).Return(nil)

		// Act
		res, err := client.Delete(context.Background(), &seminarpb.DeleteRequest{Id: seminarID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), seminarID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Delete(gomock.Any(), invalidID).Return(seminarservice.ErrInvalidArgument)

		// Act
		res, err := client.Delete(context.Background(), &seminarpb.DeleteRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid seminar ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Delete(gomock.Any(), seminarID).Return(seminarservice.ErrNotFound)

		// Act
		res, err := client.Delete(context.Background(), &seminarpb.DeleteRequest{Id: seminarID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "seminar not found")
	})
}

func TestServer_DeletePermanent(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().DeletePermanent(gomock.Any(), seminarID).Return(nil)

		// Act
		res, err := client.DeletePermanent(context.Background(), &seminarpb.DeletePermanentRequest{Id: seminarID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), seminarID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().DeletePermanent(gomock.Any(), invalidID).Return(seminarservice.ErrInvalidArgument)

		// Act
		res, err := client.DeletePermanent(context.Background(), &seminarpb.DeletePermanentRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid seminar ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().DeletePermanent(gomock.Any(), seminarID).Return(seminarservice.ErrNotFound)

		// Act
		res, err := client.DeletePermanent(context.Background(), &seminarpb.DeletePermanentRequest{Id: seminarID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "seminar not found")
	})
}

func TestServer_Restore(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	seminarID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Restore(gomock.Any(), seminarID).Return(nil)

		// Act
		res, err := client.Restore(context.Background(), &seminarpb.RestoreRequest{Id: seminarID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), seminarID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Restore(gomock.Any(), invalidID).Return(seminarservice.ErrInvalidArgument)

		// Act
		res, err := client.Restore(context.Background(), &seminarpb.RestoreRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid seminar ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Restore(gomock.Any(), seminarID).Return(seminarservice.ErrNotFound)

		// Act
		res, err := client.Restore(context.Background(), &seminarpb.RestoreRequest{Id: seminarID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "seminar not found")
	})
}
