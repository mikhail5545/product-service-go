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

package product

import (
	"context"

	productmodel "github.com/mikhail5545/product-service-go/internal/models/product"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../test/database/product_mock/repo_mock.go -package=product_mock github.com/mikhail5545/product-service-go/internal/database/product Repository

// Repository defines the interface for product data operations.
type Repository interface {
	// --- Only published and not soft-deleted ---

	// Get retrieves single published Product record from the database by it's ID.
	Get(ctx context.Context, id string) (*productmodel.Product, error)
	// GetByDetailsID retrieves single published Product record from the database by it's DetailsID.
	GetByDetailsID(ctx context.Context, detailsID string) (*productmodel.Product, error)
	// Select retrieves only specific fields from Product record in the database by it's ID.
	Select(ctx context.Context, id string, fields []string) (*productmodel.Product, error)
	// SelectByDetailsID retrieves only specific fields from Product record in the database by it's DetailsID.
	SelectByDetailsID(ctx context.Context, detailsID string, fields ...string) (*productmodel.Product, error)
	// SelectByDetailsID retrieves only specific fields from Product record in the database by it's DetailsID.
	SelectByDetailsIDs(ctx context.Context, detailsIDs []string, fields ...string) ([]productmodel.Product, error)
	// List retrieves all Product records from the database.
	List(ctx context.Context, limit, offset int) ([]productmodel.Product, error)
	// ListByDetailsType retrieves all Product records from the database that have specific DetailsType.
	ListByDetailsType(ctx context.Context, detailsType string, limit, offset int) ([]productmodel.Product, error)
	// ListByIDs retrieves all Product records from the database by a slice of IDs.
	ListByIDs(ctx context.Context, ids []string) ([]productmodel.Product, error)
	// SelectByDIs retrieves all Product specified fields from the database by a slice of IDs.
	SelectByIDs(ctx context.Context, ids []string, fields ...string) ([]productmodel.Product, error)
	// Count returns total amount of the Product records in the database
	Count(ctx context.Context) (int64, error)
	// CountByType returns the total amount of the Product records in the database that have specific DetailsType.
	CountByDetailsType(ctx context.Context, detailsType string) (int64, error)

	// --- With soft-deleted, if soft-deleted then also unpublished ---

	// GetWithDeleted retrieves single Product record including soft-deleted from the database by it's ID.
	GetWithDeleted(ctx context.Context, id string) (*productmodel.Product, error)
	// SelectByDetailsID retrieves only specific fields from soft-deleted Product record in the database by it's DetailsID.
	SelectWithDeletedByDetailsID(ctx context.Context, detailsID string, fields ...string) (*productmodel.Product, error)
	// SelectByDetailsID retrieves only specific fields from soft-deleted Product records in the database by DetailsIDs.
	SelectWithDeletedByDetailsIDs(ctx context.Context, detailsIDs []string, fields ...string) ([]productmodel.Product, error)
	// SelectByDIs retrieves all Product specified fields from the database by a slice of IDs including soft-deleted products.
	SelectWithDeletedByIDs(ctx context.Context, ids []string, fields ...string) ([]productmodel.Product, error)
	// GetWithDeletedByDetailsID retrieves single Product record from the database by it's DetailsID including soft-deleted ones.
	GetWithDeletedByDetailsID(ctx context.Context, detailsID string) (*productmodel.Product, error)
	// ListDeleted retrieves all soft-deleted Product records from the database.
	ListDeleted(ctx context.Context, limit, offset int) ([]productmodel.Product, error)
	// CountDeleted returns total amount of soft-deleted Product records in the database
	CountDeleted(ctx context.Context) (int64, error)

	// --- With unpublished, but not soft-deleted ---

	// GetWithUnpublished retrieves single Product record including unpublished from the database by it's ID.
	GetWithUnpublished(ctx context.Context, id string) (*productmodel.Product, error)
	// GetWithUnpublishedByDetailsID retrieves single unpublished Product record from the database by it's DetailsID.
	GetWithUnpublishedByDetailsID(ctx context.Context, detailsID string) (*productmodel.Product, error)
	// SelectWithUnpublishedByDetailsID retrieves only specific fields from unpublished Product record in the database by it's DetailsID.
	SelectWithUnpublishedByDetailsID(ctx context.Context, detailsID string, fields ...string) (*productmodel.Product, error)
	// SelectWithUnpublishedByIDs retrieves only specific fields from unpublished Product record in the database.
	SelectWithUnpublishedByIDs(ctx context.Context, ids []string, fields ...string) ([]productmodel.Product, error)
	// SelectWithUnpublishedByDetailsIDs retrieves only specific fields from unpublished Product record in the database by it's DetailsID.
	SelectWithUnpublishedByDetailsIDs(ctx context.Context, detailsIDs []string, fields ...string) ([]productmodel.Product, error)
	// CountUnpublished retrieves all unpublished Product records from the database.
	ListUnpublished(ctx context.Context, limit, offset int) ([]productmodel.Product, error)
	// CountUnpublished returns total amount of unpublished Product records in the database
	CountUnpublished(ctx context.Context) (int64, error)

	// -- Common --

	// Create creates new Product record in the database.
	Create(ctx context.Context, product *productmodel.Product) error
	// CreateBatch creates multiple new Product records in the database.
	CreateBatch(ctx context.Context, products ...*productmodel.Product) error
	// SetInStock sets new value for product's InStock field.
	SetInStock(ctx context.Context, id string, inStock bool) (int64, error)
	// SetInStockByDetailsID sets new value for product's InStock field by it's detailsID.
	SetInStockByDetailsID(ctx context.Context, detailsID string, inStock bool) (int64, error)
	// Update partually updates Product record using updates.
	Update(ctx context.Context, product *productmodel.Product, updates any) (int64, error)
	// Delete performs a soft-delete.
	Delete(ctx context.Context, id string) (int64, error)
	// DeleteByDetailsID performs a soft-delete of product records by details id.
	DeleteByDetailsID(ctx context.Context, detailsID string) (int64, error)
	// DeletePermanent removes product from the database completely.
	DeletePermanent(ctx context.Context, id string) (int64, error)
	// DeletePermanent removes products from the database completely by details id.
	DeletePermanentByDetailsID(ctx context.Context, detailsID string) (int64, error)
	// Restore restores soft-deleted product.
	Restore(ctx context.Context, id string) (int64, error)
	// Restore restores soft-deleted products by details id.
	RestoreByDetailsID(ctx context.Context, detailsID string) (int64, error)

	// DB returns the underlying gorm.DB instance.
	DB() *gorm.DB
	// WithTx returns a new repository instance with the given transaction.
	WithTx(tx *gorm.DB) Repository
}

// gormRepository holds *gorm.DB instance.
type gormRepository struct {
	db *gorm.DB
}

// New creates a new GORM-based product repository.
func New(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// DB returns the underlying gorm.DB instance.
func (r *gormRepository) DB() *gorm.DB {
	return r.db
}

// WithTx returns a new repository instance with the given transaction.
func (r *gormRepository) WithTx(tx *gorm.DB) Repository {
	return &gormRepository{db: tx}
}

// --- Only published and not soft-deleted ---

// Find retrieves a single product by it's ID.
func (r *gormRepository) Get(ctx context.Context, id string) (*productmodel.Product, error) {
	var product productmodel.Product
	err := r.db.WithContext(ctx).Where("in_stock = ?", true).First(&product, id).Error
	return &product, err
}

// GetByDetailsID retrieves single Product record from the database by it's DetailsID.
func (r *gormRepository) GetByDetailsID(ctx context.Context, detailsID string) (*productmodel.Product, error) {
	var product productmodel.Product
	err := r.db.WithContext(ctx).Where("in_stock = ?", true).Where("details_id = ?", detailsID).First(&product).Error
	return &product, err
}

// Select retrieves only specific fields from Product record in the database by it's ID.
func (r *gormRepository) Select(ctx context.Context, id string, fields []string) (*productmodel.Product, error) {
	var product productmodel.Product
	err := r.db.WithContext(ctx).Model(&productmodel.Product{}).Select("id", fields).Where("in_stock = ?", true).Where("id = ?", id).First(&product).Error
	return &product, err
}

// SelectByDetailsID retrieves only specific fields from Product record in the database by it's DetailsID.
func (r *gormRepository) SelectByDetailsID(ctx context.Context, detailsID string, fields ...string) (*productmodel.Product, error) {
	var product productmodel.Product
	err := r.db.WithContext(ctx).Model(&productmodel.Product{}).Select(fields).Where("in_stock = ?", true).Where("details_id = ?", detailsID).First(&product).Error
	return &product, err
}

// SelectByDetailsID retrieves only specific fields from Product records in the database by DetailsIDs.
func (r *gormRepository) SelectByDetailsIDs(ctx context.Context, detailsIDs []string, fields ...string) ([]productmodel.Product, error) {
	var products []productmodel.Product
	err := r.db.WithContext(ctx).Select(fields).Where("in_stock = ?", true).Where("details_id IN ?", detailsIDs).First(&products).Error
	return products, err
}

// List retrieves all Product records from the database.
func (r *gormRepository) List(ctx context.Context, limit, offset int) ([]productmodel.Product, error) {
	var products []productmodel.Product
	err := r.db.WithContext(ctx).Where("in_stock = ?", true).Limit(limit).Offset(offset).Order("created_at desc").Find(&products).Error
	return products, err
}

// ListByDetailsType retrieves all Product records from the database that have specific DetailsType.
func (r *gormRepository) ListByDetailsType(ctx context.Context, detailsType string, limit, offset int) ([]productmodel.Product, error) {
	var products []productmodel.Product
	err := r.db.WithContext(ctx).Where("in_stock = ?", true).Where("details_type = ?", detailsType).Limit(limit).Offset(offset).Order("created_at desc").Find(&products).Error
	return products, err
}

// ListByIDs retrieves all Product records from the database by a slice of IDs.
func (r *gormRepository) ListByIDs(ctx context.Context, ids []string) ([]productmodel.Product, error) {
	var products []productmodel.Product
	err := r.db.WithContext(ctx).Where("in_stock = ?", true).Where("id IN ?", ids).Find(&products).Error
	return products, err
}

// SelectByDIs retrieves all Product specified fields from the database by a slice of IDs.
func (r *gormRepository) SelectByIDs(ctx context.Context, ids []string, fields ...string) ([]productmodel.Product, error) {
	var products []productmodel.Product
	err := r.db.WithContext(ctx).Select(fields).Where("in_stock = ?", true).Where("id IN ?", ids).Find(&products).Error
	return products, err
}

// Count returns total amount of the Product records in the database
func (r *gormRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&productmodel.Product{}).Where("in_stock = ?", true).Count(&count).Error
	return count, err
}

// CountByType returns the total amount of the Product records in the database that have specific DetailsType.
func (r *gormRepository) CountByDetailsType(ctx context.Context, detailsType string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&productmodel.Product{}).Where("in_stock = ?", true).Where("details_type = ?", detailsType).Count(&count).Error
	return count, err
}

// --- With soft-deleted, if soft-deleted then also unpublished ---

// Get retrieves single Product record including soft-deleted from the database by it's ID.
func (r *gormRepository) GetWithDeleted(ctx context.Context, id string) (*productmodel.Product, error) {
	var product productmodel.Product
	err := r.db.WithContext(ctx).Unscoped().First(&product, id).Error
	return &product, err
}

// GetWithDeletedByDetailsID retrieves single Product record from the database by it's DetailsID including soft-deleted ones.
func (r *gormRepository) GetWithDeletedByDetailsID(ctx context.Context, detailsID string) (*productmodel.Product, error) {
	var product productmodel.Product
	err := r.db.WithContext(ctx).Unscoped().Where("details_id = ?", detailsID).First(&product).Error
	return &product, err
}

// SelectByDetailsID retrieves only specific fields from soft-deleted Product record in the database by it's DetailsID.
func (r *gormRepository) SelectWithDeletedByDetailsID(ctx context.Context, detailsID string, fields ...string) (*productmodel.Product, error) {
	var product productmodel.Product
	err := r.db.WithContext(ctx).Unscoped().Model(&productmodel.Product{}).Select(fields).Where("details_id = ?", detailsID).First(&product).Error
	return &product, err
}

// SelectByDetailsID retrieves only specific fields from soft-deleted Product records in the database by DetailsIDs.
func (r *gormRepository) SelectWithDeletedByDetailsIDs(ctx context.Context, detailsIDs []string, fields ...string) ([]productmodel.Product, error) {
	var products []productmodel.Product
	err := r.db.WithContext(ctx).Unscoped().Model(&productmodel.Product{}).Select(fields).Where("details_id IN ?", detailsIDs).Find(products).Error
	return products, err
}

// SelectByDIs retrieves all Product specified fields from the database by a slice of IDs including soft-deleted products.
func (r *gormRepository) SelectWithDeletedByIDs(ctx context.Context, ids []string, fields ...string) ([]productmodel.Product, error) {
	var products []productmodel.Product
	err := r.db.WithContext(ctx).Unscoped().Select(fields).Where("id IN ?", ids).Find(&products).Error
	return products, err
}

// ListDeleted retrieves all soft-deleted Product records from the database.
func (r *gormRepository) ListDeleted(ctx context.Context, limit, offset int) ([]productmodel.Product, error) {
	var products []productmodel.Product
	err := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Limit(limit).Offset(offset).Order("deleted_at desc").Find(&products).Error
	return products, err
}

// CountDeleted returns total amount of soft-deleted Product records in the database
func (r *gormRepository) CountDeleted(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Unscoped().Model(&productmodel.Product{}).Where("deleted_at IS NOT NULL").Count(&count).Error
	return count, err
}

// --- With unpublished, but not soft-deleted ---

// GetWithUnpublished retrieves single Product record including unpublished from the database by it's ID.
func (r *gormRepository) GetWithUnpublished(ctx context.Context, id string) (*productmodel.Product, error) {
	var product productmodel.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	return &product, err
}

// GetWithUnpublishedByDetailsID retrieves single unpublished Product record from the database by it's DetailsID.
func (r *gormRepository) GetWithUnpublishedByDetailsID(ctx context.Context, detailsID string) (*productmodel.Product, error) {
	var product productmodel.Product
	err := r.db.WithContext(ctx).Where("details_id = ?", detailsID).First(&product).Error
	return &product, err
}

// SelectWithUnpublishedByIDs retrieves only specific fields from unpublished Product record in the database.
func (r *gormRepository) SelectWithUnpublishedByIDs(ctx context.Context, ids []string, fields ...string) ([]productmodel.Product, error) {
	var products []productmodel.Product
	err := r.db.WithContext(ctx).Model(&productmodel.Product{}).Select(fields).Where("id IN ?", ids).Find(&products).Error
	return products, err
}

// SelectWithUnpublishedByDetailsID retrieves only specific fields from unpublished Product record in the database by it's DetailsID.
func (r *gormRepository) SelectWithUnpublishedByDetailsID(ctx context.Context, detailsID string, fields ...string) (*productmodel.Product, error) {
	var product productmodel.Product
	err := r.db.WithContext(ctx).Select(fields).Where("details_id = ?", detailsID).First(&product).Error
	return &product, err
}

// SelectWithUnpublishedByDetailsIDs retrieves only specific fields from unpublished Product record in the database by it's DetailsID.
func (r *gormRepository) SelectWithUnpublishedByDetailsIDs(ctx context.Context, detailsIDs []string, fields ...string) ([]productmodel.Product, error) {
	var products []productmodel.Product
	err := r.db.WithContext(ctx).Model(&productmodel.Product{}).Select(fields).Where("details_id IN ?", detailsIDs).Find(&products).Error
	return products, err
}

// CountUnpublished retrieves all unpublished Product records from the database.
func (r *gormRepository) ListUnpublished(ctx context.Context, limit, offset int) ([]productmodel.Product, error) {
	var products []productmodel.Product
	err := r.db.WithContext(ctx).
		Model(&productmodel.Product{}).
		Where("in_stock = ?", false).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&products).Error
	return products, err
}

// CountUnpublished returns total amount of unpublished Product records in the database
func (r *gormRepository) CountUnpublished(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&productmodel.Product{}).
		Where("in_stock = ?", false).
		Count(&count).Error
	return count, err
}

// --- Common ---

// Create creates new Product record in the database.
func (r *gormRepository) Create(ctx context.Context, product *productmodel.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// CreateBatch creates multiple new Product records in the database.
func (r *gormRepository) CreateBatch(ctx context.Context, products ...*productmodel.Product) error {
	return r.db.WithContext(ctx).Create(&products).Error
}

// SetInStock sets new value for product's InStock field.
func (r *gormRepository) SetInStock(ctx context.Context, id string, inStock bool) (int64, error) {
	res := r.db.WithContext(ctx).Model(&productmodel.Product{}).Where("id = ?", id).Update("in_stock", inStock)
	return res.RowsAffected, res.Error
}

// SetInStockByDetailsID sets new value for product's InStock field by it's detailsID.
func (r *gormRepository) SetInStockByDetailsID(ctx context.Context, detailsID string, inStock bool) (int64, error) {
	res := r.db.WithContext(ctx).Model(&productmodel.Product{}).Where("details_id = ?", detailsID).Update("in_stock", inStock)
	return res.RowsAffected, res.Error
}

// Update partually updates Product record using updates.
func (r *gormRepository) Update(ctx context.Context, product *productmodel.Product, updates any) (int64, error) {
	res := r.db.WithContext(ctx).Model(product).Updates(updates)
	return res.RowsAffected, res.Error
}

// Delete performs a soft-delete.
func (r *gormRepository) Delete(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Delete(&productmodel.Product{}, id)
	return res.RowsAffected, res.Error
}

// DeleteByDetailsID performs a soft-delete of product records by details id.
func (r *gormRepository) DeleteByDetailsID(ctx context.Context, detailsID string) (int64, error) {
	res := r.db.WithContext(ctx).Where("details_id = ?", detailsID).Delete(&productmodel.Product{})
	return res.RowsAffected, res.Error
}

// DeletePermanent removes product from the database completely.
func (r *gormRepository) DeletePermanent(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Delete(&productmodel.Product{}, id)
	return res.RowsAffected, res.Error
}

// DeletePermanent removes products from the database completely by details id.
func (r *gormRepository) DeletePermanentByDetailsID(ctx context.Context, detailsID string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Where("details_id = ?", detailsID).Delete(&productmodel.Product{})
	return res.RowsAffected, res.Error
}

// Restore restores soft-deleted product.
func (r *gormRepository) Restore(ctx context.Context, id string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Model(&productmodel.Product{}).Where("id = ?", id).Update("deleted_at", nil)
	return res.RowsAffected, res.Error
}

// Restore restores soft-deleted products by details id.
func (r *gormRepository) RestoreByDetailsID(ctx context.Context, detailsID string) (int64, error) {
	res := r.db.WithContext(ctx).Unscoped().Model(&productmodel.Product{}).Where("id = ?", detailsID).Update("deleted_at", nil)
	return res.RowsAffected, res.Error
}
