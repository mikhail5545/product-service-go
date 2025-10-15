package services

import (
	"context"
	"fmt"

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
			Code: 500,
		}
	}

	total, err := s.Repo.Count(ctx)
	if err != nil {
		return nil, 0, &ProductServiceError{
			Msg:  "Failed to get products count",
			Err:  err,
			Code: 500,
		}
	}
	return products, total, nil
}
