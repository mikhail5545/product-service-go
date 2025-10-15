package database

import (
	"context"

	"gorm.io/gorm"
	"vitainmove.com/product-service-go/internal/models"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	// Read operations
	Find(ctx context.Context, id string) (*models.Product, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.Product, error)
	FindByType(ctx context.Context, productType string, limit, offset int) ([]models.Product, error)
	Count(ctx context.Context) (int64, error)

	// Write operations
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int64) error

	DB() *gorm.DB
	WithTx(tx *gorm.DB) ProductRepository
}

type gormProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new GORM-based product repository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &gormProductRepository{db: db}
}

// DB returns the underlying gorm.DB instance.
func (r *gormProductRepository) DB() *gorm.DB {
	return r.db
}

// WithTx returns a new repository instance with the given transaction.
func (r *gormProductRepository) WithTx(tx *gorm.DB) ProductRepository {
	return &gormProductRepository{db: tx}
}

// Find retrieves a single product by it's ID.
func (r *gormProductRepository) Find(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	return &product, err
}

func (r *gormProductRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at desc").Find(&products).Error
	return products, err
}

func (r *gormProductRepository) FindByType(ctx context.Context, productType string, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).Where("type = ?", productType).Limit(limit).Offset(offset).Order("created_at desc").Find(&products).Error
	return products, err
}

func (r *gormProductRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&count).Error
	return count, err
}

func (r *gormProductRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *gormProductRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *gormProductRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}
