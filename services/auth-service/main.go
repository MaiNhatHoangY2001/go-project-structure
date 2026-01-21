package main

import (
	"fmt"
	"log"
	"net"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/repository"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/usecase"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/config"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/database"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/logger"
	pb "github.com/MaiNhatHoangY2001/go-project-structure/proto/auth"
	grpcServer "github.com/MaiNhatHoangY2001/go-project-structure/services/auth-service/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	userRepo := repository.NewUserRepository(mongodb.Database)

	// Initialize use cases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWT.Secret, cfg.JWT.Expiration)

	// Create gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC.AuthService.Port))
	if err != nil {
		logger.Log.Fatal("Failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, grpcServer.NewAuthServer(authUsecase))

	logger.Log.Info(fmt.Sprintf("Auth Service started on port %s", cfg.GRPC.AuthService.Port))
	if err := s.Serve(lis); err != nil {
		logger.Log.Fatal("Failed to serve", zap.Error(err))
	}
}
