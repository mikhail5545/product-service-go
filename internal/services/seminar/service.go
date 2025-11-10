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

// Package seminar provides service-layer business logic for seminars.
package seminar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	productrepo "github.com/mikhail5545/product-service-go/internal/database/product"
	seminarrepo "github.com/mikhail5545/product-service-go/internal/database/seminar"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	productmodel "github.com/mikhail5545/product-service-go/internal/models/product"
	seminarmodel "github.com/mikhail5545/product-service-go/internal/models/seminar"
	imageservice "github.com/mikhail5545/product-service-go/internal/services/image_manager"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/services/seminar_mock/service_mock.go -package=seminar_mock . Service

// Service provides service-layer business logic for seminar models.
type Service interface {
	// Get retrieves a single published and not soft-deleted seminar record from the database,
	// along with all of its associated products details (prices and product IDs).
	//
	// Returns a SeminarDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	Get(ctx context.Context, id string) (*seminarmodel.SeminarDetails, error)
	// GetWithDeleted retrieves a single seminar record from the database, including soft-deleted ones,
	// along with all of its associated products details.
	//
	// Returns a SeminarDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	GetWithDeleted(ctx context.Context, id string) (*seminarmodel.SeminarDetails, error)
	// GetWithUnpublished retrieves a single seminar record from the database, including unpublished ones (but not soft-deleted),
	// along with all of its associated products details.
	//
	// Returns a SeminarDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	GetWithUnpublished(ctx context.Context, id string) (*seminarmodel.SeminarDetails, error)
	// List retrieves a paginated list of all published and not soft-deleted seminar records.
	// Each record is returned with its associated products details.
	// It will skip seminars with missing product IDs or with incomplete product data from
	// the database.
	//
	// Returns a slice of SeminarDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs.
	List(ctx context.Context, limit, offset int) ([]seminarmodel.SeminarDetails, int64, error)
	// ListDeleted retrieves a paginated list of all soft-deleted seminar records.
	// Each record is returned with its associated products details.
	// It will skip seminars with missing product IDs or with incomplete product data from
	// the database.
	//
	// Returns a slice of SeminarDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs.
	ListDeleted(ctx context.Context, limit, offset int) ([]seminarmodel.SeminarDetails, int64, error)
	// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) seminar records.
	// Each record is returned with its associated products details.
	// It will skip seminars with missing product IDs or with incomplete product data from
	// the database.
	//
	// Returns a slice of SeminarDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs.
	ListUnpublished(ctx context.Context, limit, offset int) ([]seminarmodel.SeminarDetails, int64, error)
	// Create creates a new Seminar record and all of its associated Product records in the database.
	// It validates the request payload to ensure all required fields are present.
	// The seminar and all of the associated products are created in an unpublished state (`InStock: false`).
	//
	// Returns a CreateResponse containing the newly created SeminarID, ReservationProductID, EarlyProductID,
	// LateProductID, EarlySurchargeProductID, LateSurchargeProductID.
	// Returns an error if the request payload is invalid (ErrInvalidArgument) or a database/internal error occurs.
	Create(ctx context.Context, req *seminarmodel.CreateRequest) (*seminarmodel.CreateResponse, error)
	// Publish sets the `InStock` field to true for a seminar and all of its associated products,
	// making it available in the catalog.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Publish(ctx context.Context, id string) error
	// Unpublish sets the `InStock` field to false for a seminar and all of its associated products,
	// archiving it from the catalog.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Unpublish(ctx context.Context, id string) error
	// Update performs a partial update of a seminar and all of its related products.
	// The request should contain the seminar's ID and the fields to be updated.
	// At least one field must be provided for an update to occur.
	//
	// Returns a map containing the fields that were actually changed, nested under "seminar", "reservation_product",
	// "early_product", "late_product", "early_surcharge_product", "late_surcharge_product" keys.
	// Example: `{"seminar": {"name": "new name"}, "early_product": {"price": 99.99}}`
	// Returns an error if the request payload is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Update(ctx context.Context, req *seminarmodel.UpdateRequest) (map[string]any, error)
	// Delete performs a soft-delete of a seminar and all of its related product records.
	// It also unpublishes all records, meaning they must be manually published again after restoration.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Delete(ctx context.Context, id string) error
	// DeletePermanent performs a complete delete of a seminar and its related product records.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	DeletePermanent(ctx context.Context, id string) error
	// Restore performs a restore of a seminar and its related product records.
	// Seminar and its related product records are not being published. This should be
	// done manually.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Restore(ctx context.Context, id string) error
	// AddImage adds a new image to a seminar. It's called by media-service-go upon successful image upload.
	// It uses seminarOwnerRepoAdapter to call [imageservice.AddImage] and add an image to the seminar.
	//
	// Returns an error if:
	// 	- The request payload is invalid ([imageservice.ErrInvalidArgument]).
	// 	- The seminar (owner) is not found ([imageservice.ErrOwnerNotFound]).
	// 	- The image limit (5) is exceeded ([imageservice.ErrImageLimitExceeded]).
	// 	- A database/internal error occurs.
	//
	// Deprecated: use generic [image.Add] instead. It handles add operations for all services.
	AddImage(ctx context.Context, req *imagemodel.AddRequest) error
	// DeleteImage removes an image from a seminar. It's called by media-service-go upon successful image deletion.
	// It uses seminarOwnerRepoAdapter to call [imageservice.DeleteImage] and delete an image from the seminar.
	//
	// Returns an error if:
	//   - The request payload is invalid ([imageservice.ErrInvalidArgument]).
	//   - The seminar (owner) is not found ([imageservice.ErrOwnerNotFound]).
	//   - The image is not found on seminar (owner) ([imageservice.ErrImageNotFoundOnOwner]).
	//   - A database/internal error occurs.
	//
	// Deprecated: use generic [image.Delete] instead. It handles delete operations for all services.
	DeleteImage(ctx context.Context, req *imagemodel.DeleteRequest) error
	// AddImageBatch adds an image for a batch of seminars. It uses seminarOwnerRepoAdapter
	// to call [imageservice.AddImageBatch] and append images to the seminar. It's called by media-service-go
	// upon successfull context change.
	//
	// It returns the number of affected seminars.
	// Returns an error if the request is invalid ([imageservice.ErrInvalidArgument]), no seminars (owners) are not found ([imageservice.ErrOwnersNotFound])
	// or a database/internal error occurs.
	//
	// Deprecated: use generic [image.AddBatch] instead. It handles batch add operations for all services.
	AddImageBatch(ctx context.Context, req *imagemodel.AddBatchRequest) (int, error)
	// DeleteImageBatch removes an image from a batch of seminars. It uses seminarOwnerRepoAdapter
	// to call [imageservice.DeleteImageBatch] and append images to the seminar.
	//
	// It returns the number of affected seminars.
	// Returns an error if the request is invalid ([imageservice.ErrInvalidArgument]), no seminars (owners) are not found ([imageservice.ErrOwnersNotFound]),
	// no associations were found ([imageservice.ErrAssociationsNotFound]) or a database/internal error occurs.
	//
	// Deprecated: use generic [image.AddBatch] instead. It handles batch add operations for all services.
	DeleteImageBatch(ctx context.Context, req *imagemodel.DeleteBatchRequst) (int, error)
}

// service provides service-layer business logic for seminar models.
// It holds [seminarrepo.Repository] and [productrepo.Repository] instances
// to perform database operations.
type service struct {
	SeminarRepo seminarrepo.Repository
	ProductRepo productrepo.Repository
	ImageSvc    imageservice.Service
}

// New creates a new service instance with provided seminar and product repositories.
func New(sr seminarrepo.Repository, pr productrepo.Repository, is imageservice.Service) Service {
	return &service{
		SeminarRepo: sr,
		ProductRepo: pr,
		ImageSvc:    is,
	}
}

// Get retrieves a single published and not soft-deleted seminar record from the database,
// along with all of its associated products details (prices and product IDs).
//
// Returns a SeminarDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Get(ctx context.Context, id string) (*seminarmodel.SeminarDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	seminar, err := s.SeminarRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve seminar: %w", err)
	}

	if seminar.ReservationProductID == nil || seminar.EarlyProductID == nil || seminar.LateProductID == nil || seminar.EarlySurchargeProductID == nil || seminar.LateSurchargeProductID == nil {
		return nil, ErrIncompleteData
	}

	productIDs := []string{
		*seminar.ReservationProductID,
		*seminar.EarlyProductID,
		*seminar.LateProductID,
		*seminar.EarlySurchargeProductID,
		*seminar.LateSurchargeProductID,
	}

	products, err := s.ProductRepo.SelectByIDs(ctx, productIDs, "price")
	if err != nil {
		return nil, fmt.Errorf("failed to get seminar products: %w", err)
	}
	if len(products) != 5 {
		return nil, ErrProductsNotFound
	}

	productMap := make(map[string]*productmodel.Product, len(products))
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	details := seminarmodel.SeminarDetails{
		Seminar:             seminar,
		ReservationPrice:    productMap[*seminar.ReservationProductID].Price,
		EarlyPrice:          productMap[*seminar.EarlyProductID].Price,
		LatePrice:           productMap[*seminar.LateProductID].Price,
		EarlySurchargePrice: productMap[*seminar.EarlySurchargeProductID].Price,
		LateSurchargePrice:  productMap[*seminar.LateSurchargeProductID].Price,
	}
	details.Current()

	return &details, nil
}

// GetWithDeleted retrieves a single seminar record from the database, including soft-deleted ones,
// along with all of its associated products details.
//
// Returns a SeminarDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) GetWithDeleted(ctx context.Context, id string) (*seminarmodel.SeminarDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	seminar, err := s.SeminarRepo.GetWithDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve seminar: %w", err)
	}

	if seminar.ReservationProductID == nil || seminar.EarlyProductID == nil || seminar.LateProductID == nil || seminar.EarlySurchargeProductID == nil || seminar.LateSurchargeProductID == nil {
		return nil, ErrIncompleteData
	}

	productIDs := []string{
		*seminar.ReservationProductID,
		*seminar.EarlyProductID,
		*seminar.LateProductID,
		*seminar.EarlySurchargeProductID,
		*seminar.LateSurchargeProductID,
	}

	products, err := s.ProductRepo.SelectWithDeletedByIDs(ctx, productIDs, "price")
	if err != nil {
		return nil, fmt.Errorf("failed to get seminar products: %w", err)
	}
	if len(products) != 5 {
		return nil, ErrProductsNotFound
	}

	productMap := make(map[string]*productmodel.Product, len(products))
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	details := seminarmodel.SeminarDetails{
		Seminar:             seminar,
		ReservationPrice:    productMap[*seminar.ReservationProductID].Price,
		EarlyPrice:          productMap[*seminar.EarlyProductID].Price,
		LatePrice:           productMap[*seminar.LateProductID].Price,
		EarlySurchargePrice: productMap[*seminar.EarlySurchargeProductID].Price,
		LateSurchargePrice:  productMap[*seminar.LateSurchargeProductID].Price,
	}
	details.Current()

	return &details, nil
}

// GetWithUnpublished retrieves a single seminar record from the database, including unpublished ones (but not soft-deleted),
// along with all of its associated products details.
//
// Returns a SeminarDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) GetWithUnpublished(ctx context.Context, id string) (*seminarmodel.SeminarDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	seminar, err := s.SeminarRepo.GetWithUnpublished(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve seminar: %w", err)
	}

	if seminar.ReservationProductID == nil || seminar.EarlyProductID == nil || seminar.LateProductID == nil || seminar.EarlySurchargeProductID == nil || seminar.LateSurchargeProductID == nil {
		return nil, ErrIncompleteData
	}

	productIDs := []string{
		*seminar.ReservationProductID,
		*seminar.EarlyProductID,
		*seminar.LateProductID,
		*seminar.EarlySurchargeProductID,
		*seminar.LateSurchargeProductID,
	}

	products, err := s.ProductRepo.SelectWithUnpublishedByIDs(ctx, productIDs, "price")
	if err != nil {
		return nil, fmt.Errorf("failed to get seminar products: %w", err)
	}
	if len(products) != 5 {
		return nil, ErrProductsNotFound
	}

	productMap := make(map[string]*productmodel.Product, len(products))
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	details := seminarmodel.SeminarDetails{
		Seminar:             seminar,
		ReservationPrice:    productMap[*seminar.ReservationProductID].Price,
		EarlyPrice:          productMap[*seminar.EarlyProductID].Price,
		LatePrice:           productMap[*seminar.LateProductID].Price,
		EarlySurchargePrice: productMap[*seminar.EarlySurchargeProductID].Price,
		LateSurchargePrice:  productMap[*seminar.LateSurchargeProductID].Price,
	}
	details.Current()

	return &details, nil
}

// safeGetPrice retrieves a product's price from the map, returning 0 if the ID pointer is nil or the product is not found.
func safeGetPrice(productMap map[string]*productmodel.Product, id *string) float32 {
	if id == nil {
		return 0
	}
	if p, ok := productMap[*id]; ok {
		return p.Price
	}
	return 0
}

// hasMissingProducts checks if any of the required product IDs are missing from the product map.
func hasMissingProducts(productMap map[string]*productmodel.Product, seminar *seminarmodel.Seminar) bool {
	_, ok1 := productMap[*seminar.ReservationProductID]
	_, ok2 := productMap[*seminar.EarlyProductID]
	_, ok3 := productMap[*seminar.LateProductID]
	_, ok4 := productMap[*seminar.EarlySurchargeProductID]
	_, ok5 := productMap[*seminar.LateSurchargeProductID]
	return !ok1 || !ok2 || !ok3 || !ok4 || !ok5
}

// List retrieves a paginated list of all published and not soft-deleted seminar records.
// Each record is returned with its associated products details.
// It will skip seminars with missing product IDs or with incomplete product data from
// the database.
//
// Returns a slice of SeminarDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs.
func (s *service) List(ctx context.Context, limit, offset int) ([]seminarmodel.SeminarDetails, int64, error) {
	seminars, err := s.SeminarRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve seminars: %w", err)
	}

	// Collect all product IDs from all seminars
	var productIDs []string
	for _, seminar := range seminars {
		if seminar.ReservationProductID != nil {
			productIDs = append(productIDs, *seminar.ReservationProductID)
		}
		if seminar.EarlyProductID != nil {
			productIDs = append(productIDs, *seminar.EarlyProductID)
		}
		if seminar.LateProductID != nil {
			productIDs = append(productIDs, *seminar.LateProductID)
		}
		if seminar.EarlySurchargeProductID != nil {
			productIDs = append(productIDs, *seminar.EarlySurchargeProductID)
		}
		if seminar.LateSurchargeProductID != nil {
			productIDs = append(productIDs, *seminar.LateSurchargeProductID)
		}
	}

	// Fetch all products in a single query
	products, err := s.ProductRepo.SelectByIDs(ctx, productIDs, "price")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}

	// Create a map for quick product lookup by ID
	productMap := make(map[string]*productmodel.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = &p
	}

	var allDetails []seminarmodel.SeminarDetails
	for _, seminar := range seminars {
		// Skip seminars that have missing product IDs or if their products weren't found.
		if seminar.ReservationProductID == nil || seminar.EarlyProductID == nil || seminar.LateProductID == nil || seminar.EarlySurchargeProductID == nil || seminar.LateSurchargeProductID == nil || hasMissingProducts(productMap, &seminar) {
			continue
		}

		details := seminarmodel.SeminarDetails{
			Seminar:             &seminar,
			ReservationPrice:    safeGetPrice(productMap, seminar.ReservationProductID),
			EarlyPrice:          safeGetPrice(productMap, seminar.EarlyProductID),
			LatePrice:           safeGetPrice(productMap, seminar.LateProductID),
			EarlySurchargePrice: safeGetPrice(productMap, seminar.EarlySurchargeProductID),
			LateSurchargePrice:  safeGetPrice(productMap, seminar.LateSurchargeProductID),
		}
		details.Current()
		allDetails = append(allDetails, details)
	}
	total, err := s.SeminarRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count seminars: %w", err)
	}
	return allDetails, total, nil
}

// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) seminar records.
// Each record is returned with its associated products details.
// It will skip seminars with missing product IDs or with incomplete product data from
// the database.
//
// Returns a slice of SeminarDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs.
func (s *service) ListUnpublished(ctx context.Context, limit, offset int) ([]seminarmodel.SeminarDetails, int64, error) {
	seminars, err := s.SeminarRepo.ListUnpublished(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve seminars: %w", err)
	}

	// Collect all product IDs from all seminars
	var productIDs []string
	for _, seminar := range seminars {
		if seminar.ReservationProductID != nil {
			productIDs = append(productIDs, *seminar.ReservationProductID)
		}
		if seminar.EarlyProductID != nil {
			productIDs = append(productIDs, *seminar.EarlyProductID)
		}
		if seminar.LateProductID != nil {
			productIDs = append(productIDs, *seminar.LateProductID)
		}
		if seminar.EarlySurchargeProductID != nil {
			productIDs = append(productIDs, *seminar.EarlySurchargeProductID)
		}
		if seminar.LateSurchargeProductID != nil {
			productIDs = append(productIDs, *seminar.LateSurchargeProductID)
		}
	}

	// Fetch all products in a single query
	products, err := s.ProductRepo.SelectWithUnpublishedByIDs(ctx, productIDs, "price")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}

	// Create a map for quick product lookup by ID
	productMap := make(map[string]*productmodel.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = &p
	}

	var allDetails []seminarmodel.SeminarDetails
	for _, seminar := range seminars {
		// Skip seminars that have missing product IDs or if their products weren't found.
		if seminar.ReservationProductID == nil || seminar.EarlyProductID == nil || seminar.LateProductID == nil || seminar.EarlySurchargeProductID == nil || seminar.LateSurchargeProductID == nil || hasMissingProducts(productMap, &seminar) {
			continue
		}

		details := seminarmodel.SeminarDetails{
			Seminar:             &seminar,
			ReservationPrice:    safeGetPrice(productMap, seminar.ReservationProductID),
			EarlyPrice:          safeGetPrice(productMap, seminar.EarlyProductID),
			LatePrice:           safeGetPrice(productMap, seminar.LateProductID),
			EarlySurchargePrice: safeGetPrice(productMap, seminar.EarlySurchargeProductID),
			LateSurchargePrice:  safeGetPrice(productMap, seminar.LateSurchargeProductID),
		}
		details.Current()
		allDetails = append(allDetails, details)
	}
	total, err := s.SeminarRepo.CountUnpublished(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count seminars: %w", err)
	}
	return allDetails, total, nil
}

// ListDeleted retrieves a paginated list of all soft-deleted seminar records.
// Each record is returned with its associated products details.
// It will skip seminars with missing product IDs or with incomplete product data from
// the database.
//
// Returns a slice of SeminarDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs.
func (s *service) ListDeleted(ctx context.Context, limit, offset int) ([]seminarmodel.SeminarDetails, int64, error) {
	seminars, err := s.SeminarRepo.ListDeleted(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve seminars: %w", err)
	}

	// Collect all product IDs from all seminars
	var productIDs []string
	for _, seminar := range seminars {
		if seminar.ReservationProductID != nil {
			productIDs = append(productIDs, *seminar.ReservationProductID)
		}
		if seminar.EarlyProductID != nil {
			productIDs = append(productIDs, *seminar.EarlyProductID)
		}
		if seminar.LateProductID != nil {
			productIDs = append(productIDs, *seminar.LateProductID)
		}
		if seminar.EarlySurchargeProductID != nil {
			productIDs = append(productIDs, *seminar.EarlySurchargeProductID)
		}
		if seminar.LateSurchargeProductID != nil {
			productIDs = append(productIDs, *seminar.LateSurchargeProductID)
		}
	}

	// Fetch all products in a single query
	products, err := s.ProductRepo.SelectWithDeletedByIDs(ctx, productIDs, "price")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}

	// Create a map for quick product lookup by ID
	productMap := make(map[string]*productmodel.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = &p
	}

	var allDetails []seminarmodel.SeminarDetails
	for _, seminar := range seminars {
		// Skip seminars that have missing product IDs or if their products weren't found.
		if seminar.ReservationProductID == nil || seminar.EarlyProductID == nil || seminar.LateProductID == nil || seminar.EarlySurchargeProductID == nil || seminar.LateSurchargeProductID == nil || hasMissingProducts(productMap, &seminar) {
			continue
		}

		details := seminarmodel.SeminarDetails{
			Seminar:             &seminar,
			ReservationPrice:    safeGetPrice(productMap, seminar.ReservationProductID),
			EarlyPrice:          safeGetPrice(productMap, seminar.EarlyProductID),
			LatePrice:           safeGetPrice(productMap, seminar.LateProductID),
			EarlySurchargePrice: safeGetPrice(productMap, seminar.EarlySurchargeProductID),
			LateSurchargePrice:  safeGetPrice(productMap, seminar.LateSurchargeProductID),
		}
		details.Current()
		allDetails = append(allDetails, details)
	}
	total, err := s.SeminarRepo.CountDeleted(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count seminars: %w", err)
	}
	return allDetails, total, nil
}

// Create creates a new Seminar record and all of its associated Product records in the database.
// It validates the request payload to ensure all required fields are present.
// The seminar and all of the associated products are created in an unpublished state (`InStock: false`).
//
// Returns a CreateResponse containing the newly created SeminarID, ReservationProductID, EarlyProductID,
// LateProductID, EarlySurchargeProductID, LateSurchargeProductID.
// Returns an error if the request payload is invalid (ErrInvalidArgument) or a database/internal error occurs.
func (s *service) Create(ctx context.Context, req *seminarmodel.CreateRequest) (*seminarmodel.CreateResponse, error) {
	seminar := &seminarmodel.Seminar{}
	err := s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		if err := req.Validate(); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
		}

		seminar.ID = uuid.New().String()
		seminar.Name = req.Name
		seminar.ShortDescription = req.ShortDescription
		seminar.Date = req.Date
		seminar.EndingDate = req.EndingDate
		seminar.Place = req.Place
		seminar.LatePaymentDate = req.LatePaymentDate
		seminar.InStock = false

		products := []*productmodel.Product{
			{ID: uuid.New().String(), Price: req.ReservationPrice, InStock: false},
			{ID: uuid.New().String(), Price: req.EarlyPrice, InStock: false},
			{ID: uuid.New().String(), Price: req.LatePrice, InStock: false},
			{ID: uuid.New().String(), Price: req.EarlySurchargePrice, InStock: false},
			{ID: uuid.New().String(), Price: req.LateSurchargePrice, InStock: false},
		}

		for _, p := range products {
			p.DetailsID = seminar.ID
			p.DetailsType = "seminar"
		}

		if err := txProductRepo.CreateBatch(ctx, products...); err != nil {
			return fmt.Errorf("failed to create seminar products: %w", err)
		}

		productIDMap := map[float32]*string{
			req.ReservationPrice:    &products[0].ID,
			req.EarlyPrice:          &products[1].ID,
			req.LatePrice:           &products[2].ID,
			req.EarlySurchargePrice: &products[3].ID,
			req.LateSurchargePrice:  &products[4].ID,
		}

		seminar.ReservationProductID = productIDMap[req.ReservationPrice]
		seminar.EarlyProductID = productIDMap[req.EarlyPrice]
		seminar.LateProductID = productIDMap[req.LatePrice]
		seminar.EarlySurchargeProductID = productIDMap[req.EarlySurchargePrice]
		seminar.LateSurchargeProductID = productIDMap[req.LateSurchargePrice]

		if err := txSeminarRepo.Create(ctx, seminar); err != nil {
			return fmt.Errorf("failed to create seminar: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &seminarmodel.CreateResponse{
		ID:                      seminar.ID,
		ReservationProductID:    *seminar.ReservationProductID,
		EarlyProductID:          *seminar.EarlyProductID,
		LateProductID:           *seminar.LateProductID,
		EarlySurchargeProductID: *seminar.EarlySurchargeProductID,
		LateSurchargeProductID:  *seminar.LateSurchargeProductID,
	}, nil
}

// Publish sets the `InStock` field to true for a seminar and all of its associated products,
// making it available in the catalog.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Publish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: invalid seminar ID: %w", ErrInvalidArgument, err)
	}
	return s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txSeminarRepo.SetInStock(ctx, id, true)
		if err != nil {
			return fmt.Errorf("failed to publish seminar: %w", err)
		} else if ra == 0 {
			return ErrNotFound
		}
		ra, err = txProductRepo.SetInStockByDetailsID(ctx, id, true)
		if err != nil {
			return fmt.Errorf("failed to publish seminar products: %w", err)
		} else if ra != 5 {
			// This indicates a data integrity issue.
			return fmt.Errorf("failed to publish all 5 seminar products, only %d were updated", ra)
		}
		return nil
	})
}

// Unpublish sets the `InStock` field to false for a seminar and all of its associated products,
// archiving it from the catalog.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Unpublish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: invalid seminar ID: %w", ErrInvalidArgument, err)
	}
	return s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txSeminarRepo.SetInStock(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to unpublish seminar: %w", err)
		} else if ra == 0 {
			return ErrNotFound
		}
		ra, err = txProductRepo.SetInStockByDetailsID(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to unpublish seminar products: %w", err)
		} else if ra != 5 {
			// This indicates a data integrity issue.
			return fmt.Errorf("failed to unpublish all 5 seminar products, only %d were updated", ra)
		}
		return nil
	})
}

// Update performs a partial update of a seminar and all of its related products.
// The request should contain the seminar's ID and the fields to be updated.
// At least one field must be provided for an update to occur.
//
// Returns a map containing the fields that were actually changed, nested under "seminar", "reservation_product",
// "early_product", "late_product", "early_surcharge_product", "late_surcharge_product" keys.
// Example: `{"seminar": {"name": "new name"}, "early_product": {"price": 99.99}}`
// Returns an error if the request payload is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Update(ctx context.Context, req *seminarmodel.UpdateRequest) (map[string]any, error) {
	allUpdates := make(map[string]any)
	err := s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		if err := req.Validate(); err != nil {
			validationMsg, _ := json.Marshal(err)
			return fmt.Errorf("%w: %s", ErrInvalidArgument, string(validationMsg))
		}

		seminar, err := txSeminarRepo.Get(ctx, req.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to find seminar: %w", err)
		}

		if seminar.ReservationProductID == nil || seminar.EarlyProductID == nil || seminar.LateProductID == nil || seminar.EarlySurchargeProductID == nil || seminar.LateSurchargeProductID == nil {
			return ErrIncompleteData
		}

		productIDs := []string{
			*seminar.ReservationProductID,
			*seminar.EarlyProductID,
			*seminar.LateProductID,
			*seminar.EarlySurchargeProductID,
			*seminar.LateSurchargeProductID,
		}

		products, err := txProductRepo.SelectByIDs(ctx, productIDs, "id", "price", "details_id")
		if err != nil {
			return fmt.Errorf("failed to get seminar products: %w", err)
		}
		if len(products) != 5 {
			return ErrProductsNotFound
		}

		productMap := make(map[string]*productmodel.Product, len(products))
		for i := range products {
			productMap[products[i].ID] = &products[i]
		}

		seminarUpdates := make(map[string]any)
		if req.Name != nil && *req.Name != seminar.Name {
			seminarUpdates["name"] = *req.Name
		}
		if req.ShortDescription != nil && *req.ShortDescription != seminar.ShortDescription {
			seminarUpdates["short_description"] = *req.ShortDescription
		}
		if req.Place != nil && *req.Place != seminar.Place {
			seminarUpdates["place"] = *req.Place
		}
		if req.Date != nil && !req.Date.IsZero() && !req.Date.Equal(seminar.Date) {
			seminarUpdates["date"] = *req.Date
		}
		if req.EndingDate != nil && !req.EndingDate.IsZero() && !req.EndingDate.Equal(seminar.EndingDate) {
			seminarUpdates["ending_date"] = *req.EndingDate
		}
		if req.LatePaymentDate != nil && !req.LatePaymentDate.IsZero() && !req.LatePaymentDate.Equal(seminar.LatePaymentDate) {
			seminarUpdates["late_payment_date"] = *req.LatePaymentDate
		}
		if req.LongDescription != nil && *req.LongDescription != seminar.LongDescription {
			seminarUpdates["long_description"] = *req.LongDescription
		}
		if len(req.Tags) > 0 {
			seminarUpdates["tags"] = req.Tags
		}

		// helper function to update products
		updateProduct := func(
			reqPrice *float32,
			currentProduct *productmodel.Product,
		) (map[string]any, error) {
			if currentProduct == nil {
				// This case should be prevented by earlier checks, but as a safeguard:
				return nil, fmt.Errorf("%w: product to update not found", ErrNotFound)
			}

			productUpdates := make(map[string]any)
			if reqPrice != nil && *reqPrice != currentProduct.Price {
				productUpdates["price"] = *reqPrice
			}

			if len(productUpdates) > 0 {
				if _, err := txProductRepo.Update(ctx, currentProduct, productUpdates); err != nil {
					return nil, err
				}
			}
			return productUpdates, nil
		}

		// Check if seminar has missing products
		if hasMissingProducts(productMap, seminar) {
			return ErrProductsNotFound
		}

		// productReq represents product type as key and struct of new product price, product retrieved from the database
		productReq := map[string]struct {
			price   *float32
			product *productmodel.Product
		}{
			"reservation_product": {
				price:   req.ReservationPrice,
				product: productMap[*seminar.ReservationProductID],
			},
			"early_product": {
				price:   req.EarlyPrice,
				product: productMap[*seminar.EarlyProductID],
			},
			"late_product": {
				price:   req.LatePrice,
				product: productMap[*seminar.LateProductID],
			},
			"early_surcharge_product": {
				price:   req.EarlySurchargePrice,
				product: productMap[*seminar.EarlySurchargeProductID],
			},
			"late_surcharge_product": {
				price:   req.LateSurchargePrice,
				product: productMap[*seminar.LateSurchargeProductID],
			},
		}

		// update products
		for key, p := range productReq {
			pu, err := updateProduct(p.price, p.product)
			if err != nil {
				return fmt.Errorf("failed to update %s: %w", key, err)
			}
			if len(pu) > 0 {
				allUpdates[key] = pu
			}
		}

		if len(seminarUpdates) > 0 {
			if _, err := txSeminarRepo.Update(ctx, seminar, seminarUpdates); err != nil {
				return fmt.Errorf("failed to update seminar: %w", err)
			}
			allUpdates["seminar"] = seminarUpdates
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return allUpdates, nil
}

// AddImage adds a new image to a seminar. It's called by media-service-go upon successful image upload.
// It uses seminarOwnerRepoAdapter to call [imageservice.AddImage] and add an image to the seminar.
//
// Returns an error if:
//   - The request payload is invalid ([imageservice.ErrInvalidArgument]).
//   - The seminar (owner) is not found ([imageservice.ErrOwnerNotFound]).
//   - The image limit (5) is exceeded ([imageservice.ErrImageLimitExceeded]).
//   - A database/internal error occurs.
//
// Deprecated: use generic [image.Add] instead. It handles add operations for all services.
func (s *service) AddImage(ctx context.Context, req *imagemodel.AddRequest) error {
	ownerRepoAdapter := NewSeminarOwnerRepoAdapter(s.SeminarRepo)
	return s.ImageSvc.AddImage(ctx, req, ownerRepoAdapter)
}

// DeleteImage removes an image from a seminar. It's called by media-service-go upon successful image deletion.
// It uses seminarOwnerRepoAdapter to call [imageservice.DeleteImage] and delete an image from the seminar.
//
// Returns an error if:
//   - The request payload is invalid ([imageservice.ErrInvalidArgument]).
//   - The seminar (owner) is not found ([imageservice.ErrOwnerNotFound]).
//   - The image is not found on seminar (owner) ([imageservice.ErrImageNotFoundOnOwner]).
//   - A database/internal error occurs.
//
// Deprecated: use generic [image.Delete] instead. It handles delete operations for all services.
func (s *service) DeleteImage(ctx context.Context, req *imagemodel.DeleteRequest) error {
	ownerRepoAdapter := NewSeminarOwnerRepoAdapter(s.SeminarRepo)
	return s.ImageSvc.DeleteImage(ctx, req, ownerRepoAdapter)
}

// AddImageBatch adds an image for a batch of seminars. It uses seminarOwnerRepoAdapter
// to call [imageservice.AddImageBatch] and append images to the seminar. It's called by media-service-go
// upon successfull context change.
//
// It returns the number of affected seminars.
// Returns an error if the request is invalid ([imageservice.ErrInvalidArgument]), no seminars (owners) are not found ([imageservice.ErrOwnersNotFound])
// or a database/internal error occurs.
//
// Deprecated: use generic [image.AddBatch] instead. It handles batch add operations for all services.
func (s *service) AddImageBatch(ctx context.Context, req *imagemodel.AddBatchRequest) (int, error) {
	ownerRepoAdapter := NewSeminarOwnerRepoAdapter(s.SeminarRepo)
	return s.ImageSvc.AddImageBatch(ctx, req, ownerRepoAdapter)
}

// DeleteImageBatch removes an image from a batch of seminars. It uses seminarOwnerRepoAdapter
// to call [imageservice.DeleteImageBatch] and append images to the seminar.
//
// It returns the number of affected seminars.
// Returns an error if the request is invalid ([imageservice.ErrInvalidArgument]), no seminars (owners) are not found ([imageservice.ErrOwnersNotFound]),
// no associations were found ([imageservice.ErrAssociationsNotFound]) or a database/internal error occurs.
//
// Deprecated: use generic [image.DeleteBatch] instead. It handles batch delete operations for all services.
func (s *service) DeleteImageBatch(ctx context.Context, req *imagemodel.DeleteBatchRequst) (int, error) {
	ownerRepoAdapter := NewSeminarOwnerRepoAdapter(s.SeminarRepo)
	return s.ImageSvc.DeleteImageBatch(ctx, req, ownerRepoAdapter)
}

// Delete performs a soft-delete of a seminar and all of its related product records.
// It also unpublishes all records, meaning they must be manually published again after restoration.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: invalid seminar ID: %w", ErrInvalidArgument, err)
	}
	return s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		// Check if seminar exists
		if _, err := txSeminarRepo.GetWithUnpublished(ctx, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to get seminar: %w", err)
		}

		// Unpublish all instances
		if _, err := txSeminarRepo.SetInStock(ctx, id, false); err != nil {
			return fmt.Errorf("failed to unpublish seminar: %w", err)
		}
		ra, err := txProductRepo.SetInStockByDetailsID(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to unpublish seminar products: %w", err)
		} else if ra != 5 {
			return fmt.Errorf("failed to unpublish all 5 seminar products, only %d were updated", ra)
		}

		// Delete all instances
		if _, err = txSeminarRepo.Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete seminar: %w", err)
		}
		if _, err = txProductRepo.DeleteByDetailsID(ctx, id); err != nil {
			return fmt.Errorf("failed to delete seminar products: %w", err)
		}
		return nil
	})
}

// DeletePermanent performs a complete delete of a seminar and its related product records.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) DeletePermanent(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: invalid seminar ID: %w", ErrInvalidArgument, err)
	}
	return s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		ra, err := txSeminarRepo.DeletePermanent(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete seminar: %w", err)
		} else if ra == 0 {
			return ErrNotFound
		}

		ra, err = txProductRepo.DeletePermanentByDetailsID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete seminar products: %w", err)
		} else if ra != 5 {
			return fmt.Errorf("failed to delete all 5 seminar products, only %d were updated", ra)
		}
		return nil
	})
}

// Restore performs a restore of a seminar and its related product records.
// Seminar and its related product records are not being published. This should be
// done manually.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Restore(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: invalid seminar ID: %w", ErrInvalidArgument, err)
	}
	return s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txSeminarRepo.Restore(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to restore seminar: %w", err)
		} else if ra == 0 {
			return ErrNotFound
		}
		ra, err = txProductRepo.RestoreByDetailsID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to restore seminar products: %w", err)
		} else if ra != 5 {
			return fmt.Errorf("failed to restore all 5 seminar products, only %d were updated", ra)
		}
		return nil
	})
}
