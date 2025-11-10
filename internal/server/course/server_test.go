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

package course

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"github.com/google/uuid"
	coursemodel "github.com/mikhail5545/product-service-go/internal/models/course"
	courseservice "github.com/mikhail5545/product-service-go/internal/services/course"
	coursemock "github.com/mikhail5545/product-service-go/internal/test/services/course_mock"
	coursepb "github.com/mikhail5545/proto-go/proto/product_service/course/v0"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func setupTestServer(t *testing.T) (coursepb.CourseServiceClient, *coursemock.MockService, func()) {
	t.Helper()

	// 1. Create mock controller and mock service
	ctrl := gomock.NewController(t)
	mockService := coursemock.NewMockService(ctrl)

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

	client := coursepb.NewCourseServiceClient(conn)

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

	courseID := uuid.New().String()
	productID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &coursemodel.CourseDetails{
			Course: &coursemodel.Course{
				ID:   courseID,
				Name: "Course name",
			},
			Price:     99.99,
			ProductID: productID,
		}

		mockService.EXPECT().Get(gomock.Any(), courseID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.Get(context.Background(), &coursepb.GetRequest{Id: courseID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, courseID, res.GetCourseDetails().GetCourse().GetId())
		assert.Equal(t, productID, res.GetCourseDetails().GetProductId())
		assert.Equal(t, float32(99.99), res.GetCourseDetails().GetPrice())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Get(gomock.Any(), courseID).Return(nil, courseservice.ErrNotFound)

		// Act
		res, err := client.Get(context.Background(), &coursepb.GetRequest{Id: courseID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Get(gomock.Any(), invalidID).Return(nil, courseservice.ErrInvalidArgument)

		// Act
		res, err := client.Get(context.Background(), &coursepb.GetRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().Get(gomock.Any(), courseID).Return(nil, svcErr)

		// Act
		res, err := client.Get(context.Background(), &coursepb.GetRequest{Id: courseID})

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

	courseID := uuid.New().String()
	productID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &coursemodel.CourseDetails{
			Course: &coursemodel.Course{
				ID:   courseID,
				Name: "Course name",
			},
			Price:     99.99,
			ProductID: productID,
		}

		mockService.EXPECT().GetWithDeleted(gomock.Any(), courseID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &coursepb.GetWithDeletedRequest{Id: courseID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, courseID, res.GetCourseDetails().GetCourse().GetId())
		assert.Equal(t, productID, res.GetCourseDetails().GetProductId())
		assert.Equal(t, float32(99.99), res.GetCourseDetails().GetPrice())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().GetWithDeleted(gomock.Any(), courseID).Return(nil, courseservice.ErrNotFound)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &coursepb.GetWithDeletedRequest{Id: courseID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().GetWithDeleted(gomock.Any(), invalidID).Return(nil, courseservice.ErrInvalidArgument)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &coursepb.GetWithDeletedRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course ID")
	})

	t.Run("internal service error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().GetWithDeleted(gomock.Any(), courseID).Return(nil, svcErr)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &coursepb.GetWithDeletedRequest{Id: courseID})

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

	courseID := uuid.New().String()
	productID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &coursemodel.CourseDetails{
			Course: &coursemodel.Course{
				ID:   courseID,
				Name: "Course name",
			},
			Price:     99.99,
			ProductID: productID,
		}

		mockService.EXPECT().GetWithUnpublished(gomock.Any(), courseID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &coursepb.GetWithUnpublishedRequest{Id: courseID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, courseID, res.GetCourseDetails().GetCourse().GetId())
		assert.Equal(t, productID, res.GetCourseDetails().GetProductId())
		assert.Equal(t, float32(99.99), res.GetCourseDetails().GetPrice())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), courseID).Return(nil, courseservice.ErrNotFound)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &coursepb.GetWithUnpublishedRequest{Id: courseID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), invalidID).Return(nil, courseservice.ErrInvalidArgument)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &coursepb.GetWithUnpublishedRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), courseID).Return(nil, svcErr)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &coursepb.GetWithUnpublishedRequest{Id: courseID})

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

	courseID_1 := uuid.New().String()
	productID_1 := uuid.New().String()
	courseID_2 := uuid.New().String()
	productID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := []coursemodel.CourseDetails{
			{
				Course: &coursemodel.Course{
					ID:   courseID_1,
					Name: "Course 1 name",
				},
				Price:     99.99,
				ProductID: productID_1,
			},
			{
				Course: &coursemodel.Course{
					ID:   courseID_2,
					Name: "Course 2 name",
				},
				Price:     199.99,
				ProductID: productID_2,
			},
		}

		mockService.EXPECT().List(gomock.Any(), limit, offset).Return(expectedDetails, int64(2), nil).Times(1)

		// Act
		res, err := client.List(context.Background(), &coursepb.ListRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetCourseDetails()[0].GetCourse().GetName(), expectedDetails[0].Course.Name)
		assert.Equal(t, res.GetCourseDetails()[0].GetCourse().GetId(), expectedDetails[0].Course.ID)
		assert.Equal(t, res.GetCourseDetails()[0].GetPrice(), expectedDetails[0].Price)
		assert.Equal(t, res.GetCourseDetails()[0].GetProductId(), expectedDetails[0].ProductID)
		assert.Equal(t, res.GetCourseDetails()[1].GetCourse().GetName(), expectedDetails[1].Course.Name)
		assert.Equal(t, res.GetCourseDetails()[1].GetCourse().GetId(), expectedDetails[1].Course.ID)
		assert.Equal(t, res.GetCourseDetails()[1].GetPrice(), expectedDetails[1].Price)
		assert.Equal(t, res.GetCourseDetails()[1].GetProductId(), expectedDetails[1].ProductID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().List(gomock.Any(), limit, offset).Return(nil, int64(0), svcErr).Times(1)

		// Act
		res, err := client.List(context.Background(), &coursepb.ListRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get courses")
	})
}

func TestServer_ListDeleted(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	courseID_1 := uuid.New().String()
	productID_1 := uuid.New().String()
	courseID_2 := uuid.New().String()
	productID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := []coursemodel.CourseDetails{
			{
				Course: &coursemodel.Course{
					ID:   courseID_1,
					Name: "Course 1 name",
				},
				Price:     99.99,
				ProductID: productID_1,
			},
			{
				Course: &coursemodel.Course{
					ID:   courseID_2,
					Name: "Course 2 name",
				},
				Price:     199.99,
				ProductID: productID_2,
			},
		}

		mockService.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(expectedDetails, int64(2), nil).Times(1)

		// Act
		res, err := client.ListDeleted(context.Background(), &coursepb.ListDeletedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetCourseDetails()[0].GetCourse().GetName(), expectedDetails[0].Course.Name)
		assert.Equal(t, res.GetCourseDetails()[0].GetCourse().GetId(), expectedDetails[0].Course.ID)
		assert.Equal(t, res.GetCourseDetails()[0].GetPrice(), expectedDetails[0].Price)
		assert.Equal(t, res.GetCourseDetails()[0].GetProductId(), expectedDetails[0].ProductID)
		assert.Equal(t, res.GetCourseDetails()[1].GetCourse().GetName(), expectedDetails[1].Course.Name)
		assert.Equal(t, res.GetCourseDetails()[1].GetCourse().GetId(), expectedDetails[1].Course.ID)
		assert.Equal(t, res.GetCourseDetails()[1].GetPrice(), expectedDetails[1].Price)
		assert.Equal(t, res.GetCourseDetails()[1].GetProductId(), expectedDetails[1].ProductID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().ListDeleted(gomock.Any(), limit, offset).Return(nil, int64(0), svcErr).Times(1)

		// Act
		res, err := client.ListDeleted(context.Background(), &coursepb.ListDeletedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get courses")
	})
}

func TestServer_ListUnpublished(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	courseID_1 := uuid.New().String()
	productID_1 := uuid.New().String()
	courseID_2 := uuid.New().String()
	productID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := []coursemodel.CourseDetails{
			{
				Course: &coursemodel.Course{
					ID:   courseID_1,
					Name: "Course 1 name",
				},
				Price:     99.99,
				ProductID: productID_1,
			},
			{
				Course: &coursemodel.Course{
					ID:   courseID_2,
					Name: "Course 2 name",
				},
				Price:     199.99,
				ProductID: productID_2,
			},
		}

		mockService.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(expectedDetails, int64(2), nil).Times(1)

		// Act
		res, err := client.ListUnpublished(context.Background(), &coursepb.ListUnpublishedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetCourseDetails()[0].GetCourse().GetName(), expectedDetails[0].Course.Name)
		assert.Equal(t, res.GetCourseDetails()[0].GetCourse().GetId(), expectedDetails[0].Course.ID)
		assert.Equal(t, res.GetCourseDetails()[0].GetPrice(), expectedDetails[0].Price)
		assert.Equal(t, res.GetCourseDetails()[0].GetProductId(), expectedDetails[0].ProductID)
		assert.Equal(t, res.GetCourseDetails()[1].GetCourse().GetName(), expectedDetails[1].Course.Name)
		assert.Equal(t, res.GetCourseDetails()[1].GetCourse().GetId(), expectedDetails[1].Course.ID)
		assert.Equal(t, res.GetCourseDetails()[1].GetPrice(), expectedDetails[1].Price)
		assert.Equal(t, res.GetCourseDetails()[1].GetProductId(), expectedDetails[1].ProductID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().ListUnpublished(gomock.Any(), limit, offset).Return(nil, int64(0), svcErr).Times(1)

		// Act
		res, err := client.ListUnpublished(context.Background(), &coursepb.ListUnpublishedRequest{Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get courses")
	})
}

func TestServer_Create(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	courseID := uuid.New().String()
	productID := uuid.New().String()
	createReq := coursemodel.CreateRequest{
		Name:             "Course name",
		Topic:            "Course topic",
		ShortDescription: "Short description",
		AccessDuration:   30,
		Price:            99.99,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Create(gomock.Any(), &createReq).Return(&coursemodel.CreateResponse{ID: courseID, ProductID: productID}, nil)

		// Act
		res, err := client.Create(context.Background(), &coursepb.CreateRequest{
			Name:             createReq.Name,
			ShortDescription: createReq.ShortDescription,
			Topic:            createReq.Topic,
			AccessDuration:   int32(createReq.AccessDuration),
			Price:            createReq.Price,
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), courseID)
		assert.Equal(t, res.GetProductId(), productID)
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().Create(gomock.Any(), &createReq).Return(nil, svcErr)

		// Act
		res, err := client.Create(context.Background(), &coursepb.CreateRequest{
			Name:             createReq.Name,
			ShortDescription: createReq.ShortDescription,
			Topic:            createReq.Topic,
			AccessDuration:   int32(createReq.AccessDuration),
			Price:            createReq.Price,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to create course")
	})
}

func TestServer_Publish(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	courseID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Publish(gomock.Any(), courseID).Return(nil)

		// Act
		res, err := client.Publish(context.Background(), &coursepb.PublishRequest{Id: courseID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), courseID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Publish(gomock.Any(), invalidID).Return(courseservice.ErrInvalidArgument)

		// Act
		res, err := client.Publish(context.Background(), &coursepb.PublishRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Publish(gomock.Any(), courseID).Return(courseservice.ErrNotFound)

		// Act
		res, err := client.Publish(context.Background(), &coursepb.PublishRequest{Id: courseID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course not found")
	})
}

func TestServer_Unpublish(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	courseID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Unpublish(gomock.Any(), courseID).Return(nil)

		// Act
		res, err := client.Unpublish(context.Background(), &coursepb.UnpublishRequest{Id: courseID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), courseID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Unpublish(gomock.Any(), invalidID).Return(courseservice.ErrInvalidArgument)

		// Act
		res, err := client.Unpublish(context.Background(), &coursepb.UnpublishRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Unpublish(gomock.Any(), courseID).Return(courseservice.ErrNotFound)

		// Act
		res, err := client.Unpublish(context.Background(), &coursepb.UnpublishRequest{Id: courseID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course not found")
	})
}

func TestServer_Update(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	courseID := uuid.New().String()
	newName := "New course name"
	newShortDescription := "New course short description"
	newPrice := float32(99.99)

	t.Run("success", func(t *testing.T) {
		// Arrange
		updates := make(map[string]any)
		courseUpdates := map[string]any{"name": newName, "short_description": newShortDescription}
		productUpdates := map[string]any{"price": newPrice}
		updates["course"] = courseUpdates
		updates["product"] = productUpdates
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(updates, nil).Times(1)

		// Act
		res, err := client.Update(context.Background(), &coursepb.UpdateRequest{
			Id:               courseID,
			Name:             &newName,
			ShortDescription: &newShortDescription,
			Price:            &newPrice,
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.Updated.Paths[0], "updateresponse.name")
		assert.Equal(t, res.Updated.Paths[1], "updateresponse.short_description")
		assert.Equal(t, res.GetName(), courseUpdates["name"])
		assert.Equal(t, res.GetShortDescription(), courseUpdates["short_description"])
		assert.Equal(t, res.GetPrice(), productUpdates["price"])
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, courseservice.ErrInvalidArgument).Times(1)

		// Act
		res, err := client.Update(context.Background(), &coursepb.UpdateRequest{
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
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, courseservice.ErrNotFound)

		// Act
		res, err := client.Update(context.Background(), &coursepb.UpdateRequest{
			Id:               courseID,
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
		assert.Contains(t, st.Message(), "Course not found")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		svcErr := errors.New("unexpected error")
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, svcErr)

		// Act
		res, err := client.Update(context.Background(), &coursepb.UpdateRequest{
			Id:               courseID,
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

	courseID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Delete(gomock.Any(), courseID).Return(nil)

		// Act
		res, err := client.Delete(context.Background(), &coursepb.DeleteRequest{Id: courseID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), courseID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Delete(gomock.Any(), invalidID).Return(courseservice.ErrInvalidArgument)

		// Act
		res, err := client.Delete(context.Background(), &coursepb.DeleteRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Delete(gomock.Any(), courseID).Return(courseservice.ErrNotFound)

		// Act
		res, err := client.Delete(context.Background(), &coursepb.DeleteRequest{Id: courseID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course not found")
	})
}

func TestServer_DeletePermanent(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	courseID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().DeletePermanent(gomock.Any(), courseID).Return(nil)

		// Act
		res, err := client.DeletePermanent(context.Background(), &coursepb.DeletePermanentRequest{Id: courseID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), courseID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().DeletePermanent(gomock.Any(), invalidID).Return(courseservice.ErrInvalidArgument)

		// Act
		res, err := client.DeletePermanent(context.Background(), &coursepb.DeletePermanentRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().DeletePermanent(gomock.Any(), courseID).Return(courseservice.ErrNotFound)

		// Act
		res, err := client.DeletePermanent(context.Background(), &coursepb.DeletePermanentRequest{Id: courseID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course not found")
	})
}

func TestServer_Restore(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	courseID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Restore(gomock.Any(), courseID).Return(nil)

		// Act
		res, err := client.Restore(context.Background(), &coursepb.RestoreRequest{Id: courseID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), courseID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		mockService.EXPECT().Restore(gomock.Any(), invalidID).Return(courseservice.ErrInvalidArgument)

		// Act
		res, err := client.Restore(context.Background(), &coursepb.RestoreRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Restore(gomock.Any(), courseID).Return(courseservice.ErrNotFound)

		// Act
		res, err := client.Restore(context.Background(), &coursepb.RestoreRequest{Id: courseID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course not found")
	})
}
