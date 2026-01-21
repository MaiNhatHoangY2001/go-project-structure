package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/repository"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/usecase"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/config"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/database"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/logger"
	pb "github.com/MaiNhatHoangY2001/go-project-structure/proto/todo"
	grpcServer "github.com/MaiNhatHoangY2001/go-project-structure/services/todo-service/grpc"
	httpServer "github.com/MaiNhatHoangY2001/go-project-structure/services/todo-service/http"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("configs/app.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	if err := logger.InitLogger(cfg.Logger.Level, cfg.Logger.Encoding); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Connect to MongoDB
	mongodb, err := database.NewMongoDB(cfg.Database.MongoDB.URI, cfg.Database.MongoDB.Database, cfg.Database.MongoDB.Timeout)
	if err != nil {
		logger.Log.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		if err := mongodb.Disconnect(); err != nil {
			logger.Log.Error("Failed to disconnect MongoDB", zap.Error(err))
		}
	}()

	// Initialize repositories
	todoRepo := repository.NewTodoRepository(mongodb.Database)

	// Initialize use cases
	todoUsecase := usecase.NewTodoUsecase(todoRepo)

	// Initialize Auth gRPC client
	authClient, err := grpcServer.NewAuthClient(cfg.GRPC.AuthService.Host, cfg.GRPC.AuthService.Port)
	if err != nil {
		logger.Log.Fatal("Failed to connect to Auth Service", zap.Error(err))
	}
	defer authClient.Close()

	// Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC.TodoService.Port))
		if err != nil {
			logger.Log.Fatal("Failed to listen", zap.Error(err))
		}

		s := grpc.NewServer()
		pb.RegisterTodoServiceServer(s, grpcServer.NewTodoServer(todoUsecase))

		logger.Log.Info(fmt.Sprintf("Todo gRPC Service started on port %s", cfg.GRPC.TodoService.Port))
		if err := s.Serve(lis); err != nil {
			logger.Log.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	// Create local gRPC client for HTTP handlers with retry logic
	todoGrpcClient, conn, err := createLocalTodoClient(cfg.GRPC.TodoService.Host, cfg.GRPC.TodoService.Port)
	if err != nil {
		logger.Log.Fatal("Failed to create local todo gRPC client", zap.Error(err))
	}
	defer conn.Close()

	// Setup and start HTTP server
	router := httpServer.SetupRouter(todoGrpcClient, authClient)
	
	logger.Log.Info(fmt.Sprintf("Todo HTTP Service started on port %s", cfg.Server.Port))
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		logger.Log.Fatal("Failed to start HTTP server", zap.Error(err))
	}
}

func createLocalTodoClient(host, port string) (pb.TodoServiceClient, *grpc.ClientConn, error) {
	address := fmt.Sprintf("%s:%s", host, port)
	
	// Retry logic with exponential backoff
	maxRetries := 5
	baseDelay := 200 * time.Millisecond
	
	var conn *grpc.ClientConn
	var err error
	
	for i := 0; i < maxRetries; i++ {
		conn, err = grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			return pb.NewTodoServiceClient(conn), conn, nil
		}
		
		if i < maxRetries-1 {
			delay := baseDelay * time.Duration(1<<uint(i))
			logger.Log.Warn(fmt.Sprintf("Failed to connect to local todo service, retrying in %v", delay), zap.Error(err))
			time.Sleep(delay)
		}
	}
	
	return nil, nil, fmt.Errorf("failed to connect to local todo service after %d retries: %w", maxRetries, err)
}
