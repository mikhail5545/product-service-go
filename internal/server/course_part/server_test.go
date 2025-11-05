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

package coursepart

import (
	"context"
	"log"
	"net"
	"net/http"
	"testing"

	"github.com/google/uuid"
	coursepartmodel "github.com/mikhail5545/product-service-go/internal/models/course_part"
	coursepartservice "github.com/mikhail5545/product-service-go/internal/services/course_part"
	coursepartmock "github.com/mikhail5545/product-service-go/internal/test/services/course_part_mock"
	coursepartpb "github.com/mikhail5545/proto-go/proto/course_part/v0"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func setupTestServer(t *testing.T) (coursepartpb.CoursePartServiceClient, *coursepartmock.MockService, func()) {
	t.Helper()

	// 1. Create mock controller and mock service
	ctrl := gomock.NewController(t)
	mockService := coursepartmock.NewMockService(ctrl)

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

	client := coursepartpb.NewCoursePartServiceClient(conn)

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

	partID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &coursepartmodel.CoursePart{
			ID:               partID,
			Name:             "Course part name",
			ShortDescription: "Course part short description",
			Number:           33,
		}

		mockService.EXPECT().Get(gomock.Any(), partID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.Get(context.Background(), &coursepartpb.GetRequest{Id: partID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, partID, res.GetCoursePart().GetId())
		assert.Equal(t, expectedDetails.Name, res.GetCoursePart().GetName())
		assert.Equal(t, expectedDetails.ShortDescription, res.GetCoursePart().GetShortDescription())
		assert.Equal(t, int32(expectedDetails.Number), res.GetCoursePart().GetNumber())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Course part not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Get(gomock.Any(), partID).Return(nil, serviceErr)

		// Act
		res, err := client.Get(context.Background(), &coursepartpb.GetRequest{Id: partID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course part not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &coursepartservice.Error{
			Msg:  "Invalid course part ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().Get(gomock.Any(), invalidID).Return(nil, serviceErr)

		// Act
		res, err := client.Get(context.Background(), &coursepartpb.GetRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course part ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Database error",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().Get(gomock.Any(), partID).Return(nil, serviceErr)

		// Act
		res, err := client.Get(context.Background(), &coursepartpb.GetRequest{Id: partID})

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

	partID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &coursepartmodel.CoursePart{
			ID:               partID,
			Name:             "Course part name",
			ShortDescription: "Course part short description",
			Number:           33,
		}

		mockService.EXPECT().GetWithDeleted(gomock.Any(), partID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &coursepartpb.GetWithDeletedRequest{Id: partID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, partID, res.GetCoursePart().GetId())
		assert.Equal(t, expectedDetails.Name, res.GetCoursePart().GetName())
		assert.Equal(t, expectedDetails.ShortDescription, res.GetCoursePart().GetShortDescription())
		assert.Equal(t, int32(expectedDetails.Number), res.GetCoursePart().GetNumber())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Course part not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().GetWithDeleted(gomock.Any(), partID).Return(nil, serviceErr)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &coursepartpb.GetWithDeletedRequest{Id: partID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course part not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &coursepartservice.Error{
			Msg:  "Invalid course part ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().GetWithDeleted(gomock.Any(), invalidID).Return(nil, serviceErr)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &coursepartpb.GetWithDeletedRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course part ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Database error",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().GetWithDeleted(gomock.Any(), partID).Return(nil, serviceErr)

		// Act
		res, err := client.GetWithDeleted(context.Background(), &coursepartpb.GetWithDeletedRequest{Id: partID})

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

	partID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedDetails := &coursepartmodel.CoursePart{
			ID:               partID,
			Name:             "Course part name",
			ShortDescription: "Course part short description",
			Number:           33,
		}

		mockService.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(expectedDetails, nil).Times(1)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &coursepartpb.GetWithUnpublishedRequest{Id: partID})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, partID, res.GetCoursePart().GetId())
		assert.Equal(t, expectedDetails.Name, res.GetCoursePart().GetName())
		assert.Equal(t, expectedDetails.ShortDescription, res.GetCoursePart().GetShortDescription())
		assert.Equal(t, int32(expectedDetails.Number), res.GetCoursePart().GetNumber())
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Course part not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(nil, serviceErr)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &coursepartpb.GetWithUnpublishedRequest{Id: partID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course part not found")
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &coursepartservice.Error{
			Msg:  "Invalid course part ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), invalidID).Return(nil, serviceErr)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &coursepartpb.GetWithUnpublishedRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course part ID")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Database error",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().GetWithUnpublished(gomock.Any(), partID).Return(nil, serviceErr)

		// Act
		res, err := client.GetWithUnpublished(context.Background(), &coursepartpb.GetWithUnpublishedRequest{Id: partID})

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

	courseID := uuid.New().String()
	partID_1 := uuid.New().String()
	partID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedParts := []coursepartmodel.CoursePart{
			{
				ID:               partID_1,
				CourseID:         courseID,
				Name:             "Course part 1 name",
				ShortDescription: "Course part 1 short description",
				Number:           22,
			},
			{
				ID:               partID_2,
				CourseID:         courseID,
				Name:             "Course part 2 name",
				ShortDescription: "Course part  short description",
				Number:           33,
			},
		}

		mockService.EXPECT().List(gomock.Any(), courseID, limit, offset).Return(expectedParts, int64(2), nil).Times(1)

		// Act
		res, err := client.List(context.Background(), &coursepartpb.ListRequest{CourseId: courseID, Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetCourseParts()[0].GetName(), expectedParts[0].Name)
		assert.Equal(t, res.GetCourseParts()[0].GetId(), expectedParts[0].ID)
		assert.Equal(t, res.GetCourseParts()[0].GetNumber(), int32(expectedParts[0].Number))
		assert.Equal(t, res.GetCourseParts()[0].GetCourseId(), expectedParts[0].CourseID)
		assert.Equal(t, res.GetCourseParts()[1].GetName(), expectedParts[1].Name)
		assert.Equal(t, res.GetCourseParts()[1].GetId(), expectedParts[1].ID)
		assert.Equal(t, res.GetCourseParts()[1].GetNumber(), int32(expectedParts[1].Number))
		assert.Equal(t, res.GetCourseParts()[1].GetCourseId(), expectedParts[1].CourseID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Failed to get course parts",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().List(gomock.Any(), courseID, limit, offset).Return(nil, int64(0), serviceErr).Times(1)

		// Act
		res, err := client.List(context.Background(), &coursepartpb.ListRequest{CourseId: courseID, Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get course parts")
	})

	t.Run("invalid course ID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &coursepartservice.Error{
			Msg:  "Invalid course ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().List(gomock.Any(), invalidID, limit, offset).Return(nil, int64(0), serviceErr).Times(1)

		// Act
		res, err := client.List(context.Background(), &coursepartpb.ListRequest{CourseId: invalidID, Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course ID")
	})
}

func TestServer_ListDeleted(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	courseID := uuid.New().String()
	partID_1 := uuid.New().String()
	partID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedParts := []coursepartmodel.CoursePart{
			{
				ID:               partID_1,
				CourseID:         courseID,
				Name:             "Course part 1 name",
				ShortDescription: "Course part 1 short description",
				Number:           22,
			},
			{
				ID:               partID_2,
				CourseID:         courseID,
				Name:             "Course part 2 name",
				ShortDescription: "Course part  short description",
				Number:           33,
			},
		}

		mockService.EXPECT().ListDeleted(gomock.Any(), courseID, limit, offset).Return(expectedParts, int64(2), nil).Times(1)

		// Act
		res, err := client.ListDeleted(context.Background(), &coursepartpb.ListDeletedRequest{CourseId: courseID, Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetCourseParts()[0].GetName(), expectedParts[0].Name)
		assert.Equal(t, res.GetCourseParts()[0].GetId(), expectedParts[0].ID)
		assert.Equal(t, res.GetCourseParts()[0].GetNumber(), int32(expectedParts[0].Number))
		assert.Equal(t, res.GetCourseParts()[0].GetCourseId(), expectedParts[0].CourseID)
		assert.Equal(t, res.GetCourseParts()[1].GetName(), expectedParts[1].Name)
		assert.Equal(t, res.GetCourseParts()[1].GetId(), expectedParts[1].ID)
		assert.Equal(t, res.GetCourseParts()[1].GetNumber(), int32(expectedParts[1].Number))
		assert.Equal(t, res.GetCourseParts()[1].GetCourseId(), expectedParts[1].CourseID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Failed to get course parts",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().ListDeleted(gomock.Any(), courseID, limit, offset).Return(nil, int64(0), serviceErr).Times(1)

		// Act
		res, err := client.ListDeleted(context.Background(), &coursepartpb.ListDeletedRequest{CourseId: courseID, Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get course parts")
	})

	t.Run("invalid course ID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &coursepartservice.Error{
			Msg:  "Invalid course ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().ListDeleted(gomock.Any(), invalidID, limit, offset).Return(nil, int64(0), serviceErr).Times(1)

		// Act
		res, err := client.ListDeleted(context.Background(), &coursepartpb.ListDeletedRequest{CourseId: invalidID, Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course ID")
	})
}

func TestServer_ListUnpublished(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	courseID := uuid.New().String()
	partID_1 := uuid.New().String()
	partID_2 := uuid.New().String()

	limit, offset := 2, 0

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedParts := []coursepartmodel.CoursePart{
			{
				ID:               partID_1,
				CourseID:         courseID,
				Name:             "Course part 1 name",
				ShortDescription: "Course part 1 short description",
				Number:           22,
			},
			{
				ID:               partID_2,
				CourseID:         courseID,
				Name:             "Course part 2 name",
				ShortDescription: "Course part  short description",
				Number:           33,
			},
		}

		mockService.EXPECT().ListUnpublished(gomock.Any(), courseID, limit, offset).Return(expectedParts, int64(2), nil).Times(1)

		// Act
		res, err := client.ListUnpublished(context.Background(), &coursepartpb.ListUnpublishedRequest{CourseId: courseID, Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetCourseParts()[0].GetName(), expectedParts[0].Name)
		assert.Equal(t, res.GetCourseParts()[0].GetId(), expectedParts[0].ID)
		assert.Equal(t, res.GetCourseParts()[0].GetNumber(), int32(expectedParts[0].Number))
		assert.Equal(t, res.GetCourseParts()[0].GetCourseId(), expectedParts[0].CourseID)
		assert.Equal(t, res.GetCourseParts()[1].GetName(), expectedParts[1].Name)
		assert.Equal(t, res.GetCourseParts()[1].GetId(), expectedParts[1].ID)
		assert.Equal(t, res.GetCourseParts()[1].GetNumber(), int32(expectedParts[1].Number))
		assert.Equal(t, res.GetCourseParts()[1].GetCourseId(), expectedParts[1].CourseID)
		assert.Equal(t, res.GetTotal(), int64(2))
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Failed to get course parts",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().ListUnpublished(gomock.Any(), courseID, limit, offset).Return(nil, int64(0), serviceErr).Times(1)

		// Act
		res, err := client.ListUnpublished(context.Background(), &coursepartpb.ListUnpublishedRequest{CourseId: courseID, Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to get course parts")
	})

	t.Run("invalid course ID", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &coursepartservice.Error{
			Msg:  "Invalid course ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().ListUnpublished(gomock.Any(), invalidID, limit, offset).Return(nil, int64(0), serviceErr).Times(1)

		// Act
		res, err := client.ListUnpublished(context.Background(), &coursepartpb.ListUnpublishedRequest{CourseId: invalidID, Limit: int32(limit), Offset: int32(offset)})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course ID")
	})
}

func TestServer_Create(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	courseID := uuid.New().String()
	partID := uuid.New().String()

	createReq := coursepartmodel.CreateRequest{
		CourseID:         courseID,
		Name:             "Course part name",
		ShortDescription: "Course part short description",
		Number:           2,
	}

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Create(gomock.Any(), &createReq).Return(&coursepartmodel.CreateResponse{ID: partID, CourseID: courseID}, nil)

		// Act
		res, err := client.Create(context.Background(), &coursepartpb.CreateRequest{
			CourseId:         createReq.CourseID,
			Name:             createReq.Name,
			ShortDescription: createReq.ShortDescription,
			Number:           int32(createReq.Number),
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), partID)
		assert.Equal(t, res.GetCourseId(), courseID)
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Failed to create course part",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().Create(gomock.Any(), &createReq).Return(nil, serviceErr)

		// Act
		res, err := client.Create(context.Background(), &coursepartpb.CreateRequest{
			CourseId:         createReq.CourseID,
			Name:             createReq.Name,
			ShortDescription: createReq.ShortDescription,
			Number:           int32(createReq.Number),
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "Failed to create course part")
	})
}

func TestServer_Publish(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	partID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Publish(gomock.Any(), partID).Return(nil)

		// Act
		res, err := client.Publish(context.Background(), &coursepartpb.PublishRequest{Id: partID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), partID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &coursepartservice.Error{
			Msg:  "Invalid course part ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().Publish(gomock.Any(), invalidID).Return(serviceErr)

		// Act
		res, err := client.Publish(context.Background(), &coursepartpb.PublishRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course part ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Course part not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Publish(gomock.Any(), partID).Return(serviceErr)

		// Act
		res, err := client.Publish(context.Background(), &coursepartpb.PublishRequest{Id: partID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course part not found")
	})
}

func TestServer_Unpublish(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	partID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Unpublish(gomock.Any(), partID).Return(nil)

		// Act
		res, err := client.Unpublish(context.Background(), &coursepartpb.UnpublishRequest{Id: partID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), partID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &coursepartservice.Error{
			Msg:  "Invalid course part ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().Unpublish(gomock.Any(), invalidID).Return(serviceErr)

		// Act
		res, err := client.Unpublish(context.Background(), &coursepartpb.UnpublishRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course part ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Course part not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Unpublish(gomock.Any(), partID).Return(serviceErr)

		// Act
		res, err := client.Unpublish(context.Background(), &coursepartpb.UnpublishRequest{Id: partID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course part not found")
	})
}

func TestServer_Update(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	partID := uuid.New().String()
	newName := "New course name"
	newShortDescription := "New course short description"
	newNumber := int32(33)

	t.Run("success", func(t *testing.T) {
		// Arrange
		updates := map[string]any{"name": newName, "short_description": newShortDescription, "number": newNumber}
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(updates, nil).Times(1)

		// Act
		res, err := client.Update(context.Background(), &coursepartpb.UpdateRequest{
			Id:               partID,
			Name:             &newName,
			ShortDescription: &newShortDescription,
			Number:           &newNumber,
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.Updated.Paths[0], "updateresponse.name")
		assert.Equal(t, res.Updated.Paths[1], "updateresponse.short_description")
		assert.Equal(t, res.GetName(), updates["name"])
		assert.Equal(t, res.GetShortDescription(), updates["short_description"])
		assert.Equal(t, res.GetNumber(), updates["number"])
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &coursepartservice.Error{
			Msg:  "Invalid request payload",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, serviceErr).Times(1)

		// Act
		res, err := client.Update(context.Background(), &coursepartpb.UpdateRequest{
			Id:               invalidID,
			Name:             &newName,
			ShortDescription: &newShortDescription,
			Number:           &newNumber,
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
		serviceErr := &coursepartservice.Error{
			Msg:  "Course part not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, serviceErr)

		// Act
		res, err := client.Update(context.Background(), &coursepartpb.UpdateRequest{
			Id:               partID,
			Name:             &newName,
			ShortDescription: &newShortDescription,
			Number:           &newNumber,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course part not found")
	})

	t.Run("internal server error", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Database error",
			Code: http.StatusInternalServerError,
		}
		mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, serviceErr)

		// Act
		res, err := client.Update(context.Background(), &coursepartpb.UpdateRequest{
			Id:               partID,
			Name:             &newName,
			ShortDescription: &newShortDescription,
			Number:           &newNumber,
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

	partID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Delete(gomock.Any(), partID).Return(nil)

		// Act
		res, err := client.Delete(context.Background(), &coursepartpb.DeleteRequest{Id: partID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), partID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &coursepartservice.Error{
			Msg:  "Invalid course part ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().Delete(gomock.Any(), invalidID).Return(serviceErr)

		// Act
		res, err := client.Delete(context.Background(), &coursepartpb.DeleteRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course part ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Course part not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Delete(gomock.Any(), partID).Return(serviceErr)

		// Act
		res, err := client.Delete(context.Background(), &coursepartpb.DeleteRequest{Id: partID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course part not found")
	})
}

func TestServer_DeletePermanent(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	partID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().DeletePermanent(gomock.Any(), partID).Return(nil)

		// Act
		res, err := client.DeletePermanent(context.Background(), &coursepartpb.DeletePermanentRequest{Id: partID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), partID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &coursepartservice.Error{
			Msg:  "Invalid course part ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().DeletePermanent(gomock.Any(), invalidID).Return(serviceErr)

		// Act
		res, err := client.DeletePermanent(context.Background(), &coursepartpb.DeletePermanentRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course part ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Course part not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().DeletePermanent(gomock.Any(), partID).Return(serviceErr)

		// Act
		res, err := client.DeletePermanent(context.Background(), &coursepartpb.DeletePermanentRequest{Id: partID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course part not found")
	})
}

func TestServer_Restore(t *testing.T) {
	client, mockService, cleanup := setupTestServer(t)
	defer cleanup()

	partID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockService.EXPECT().Restore(gomock.Any(), partID).Return(nil)

		// Act
		res, err := client.Restore(context.Background(), &coursepartpb.RestoreRequest{Id: partID})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, res.GetId(), partID)
	})

	t.Run("invalid argument", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"
		serviceErr := &coursepartservice.Error{
			Msg:  "Invalid course part ID",
			Code: http.StatusBadRequest,
		}
		mockService.EXPECT().Restore(gomock.Any(), invalidID).Return(serviceErr)

		// Act
		res, err := client.Restore(context.Background(), &coursepartpb.RestoreRequest{Id: invalidID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "Invalid course part ID")
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		serviceErr := &coursepartservice.Error{
			Msg:  "Course part not found",
			Code: http.StatusNotFound,
		}
		mockService.EXPECT().Restore(gomock.Any(), partID).Return(serviceErr)

		// Act
		res, err := client.Restore(context.Background(), &coursepartpb.RestoreRequest{Id: partID})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "Course part not found")
	})
}
