package coursepart

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/mikhail5545/product-service-go/internal/database/course"
	coursepart "github.com/mikhail5545/product-service-go/internal/database/course_part"
	"github.com/mikhail5545/product-service-go/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	partRepo   coursepart.Repository
	courseRepo course.Repository
}

type Error struct {
	Msg  string
	Err  error
	Code int
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) GetCode() int {
	return e.Code
}

func New(pr coursepart.Repository, cr course.Repository) *Service {
	return &Service{
		partRepo:   pr,
		courseRepo: cr,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*models.CoursePart, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{
			Msg:  "Invalid Course part ID",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	part, err := s.partRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{
				Msg:  "Course part not found",
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return nil, &Error{
			Msg:  "Failed to get course part",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return part, nil
}

func (s *Service) GetReduced(ctx context.Context, id string) (*models.CoursePart, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{
			Msg:  "Invalid Course part ID",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	part, err := s.partRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{
				Msg:  "Course part not found",
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return nil, &Error{
			Msg:  "Failed to get course part",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return part, nil
}

func (s *Service) List(ctx context.Context, courseID string, limit, offset int) ([]models.CoursePart, int64, error) {
	if _, err := uuid.Parse(courseID); err != nil {
		return nil, 0, &Error{
			Msg:  "Invalid course ID",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	parts, err := s.partRepo.List(ctx, courseID, limit, offset)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to get course parts",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	total, err := s.partRepo.Count(ctx, courseID)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to count course parts",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return parts, total, nil
}

func (s *Service) ListReduced(ctx context.Context, courseID string, limit, offset int) ([]models.CoursePart, int64, error) {
	if _, err := uuid.Parse(courseID); err != nil {
		return nil, 0, &Error{
			Msg:  "Invalid course ID",
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	parts, err := s.partRepo.List(ctx, courseID, limit, offset)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to get course parts",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	total, err := s.partRepo.Count(ctx, courseID)
	if err != nil {
		return nil, 0, &Error{
			Msg:  "Failed to count course parts",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return parts, total, nil
}

func (s *Service) Create(ctx context.Context, part *models.CoursePart) (*models.CoursePart, error) {
	err := s.partRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPartRepo := s.partRepo.WithTx(tx)
		txCourseRepo := s.courseRepo.WithTx(tx)

		if _, err := uuid.Parse(part.CourseID); err != nil {
			return &Error{
				Msg:  "Invalid course ID",
				Err:  err,
				Code: http.StatusBadRequest,
			}
		}
		if part.Name == "" {
			return &Error{
				Msg:  "Course part name is required",
				Err:  nil,
				Code: http.StatusBadRequest,
			}
		}
		if part.Description == "" {
			return &Error{
				Msg:  "Course part description is required",
				Err:  nil,
				Code: http.StatusBadRequest,
			}
		}
		if part.Number <= 0 {
			return &Error{
				Msg:  "Course part number cannot be negative or null",
				Err:  nil,
				Code: http.StatusBadRequest,
			}
		}
		part.ID = uuid.New().String()

		course, err := txCourseRepo.Get(ctx, part.CourseID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{
					Msg:  "Course not found",
					Err:  err,
					Code: http.StatusNotFound,
				}
			}
			return &Error{
				Msg:  "Failed to get course",
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}

		if err := txPartRepo.Create(ctx, part); err != nil {
			return &Error{
				Msg:  "Failed to create course part",
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}
		updates := make(map[string]any)
		course.CourseParts = append(course.CourseParts, part)
		updates["course_parts"] = course.CourseParts
		if _, err := txCourseRepo.Update(ctx, course, updates); err != nil {
			return &Error{
				Msg:  "Failed to add new course part to the course",
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return part, err
}

// func (s *CourseService) DeleteCoursePart(ctx context.Context, id string) error {
// 	return s.CourseRepo.DB().Transaction(func(tx *gorm.DB) error {
// 		txCourseRepo := s.CourseRepo.WithTx(tx)

// 		if _, err := uuid.Parse(id); err != nil {
// 			return &CourseServiceError{
// 				Msg:  "Invalid course part ID",
// 				Err:  nil,
// 				Code: http.StatusBadRequest,
// 			}
// 		}

// 		part, err := txCourseRepo.FindCoursePart(ctx, id)
// 		if err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return &CourseServiceError{
// 					Msg:  "Course part not found",
// 					Err:  err,
// 					Code: http.StatusNotFound,
// 				}
// 			}
// 			return &CourseServiceError{
// 				Msg:  "Failed to get course part",
// 				Err:  err,
// 				Code: http.StatusNotFound,
// 			}
// 		}

// 		// If CoursePart has uploaded MUXVideo, call github.com/mikhail5545/media-service-go
// 		// mux client to delete MUX asset.
// 		if part.MUXVideoID != nil {
// 			_, err := s.MediaClient.DeleteMuxUpload(ctx, &muxpb.DeleteMuxUploadRequest{Id: *part.MUXVideoID})
// 			if err != nil {
// 				return &CourseServiceError{
// 					Msg:  "Failed to delete mux upload",
// 					Err:  err,
// 					Code: http.StatusServiceUnavailable,
// 				}
// 			}
// 		}

// 		course, err := txCourseRepo.FindWithParts(ctx, part.CourseID)
// 		if err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return &CourseServiceError{
// 					Msg:  "Course not found",
// 					Err:  err,
// 					Code: http.StatusNotFound,
// 				}
// 			}
// 			return &CourseServiceError{
// 				Msg:  "Failed to get course",
// 				Err:  err,
// 				Code: http.StatusNotFound,
// 			}
// 		}

// 		if err := txCourseRepo.DeleteCoursePart(ctx, id); err != nil {
// 			return &CourseServiceError{
// 				Msg:  "Failed to delete course part",
// 				Err:  err,
// 				Code: http.StatusInternalServerError,
// 			}
// 		}

// 		// Update CourseParts field in the course record.
// 		for i, part := range course.CourseParts {
// 			if part.ID == id {
// 				course.CourseParts = append(course.CourseParts[:i], course.CourseParts[i+1:]...)
// 				break
// 			}
// 		}

// 		if err := txCourseRepo.Update(ctx, course); err != nil {
// 			return &CourseServiceError{
// 				Msg:  "Failed to update course",
// 				Err:  err,
// 				Code: http.StatusInternalServerError,
// 			}
// 		}
// 		return nil
// 	})
// }

// func (s *CourseService) AddMuxVideoToCoursePart(ctx context.Context, id string, uploadID string) error {
// 	return s.CourseRepo.DB().Transaction(func(tx *gorm.DB) error {
// 		txCourseRepo := s.CourseRepo.WithTx(tx)

// 		if _, err := uuid.Parse(id); err != nil {
// 			return &CourseServiceError{
// 				Msg:  "Invalid course part ID",
// 				Err:  nil,
// 				Code: http.StatusBadRequest,
// 			}
// 		}
// 		if _, err := uuid.Parse(uploadID); err != nil {
// 			return &CourseServiceError{
// 				Msg:  "Invalid MUX upload ID",
// 				Err:  nil,
// 				Code: http.StatusBadRequest,
// 			}
// 		}

// 		part, err := txCourseRepo.FindCoursePart(ctx, id)
// 		if err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return &CourseServiceError{
// 					Msg:  "Course part not found",
// 					Err:  err,
// 					Code: http.StatusNotFound,
// 				}
// 			}
// 			return &CourseServiceError{
// 				Msg:  "Failed to get course part",
// 				Err:  err,
// 				Code: http.StatusNotFound,
// 			}
// 		}

// 		if part.MUXVideoID != nil {
// 			return &CourseServiceError{
// 				Msg:  "Course part already has MUX video",
// 				Err:  nil,
// 				Code: http.StatusBadRequest,
// 			}
// 		}

// 		part.MUXVideoID = &uploadID
// 		if err := txCourseRepo.UpdateCoursePart(ctx, part); err != nil {
// 			return &CourseServiceError{
// 				Msg:  "Failed to update course part",
// 				Err:  err,
// 				Code: http.StatusInternalServerError,
// 			}
// 		}
// 		return nil
// 	})
// }
