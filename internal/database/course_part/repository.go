package coursepart

import (
	"context"

	"github.com/mikhail5545/product-service-go/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	// Find retrieves single Course Part from the database.
	Get(ctx context.Context, id string) (*models.CoursePart, error)
	// List retrieves Course Parts data from the database.
	List(ctx context.Context, courseID string, limit, offset int) ([]models.CoursePart, error)
	// Count counts the number of Course Parts related to Course with provided ID records in the database.
	Count(ctx context.Context, courseID string) (int64, error)
	// Create creates a new CoursePart record in the database.
	Create(ctx context.Context, coursePart *models.CoursePart) error
	// Update updates CoursePart record.
	Update(ctx context.Context, coursePart *models.CoursePart) error
	// Delete deletes CoursePart record.
	Delete(ctx context.Context, id string) error

	DB() *gorm.DB
	WithTx(tx *gorm.DB) Repository
}

type gormRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &gormRepository{
		db: db,
	}
}

func (r *gormRepository) DB() *gorm.DB {
	return r.db
}

func (r *gormRepository) WithTx(tx *gorm.DB) Repository {
	return &gormRepository{
		db: tx,
	}
}

func (r *gormRepository) Get(ctx context.Context, id string) (*models.CoursePart, error) {
	var coursePart models.CoursePart
	err := r.db.WithContext(ctx).First(&coursePart, "id = ?", id).Error
	return &coursePart, err
}

func (r *gormRepository) List(ctx context.Context, courseID string, limit, offset int) ([]models.CoursePart, error) {
	var courseParts []models.CoursePart
	err := r.db.WithContext(ctx).Order("created_at desc").Limit(limit).Offset(offset).Find(&courseParts, "course_id = ?", courseID).Error
	return courseParts, err
}

func (r *gormRepository) Count(ctx context.Context, courseID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.CoursePart{}).Where("course_id = ?", courseID).Count(&count).Error
	return count, err
}

func (r *gormRepository) Create(ctx context.Context, coursePart *models.CoursePart) error {
	return r.db.WithContext(ctx).Create(coursePart).Error
}

func (r *gormRepository) Update(ctx context.Context, coursePart *models.CoursePart) error {
	return r.db.WithContext(ctx).Save(coursePart).Error
}

func (r *gormRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.CoursePart{}, id).Error
}
