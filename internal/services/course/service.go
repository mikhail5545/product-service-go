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

// Package course provides service-layer business logic for courses.
package course

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	courserepo "github.com/mikhail5545/product-service-go/internal/database/course"
	coursepartrepo "github.com/mikhail5545/product-service-go/internal/database/course_part"
	productrepo "github.com/mikhail5545/product-service-go/internal/database/product"
	coursemodel "github.com/mikhail5545/product-service-go/internal/models/course"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	"github.com/mikhail5545/product-service-go/internal/models/product"
	imageservice "github.com/mikhail5545/product-service-go/internal/services/image_manager"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/services/course_mock/service_mock.go -package=course_mock . Service

// Service provides service-layer business logic for course models.
type Service interface {
	// Get retrieves a single published and not soft-deleted course record from the database,
	// along with its associated product details (price and product ID). Also it preloads all
	// its associated course part records.
	//
	// Returns a CourseDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	Get(ctx context.Context, id string) (*coursemodel.CourseDetails, error)
	// GetWithDeleted retrieves a single course record from the database, including soft-deleted ones,
	// along with its associated product details. Also it preloads all its associated course part records.
	//
	// Returns a CourseDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	GetWithDeleted(ctx context.Context, id string) (*coursemodel.CourseDetails, error)
	// GetWithUnpublished retrieves a single course record from the database, including unpublished ones (but not soft-deleted),
	// along with its associated product details. Also it preloads all its associated course part records.
	//
	// Returns a CourseDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	GetWithUnpublished(ctx context.Context, id string) (*coursemodel.CourseDetails, error)
	// GetReduced retrieves a single published and not soft-deleted course record from the database,
	// along with its associated product details (price and product ID).
	//
	// Returns a CourseDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	GetReduced(ctx context.Context, id string) (*coursemodel.CourseDetails, error)
	// GetReducedWithDeleted retrieves a single course record from the database, including soft-deleted ones,
	// along with its associated product details.
	//
	// Returns a CourseDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	GetReducedWithDeleted(ctx context.Context, id string) (*coursemodel.CourseDetails, error)
	// List retrieves a paginated list of all published and not soft-deleted course records.
	// Each record is returned with its associated product details and preloaded course part records.
	//
	// Returns a slice of CourseDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs.
	List(ctx context.Context, limit, offset int) ([]coursemodel.CourseDetails, int64, error)
	// ListDeleted retrieves a paginated list of all soft-deleted course records.
	// Each record is returned with its associated product details.
	//
	// Returns a slice of CourseDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs.
	ListDeleted(ctx context.Context, limit, offset int) ([]coursemodel.CourseDetails, int64, error)
	// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) course records.
	// Each record is returned with its associated product details.
	//
	// Returns a slice of CourseDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs.
	ListUnpublished(ctx context.Context, limit, offset int) ([]coursemodel.CourseDetails, int64, error)
	// Create creates a new Course record and its associated Product record in the database.
	// It validates the request payload to ensure all required fields are present.
	// Both the course and the product are created in an unpublished state (`InStock: false`).
	//
	// Returns a CreateResponse containing the newly created CourseID and ProductID.
	// Returns an error if the request payload is invalid (ErrInvalidArgument) or a database/internal error occurs.
	Create(ctx context.Context, req *coursemodel.CreateRequest) (*coursemodel.CreateResponse, error)
	// Publish sets the `InStock` field to true for a course and its associated product,
	// making it available in the catalog. All of its associated course parts (if they exist)
	// should be unpublished separately.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Publish(ctx context.Context, id string) error
	// Unpublish sets the `InStock` field to false for a course, its associated course parts
	// and its associated product, archiving it from the catalog.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Unpublish(ctx context.Context, id string) error
	// Update performs a partial update of a course and its related product.
	// The request should contain the course's ID and the fields to be updated.
	// At least one field must be provided for an update to occur.
	//
	// Returns a map containing the fields that were actually changed, nested under "course" and "product" keys.
	// Example: `{"course": {"name": "new name"}, "product": {"price": 99.99}}`
	// Returns an error if the request payload is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Update(ctx context.Context, req *coursemodel.UpdateRequest) (map[string]any, error)
	// Delete performs a soft-delete of a course, its associated course parts
	// and its associated product record.
	// It also unpublishes all records, meaning they must be manually published again after restoration.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Delete(ctx context.Context, id string) error
	// DeletePermanent performs a complete delete of a course, its associated course parts
	// and its associated product record.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	DeletePermanent(ctx context.Context, id string) error
	// Restore performs a restore of a course, its associated course parts
	// and its related product record.
	// Course record, its associated course part records and its related product record
	// are not being published. This should be done manually.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Restore(ctx context.Context, id string) error
	// AddImage adds a new image to a course. It's called by media-service-go upon successful image upload.
	// It uses courseOwnerRepoAdapter to call [imageservice.AddImage] and add an image to the course.
	//
	// Returns an error if:
	//   - The request payload is invalid ([imageservice.ErrInvalidArgument]).
	//   - The course (owner) is not found ([imageservice.ErrOwnerNotFound]).
	//   - The image limit (5) is exceeded ([imageservice.ErrImageLimitExceeded]).
	//   - A database/internal error occurs.
	//
	// Deprecated: use [image.Add] instead, it performs adding of an image for all services (generic method).
	AddImage(ctx context.Context, req *imagemodel.AddRequest) error
	// DeleteImage removes an image from a course. It's called by media-service-go upon successful image deletion.
	// It uses courseOwnerRepoAdapter to call [imageservice.DeleteImage] and delete an image from the course.
	//
	// Returns an error if:
	//   - The request payload is invalid ([imageservice.ErrInvalidArgument]).
	//   - The course (owner) is not found ([imageservice.ErrOwnerNotFound]).
	//   - The image is not found on course (owner) ([imageservice.ErrImageNotFoundOnOwner]).
	//   - A database/internal error occurs.
	//
	// Deprecated: use [image.Delete] instead, it handles image deletion for all services (generic method).
	DeleteImage(ctx context.Context, req *imagemodel.DeleteRequest) error
	// AddImageBatch adds an image for a batch of courses. It uses courseOwnerRepoAdapter
	// to call [imageservice.AddImageBatch] and append images to the course. It's called by media-service-go
	// upon successfull context change.
	//
	// It returns a number of affected courses.
	// Returns an error if the request is invalid ([imageservice.ErrInvalidArgument]), no courses (owners) are not found ([imageservice.ErrOwnersNotFound])
	// or a database/internal error occurs.
	//
	// Deprecated: use generic function instead
	//
	// Deprecated: use [image.AddBatch] instead, it handles image batch add for all services (generic method).
	AddImageBatch(ctx context.Context, req *imagemodel.AddBatchRequest) (int, error)
	// DeleteImageBatch removes an image from a batch of courses. It uses courseOwnerRepoAdapter
	// to call [imageservice.DeleteImageBatch] and append images to the course.
	//
	// It returns a number of affected courses.
	// Returns an error if the request is invalid ([imageservice.ErrInvalidArgument]), no courses (owners) are not found ([imageservice.ErrOwnersNotFound]),
	// no associations were found ([imageservice.ErrAssociationsNotFound]) or a database/internal error occurs.
	//
	// Deprecated: use [image.DeleteBatch] instead, it handles image batch delete for all services (generic method).
	DeleteImageBatch(ctx context.Context, req *imagemodel.DeleteBatchRequst) (int, error)
}

// service provides service-layer business logic for course models.
// It holds [course.Repository],
// [product.Repository] and [coursepart.Repository]
// instances to perform database operations.
type service struct {
	CourseRepo  courserepo.Repository
	ProductRepo productrepo.Repository
	PartRepo    coursepartrepo.Repository
	ImageSvc    imageservice.Service
}

// New creates a new Service instance with provided
// course, product and course part repositories.
func New(
	cr courserepo.Repository,
	pr productrepo.Repository,
	cpr coursepartrepo.Repository,
	is imageservice.Service,
) Service {
	return &service{
		CourseRepo:  cr,
		ProductRepo: pr,
		PartRepo:    cpr,
		ImageSvc:    is,
	}
}

// Get retrieves a single published and not soft-deleted course record from the database,
// along with its associated product details (price and product ID). Also it preloads all
// its associated course part records.
//
// Returns a CourseDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Get(ctx context.Context, id string) (*coursemodel.CourseDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	courseRec, err := s.CourseRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve course: %w", err)
	}
	productRec, err := s.ProductRepo.GetByDetailsID(ctx, courseRec.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve product for course: %w", err)
	}
	return &coursemodel.CourseDetails{
		Course:    courseRec,
		Price:     productRec.Price,
		ProductID: productRec.ID,
	}, nil
}

// GetWithDeleted retrieves a single course record from the database, including soft-deleted ones,
// along with its associated product details. Also it preloads all its associated course part records.
//
// Returns a CourseDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) GetWithDeleted(ctx context.Context, id string) (*coursemodel.CourseDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	courseRec, err := s.CourseRepo.GetWithDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve course: %w", err)
	}
	productRec, err := s.ProductRepo.GetWithDeletedByDetailsID(ctx, courseRec.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve product for course: %w", err)
	}
	return &coursemodel.CourseDetails{
		Course:    courseRec,
		Price:     productRec.Price,
		ProductID: productRec.ID,
	}, nil
}

// GetWithUnpublished retrieves a single course record from the database, including unpublished ones (but not soft-deleted),
// along with its associated product details. Also it preloads all its associated course part records.
//
// Returns a CourseDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) GetWithUnpublished(ctx context.Context, id string) (*coursemodel.CourseDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	courseRec, err := s.CourseRepo.GetWithUnpublished(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve course: %w", err)
	}
	productRec, err := s.ProductRepo.GetWithUnpublishedByDetailsID(ctx, courseRec.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve product for course: %w", err)
	}
	return &coursemodel.CourseDetails{
		Course:    courseRec,
		Price:     productRec.Price,
		ProductID: productRec.ID,
	}, nil
}

// GetReduced retrieves a single published and not soft-deleted course record from the database,
// along with its associated product details (price and product ID).
//
// Returns a CourseDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) GetReduced(ctx context.Context, id string) (*coursemodel.CourseDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	courseRec, err := s.CourseRepo.GetReduced(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve course: %w", err)
	}
	productRec, err := s.ProductRepo.GetByDetailsID(ctx, courseRec.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve product for course: %w", err)
	}

	return &coursemodel.CourseDetails{
		Course:    courseRec,
		Price:     productRec.Price,
		ProductID: productRec.ID,
	}, nil
}

// GetReducedWithDeleted retrieves a single course record from the database, including soft-deleted ones,
// along with its associated product details.
//
// Returns a CourseDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) GetReducedWithDeleted(ctx context.Context, id string) (*coursemodel.CourseDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	courseRec, err := s.CourseRepo.GetReducedWithDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve course: %w", err)
	}
	productRec, err := s.ProductRepo.GetWithDeletedByDetailsID(ctx, courseRec.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve product for course: %w", err)
	}

	return &coursemodel.CourseDetails{
		Course:    courseRec,
		Price:     productRec.Price,
		ProductID: productRec.ID,
	}, nil
}

// List retrieves a paginated list of all published and not soft-deleted course records.
// Each record is returned with its associated product details and preloaded course part records.
//
// Returns a slice of CourseDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs.
func (s *service) List(ctx context.Context, limit, offset int) ([]coursemodel.CourseDetails, int64, error) {
	courses, err := s.CourseRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve courses: %w", err)
	}
	total, err := s.CourseRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count courses: %w", err)
	}

	coursesMap := make(map[string]*coursemodel.Course, len(courses))
	var courseIDs []string
	for i := range courses {
		coursesMap[courses[i].ID] = &courses[i]
		courseIDs = append(courseIDs, courses[i].ID)
	}

	products, err := s.ProductRepo.SelectByDetailsIDs(ctx, courseIDs, "id", "price", "details_id")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}
	var allDetails []coursemodel.CourseDetails
	for _, p := range products {
		allDetails = append(allDetails, coursemodel.CourseDetails{
			Course:    coursesMap[p.DetailsID],
			Price:     p.Price,
			ProductID: p.ID,
		})
	}
	return allDetails, total, nil
}

// ListDeleted retrieves a paginated list of all soft-deleted course records.
// Each record is returned with its associated product details.
//
// Returns a slice of CourseDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs.
func (s *service) ListDeleted(ctx context.Context, limit, offset int) ([]coursemodel.CourseDetails, int64, error) {
	courses, err := s.CourseRepo.ListDeleted(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve courses: %w", err)
	}
	total, err := s.CourseRepo.CountDeleted(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count courses: %w", err)
	}

	coursesMap := make(map[string]*coursemodel.Course, len(courses))
	var courseIDs []string
	for i := range courses {
		coursesMap[courses[i].ID] = &courses[i]
		courseIDs = append(courseIDs, courses[i].ID)
	}

	products, err := s.ProductRepo.SelectWithDeletedByDetailsIDs(ctx, courseIDs, "id", "price", "details_id")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}
	var allDetails []coursemodel.CourseDetails
	for _, p := range products {
		allDetails = append(allDetails, coursemodel.CourseDetails{
			Course:    coursesMap[p.DetailsID],
			Price:     p.Price,
			ProductID: p.ID,
		})
	}
	return allDetails, total, nil
}

// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) course records.
// Each record is returned with its associated product details.
//
// Returns a slice of CourseDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs.
func (s *service) ListUnpublished(ctx context.Context, limit, offset int) ([]coursemodel.CourseDetails, int64, error) {
	courses, err := s.CourseRepo.ListUnpublished(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve courses: %w", err)
	}
	total, err := s.CourseRepo.CountUnpublished(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count courses: %w", err)
	}

	coursesMap := make(map[string]*coursemodel.Course, len(courses))
	var courseIDs []string
	for i := range courses {
		coursesMap[courses[i].ID] = &courses[i]
		courseIDs = append(courseIDs, courses[i].ID)
	}

	products, err := s.ProductRepo.SelectWithUnpublishedByDetailsIDs(ctx, courseIDs, "id", "price", "details_id")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}
	var allDetails []coursemodel.CourseDetails
	for _, p := range products {
		allDetails = append(allDetails, coursemodel.CourseDetails{
			Course:    coursesMap[p.DetailsID],
			Price:     p.Price,
			ProductID: p.ID,
		})
	}
	return allDetails, total, nil
}

// Create creates a new Course record and its associated Product record in the database.
// It validates the request payload to ensure all required fields are present.
// Both the course and the product are created in an unpublished state (`InStock: false`).
//
// Returns a CreateResponse containing the newly created CourseID and ProductID.
// Returns an error if the request payload is invalid (ErrInvalidArgument) or a database/internal error occurs.
func (s *service) Create(ctx context.Context, req *coursemodel.CreateRequest) (*coursemodel.CreateResponse, error) {
	var courseID, productID string
	err := s.CourseRepo.DB().Transaction(func(tx *gorm.DB) error {
		txCourseRepo := s.CourseRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		if err := req.Validate(); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
		}

		course := &coursemodel.Course{
			ID:               uuid.New().String(),
			Name:             req.Name,
			ShortDescription: req.ShortDescription,
			Topic:            req.Topic,
			AccessDuration:   req.AccessDuration,
			InStock:          false,
		}

		product := &product.Product{
			ID:          uuid.New().String(),
			Price:       req.Price,
			DetailsID:   course.ID,
			DetailsType: "course",
			InStock:     false,
		}

		if err := txCourseRepo.Create(ctx, course); err != nil {
			return fmt.Errorf("failed to create course: %w", err)
		}
		if err := txProductRepo.Create(ctx, product); err != nil {
			return fmt.Errorf("failed to create course product: %w", err)
		}

		courseID = course.ID
		productID = product.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &coursemodel.CreateResponse{ID: courseID, ProductID: productID}, nil
}

// Publish sets the `InStock` field to true for a course and its associated product,
// making it available in the catalog. All of its associated course parts (if they exist)
// should be unpublished separately.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Publish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.CourseRepo.DB().Transaction(func(tx *gorm.DB) error {
		ra, err := s.CourseRepo.WithTx(tx).SetInStock(ctx, id, true)
		if err != nil {
			return fmt.Errorf("failed to publish course: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		ra, err = s.ProductRepo.WithTx(tx).SetInStockByDetailsID(ctx, id, true)
		if err != nil {
			return fmt.Errorf("failed to publish course product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil
	})
}

// Unpublish sets the `InStock` field to false for a course, its associated course parts
// and its associated product, archiving it from the catalog.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Unpublish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.CourseRepo.DB().Transaction(func(tx *gorm.DB) error {
		ra, err := s.CourseRepo.WithTx(tx).SetInStock(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to unpublish course: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		ra, err = s.ProductRepo.WithTx(tx).SetInStockByDetailsID(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to unpublish course product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		if _, err := s.PartRepo.WithTx(tx).SetPublishedByCourseID(ctx, id, false); err != nil {
			return fmt.Errorf("failed to unpublish course parts: %w", err)
		}
		return nil
	})
}

// Update performs a partial update of a course and its related product.
// The request should contain the course's ID and the fields to be updated.
// At least one field must be provided for an update to occur.
//
// Returns a map containing the fields that were actually changed, nested under "course" and "product" keys.
// Example: `{"course": {"name": "new name"}, "product": {"price": 99.99}}`
// Returns an error if the request payload is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Update(ctx context.Context, req *coursemodel.UpdateRequest) (map[string]any, error) {
	updates := make(map[string]any)
	err := s.CourseRepo.DB().Transaction(func(tx *gorm.DB) error {
		txCourseRepo := s.CourseRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		if err := req.Validate(); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
		}

		course, err := txCourseRepo.Get(ctx, req.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to retrieve course: %w", err)
		}
		product, err := txProductRepo.GetByDetailsID(ctx, course.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to retrieve product for course: %w", err)
		}

		courseUpdates := make(map[string]any)
		productUpdates := make(map[string]any)

		if req.Name != nil && *req.Name != course.Name {
			// Pass the actual value instead of pointer to it
			courseUpdates["name"] = *req.Name
		}
		if req.ShortDescription != nil && *req.ShortDescription != course.ShortDescription {
			courseUpdates["short_description"] = *req.ShortDescription
		}
		if req.LongDescription != nil && *req.LongDescription != course.LongDescription {
			courseUpdates["long_description"] = *req.LongDescription
		}
		if req.AccessDuration != nil && *req.AccessDuration != course.AccessDuration {
			courseUpdates["access_duration"] = *req.AccessDuration
		}
		if req.Price != nil && *req.Price != product.Price {
			productUpdates["price"] = *req.Price
		}
		if req.Topic != nil && *req.Topic != course.Topic {
			courseUpdates["topic"] = *req.Topic
		}
		if len(req.Tags) > 0 {
			courseUpdates["tags"] = req.Tags
		}

		if len(productUpdates) > 0 {
			productUpdates["updated_at"] = time.Now()
			if _, err := txProductRepo.Update(ctx, product, productUpdates); err != nil {
				return fmt.Errorf("failed to update course: %w", err)
			}
		}
		if len(courseUpdates) > 0 {
			courseUpdates["updated_at"] = time.Now()
			if _, err := txCourseRepo.Update(ctx, course, courseUpdates); err != nil {
				return fmt.Errorf("failed to update course product: %w", err)
			}
		}
		updates["course"] = courseUpdates
		updates["product"] = productUpdates
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updates, nil
}

// AddImage adds a new image to a course. It's called by media-service-go upon successful image upload.
// It uses courseOwnerRepoAdapter to call [imageservice.AddImage] and add an image to the course.
//
// Returns an error if:
//   - The request payload is invalid ([imageservice.ErrInvalidArgument]).
//   - The course (owner) is not found ([imageservice.ErrOwnerNotFound]).
//   - The image limit (5) is exceeded ([imageservice.ErrImageLimitExceeded]).
//   - A database/internal error occurs.
//
// Deprecated: use [image_router.AddImage] instead, it performs adding of an image for all services (generic method).
func (s *service) AddImage(ctx context.Context, req *imagemodel.AddRequest) error {
	ownerRepoAdapter := NewOwnerRepoAdapter(s.CourseRepo)
	return s.ImageSvc.AddImage(ctx, req, ownerRepoAdapter)
}

// DeleteImage removes an image from a course. It's called by media-service-go upon successful image deletion.
// It uses courseOwnerRepoAdapter to call [imageservice.DeleteImage] and delete an image from the course.
//
// Returns an error if:
//   - The request payload is invalid ([imageservice.ErrInvalidArgument]).
//   - The course (owner) is not found ([imageservice.ErrOwnerNotFound]).
//   - The image is not found on course (owner) ([imageservice.ErrImageNotFoundOnOwner]).
//   - A database/internal error occurs.
//
// Deprecated: use [image_router.DeleteImage] instead, it handles image deletion for all services (generic method).
func (s *service) DeleteImage(ctx context.Context, req *imagemodel.DeleteRequest) error {
	ownerRepoAdapter := NewOwnerRepoAdapter(s.CourseRepo)
	return s.ImageSvc.DeleteImage(ctx, req, ownerRepoAdapter)
}

// AddImageBatch adds an image for a batch of courses. It uses courseOwnerRepoAdapter
// to call [imageservice.AddImageBatch] and append images to the course. It's called by media-service-go
// upon successfull context change.
//
// It returns a number of affected courses.
// Returns an error if the request is invalid ([imageservice.ErrInvalidArgument]), no courses (owners) are not found ([imageservice.ErrOwnersNotFound])
// or a database/internal error occurs.\
//
// Deprecated: use [image_router.AddImageBatch] instead, it handles image batch add for all services (generic method).
func (s *service) AddImageBatch(ctx context.Context, req *imagemodel.AddBatchRequest) (int, error) {
	ownerRepoAdapter := NewOwnerRepoAdapter(s.CourseRepo)
	return s.ImageSvc.AddImageBatch(ctx, req, ownerRepoAdapter)
}

// DeleteImageBatch removes an image from a batch of courses. It uses courseOwnerRepoAdapter
// to call [imageservice.DeleteImageBatch] and append images to the course.
//
// It returns a number of affected courses.
// Returns an error if the request is invalid ([imageservice.ErrInvalidArgument]), no courses (owners) are not found ([imageservice.ErrOwnersNotFound]),
// no associations were found ([imageservice.ErrAssociationsNotFound]) or a database/internal error occurs.
//
// Deprecated: use [image_router.DeleteImageBatch] instead, it handles image batch delete for all services (generic method).
func (s *service) DeleteImageBatch(ctx context.Context, req *imagemodel.DeleteBatchRequst) (int, error) {
	ownerRepoAdapter := NewOwnerRepoAdapter(s.CourseRepo)
	return s.ImageSvc.DeleteImageBatch(ctx, req, ownerRepoAdapter)
}

// Delete performs a soft-delete of a course, its associated course parts
// and its associated product record.
// It also unpublishes all records, meaning they must be manually published again after restoration.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.CourseRepo.DB().Transaction(func(tx *gorm.DB) error {
		txCourseRepo := s.CourseRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		txPartRepo := s.PartRepo.WithTx(tx)

		// Check if the record exists first (including unpublished, but not soft-deleted)
		if _, err := txCourseRepo.GetWithUnpublished(ctx, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to retrieve course: %w", err)
		}

		if _, err := txCourseRepo.SetInStock(ctx, id, false); err != nil {
			return fmt.Errorf("failed to unpublish course: %w", err)
		}

		// Course may not have any parts
		if _, err := txPartRepo.SetPublishedByCourseID(ctx, id, false); err != nil {
			return fmt.Errorf("failed to unpublish course parts: %w", err)
		}

		ra, err := txProductRepo.SetInStockByDetailsID(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to unpublish course product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%s: %w", ErrNotFound, err)
		}

		// Delete all instances
		if _, err = txCourseRepo.Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete course: %w", err)
		}

		if _, err = txProductRepo.DeleteByDetailsID(ctx, id); err != nil {
			return fmt.Errorf("failed to delete course product: %w", err)
		}

		if _, err = txPartRepo.DeleteByCourseID(ctx, id); err != nil {
			return fmt.Errorf("failed to delete course parts: %w", err)
		}
		return nil
	})
}

// DeletePermanent performs a complete delete of a course, its associated course parts
// and its associated product record.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) DeletePermanent(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.CourseRepo.DB().Transaction(func(tx *gorm.DB) error {
		txCourseRepo := s.CourseRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		txPartRepo := s.PartRepo.WithTx(tx)

		ra, err := txCourseRepo.DeletePermanent(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete course: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}

		ra, err = txProductRepo.DeletePermanentByDetailsID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete course product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		// Course may not have any parts
		if _, err = txPartRepo.DeletePermanentByCourseID(ctx, id); err != nil {
			return fmt.Errorf("failed to delete course parts: %w", err)
		}
		return nil
	})
}

// Restore performs a restore of a course, its associated course parts
// and its related product record.
// Course record, its associated course part records and its related product record
// are not being published. This should be done manually.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Restore(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.CourseRepo.DB().Transaction(func(tx *gorm.DB) error {
		txCourseRepo := s.CourseRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		txPartRepo := s.PartRepo.WithTx(tx)

		ra, err := txCourseRepo.Restore(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to restore course: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}

		ra, err = txProductRepo.RestoreByDetailsID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to restore course product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}

		// Course may not have any parts
		if _, err := txPartRepo.RestoreByCourseID(ctx, id); err != nil {
			return fmt.Errorf("failed to restore course parts: %w", err)
		}
		return nil
	})
}
