package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"vitainmove.com/product-service-go/internal/database"
	"vitainmove.com/product-service-go/internal/models"
)

type ProductService struct {
	Repo database.ProductRepository
}

type ProductServiceError struct {
	Msg  string
	Err  error
	Code int
}

func (e *ProductServiceError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func (e *ProductServiceError) Unwrap() error {
	return e.Err
}

func (e *ProductServiceError) GetCode() int {
	return e.Code
}

func NewProductService(pr database.ProductRepository) *ProductService {
	return &ProductService{Repo: pr}
}

func (s *ProductService) GetProducts(ctx context.Context, limit, offset int) ([]models.Product, int64, error) {
	products, err := s.Repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, &ProductServiceError{
			Msg:  "Failed to get products",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	total, err := s.Repo.Count(ctx)
	if err != nil {
		return nil, 0, &ProductServiceError{
			Msg:  "Failed to get products count",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return products, total, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	product, err := s.Repo.Find(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &ProductServiceError{
				Msg:  "Product not found",
				Err:  err,
				Code: http.StatusNotFound,
			}
		}
		return nil, &ProductServiceError{
			Msg:  "Failed to get product",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return product, nil
}

func (s *ProductService) GetProductByType(ctx context.Context, productType string, limit, offset int) ([]models.Product, int64, error) {
	products, err := s.Repo.FindByType(ctx, productType, limit, offset)
	if err != nil {
		return nil, 0, &ProductServiceError{
			Msg:  "Failed to get product",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	total, err := s.Repo.CountByType(ctx, productType)
	if err != nil {
		return nil, 0, &ProductServiceError{
			Msg:  "Failed to count products",
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return products, total, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, req *models.AddProductRequest) (*models.Product, error) {
	var product *models.Product

	err := s.Repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.Repo.WithTx(tx)

		if req.Name == "" {
			return &ProductServiceError{Msg: "product name cannot be empty", Err: nil, Code: http.StatusBadRequest}
		}
		if req.Description == "" {
			return &ProductServiceError{Msg: "product description cannot be empty", Err: nil, Code: http.StatusBadRequest}
		}
		if req.Amount < 0 {
			return &ProductServiceError{Msg: "amount cannot be negative", Err: nil, Code: http.StatusBadRequest}
		}
		if req.Price <= 0 {
			return &ProductServiceError{Msg: "price cannot be negative or null", Err: nil, Code: http.StatusBadRequest}
		}

		product = &models.Product{
			ID:               uuid.New().String(),
			Name:             req.Name,
			Description:      req.Description,
			Price:            req.Price,
			Amount:           req.Amount,
			ShippingRequired: req.ShippingRequired,
			ProductType:      "physical",
		}

		if err := txRepo.Create(ctx, product); err != nil {
			return &ProductServiceError{Msg: "failed to create product", Err: err, Code: http.StatusInternalServerError}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *models.EditProductRequest, id string) (*models.Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, &ProductServiceError{Msg: "invalid product id", Err: err, Code: http.StatusBadRequest}
	}

	var product *models.Product

	err := s.Repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.Repo.WithTx(tx)

		var findErr error
		product, findErr = txRepo.Find(ctx, id)
		if findErr != nil {
			if errors.Is(findErr, gorm.ErrRecordNotFound) {
				return &ProductServiceError{Msg: "product not found", Err: findErr, Code: http.StatusNotFound}
			}
			return &ProductServiceError{Msg: "failed to find product", Err: findErr, Code: http.StatusInternalServerError}
		}

		var updated bool
		if req.Name != "" && req.Name != product.Name {
			product.Name = req.Name
			updated = true
		}
		if req.Price != 0 && req.Price != product.Price {
			product.Price = req.Price
			updated = true
		}
		if req.Amount >= 0 && req.Amount != product.Amount { // amount can be null
			product.Amount = req.Amount
			updated = true
		}
		// Only update if a field has changed to avoid unnecessary DB calls.
		if updated {
			product.UpdatedAt = time.Now()
			if err := txRepo.Update(ctx, product); err != nil {
				return &ProductServiceError{Msg: "failed to update product", Err: err, Code: http.StatusInternalServerError}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return product, nil
}
