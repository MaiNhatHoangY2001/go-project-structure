package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/config"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/logger"
	"github.com/MaiNhatHoangY2001/go-project-structure/services/api-gateway/grpc/clients"
	"github.com/MaiNhatHoangY2001/go-project-structure/services/api-gateway/http/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := logger.InitLogger("info", "json"); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := logger.Log.Sync(); err != nil {
			fmt.Printf("Error syncing logger: %v\n", err)
		}
	}()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal("Failed to load config", zap.Error(err))
	}

	// Initialize gRPC clients
	authClient, err := clients.NewAuthClient(cfg.AuthService.Host, cfg.AuthService.Port)
	if err != nil {
		logger.Log.Fatal("Failed to connect to Auth Service", zap.Error(err))
	}
	defer func() {
		if err := authClient.Close(); err != nil {
			logger.Log.Error("Failed to close Auth client", zap.Error(err))
		}
	}()

	todoClient, err := clients.NewTodoClient(cfg.TodoService.Host, cfg.TodoService.Port)
	if err != nil {
		logger.Log.Fatal("Failed to connect to Todo Service", zap.Error(err))
	}
	defer func() {
		if err := todoClient.Close(); err != nil {
			logger.Log.Error("Failed to close Todo client", zap.Error(err))
		}
	}()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := router.SetupRouter(authClient, todoClient)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Log.Info("Starting API Gateway", 
			zap.String("port", cfg.Server.Port),
			zap.String("auth_service", fmt.Sprintf("%s:%s", cfg.AuthService.Host, cfg.AuthService.Port)),
			zap.String("todo_service", fmt.Sprintf("%s:%s", cfg.TodoService.Host, cfg.TodoService.Port)),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down API Gateway...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("API Gateway stopped")
}
