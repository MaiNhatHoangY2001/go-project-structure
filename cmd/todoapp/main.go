// Package main Todo App
// @title Todo App API
// @version 1.0
// @description A Todo application with JWT authentication and MongoDB
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/MaiNhatHoangY2001/go-project-structure/api/docs"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/delivery/http/handler"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/delivery/http/router"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/repository"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/usecase"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/config"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/database"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("configs/app.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	if err := logger.InitLogger(cfg.Logger.Level, cfg.Logger.Encoding); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Log.Info("Starting Todo App")

	// Connect to MongoDB
	mongoDB, err := database.NewMongoDB(
		cfg.Database.MongoDB.URI,
		cfg.Database.MongoDB.Database,
		cfg.Database.MongoDB.Timeout,
	)
	if err != nil {
		logger.Log.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		if err := mongoDB.Disconnect(); err != nil {
			logger.Log.Error("Failed to disconnect MongoDB", zap.Error(err))
		}
	}()

	logger.Log.Info("Connected to MongoDB")

	// Initialize repositories
	userRepo := repository.NewUserRepository(mongoDB.Database)
	todoRepo := repository.NewTodoRepository(mongoDB.Database)

	// Initialize use cases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWT.Secret, cfg.JWT.Expiration)
	todoUsecase := usecase.NewTodoUsecase(todoRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	todoHandler := handler.NewTodoHandler(todoUsecase)

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Setup router
	r := gin.Default()

	// Setup routes
	router.SetupRoutes(r, authUsecase, authHandler, todoHandler)

	// Start server
	srv := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Log.Info("Server starting", zap.String("address", srv))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(srv); err != nil {
			logger.Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	<-quit
	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Disconnect MongoDB
	if err := mongoDB.Disconnect(); err != nil {
		logger.Log.Error("Error disconnecting MongoDB", zap.Error(err))
	}

	<-ctx.Done()
	logger.Log.Info("Server stopped")
}
