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
	"net/http"

	"github.com/google/uuid"
	productrepo "github.com/mikhail5545/product-service-go/internal/database/product"
	seminarrepo "github.com/mikhail5545/product-service-go/internal/database/seminar"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	productmodel "github.com/mikhail5545/product-service-go/internal/models/product"
	seminarmodel "github.com/mikhail5545/product-service-go/internal/models/seminar"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/services/seminar_mock/service_mock.go -package=seminar_mock . Service

// Service provides service-layer business logic for seminar models.
type Service interface {
	// Get retrieves a single published and not soft-deleted seminar record from the database,
	// along with all of its associated products details (prices and product IDs).
	//
	// Returns a SeminarDetails struct containing the combined information.
	// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Get(ctx context.Context, id string) (*seminarmodel.SeminarDetails, error)
	// GetWithDeleted retrieves a single seminar record from the database, including soft-deleted ones,
	// along with all of its associated products details.
	//
	// Returns a SeminarDetails struct containing the combined information.
	// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	GetWithDeleted(ctx context.Context, id string) (*seminarmodel.SeminarDetails, error)
	// GetWithUnpublished retrieves a single seminar record from the database, including unpublished ones (but not soft-deleted),
	// along with all of its associated products details.
	//
	// Returns a SeminarDetails struct containing the combined information.
	// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	GetWithUnpublished(ctx context.Context, id string) (*seminarmodel.SeminarDetails, error)
	// List retrieves a paginated list of all published and not soft-deleted seminar records.
	// Each record is returned with its associated products details.
	//
	// Returns a slice of SeminarDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs (http.StatusInternalServerError).
	List(ctx context.Context, limit, offset int) ([]seminarmodel.SeminarDetails, int64, error)
	// ListDeleted retrieves a paginated list of all soft-deleted seminar records.
	// Each record is returned with its associated products details.
	//
	// Returns a slice of SeminarDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs (http.StatusInternalServerError).
	ListDeleted(ctx context.Context, limit, offset int) ([]seminarmodel.SeminarDetails, int64, error)
	// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) seminar records.
	// Each record is returned with its associated products details.
	//
	// Returns a slice of SeminarDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs (http.StatusInternalServerError).
	ListUnpublished(ctx context.Context, limit, offset int) ([]seminarmodel.SeminarDetails, int64, error)
	// Create creates a new Seminar record and all of its associated Product records in the database.
	// It validates the request payload to ensure all required fields are present.
	// The seminar and all of the associated products are created in an unpublished state (`InStock: false`).
	//
	// Returns a CreateResponse containing the newly created SeminarID, ReservationProductID, EarlyProductID,
	// LateProductID, EarlySurchargeProductID, LateSurchargeProductID.
	// Returns an error if the request payload is invalid (http.StatusBadRequest) or a database/internal error occurs (http.StatusInternalServerError).
	Create(ctx context.Context, req *seminarmodel.CreateRequest) (*seminarmodel.CreateResponse, error)
	// Publish sets the `InStock` field to true for a seminar and all of its associated products,
	// making it available in the catalog.
	//
	// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Publish(ctx context.Context, id string) error
	// Unpublish sets the `InStock` field to false for a seminar and all of its associated products,
	// archiving it from the catalog.
	//
	// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Unpublish(ctx context.Context, id string) error
	// Update performs a partial update of a seminar and all of its related products.
	// The request should contain the seminar's ID and the fields to be updated.
	// At least one field must be provided for an update to occur.
	//
	// Returns a map containing the fields that were actually changed, nested under "seminar", "reservation_product",
	// "early_product", "late_product", "early_surcharge_product", "late_surcharge_product" keys.
	// Example: `{"seminar": {"name": "new name"}, "early_product": {"price": 99.99}}`
	// Returns an error if the request payload is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Update(ctx context.Context, req *seminarmodel.UpdateRequest) (map[string]any, error)
	// Delete performs a soft-delete of a seminar and all of its related product records.
	// It also unpublishes all records, meaning they must be manually published again after restoration.
	//
	// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Delete(ctx context.Context, id string) error
	// DeletePermanent performs a complete delete of a seminar and its related product records.
	//
	// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	DeletePermanent(ctx context.Context, id string) error
	// Restore performs a restore of a seminar and its related product records.
	// Seminar and its related product records are not being published. This should be
	// done manually.
	//
	// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
	// or a database/internal error occurs (http.StatusInternalServerError).
	Restore(ctx context.Context, id string) error
	// AddImage adds a new image to a seminar. It's called by media-service-go upon successful image upload.
	// The function validates the request, checks the image limit, and appends the new image information.
	//
	// Returns an AddResponse with the MediaServiceID on success.
	// Returns an error if:
	// - The request payload is invalid (http.StatusBadRequest).
	// - The seminar (owner) is not found (http.StatusNotFound).
	// - The image limit (5) is exceeded (http.StatusBadRequest).
	// - A database/internal error occurs (http.StatusInternalServerError).
	AddImage(ctx context.Context, req *imagemodel.AddRequest) (*imagemodel.AddResponse, error)
	// DeleteImage removes an image from a seminar. It's called by media-service-go upon successful image deletion.
	// The function validates the request and removes the image information from the seminar.
	//
	// Returns an error if:
	// - The request payload is invalid (http.StatusBadRequest).
	// - The seminar (owner) or image is not found (http.StatusNotFound).
	// - A database/internal error occurs (http.StatusInternalServerError).
	DeleteImage(ctx context.Context, req *imagemodel.DeleteRequest) error
}

// service provides service-layer business logic for seminar models.
// It holds [seminarrepo.Repository] and [productrepo.Repository] instances
// to perform database operations.
type service struct {
	SeminarRepo seminarrepo.Repository
	ProductRepo productrepo.Repository
}

// Error represents seminar service error.
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

// New creates a new service instance with provided seminar and product repositories.
func New(sr seminarrepo.Repository, pr productrepo.Repository) Service {
	return &service{
		SeminarRepo: sr,
		ProductRepo: pr,
	}
}

// Get retrieves a single published and not soft-deleted seminar record from the database,
// along with all of its associated products details (prices and product IDs).
//
// Returns a SeminarDetails struct containing the combined information.
// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Get(ctx context.Context, id string) (*seminarmodel.SeminarDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{Msg: "Invalid seminar ID", Err: err, Code: http.StatusBadRequest}
	}
	seminar, err := s.SeminarRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{Msg: "Seminar not found", Err: err, Code: http.StatusNotFound}
		}
		return nil, &Error{Msg: "Failed to retrieve seminar", Err: err, Code: http.StatusInternalServerError}
	}

	if seminar.ReservationProductID == nil || seminar.EarlyProductID == nil || seminar.LateProductID == nil || seminar.EarlySurchargeProductID == nil || seminar.LateSurchargeProductID == nil {
		return nil, &Error{
			Msg:  "seminar record is missing one or more required product IDs",
			Err:  errors.New("incomplete seminar data"),
			Code: http.StatusInternalServerError,
		}
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
		return nil, &Error{Msg: "failed to get seminar products", Err: err, Code: http.StatusInternalServerError}
	}
	if len(products) != 5 {
		return nil, &Error{Msg: "could not find all products for seminar", Err: err, Code: http.StatusNotFound}
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
// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) GetWithDeleted(ctx context.Context, id string) (*seminarmodel.SeminarDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{Msg: "Invalid seminar ID", Err: err, Code: http.StatusBadRequest}
	}
	seminar, err := s.SeminarRepo.GetWithDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{Msg: "Seminar not found", Err: err, Code: http.StatusNotFound}
		}
		return nil, &Error{Msg: "Failed to retrieve seminar", Err: err, Code: http.StatusInternalServerError}
	}

	if seminar.ReservationProductID == nil || seminar.EarlyProductID == nil || seminar.LateProductID == nil || seminar.EarlySurchargeProductID == nil || seminar.LateSurchargeProductID == nil {
		return nil, &Error{
			Msg:  "seminar record is missing one or more required product IDs",
			Err:  errors.New("incomplete seminar data"),
			Code: http.StatusInternalServerError,
		}
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
		return nil, &Error{Msg: "failed to get seminar products", Err: err, Code: http.StatusInternalServerError}
	}
	if len(products) != 5 {
		return nil, &Error{Msg: "could not find all products for seminar", Err: err, Code: http.StatusNotFound}
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
// Returns an error if the ID is invalid (http.StatusBadRequest), the record is not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) GetWithUnpublished(ctx context.Context, id string) (*seminarmodel.SeminarDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &Error{Msg: "Invalid seminar ID", Err: err, Code: http.StatusBadRequest}
	}
	seminar, err := s.SeminarRepo.GetWithUnpublished(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &Error{Msg: "Seminar not found", Err: err, Code: http.StatusNotFound}
		}
		return nil, &Error{Msg: "Failed to retrieve seminar", Err: err, Code: http.StatusInternalServerError}
	}

	if seminar.ReservationProductID == nil || seminar.EarlyProductID == nil || seminar.LateProductID == nil || seminar.EarlySurchargeProductID == nil || seminar.LateSurchargeProductID == nil {
		return nil, &Error{
			Msg:  "seminar record is missing one or more required product IDs",
			Err:  errors.New("incomplete seminar data"),
			Code: http.StatusInternalServerError,
		}
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
		return nil, &Error{Msg: "failed to get seminar products", Err: err, Code: http.StatusInternalServerError}
	}
	if len(products) != 5 {
		return nil, &Error{Msg: "could not find all products for seminar", Err: err, Code: http.StatusNotFound}
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
//
// Returns a slice of SeminarDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs (http.StatusInternalServerError).
func (s *service) List(ctx context.Context, limit, offset int) ([]seminarmodel.SeminarDetails, int64, error) {
	seminars, err := s.SeminarRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, &Error{Msg: "Failed to retrieve seminars", Err: err, Code: http.StatusInternalServerError}
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
		return nil, 0, &Error{Msg: "Failed to retrieve products", Err: err, Code: http.StatusInternalServerError}
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
		return nil, 0, &Error{Msg: "Failed to count seminars", Err: err, Code: http.StatusInternalServerError}
	}
	return allDetails, total, nil
}

// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) seminar records.
// Each record is returned with its associated products details.
//
// Returns a slice of SeminarDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs (http.StatusInternalServerError).
func (s *service) ListUnpublished(ctx context.Context, limit, offset int) ([]seminarmodel.SeminarDetails, int64, error) {
	seminars, err := s.SeminarRepo.ListUnpublished(ctx, limit, offset)
	if err != nil {
		return nil, 0, &Error{Msg: "Failed to retrieve seminars", Err: err, Code: http.StatusInternalServerError}
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
		return nil, 0, &Error{Msg: "Failed to retrieve products", Err: err, Code: http.StatusInternalServerError}
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
		return nil, 0, &Error{Msg: "Failed to count seminars", Err: err, Code: http.StatusInternalServerError}
	}
	return allDetails, total, nil
}

// ListDeleted retrieves a paginated list of all soft-deleted seminar records.
// Each record is returned with its associated products details.
//
// Returns a slice of SeminarDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs (http.StatusInternalServerError).
func (s *service) ListDeleted(ctx context.Context, limit, offset int) ([]seminarmodel.SeminarDetails, int64, error) {
	seminars, err := s.SeminarRepo.ListDeleted(ctx, limit, offset)
	if err != nil {
		return nil, 0, &Error{Msg: "Failed to retrieve seminars", Err: err, Code: http.StatusInternalServerError}
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
		return nil, 0, &Error{Msg: "Failed to retrieve products", Err: err, Code: http.StatusInternalServerError}
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
		return nil, 0, &Error{Msg: "Failed to count seminars", Err: err, Code: http.StatusInternalServerError}
	}
	return allDetails, total, nil
}

// Create creates a new Seminar record and all of its associated Product records in the database.
// It validates the request payload to ensure all required fields are present.
// The seminar and all of the associated products are created in an unpublished state (`InStock: false`).
//
// Returns a CreateResponse containing the newly created SeminarID, ReservationProductID, EarlyProductID,
// LateProductID, EarlySurchargeProductID, LateSurchargeProductID.
// Returns an error if the request payload is invalid (http.StatusBadRequest) or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Create(ctx context.Context, req *seminarmodel.CreateRequest) (*seminarmodel.CreateResponse, error) {
	seminar := &seminarmodel.Seminar{}
	err := s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		if err := req.Validate(); err != nil {
			return &Error{
				Msg:  "Invalid request payload",
				Err:  err,
				Code: http.StatusBadRequest,
			}
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
			return &Error{Msg: "Failed to create seminar products", Err: err, Code: http.StatusInternalServerError}
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
			return &Error{Msg: "Failed to create seminar", Err: err, Code: http.StatusInternalServerError}
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
// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Publish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &Error{Msg: "Invalid seminar ID", Err: err, Code: http.StatusBadRequest}
	}
	return s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txSeminarRepo.SetInStock(ctx, id, true)
		if err != nil {
			return &Error{Msg: "Failed to publish seminar", Err: err, Code: http.StatusInternalServerError}
		} else if ra == 0 {
			return &Error{Msg: "Seminar not found", Err: err, Code: http.StatusNotFound}
		}
		ra, err = txProductRepo.SetInStockByDetailsID(ctx, id, true)
		if err != nil {
			return &Error{Msg: "Failed to publish seminar products", Err: err, Code: http.StatusInternalServerError}
		} else if ra != 5 {
			return &Error{Msg: "Failed to publish all seminar products", Err: err, Code: http.StatusInternalServerError}
		}
		return nil
	})
}

// Unpublish sets the `InStock` field to false for a seminar and all of its associated products,
// archiving it from the catalog.
//
// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Unpublish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &Error{Msg: "Invalid seminar ID", Err: err, Code: http.StatusBadRequest}
	}
	return s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txSeminarRepo.SetInStock(ctx, id, false)
		if err != nil {
			return &Error{Msg: "Failed to unpublish seminar", Err: err, Code: http.StatusInternalServerError}
		} else if ra == 0 {
			return &Error{Msg: "Seminar not found", Err: err, Code: http.StatusNotFound}
		}
		ra, err = txProductRepo.SetInStockByDetailsID(ctx, id, false)
		if err != nil {
			return &Error{Msg: "Failed to unpublish seminar products", Err: err, Code: http.StatusInternalServerError}
		} else if ra != 5 {
			return &Error{Msg: "Failed to unpublish all seminar products", Err: err, Code: http.StatusInternalServerError}
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
// Returns an error if the request payload is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Update(ctx context.Context, req *seminarmodel.UpdateRequest) (map[string]any, error) {
	allUpdates := make(map[string]any)
	err := s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		if err := req.Validate(); err != nil {
			validationMsg, err := json.Marshal(err)
			return &Error{Msg: string(validationMsg), Err: err, Code: http.StatusBadRequest}
		}

		seminar, err := txSeminarRepo.Get(ctx, req.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{Msg: "seminar not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "failed to find seminar", Err: err, Code: http.StatusInternalServerError}
		}

		if seminar.ReservationProductID == nil || seminar.EarlyProductID == nil || seminar.LateProductID == nil || seminar.EarlySurchargeProductID == nil || seminar.LateSurchargeProductID == nil {
			return &Error{
				Msg:  "seminar record is missing one or more required product IDs",
				Err:  errors.New("incomplete seminar data"),
				Code: http.StatusInternalServerError,
			}
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
			return &Error{Msg: "failed to get seminar products", Err: err, Code: http.StatusInternalServerError}
		}
		if len(products) != 5 {
			return &Error{Msg: "could not find all products for seminar", Err: err, Code: http.StatusNotFound}
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
				return nil, &Error{Msg: "product to update not found", Err: errors.New("nil product pointer"), Code: http.StatusNotFound}
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
			return &Error{Msg: "Cannot find all seminar products", Err: nil, Code: http.StatusNotFound}
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
				return &Error{
					Msg:  fmt.Sprintf("Failed to update %s", key),
					Err:  err,
					Code: http.StatusInternalServerError,
				}
			}
			if len(pu) > 0 {
				allUpdates[key] = pu
			}
		}

		if len(seminarUpdates) > 0 {
			if _, err := txSeminarRepo.Update(ctx, seminar, seminarUpdates); err != nil {
				return &Error{
					Msg:  "Failed to update seminar",
					Err:  err,
					Code: http.StatusInternalServerError,
				}
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
// The function validates the request, checks the image limit, and appends the new image information.
//
// Returns an AddResponse with the MediaServiceID on success.
// Returns an error if:
// - The request payload is invalid (http.StatusBadRequest).
// - The seminar (owner) is not found (http.StatusNotFound).
// - The image limit (5) is exceeded (http.StatusBadRequest).
// - A database/internal error occurs (http.StatusInternalServerError).
func (s *service) AddImage(ctx context.Context, req *imagemodel.AddRequest) (*imagemodel.AddResponse, error) {
	if err := req.Validate(); err != nil {
		validationMsg, _ := json.Marshal(err)
		return nil, &Error{Msg: string(validationMsg), Err: err, Code: http.StatusBadRequest}
	}

	err := s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)

		seminarRec, err := txSeminarRepo.GetWithUnpublished(ctx, req.OwnerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{Msg: "Seminar not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "Failed to get seminar", Err: err, Code: http.StatusInternalServerError}
		}

		if seminarRec.UploadedImageAmount >= 5 {
			return &Error{Msg: "Maximum number of uploaded images is 5 per item", Err: nil, Code: http.StatusBadRequest}
		}

		newImage := &imagemodel.Image{
			PublicID:       req.PublicID,
			URL:            req.URL,
			SecureURL:      req.SecureURL,
			MediaServiceID: req.MediaServiceID,
		}

		if err := txSeminarRepo.AddImage(ctx, seminarRec, newImage); err != nil {
			return &Error{Msg: "Failed to add image to seminar", Err: err, Code: http.StatusInternalServerError}
		}

		// Increment the image count and save
		seminarRec.UploadedImageAmount++
		if _, err := txSeminarRepo.Update(ctx, seminarRec, map[string]any{"uploaded_image_amount": seminarRec.UploadedImageAmount}); err != nil {
			return &Error{Msg: "Failed to update seminar", Err: err, Code: http.StatusInternalServerError}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &imagemodel.AddResponse{MediaServiceID: req.MediaServiceID}, nil
}

// DeleteImage removes an image from a seminar. It's called by media-service-go upon successful image deletion.
// The function validates the request and removes the image information from the seminar.
//
// Returns an error if:
// - The request payload is invalid (http.StatusBadRequest).
// - The seminar (owner) or image is not found (http.StatusNotFound).
// - A database/internal error occurs (http.StatusInternalServerError).
func (s *service) DeleteImage(ctx context.Context, req *imagemodel.DeleteRequest) error {
	if err := req.Validate(); err != nil {
		validationMsg, _ := json.Marshal(err)
		return &Error{Msg: string(validationMsg), Err: err, Code: http.StatusBadRequest}
	}

	return s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)

		seminarRec, err := txSeminarRepo.GetWithUnpublished(ctx, req.OwnerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{Msg: "Seminar not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "Failed to get seminar", Err: err, Code: http.StatusInternalServerError}
		}

		var imageFound bool
		for i, img := range seminarRec.Images {
			if img.MediaServiceID == req.MediaServiceID {
				imageFound = true
				// Remove the image from the slice
				seminarRec.Images = append(seminarRec.Images[:i], seminarRec.Images[i+1:]...)
				break
			}
		}

		if !imageFound {
			return &Error{Msg: "Image not found on seminar", Err: gorm.ErrRecordNotFound, Code: http.StatusNotFound}
		}

		if err := txSeminarRepo.DeleteImage(ctx, seminarRec, req.MediaServiceID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{Msg: "Image not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "Failed to delete seminar image", Err: err, Code: http.StatusInternalServerError}
		}

		seminarRec.UploadedImageAmount--

		// Save the updated image list and count
		if _, err := txSeminarRepo.Update(ctx, seminarRec, map[string]any{"images": seminarRec.Images, "uploaded_image_amount": seminarRec.UploadedImageAmount}); err != nil {
			return &Error{Msg: "Failed to update seminar", Err: err, Code: http.StatusInternalServerError}
		}
		return nil
	})
}

// Delete performs a soft-delete of a seminar and all of its related product records.
// It also unpublishes all records, meaning they must be manually published again after restoration.
//
// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &Error{Msg: "invalid seminar id", Err: err, Code: http.StatusBadRequest}
	}
	return s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		// Check if seminar exists
		if _, err := txSeminarRepo.GetWithUnpublished(ctx, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &Error{Msg: "Seminar not found", Err: err, Code: http.StatusNotFound}
			}
			return &Error{Msg: "Failed to get seminar", Err: err, Code: http.StatusInternalServerError}
		}

		// Unpublish all instances
		if _, err := txSeminarRepo.SetInStock(ctx, id, false); err != nil {
			return &Error{Msg: "failed to unpublish seminar", Err: err, Code: http.StatusInternalServerError}
		}
		ra, err := txProductRepo.SetInStockByDetailsID(ctx, id, false)
		if err != nil {
			return &Error{Msg: "failed to unpublish seminar products", Err: err, Code: http.StatusInternalServerError}
		} else if ra != 5 {
			return &Error{Msg: "failed to unpublish all seminar products", Err: err, Code: http.StatusInternalServerError}
		}

		// Delete all instances
		if _, err = txSeminarRepo.Delete(ctx, id); err != nil {
			return &Error{Msg: "failed to delete seminar", Err: err, Code: http.StatusInternalServerError}
		}
		if _, err = txProductRepo.DeleteByDetailsID(ctx, id); err != nil {
			return &Error{Msg: "failed to delete seminar products", Err: err, Code: http.StatusInternalServerError}
		}
		return nil
	})
}

// DeletePermanent performs a complete delete of a seminar and its related product records.
//
// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) DeletePermanent(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &Error{Msg: "invalid seminar id", Err: err, Code: http.StatusBadRequest}
	}
	return s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		ra, err := txSeminarRepo.DeletePermanent(ctx, id)
		if err != nil {
			return &Error{Msg: "failed to delete seminar", Err: err, Code: http.StatusInternalServerError}
		} else if ra == 0 {
			return &Error{Msg: "seminar not found", Err: err, Code: http.StatusNotFound}
		}

		ra, err = txProductRepo.DeletePermanentByDetailsID(ctx, id)
		if err != nil {
			return &Error{Msg: "failed to delete seminar products", Err: err, Code: http.StatusInternalServerError}
		} else if ra != 5 {
			return &Error{Msg: "failed to delete all seminar products", Err: err, Code: http.StatusInternalServerError}
		}
		return nil
	})
}

// Restore performs a restore of a seminar and its related product records.
// Seminar and its related product records are not being published. This should be
// done manually.
//
// Returns an error if the ID is invalid (http.StatusBadRequest), the records are not found (http.StatusNotFound),
// or a database/internal error occurs (http.StatusInternalServerError).
func (s *service) Restore(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return &Error{Msg: "invalid seminar id", Err: err, Code: http.StatusBadRequest}
	}
	return s.SeminarRepo.DB().Transaction(func(tx *gorm.DB) error {
		txSeminarRepo := s.SeminarRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txSeminarRepo.Restore(ctx, id)
		if err != nil {
			return &Error{Msg: "failed to delete seminar", Err: err, Code: http.StatusInternalServerError}
		} else if ra == 0 {
			return &Error{Msg: "seminar not found", Err: err, Code: http.StatusNotFound}
		}
		ra, err = txProductRepo.RestoreByDetailsID(ctx, id)
		if err != nil {
			return &Error{Msg: "failed to delete seminar products", Err: err, Code: http.StatusInternalServerError}
		} else if ra != 5 {
			return &Error{Msg: "failed to delete all seminar products", Err: err, Code: http.StatusInternalServerError}
		}
		return nil
	})
}
