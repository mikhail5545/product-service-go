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

// Package coursepart provides servive-layer business logic for course parts.
package coursepart

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	courserepo "github.com/mikhail5545/product-service-go/internal/database/course"
	coursepartrepo "github.com/mikhail5545/product-service-go/internal/database/course_part"
	coursepartmodel "github.com/mikhail5545/product-service-go/internal/models/course_part"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/services/course_part_mock/service_mock.go -package=course_part_mock . Service

// Service provides service-layer business logic for course part models.
type Service interface {
	// Get retrieves a single published and not soft-deleted course part record from the database.
	// It attempts to retrieve MUXVideo information by calling the media service.
	//
	// It returns the course part record with populated MUXVideo details if found.
	// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Get(ctx context.Context, id string) (*coursepartmodel.CoursePart, error)
	// GetWithDeleted retrieves a single course part record from the database, including soft-deleted ones.
	// It attempts to retrieve MUXVideo information by calling the media service.
	//
	// It returns the course part record with populated MUXVideo details if found.
	// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	GetWithDeleted(ctx context.Context, id string) (*coursepartmodel.CoursePart, error)
	// GetWithUnpublished retrieves a single course part record from the database, including unpublished ones, but not soft-deleted.
	// It attempts to retrieve MUXVideo information by calling the media service.
	//
	// It returns the course part record with populated MUXVideo details if found.
	// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	GetWithUnpublished(ctx context.Context, id string) (*coursepartmodel.CoursePart, error)
	// GetReduced retrieves a single published and not soft-deleted course part record from the database.
	// It does not populate MUXVideo details; the MUXVideo field in the returned struct will be nil.
	//
	// Returns the course part record if found.
	// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	GetReduced(ctx context.Context, id string) (*coursepartmodel.CoursePart, error)
	// GetWithDeletedReduced retrieves a single course part record from the database, including soft-deleted ones.
	// It does not populate MUXVideo details; the MUXVideo field in the returned struct will be nil.
	//
	// Returns the course part record if found.
	// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	GetWithDeletedReduced(ctx context.Context, id string) (*coursepartmodel.CoursePart, error)
	// GetWithUnpublishedReduced retrieves a single course part record from the database, including unpublished ones, but not soft-deleted.
	// It does not populate MUXVideo details; the MUXVideo field in the returned struct will be nil.
	//
	// Returns the course part record if found.
	// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	GetWithUnpublishedReduced(ctx context.Context, id string) (*coursepartmodel.CoursePart, error)
	// List retrieves a paginated list of all published and not soft-deleted course part records for a given course ID.
	// It attempts to retrieve MUXVideo information for each course part by calling the media service.
	//
	// Returns a slice of course part records with populated MUXVideo details and the total count of such records.
	// Returns an error if the course ID is invalid (http.StatusBadRequest) or a database/internal error occurs (http.StatusInternalServerError).
	List(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, int64, error)
	// ListReduced retrieves a paginated list of all published and not soft-deleted course part records for a given course ID.
	// It does not populate MUXVideo details for the course parts.
	//
	// Returns a slice of course part records and the total count of such records.
	// Returns an error if the course ID is invalid (http.StatusBadRequest) or a database/internal error occurs (http.StatusInternalServerError).
	ListReduced(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, int64, error)
	// ListDeleted retrieves a paginated list of all soft-deleted course part records for a given course ID.
	// It does not populate MUXVideo details for the course parts.
	//
	// Returns a slice of soft-deleted course part records and the total count of such records.
	// Returns an error if the course ID is invalid (http.StatusBadRequest) or a database/internal error occurs (http.StatusInternalServerError).
	ListDeleted(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, int64, error)
	// ListDeleted retrieves a paginated list of all soft-deleted course part records for a given course ID.
	// It does not populate MUXVideo details for the course parts.
	//
	// Returns a slice of soft-deleted course part records and the total count of such records.
	// Returns an error if the course ID is invalid (http.StatusBadRequest) or a database/internal error occurs (http.StatusInternalServerError).
	ListUnpublished(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, int64, error)
	// Create creates a new CoursePart record in the database and associates it with an existing Course.
	// It validates the request payload and ensures the Course exists.
	// It also checks for uniqueness of the part number within the course.
	//
	// Returns a CreateResponse containing the newly created CoursePartID and CourseID.
	// Returns an error if the request payload is invalid (http.StatusBadRequest), the associated course is not found (http.StatusNotFound),
	// the part number is not unique within the course (http.StatusBadRequest), or a database/internal error occurs (http.StatusInternalServerError).
	Create(ctx context.Context, req *coursepartmodel.CreateRequest) (*coursepartmodel.CreateResponse, error)
	// Publish sets the 'published' field to true for a specific course part.
	// It will fail if the parent course is not published.
	//
	// Returns an error if the course part ID is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
	// the parent course is unpublished (http.StatusBadRequest), or a database/internal error occurs (http.StatusInternalServerError).
	Publish(ctx context.Context, id string) error
	// Unpublish sets the 'published' field to false for a specific course part.
	//
	// Returns an error if the course part ID is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Unpublish(ctx context.Context, id string) error
	// Update performs a partial update of a course part's information.
	// The request should contain the course part's ID and the fields to be updated.
	// At least one field must be provided for an update to occur.
	// It also ensures that the 'Number' field, if updated, remains unique within its course.
	//
	// Returns a map of the fields that were actually changed.
	// Returns an error if the request payload is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
	// the new part number is not unique within the course (http.StatusBadRequest), or a database/internal error occurs (http.StatusInternalServerError).
	Update(ctx context.Context, req *coursepartmodel.UpdateRequest) (map[string]any, error)
	// AddVideo populates the MUXVideoID field in the course part record if it is different from the previous one.
	//
	// Returns a map representation of the changed field ("mux_video_id": val) if an update occurred.
	// Returns an error if the request payload is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	AddVideo(ctx context.Context, req *coursepartmodel.AddVideoRequest) (map[string]any, error)
	// Delete performs a soft-delete for a specific course part.
	// It also unpublishes the course part, meaning it must be manually published again after restoration.
	//
	// Returns an error if the course part ID is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
	// after restore.
	Delete(ctx context.Context, id string) error
	// DeletePermanent completely removes a course part record from the database.
	//
	// Returns an error if the course part ID is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	DeletePermanent(ctx context.Context, id string) error
	// Restore restores a soft-deleted course part record.
	// The course part will remain unpublished and must be manually published again.
	//
	// Returns an error if the course part ID is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Restore(ctx context.Context, id string) error
}

// service provides service-layer business logic for course part models.
type service struct {
	partRepo   coursepartrepo.Repository
	courseRepo courserepo.Repository
}

// New creates a new Service instance with the provided course part and course repositories.
func New(pr coursepartrepo.Repository, cr courserepo.Repository) Service {
	return &service{
		partRepo:   pr,
		courseRepo: cr,
	}
}

// Get retrieves a single published and not soft-deleted course part record from the database.
// It attempts to retrieve MUXVideo information by calling the media service.
//
// It returns the course part record with populated MUXVideo details if found.
// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Get(ctx context.Context, id string) (*coursepartmodel.CoursePart, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	part, err := s.partRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve course part: %w", err)
	}
	// TODO: Implement call to media service to retrieve mux video.
	return part, nil
}

// GetWithDeleted retrieves a single course part record from the database, including soft-deleted ones.
// It attempts to retrieve MUXVideo information by calling the media service.
//
// It returns the course part record with populated MUXVideo details if found.
// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) GetWithDeleted(ctx context.Context, id string) (*coursepartmodel.CoursePart, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	part, err := s.partRepo.GetWithDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve course part: %w", err)
	}
	// TODO: Implement call to media service to retrieve mux video.
	return part, nil
}

// GetWithUnpublished retrieves a single course part record from the database, including unpublished ones, but not soft-deleted.
// It attempts to retrieve MUXVideo information by calling the media service.
//
// It returns the course part record with populated MUXVideo details if found.
// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) GetWithUnpublished(ctx context.Context, id string) (*coursepartmodel.CoursePart, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	part, err := s.partRepo.GetWithUnpublished(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve course part: %w", err)
	}
	// TODO: Implement call to media service to retrieve mux video.
	return part, nil
}

// GetReduced retrieves a single published and not soft-deleted course part record from the database.
// It does not populate MUXVideo details; the MUXVideo field in the returned struct will be nil.
//
// Returns the course part record if found.
// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) GetReduced(ctx context.Context, id string) (*coursepartmodel.CoursePart, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	part, err := s.partRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve course part: %w", err)
	}
	return part, nil
}

// GetWithDeletedReduced retrieves a single course part record from the database, including soft-deleted ones.
// It does not populate MUXVideo details; the MUXVideo field in the returned struct will be nil.
//
// Returns the course part record if found.
// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) GetWithDeletedReduced(ctx context.Context, id string) (*coursepartmodel.CoursePart, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	part, err := s.partRepo.GetWithDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve course part: %w", err)
	}
	return part, nil
}

// GetWithUnpublishedReduced retrieves a single course part record from the database, including unpublished ones, but not soft-deleted.
// It does not populate MUXVideo details; the MUXVideo field in the returned struct will be nil.
//
// Returns the course part record if found.
// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) GetWithUnpublishedReduced(ctx context.Context, id string) (*coursepartmodel.CoursePart, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	part, err := s.partRepo.GetWithUnpublished(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve course part: %w", err)
	}
	return part, nil
}

// List retrieves a paginated list of all published and not soft-deleted course part records for a given course ID.
// It attempts to retrieve MUXVideo information for each course part by calling the media service.
//
// Returns a slice of course part records with populated MUXVideo details and the total count of such records.
// Returns an error if the course ID is invalid (http.StatusBadRequest) or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) List(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, int64, error) {
	if _, err := uuid.Parse(courseID); err != nil {
		return nil, 0, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	parts, err := s.partRepo.List(ctx, courseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve course parts: %w", err)
	}
	total, err := s.partRepo.Count(ctx, courseID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count course parts: %w", err)
	}
	// TODO: Implement call to media service to retrieve mux video.
	return parts, total, nil
}

// ListReduced retrieves a paginated list of all published and not soft-deleted course part records for a given course ID.
// It does not populate MUXVideo details for the course parts.
//
// Returns a slice of course part records and the total count of such records.
// Returns an error if the course ID is invalid (http.StatusBadRequest) or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) ListReduced(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, int64, error) {
	if _, err := uuid.Parse(courseID); err != nil {
		return nil, 0, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	parts, err := s.partRepo.List(ctx, courseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve course parts: %w", err)
	}
	total, err := s.partRepo.Count(ctx, courseID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count course parts: %w", err)
	}
	return parts, total, nil
}

// ListDeleted retrieves a paginated list of all soft-deleted course part records for a given course ID.
// It does not populate MUXVideo details for the course parts.
//
// Returns a slice of soft-deleted course part records and the total count of such records.
// Returns an error if the course ID is invalid (http.StatusBadRequest) or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) ListDeleted(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, int64, error) {
	if _, err := uuid.Parse(courseID); err != nil {
		return nil, 0, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	parts, err := s.partRepo.ListDeleted(ctx, courseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve course parts: %w", err)
	}
	total, err := s.partRepo.CountDeleted(ctx, courseID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count course parts: %w", err)
	}
	return parts, total, nil
}

// ListUnpublished retrieves a paginated list of all unpublished course part records for a given course ID.
// It does not populate MUXVideo details for the course parts.
//
// Returns a slice of unpublished course part records and the total count of such records.
// Returns an error if the course ID is invalid (http.StatusBadRequest) or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) ListUnpublished(ctx context.Context, courseID string, limit, offset int) ([]coursepartmodel.CoursePart, int64, error) {
	if _, err := uuid.Parse(courseID); err != nil {
		return nil, 0, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	parts, err := s.partRepo.ListUnpublished(ctx, courseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve course parts: %w", err)
	}
	total, err := s.partRepo.CountUnpublished(ctx, courseID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count course parts: %w", err)
	}
	return parts, total, nil
}

// Create creates a new CoursePart record in the database and associates it with an existing Course.
// It validates the request payload and ensures the Course exists.
// It also checks for uniqueness of the part number within the course.
//
// Returns a CreateResponse containing the newly created CoursePartID and CourseID.
// Returns an error if the request payload is invalid (http.StatusBadRequest), the associated course is not found (http.StatusNotFound),
// the part number is not unique within the course (http.StatusBadRequest), or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Create(ctx context.Context, req *coursepartmodel.CreateRequest) (*coursepartmodel.CreateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	var partID, courseID string
	err := s.partRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPartRepo := s.partRepo.WithTx(tx)
		txCourseRepo := s.courseRepo.WithTx(tx)

		part := &coursepartmodel.CoursePart{
			ID:               uuid.New().String(),
			Name:             req.Name,
			ShortDescription: req.ShortDescription,
			Number:           req.Number,
			CourseID:         req.CourseID,
			Published:        false,
		}

		_, err := txCourseRepo.Select(ctx, part.CourseID, "id") // Only need to check if course exists
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to retrieve course: %w", err)
		}

		// Check for unique part number within the course
		count, err := txPartRepo.CountQuery(ctx, "course_id = ? AND number = ?", req.CourseID, req.Number)
		if err != nil {
			return fmt.Errorf("failed to check for unique course part number: %w", err)
		}
		if count != 0 {
			return fmt.Errorf("%w, course part with number %d already exists in course %s: %w", ErrInvalidArgument, req.Number, req.CourseID, err)
		}

		if err := txPartRepo.Create(ctx, part); err != nil {
			return fmt.Errorf("failed to create course part: %w", err)
		}

		partID = part.ID
		courseID = part.CourseID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &coursepartmodel.CreateResponse{ID: partID, CourseID: courseID}, err
}

// Publish sets the 'published' field to true for a specific course part.
// It will fail if the parent course is not published.
//
// Returns an error if the course part ID is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
// the parent course is unpublished (http.StatusBadRequest), or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Publish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.partRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPartRepo := s.partRepo.WithTx(tx)
		txCourseRepo := s.courseRepo.WithTx(tx)

		part, err := txPartRepo.GetWithUnpublished(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to retrieve course part: %w", err)
		}

		course, err := txCourseRepo.GetReduced(ctx, part.CourseID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to retrieve course: %w", err)
		}

		if !course.InStock {
			return fmt.Errorf("%w, cannot publish course part because parent course is not published: %w", ErrInvalidArgument, err)
		}

		if _, err := txPartRepo.SetPublished(ctx, id, true); err != nil {
			return fmt.Errorf("failed to publish course part: %w", err)
		}
		return nil
	})
}

// Unpublish sets the 'published' field to false for a specific course part.
//
// Returns an error if the course part ID is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Unpublish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.partRepo.DB().Transaction(func(tx *gorm.DB) error {
		ra, err := s.partRepo.WithTx(tx).SetPublished(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to upublish course part: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil
	})
}

// Update performs a partial update of a course part's information.
// The request should contain the course part's ID and the fields to be updated.
// At least one field must be provided for an update to occur.
// It also ensures that the 'Number' field, if updated, remains unique within its course.
//
// Returns a map of the fields that were actually changed.
// Returns an error if the request payload is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
// the new part number is not unique within the course (http.StatusBadRequest), or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Update(ctx context.Context, req *coursepartmodel.UpdateRequest) (map[string]any, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	updates := make(map[string]any)
	err := s.partRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPartRepo := s.partRepo.WithTx(tx)

		part, err := txPartRepo.Get(ctx, req.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to retrieve course part: %w", err)
		}

		if req.Name != nil && *req.Name != part.Name {
			updates["name"] = *req.Name
		}
		if req.ShortDescription != nil && *req.ShortDescription != part.ShortDescription {
			updates["short_description"] = *req.ShortDescription
		}
		if req.LongDescription != nil && *req.LongDescription != part.LongDescription {
			updates["long_description"] = *req.LongDescription
		}
		if req.Number != nil && *req.Number != part.Number {
			count, err := txPartRepo.CountQuery(ctx, "course_id = ? AND number = ?", part.CourseID, *req.Number) // Use part.CourseID
			if err != nil {
				return fmt.Errorf("failed to check for unique course part number: %w", err)
			}
			if count > 0 { // If count is greater than 0, a part with this number already exists in this course
				return fmt.Errorf("%w, course part with number %d already exists in course %s: %w", ErrInvalidArgument, *req.Number, part.CourseID, err)
			}
			updates["number"] = *req.Number
		}
		if len(req.Tags) > 0 {
			updates["tags"] = req.Tags
		}

		if len(updates) > 0 {
			if _, err := txPartRepo.Update(ctx, part, updates); err != nil {
				return fmt.Errorf("failed to update course part: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updates, nil
}

// AddVideo populates the MUXVideoID field in the course part record if it is different from the previous one.
//
// Returns a map representation of the changed field ("mux_video_id": val) if an update occurred.
// Returns an error if the request payload is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) AddVideo(ctx context.Context, req *coursepartmodel.AddVideoRequest) (map[string]any, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	updates := make(map[string]any)
	err := s.partRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPartRepo := s.partRepo.WithTx(tx)

		part, err := txPartRepo.Select(ctx, req.ID, "id", "course_id", "mux_video_id")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to retrieve course part: %w", err)
		}

		// Only update if the new MUXVideoID is different from the existing one
		if part.MUXVideoID == nil || req.MUXVideoID != *part.MUXVideoID {
			updates["mux_video_id"] = req.MUXVideoID
		}

		if len(updates) > 0 { // Only perform update if there are actual changes
			if _, err := txPartRepo.Update(ctx, part, updates); err != nil {
				return fmt.Errorf("failed to update course part: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updates, err
}

// Delete performs a soft-delete for a specific course part.
// It also unpublishes the course part, meaning it must be manually published again after restoration.
//
// Returns an error if the course part ID is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
// after restore.
func (s *service) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.partRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPartRepo := s.partRepo.WithTx(tx)

		// Check if the record exists first (including unpublished, but not soft-deleted)
		if _, err := txPartRepo.GetWithUnpublished(ctx, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to retrieve course part: %w", err)
		}

		if _, err := txPartRepo.SetPublished(ctx, id, false); err != nil {
			return fmt.Errorf("failed to unpublish course part: %w", err)
		}

		// Perform soft-delete
		if _, err := txPartRepo.Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete course part: %w", err)
		}
		return nil
	})
}

// DeletePermanent completely removes a course part record from the database.
//
// Returns an error if the course part ID is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) DeletePermanent(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.partRepo.DB().Transaction(func(tx *gorm.DB) error {
		ra, err := s.partRepo.WithTx(tx).DeletePermanent(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete course part: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil
	})
}

// Restore restores a soft-deleted course part record.
// The course part will remain unpublished and must be manually published again.
//
// Returns an error if the course part ID is invalid (http.StatusBadRequest), the course part is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Restore(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.partRepo.DB().Transaction(func(tx *gorm.DB) error {
		ra, err := s.partRepo.WithTx(tx).Restore(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to restore course part: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil
	})
}
