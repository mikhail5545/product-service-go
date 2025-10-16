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
	productpb "vitainmove.com/product-service-go/proto/product/v0"
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

	// Create an instance of required services
	productService := services.NewProductService(productRepo)

	// --- Start gRPC server ---
	go func() {
		lis, err := net.Listen("tcp", grpcListenAddr)
		if err != nil {
			log.Fatalf("Failed to listen on %s: %v", grpcListenAddr, err)
		}

		grpcServer := startGRPCServer(productService, lis)
		log.Printf("gRPC server listening on %s", grpcListenAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// --- Start HTTP server ---
	e := echo.New()

	// Register HTTP handlers

	routers.SetupRouter(e, productService)
	httpListenAddr := fmt.Sprintf(":%d", httpPort)
	if err := e.Start(httpListenAddr); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func startGRPCServer(productService *services.ProductService, lis net.Listener) *grpc.Server {
	grpcServer := grpc.NewServer()

	// Create instances of required gRPC server implementations
	productServer := server.NewProductServer(productService)

	// Register gRPC services with the server
	productpb.RegisterProductServiceServer(grpcServer, productServer)

	return grpcServer
}
