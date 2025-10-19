// vitainmove.com/product-service-go
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
	"google.golang.org/grpc"
	"vitainmove.com/product-service-go/internal/database"
	"vitainmove.com/product-service-go/internal/routers"
	"vitainmove.com/product-service-go/internal/server"
	"vitainmove.com/product-service-go/internal/services"
	coursepb "vitainmove.com/product-service-go/proto/course/v0"
	productpb "vitainmove.com/product-service-go/proto/product/v0"
	trainingsessionpb "vitainmove.com/product-service-go/proto/training_session/v0"
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

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName)

	db, err := database.NewPostgresDB(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established.")

	// Create an instance of required repositories
	productRepo := database.NewProductRepository(db)
	trainingSessionRepo := database.NewTrainingSessionRepository(db)
	courseRepo := database.NewCourseRepository(db)

	// Create an instance of required services
	productService := services.NewProductService(productRepo)
	trainingSessionService := services.NewTrainingSessionService(trainingSessionRepo, productRepo)
	courseService := services.NewCourseService(courseRepo, productRepo)

	// --- Start gRPC server ---
	go func() {
		lis, err := net.Listen("tcp", grpcListenAddr)
		if err != nil {
			log.Fatalf("Failed to listen on %s: %v", grpcListenAddr, err)
		}

		grpcServer := startGRPCServer(productService, trainingSessionService, courseService, lis)
		log.Printf("gRPC server listening on %s", grpcListenAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// --- Start HTTP server ---
	e := echo.New()

	// Register HTTP handlers

	routers.SetupRouter(e, productService, trainingSessionService, courseService)
	httpListenAddr := fmt.Sprintf(":%d", httpPort)
	if err := e.Start(httpListenAddr); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func startGRPCServer(productService *services.ProductService, tsService *services.TrainingSessionService, courseService *services.CourseService, lis net.Listener) *grpc.Server {
	grpcServer := grpc.NewServer()

	// Create instances of required gRPC server implementations
	productServer := server.NewProductServer(productService)
	trainingSessionServer := server.NewTrainingSessionServer(tsService)
	courseServer := server.NewCourseServer(courseService)

	// Register gRPC services with the server
	productpb.RegisterProductServiceServer(grpcServer, productServer)
	trainingsessionpb.RegisterTrainingSessionServiceServer(grpcServer, trainingSessionServer)
	coursepb.RegisterCourseServiceServer(grpcServer, courseServer)

	return grpcServer
}
