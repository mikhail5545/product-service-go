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

// Package app bootstraps and runs the application.
package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/mikhail5545/product-service-go/internal/database"
	courserepo "github.com/mikhail5545/product-service-go/internal/database/course"
	cprepo "github.com/mikhail5545/product-service-go/internal/database/course_part"
	imagerepo "github.com/mikhail5545/product-service-go/internal/database/image"
	physicalgoodrepo "github.com/mikhail5545/product-service-go/internal/database/physical_good"
	productrepo "github.com/mikhail5545/product-service-go/internal/database/product"
	seminarrepo "github.com/mikhail5545/product-service-go/internal/database/seminar"
	tsrepo "github.com/mikhail5545/product-service-go/internal/database/training_session"
	"github.com/mikhail5545/product-service-go/internal/routers"
	courseserver "github.com/mikhail5545/product-service-go/internal/server/course"
	cpserver "github.com/mikhail5545/product-service-go/internal/server/course_part"
	physicalgoodserver "github.com/mikhail5545/product-service-go/internal/server/physical_good"
	productserver "github.com/mikhail5545/product-service-go/internal/server/product"
	seminarserver "github.com/mikhail5545/product-service-go/internal/server/seminar"
	tsserver "github.com/mikhail5545/product-service-go/internal/server/training_session"
	courseservice "github.com/mikhail5545/product-service-go/internal/services/course"
	cpservice "github.com/mikhail5545/product-service-go/internal/services/course_part"
	imageservice "github.com/mikhail5545/product-service-go/internal/services/image"
	physicalgoodservice "github.com/mikhail5545/product-service-go/internal/services/physical_good"
	productservice "github.com/mikhail5545/product-service-go/internal/services/product"
	seminarservice "github.com/mikhail5545/product-service-go/internal/services/seminar"
	tsservice "github.com/mikhail5545/product-service-go/internal/services/training_session"
	"google.golang.org/grpc"
)

const (
	grpcPort = 50052
	httpPort = 8082
)

// Run initializes and starts the application servers.
func Run(ctx context.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBHost := os.Getenv("POSTGRES_HOST")
	DBPort := os.Getenv("POSTGRES_PORT")
	DBUser := os.Getenv("POSTGRES_USER")
	DBPassword := os.Getenv("POSTGRES_PASSWORD")
	DBName := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName)

	db, err := database.NewPostgresDB(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	// Create an instance of required repositories
	productRepo := productrepo.New(db)
	trainingSessionRepo := tsrepo.New(db)
	courseRepo := courserepo.New(db)
	seminarRepo := seminarrepo.New(db)
	coursePartRepo := cprepo.New(db)
	physicalGoodRepo := physicalgoodrepo.New(db)
	imageRepo := imagerepo.New(db)

	// Create an instance of required services
	productService := productservice.New(productRepo)
	imageService := imageservice.New(imageRepo)
	trainingSessionService := tsservice.New(trainingSessionRepo, productRepo, imageService)
	courseService := courseservice.New(courseRepo, productRepo, coursePartRepo, imageService)
	seminarService := seminarservice.New(seminarRepo, productRepo, imageService)
	coursePartService := cpservice.New(coursePartRepo, courseRepo)
	physicalGoodService := physicalgoodservice.New(physicalGoodRepo, productRepo, imageService)

	// --- Start gRPC server ---
	go func() {
		grpcListenAddr := fmt.Sprintf(":%d", grpcPort)
		lis, err := net.Listen("tcp", grpcListenAddr)
		if err != nil {
			log.Fatalf("Failed to listen on %s: %v", grpcListenAddr, err)
		}

		grpcServer := grpc.NewServer()

		// --- Register gRPC services with the server ---
		courseserver.Register(grpcServer, courseService)
		tsserver.Register(grpcServer, trainingSessionService)
		cpserver.Register(grpcServer, coursePartService)
		seminarserver.Register(grpcServer, seminarService)
		productserver.Register(grpcServer, productService)
		physicalgoodserver.Register(grpcServer, physicalGoodService)

		log.Printf("gRPC server listening on %s", grpcListenAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// --- Start HTTP server ---
	e := echo.New()

	// Register HTTP handlers
	routers.Setup(e, productService, coursePartService, trainingSessionService, courseService, seminarService, physicalGoodService)
	httpListenAddr := fmt.Sprintf(":%d", httpPort)
	if err := e.Start(httpListenAddr); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
