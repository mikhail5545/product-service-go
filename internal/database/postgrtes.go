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

package database

import (
	"context"

	coursemodel "github.com/mikhail5545/product-service-go/internal/models/course"
	coursepartmodel "github.com/mikhail5545/product-service-go/internal/models/course_part"
	physicalgoodmodel "github.com/mikhail5545/product-service-go/internal/models/physical_good"
	productmodel "github.com/mikhail5545/product-service-go/internal/models/product"
	seminarmodel "github.com/mikhail5545/product-service-go/internal/models/seminar"
	trainingsessionmodel "github.com/mikhail5545/product-service-go/internal/models/training_session"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(ctx context.Context, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&productmodel.Product{},
		&trainingsessionmodel.TrainingSession{},
		&coursepartmodel.CoursePart{},
		&coursemodel.Course{},
		&trainingsessionmodel.TrainingSession{},
		&seminarmodel.Seminar{},
		&physicalgoodmodel.PhysicalGood{},
	)
	if err != nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		return nil, err
	}

	return db, nil
}
