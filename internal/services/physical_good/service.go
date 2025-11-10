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

// Package physicalgood provides service-layer business logic for physical goods.
package physicalgood

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	physicalgoodrepo "github.com/mikhail5545/product-service-go/internal/database/physical_good"
	productrepo "github.com/mikhail5545/product-service-go/internal/database/product"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	physicalgoodmodel "github.com/mikhail5545/product-service-go/internal/models/physical_good"
	productmodel "github.com/mikhail5545/product-service-go/internal/models/product"
	imageservice "github.com/mikhail5545/product-service-go/internal/services/image"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/services/physical_good_mock/service_mock.go -package=physical_good_mock . Service

// Service provides service-layer business logic for physical good models.
type Service interface {
	// Get retrieves a single published and not soft-deleted physical good record from the database,
	// along with its associated product details (price and product ID).
	//
	// Returns a PhysicalGoodDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	Get(ctx context.Context, id string) (*physicalgoodmodel.PhysicalGoodDetails, error)
	// GetWithDeleted retrieves a single physical good record from the database, including soft-deleted ones,
	// along with its associated product details (price and product ID).
	//
	// Returns a PhysicalGoodDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	GetWithDeleted(ctx context.Context, id string) (*physicalgoodmodel.PhysicalGoodDetails, error)
	// GetWithUnpublished retrieves a single physical good record from the database, including unpublished ones (but not soft-deleted),
	// along with its associated product details (price and product ID).
	//
	// Returns a PhysicalGoodDetails struct containing the combined information.
	// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
	// or a database/internal error occurs.
	GetWithUnpublished(ctx context.Context, id string) (*physicalgoodmodel.PhysicalGoodDetails, error)
	// List retrieves a paginated list of all published and not soft-deleted physical good records.
	// Each record is returned with its associated product details.
	//
	// Returns a slice of PhysicalGoodDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs.
	List(ctx context.Context, limit, offset int) ([]physicalgoodmodel.PhysicalGoodDetails, int64, error)
	// ListDeleted retrieves a paginated list of all soft-deleted physical good records.
	// Each record is returned with its associated product details.
	//
	// Returns a slice of PhysicalGoodDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs.
	ListDeleted(ctx context.Context, limit, offset int) ([]physicalgoodmodel.PhysicalGoodDetails, int64, error)
	// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) physical good records.
	// Each record is returned with its associated product details.
	//
	// Returns a slice of PhysicalGoodDetails, the total count of such records, and an error if one occurs.
	// Returns an error if a database/internal error occurs.
	ListUnpublished(ctx context.Context, limit, offset int) ([]physicalgoodmodel.PhysicalGoodDetails, int64, error)
	// Create creates a new PhysicalGood record and its associated Product record in the database.
	// It validates the request payload to ensure all required fields are present.
	// Both the physical good and the product are created in an unpublished state (`InStock: false`).
	//
	// Returns a CreateResponse containing the newly created PhysicalGoodID and ProductID.
	// Returns an error if the request payload is invalid (ErrInvalidArgument) or a database/internal error occurs.
	Create(ctx context.Context, req *physicalgoodmodel.CreateRequest) (*physicalgoodmodel.CreateResponse, error)
	// Update performs a partial update of a physical good and its related product.
	// The request should contain the physical good's ID and the fields to be updated.
	// At least one field must be provided for an update to occur.
	//
	// Returns a map containing the fields that were actually changed, nested under "physical_good" and "product" keys.
	// Example: `{"physical_good": {"name": "new name"}, "product": {"price": 99.99}}`
	// Returns an error if the request payload is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Update(ctx context.Context, req *physicalgoodmodel.UpdateRequest) (map[string]any, error)
	// Publish sets the `InStock` field to true for a physical good and its associated product,
	// making it available in the catalog.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Publish(ctx context.Context, id string) error
	// Unpublish sets the `InStock` field to false for a physical good and its associated product,
	// archiving it from the catalog.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Unpublish(ctx context.Context, id string) error
	// Delete performs a soft-delete of a physical good and its related product record.
	// It also unpublishes both records, meaning they must be manually published again after restoration.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Delete(ctx context.Context, id string) error
	// DeletePermanent performs a complete delete of a physical good and its related product record.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	DeletePermanent(ctx context.Context, id string) error
	// Restore performs a restore of a physical good and its related product record.
	// Physical good and its related product record are not being published. This should be
	// done manually.
	//
	// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
	// or a database/internal error occurs.
	Restore(ctx context.Context, id string) error
	// AddImage adds a new image to a physical good. It's called by media-service-go upon successful image upload.
	// It uses physicalGoodOwnerRepoAdapter to call [imageservice.AddImage] and add an image to the physical good.
	//
	// Returns an error if:
	// - The request payload is invalid ([imageservice.ErrInvalidArgument]).
	// - The physical good (owner) is not found ([imageservice.ErrOwnerNotFound]).
	// - The image limit (5) is exceeded ([imageservice.ErrImageLimitExceeded]).
	// - A database/internal error occurs.
	AddImage(ctx context.Context, req *imagemodel.AddRequest) error
	// DeleteImage removes an image from a physical good. It's called by media-service-go upon successful image deletion.
	// It uses physicalGoodOwnerRepoAdapter to call [imageservice.DeleteImage] and delete an image from the physical good.
	//
	// Returns an error if:
	//   - The request payload is invalid ([imageservice.ErrInvalidArgument]).
	//   - The physical good (owner) is not found ([imageservice.ErrOwnerNotFound]).
	//   - The image is not found on physical good (owner) ([imageservice.ErrImageNotFoundOnOwner]).
	//   - A database/internal error occurs.
	DeleteImage(ctx context.Context, req *imagemodel.DeleteRequest) error
	// AddImageBatch adds an image for a batch of physical goods. It uses seminarOwnerRepoAdapter
	// to call [imageservice.AddImageBatch] and append images to the seminar. It's called by media-service-go
	// upon successfull context change.
	//
	// It returns the number of affected physical goods.
	// Returns an error if the request is invalid ([imageservice.ErrInvalidArgument]), no physical goods (owners) are not found ([imageservice.ErrOwnersNotFound])
	// or a database/internal error occurs.
	AddImageBatch(ctx context.Context, req *imagemodel.AddBatchRequest) (int, error)
	// DeleteImageBatch removes an image from a batch of physical goods. It uses seminarOwnerRepoAdapter
	// to call [imageservice.DeleteImageBatch] and append images to the seminar.
	//
	// It returns the number of affected physical goods.
	// Returns an error if the request is invalid ([imageservice.ErrInvalidArgument]), no physical goods (owners) are not found ([imageservice.ErrOwnersNotFound]),
	// no associations were found ([imageservice.ErrAssociationsNotFound]) or a database/internal error occurs.
	DeleteImageBatch(ctx context.Context, req *imagemodel.DeleteBatchRequst) (int, error)
}

// service provides service-layer business logic for physical good models.
// It holds [physicalgoodrepo.Repository] and [productrepo.Repository] instances
// to perform database operations.
type service struct {
	PhysicalGoodRepo physicalgoodrepo.Repository
	ProductRepo      productrepo.Repository
	ImageSvc         imageservice.Service
}

// New creates a new service instance with provided physical good and product repositories.
func New(gr physicalgoodrepo.Repository, pr productrepo.Repository, is imageservice.Service) Service {
	return &service{
		PhysicalGoodRepo: gr,
		ProductRepo:      pr,
		ImageSvc:         is,
	}
}

// Get retrieves a single published and not soft-deleted physical good record from the database,
// along with its associated product details (price and product ID).
//
// Returns a PhysicalGoodDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Get(ctx context.Context, id string) (*physicalgoodmodel.PhysicalGoodDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	phGood, err := s.PhysicalGoodRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve physical good: %w", err)
	}
	product, err := s.ProductRepo.SelectByDetailsID(ctx, phGood.ID, "id", "price")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve physical good product: %w", err)
	}
	return &physicalgoodmodel.PhysicalGoodDetails{
		PhysicalGood: phGood,
		Price:        product.Price,
		ProductID:    product.ID,
	}, nil
}

// GetWithDeleted retrieves a single physical good record from the database, including soft-deleted ones,
// along with its associated product details (price and product ID).
//
// Returns a PhysicalGoodDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) GetWithDeleted(ctx context.Context, id string) (*physicalgoodmodel.PhysicalGoodDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	phGood, err := s.PhysicalGoodRepo.GetWithDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve physical good: %w", err)
	}
	product, err := s.ProductRepo.SelectWithDeletedByDetailsID(ctx, phGood.ID, "id", "price")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve physical good product: %w", err)
	}
	return &physicalgoodmodel.PhysicalGoodDetails{
		PhysicalGood: phGood,
		Price:        product.Price,
		ProductID:    product.ID,
	}, nil
}

// GetWithUnpublished retrieves a single physical good record from the database, including unpublished ones (but not soft-deleted),
// along with its associated product details (price and product ID).
//
// Returns a PhysicalGoodDetails struct containing the combined information.
// Returns an error if the ID is invalid (ErrInvalidArgument), the record is not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) GetWithUnpublished(ctx context.Context, id string) (*physicalgoodmodel.PhysicalGoodDetails, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	phGood, err := s.PhysicalGoodRepo.GetWithUnpublished(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve physical good: %w", err)
	}
	product, err := s.ProductRepo.SelectWithUnpublishedByDetailsID(ctx, phGood.ID, "id", "price")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("failed to retrieve physical good product: %w", err)
	}
	return &physicalgoodmodel.PhysicalGoodDetails{
		PhysicalGood: phGood,
		Price:        product.Price,
		ProductID:    product.ID,
	}, nil
}

// List retrieves a paginated list of all published and not soft-deleted physical good records.
// Each record is returned with its associated product details.
//
// Returns a slice of PhysicalGoodDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs.
func (s *service) List(ctx context.Context, limit, offset int) ([]physicalgoodmodel.PhysicalGoodDetails, int64, error) {
	phGoods, err := s.PhysicalGoodRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve physical goods: %w", err)
	}
	total, err := s.PhysicalGoodRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count physical goods: %w", err)
	}

	phGoodsMap := make(map[string]*physicalgoodmodel.PhysicalGood, len(phGoods))
	var phGoodsIDs []string
	for i := range phGoods {
		phGoodsMap[phGoods[i].ID] = &phGoods[i]
		phGoodsIDs = append(phGoodsIDs, phGoods[i].ID)
	}

	products, err := s.ProductRepo.SelectByDetailsIDs(ctx, phGoodsIDs, "id", "price", "details_id")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}
	var allDetails []physicalgoodmodel.PhysicalGoodDetails
	for _, p := range products {
		allDetails = append(allDetails, physicalgoodmodel.PhysicalGoodDetails{
			PhysicalGood: phGoodsMap[p.DetailsID],
			Price:        p.Price,
			ProductID:    p.ID,
		})
	}
	return allDetails, total, nil
}

// ListUnpublished retrieves a paginated list of all unpublished (but not soft-deleted) physical good records.
// Each record is returned with its associated product details.
//
// Returns a slice of PhysicalGoodDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs.
func (s *service) ListUnpublished(ctx context.Context, limit, offset int) ([]physicalgoodmodel.PhysicalGoodDetails, int64, error) {
	phGoods, err := s.PhysicalGoodRepo.ListUnpublished(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve physical goods: %w", err)
	}
	total, err := s.PhysicalGoodRepo.CountUnpublished(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count physical goods: %w", err)
	}

	phGoodsMap := make(map[string]*physicalgoodmodel.PhysicalGood, len(phGoods))
	var phGoodsIDs []string
	for i := range phGoods {
		phGoodsMap[phGoods[i].ID] = &phGoods[i]
		phGoodsIDs = append(phGoodsIDs, phGoods[i].ID)
	}

	products, err := s.ProductRepo.SelectWithUnpublishedByDetailsIDs(ctx, phGoodsIDs, "id", "price", "details_id")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}
	var allDetails []physicalgoodmodel.PhysicalGoodDetails
	for _, p := range products {
		allDetails = append(allDetails, physicalgoodmodel.PhysicalGoodDetails{
			PhysicalGood: phGoodsMap[p.DetailsID],
			Price:        p.Price,
			ProductID:    p.ID,
		})
	}
	return allDetails, total, nil
}

// ListDeleted retrieves a paginated list of all soft-deleted physical good records.
// Each record is returned with its associated product details.
//
// Returns a slice of PhysicalGoodDetails, the total count of such records, and an error if one occurs.
// Returns an error if a database/internal error occurs.
func (s *service) ListDeleted(ctx context.Context, limit, offset int) ([]physicalgoodmodel.PhysicalGoodDetails, int64, error) {
	phGoods, err := s.PhysicalGoodRepo.ListDeleted(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve physical goods: %w", err)
	}
	total, err := s.PhysicalGoodRepo.CountDeleted(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count physical goods: %w", err)
	}

	phGoodsMap := make(map[string]*physicalgoodmodel.PhysicalGood, len(phGoods))
	var phGoodsIDs []string
	for i := range phGoods {
		phGoodsMap[phGoods[i].ID] = &phGoods[i]
		phGoodsIDs = append(phGoodsIDs, phGoods[i].ID)
	}

	products, err := s.ProductRepo.SelectWithDeletedByDetailsIDs(ctx, phGoodsIDs, "id", "price", "details_id")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve products: %w", err)
	}
	var allDetails []physicalgoodmodel.PhysicalGoodDetails
	for _, p := range products {
		allDetails = append(allDetails, physicalgoodmodel.PhysicalGoodDetails{
			PhysicalGood: phGoodsMap[p.DetailsID],
			Price:        p.Price,
			ProductID:    p.ID,
		})
	}
	return allDetails, total, nil
}

// Create creates a new PhysicalGood record and its associated Product record in the database.
// It validates the request payload to ensure all required fields are present.
// Both the physical good and the product are created in an unpublished state (`InStock: false`).
//
// Returns a CreateResponse containing the newly created PhysicalGoodID and ProductID.
// Returns an error if the request payload is invalid (ErrInvalidArgument) or a database/internal error occurs.
func (s *service) Create(ctx context.Context, req *physicalgoodmodel.CreateRequest) (*physicalgoodmodel.CreateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	var phGoodID, productID string
	err := s.PhysicalGoodRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPhysicalGoodRepo := s.PhysicalGoodRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		phGood := &physicalgoodmodel.PhysicalGood{
			ID:               uuid.New().String(),
			Name:             req.Name,
			ShortDescription: req.ShortDescription,
			Amount:           req.Amount,
			ShippingRequired: req.ShippingRequired,
			InStock:          false,
		}

		product := &productmodel.Product{
			ID:          uuid.New().String(),
			Price:       req.Price,
			DetailsID:   phGood.ID,
			DetailsType: "physical_good",
			InStock:     false,
		}

		if err := txPhysicalGoodRepo.Create(ctx, phGood); err != nil {
			return fmt.Errorf("failed to create physical good: %w", err)
		}
		if err := txProductRepo.Create(ctx, product); err != nil {
			return fmt.Errorf("failed to create physical good product: %w", err)
		}

		phGoodID = phGood.ID
		productID = product.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &physicalgoodmodel.CreateResponse{ID: phGoodID, ProductID: productID}, nil
}

// Publish sets the `InStock` field to true for a physical good and its associated product,
// making it available in the catalog.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Publish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.PhysicalGoodRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPhysicalGoodRepo := s.PhysicalGoodRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txPhysicalGoodRepo.SetInStock(ctx, id, true)
		if err != nil {
			return fmt.Errorf("failed to publish physical good: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		ra, err = txProductRepo.SetInStockByDetailsID(ctx, id, true)
		if err != nil {
			return fmt.Errorf("failed to publish physical good product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil
	})
}

// Unpublish sets the `InStock` field to false for a physical good and its associated product,
// archiving it from the catalog.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Unpublish(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.PhysicalGoodRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPhysicalGoodRepo := s.PhysicalGoodRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txPhysicalGoodRepo.SetInStock(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to unpublish physical good: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		ra, err = txProductRepo.SetInStockByDetailsID(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to unpublish physical good product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil
	})
}

// Update performs a partial update of a physical good and its related product.
// The request should contain the physical good's ID and the fields to be updated.
// At least one field must be provided for an update to occur.
//
// Returns a map containing the fields that were actually changed, nested under "physical_good" and "product" keys.
// Example: `{"physical_good": {"name": "new name"}, "product": {"price": 99.99}}`
// Returns an error if the request payload is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Update(ctx context.Context, req *physicalgoodmodel.UpdateRequest) (map[string]any, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	allUpdates := make(map[string]any)
	err := s.PhysicalGoodRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPhysicalGoodRepo := s.PhysicalGoodRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		phGood, err := txPhysicalGoodRepo.Get(ctx, req.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to retrieve physical good: %w", err)
		}
		product, err := txProductRepo.SelectByDetailsID(ctx, req.ID, "id", "price")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to retrieve physical good product: %w", err)
		}

		updates := make(map[string]any)
		productUpdates := make(map[string]any)
		if req.Name != nil && *req.Name != phGood.Name {
			updates["name"] = *req.Name
		}
		if req.ShortDescription != nil && *req.ShortDescription != phGood.ShortDescription {
			updates["short_description"] = *req.ShortDescription
		}
		if req.LongDescription != nil && *req.LongDescription != phGood.LongDescription {
			updates["long_description"] = *req.LongDescription
		}
		if req.Amount != nil && *req.Amount != phGood.Amount {
			updates["amount"] = *req.Amount
		}
		if req.ShippingRequired != nil && *req.ShippingRequired != phGood.ShippingRequired {
			updates["shipping_required"] = *req.ShippingRequired
		}
		if len(req.Tags) > 0 {
			updates["tags"] = req.Tags
		}
		if req.Price != nil && *req.Price != product.Price {
			productUpdates["price"] = *req.Price
		}

		if len(updates) > 0 {
			if _, err := txPhysicalGoodRepo.Update(ctx, phGood, updates); err != nil {
				return fmt.Errorf("failed to update physical good: %w", err)
			}
		}
		if len(productUpdates) > 0 {
			if _, err := txProductRepo.Update(ctx, product, productUpdates); err != nil {
				return fmt.Errorf("failed to update physical good product: %w", err)
			}
		}

		allUpdates["physical_good"] = updates
		allUpdates["product"] = productUpdates
		return nil
	})
	if err != nil {
		return nil, err
	}
	return allUpdates, nil
}

// AddImage adds a new image to a physical good. It's called by media-service-go upon successful image upload.
// It uses physicalGoodOwnerRepoAdapter to call [imageservice.AddImage] and add an image to the physical good.
//
// Returns an error if:
// - The request payload is invalid ([imageservice.ErrInvalidArgument]).
// - The physical good (owner) is not found ([imageservice.ErrOwnerNotFound]).
// - The image limit (5) is exceeded ([imageservice.ErrImageLimitExceeded]).
// - A database/internal error occurs.
func (s *service) AddImage(ctx context.Context, req *imagemodel.AddRequest) error {
	ownerRepoAdapter := newPhysicalGoodOwnerRepoAdapter(s.PhysicalGoodRepo)
	return s.ImageSvc.AddImage(ctx, req, ownerRepoAdapter)
}

// DeleteImage removes an image from a physical good. It's called by media-service-go upon successful image deletion.
// It uses physicalGoodOwnerRepoAdapter to call [imageservice.DeleteImage] and delete an image from the physical good.
//
// Returns an error if:
//   - The request payload is invalid ([imageservice.ErrInvalidArgument]).
//   - The physical good (owner) is not found ([imageservice.ErrOwnerNotFound]).
//   - The image is not found on physical good (owner) ([imageservice.ErrImageNotFoundOnOwner]).
//   - A database/internal error occurs.
func (s *service) DeleteImage(ctx context.Context, req *imagemodel.DeleteRequest) error {
	ownerRepoAdapter := newPhysicalGoodOwnerRepoAdapter(s.PhysicalGoodRepo)
	return s.ImageSvc.DeleteImage(ctx, req, ownerRepoAdapter)
}

// AddImageBatch adds an image for a batch of physical goods. It uses seminarOwnerRepoAdapter
// to call [imageservice.AddImageBatch] and append images to the seminar. It's called by media-service-go
// upon successfull context change.
//
// It returns the number of affected physical goods.
// Returns an error if the request is invalid ([imageservice.ErrInvalidArgument]), no physical goods (owners) are not found ([imageservice.ErrOwnersNotFound])
// or a database/internal error occurs.
func (s *service) AddImageBatch(ctx context.Context, req *imagemodel.AddBatchRequest) (int, error) {
	// Use the adapter to bridge the specific course repository to the generic owner repository interface.
	ownerRepoAdapter := newPhysicalGoodOwnerRepoAdapter(s.PhysicalGoodRepo)
	return s.ImageSvc.AddImageBatch(ctx, req, ownerRepoAdapter)
}

// DeleteImageBatch removes an image from a batch of physical goods. It uses seminarOwnerRepoAdapter
// to call [imageservice.DeleteImageBatch] and append images to the seminar.
//
// It returns the number of affected physical goods.
// Returns an error if the request is invalid ([imageservice.ErrInvalidArgument]), no physical goods (owners) are not found ([imageservice.ErrOwnersNotFound]),
// no associations were found ([imageservice.ErrAssociationsNotFound]) or a database/internal error occurs.
func (s *service) DeleteImageBatch(ctx context.Context, req *imagemodel.DeleteBatchRequst) (int, error) {
	ownerRepoAdapter := newPhysicalGoodOwnerRepoAdapter(s.PhysicalGoodRepo)
	return s.ImageSvc.DeleteImageBatch(ctx, req, ownerRepoAdapter)
}

// Delete performs a soft-delete of a physical good and its related product record.
// It also unpublishes both records, meaning they must be manually published again after restoration.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.PhysicalGoodRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPhysicalGoodRepo := s.PhysicalGoodRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		// Check if the record exists first (including unpublished, but not soft-deleted)
		if _, err := txPhysicalGoodRepo.GetWithUnpublished(ctx, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return fmt.Errorf("failed to retrieve physical good: %w", err)
		}

		// Unpublish all instances
		if _, err := txPhysicalGoodRepo.SetInStock(ctx, id, false); err != nil {
			return fmt.Errorf("failed to unpublish physical good: %w", err)
		}
		ra, err := txProductRepo.SetInStockByDetailsID(ctx, id, false)
		if err != nil {
			return fmt.Errorf("failed to unpublish physical good product: %w", err)
		}
		if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		// Delete
		if _, err = txPhysicalGoodRepo.Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete physical good: %w", err)
		}
		if _, err = txProductRepo.DeleteByDetailsID(ctx, id); err != nil {
			return fmt.Errorf("failed to delete physical good product: %w", err)
		}
		return nil
	})
}

// DeletePermanent performs a complete delete of a physical good and its related product record.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) DeletePermanent(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.PhysicalGoodRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPhysicalGoodRepo := s.PhysicalGoodRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)

		ra, err := txPhysicalGoodRepo.DeletePermanent(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete physical good: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		ra, err = txProductRepo.DeletePermanentByDetailsID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete physical good product: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil
	})
}

// Restore performs a restore of a physical good and its related product record.
// Physical good and its related product record are not being published. This should be
// done manually.
//
// Returns an error if the ID is invalid (ErrInvalidArgument), the records are not found (ErrNotFound),
// or a database/internal error occurs.
func (s *service) Restore(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}
	return s.PhysicalGoodRepo.DB().Transaction(func(tx *gorm.DB) error {
		txPhysicalGoodRepo := s.PhysicalGoodRepo.WithTx(tx)
		txProductRepo := s.ProductRepo.WithTx(tx)
		ra, err := txPhysicalGoodRepo.Restore(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to restore physical good: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		ra, err = txProductRepo.RestoreByDetailsID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to restore physical good: %w", err)
		} else if ra == 0 {
			return fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil
	})
}
