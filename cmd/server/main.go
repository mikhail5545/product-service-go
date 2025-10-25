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

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	mediaclient "github.com/mikhail5545/product-service-go/internal/clients/mediaservice"
	"github.com/mikhail5545/product-service-go/internal/database"
	courserepo "github.com/mikhail5545/product-service-go/internal/database/course"
	cprepo "github.com/mikhail5545/product-service-go/internal/database/course_part"
	productrepo "github.com/mikhail5545/product-service-go/internal/database/product"
	seminarrepo "github.com/mikhail5545/product-service-go/internal/database/seminar"
	tsrepo "github.com/mikhail5545/product-service-go/internal/database/training_session"
	"github.com/mikhail5545/product-service-go/internal/routers"
	courseservice "github.com/mikhail5545/product-service-go/internal/services/course"
	cpservice "github.com/mikhail5545/product-service-go/internal/services/course_part"
	productservice "github.com/mikhail5545/product-service-go/internal/services/product"
	seminarservice "github.com/mikhail5545/product-service-go/internal/services/seminar"
	tsservice "github.com/mikhail5545/product-service-go/internal/services/training_session"
	"google.golang.org/grpc"
)

func main() {
	const grpcPort = 50052
	const httpPort = 8082
	grpcListenAddr := fmt.Sprintf(":%d", grpcPort)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBHost := os.Getenv("POSTGRES_HOST")
	DBPort := os.Getenv("POSTGRES_PORT")
	DBUser := os.Getenv("POSTGRES_USER")
	DBPassword := os.Getenv("POSTGRES_PASSWORD")
	DBName := os.Getenv("POSTGRES_DB")

	mediaServiceAddr := os.Getenv("MEDIA_SERVICE_ADDR")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName)

	db, err := database.NewPostgresDB(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	// --- Create gRPC Clients ---
	mediaSvcClient, err := mediaclient.NewClient(context.Background(), mediaServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create media service client: %v", err)
	}
	defer mediaSvcClient.Close()

	// Create an instance of required repositories
	productRepo := productrepo.New(db)
	trainingSessionRepo := tsrepo.New(db)
	courseRepo := courserepo.New(db)
	seminarRepo := seminarrepo.New(db)
	coursePartRepo := cprepo.New(db)

	// Create an instance of required services
	productService := productservice.New(productRepo)
	trainingSessionService := tsservice.New(trainingSessionRepo, productRepo)
	courseService := courseservice.New(courseRepo, productRepo, mediaSvcClient)
	seminarService := seminarservice.New(seminarRepo, productRepo)
	coursePartService := cpservice.New(coursePartRepo, courseRepo)

	// --- Start gRPC server ---
	go func() {
		lis, err := net.Listen("tcp", grpcListenAddr)
		if err != nil {
			log.Fatalf("Failed to listen on %s: %v", grpcListenAddr, err)
		}

		grpcServer := startGRPCServer(productService, trainingSessionService, courseService, seminarService, lis)
		log.Printf("gRPC server listening on %s", grpcListenAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// --- Start HTTP server ---
	e := echo.New()

	// Register HTTP handlers

	routers.Setup(e, productService, coursePartService, trainingSessionService, courseService, seminarService)
	httpListenAddr := fmt.Sprintf(":%d", httpPort)
	if err := e.Start(httpListenAddr); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func startGRPCServer(
	productService *productservice.Service,
	tsService *tsservice.Service,
	courseService *courseservice.Service,
	seminarService *seminarservice.Service,
	lis net.Listener,
) *grpc.Server {
	grpcServer := grpc.NewServer()

	// // Create instances of required gRPC server implementations
	// productServer := server.NewProductServer(productService)
	// trainingSessionServer := server.NewTrainingSessionServer(tsService)
	// courseServer := server.NewCourseServer(courseService)
	// seminarServer := server.NewSeminarServer(seminarService)

	// // Register gRPC services with the server
	// productpb.RegisterProductServiceServer(grpcServer, productServer)
	// trainingsessionpb.RegisterTrainingSessionServiceServer(grpcServer, trainingSessionServer)
	// coursepb.RegisterCourseServiceServer(grpcServer, courseServer)
	// seminarpb.RegisterSeminarServiceServer(grpcServer, seminarServer)

	return grpcServer
}
