package router

import (
	"net/http"

	"github.com/MaiNhatHoangY2001/go-project-structure/services/api-gateway/grpc/clients"
	"github.com/MaiNhatHoangY2001/go-project-structure/services/api-gateway/http/handler"
	"github.com/MaiNhatHoangY2001/go-project-structure/services/api-gateway/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(authClient *clients.AuthClient, todoClient *clients.TodoClient) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "api-gateway"})
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (no authentication required)
		authHandler := handler.NewAuthHandler(authClient)
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
		}

		// Todo routes (authentication required)
		todoHandler := handler.NewTodoHandler(todoClient)
		todos := v1.Group("/todos")
		todos.Use(middleware.AuthMiddleware(authClient))
		{
			todos.POST("", todoHandler.CreateTodo)
			todos.GET("", todoHandler.ListTodos)
			todos.GET("/:id", todoHandler.GetTodo)
			todos.PUT("/:id", todoHandler.UpdateTodo)
			todos.DELETE("/:id", todoHandler.DeleteTodo)
		}
	}

	return r
}
