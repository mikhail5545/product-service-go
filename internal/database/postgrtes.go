package database

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"vitainmove.com/product-service-go/internal/models"
)

func NewPostgresDB(ctx context.Context, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Product{})
	if err != nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		return nil, err
	}

	return db, nil
}
