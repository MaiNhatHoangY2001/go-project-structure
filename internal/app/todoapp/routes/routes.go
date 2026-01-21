package routes

import (
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/handlers"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/middleware"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, authService services.AuthService, authHandler *handlers.AuthHandler, todoHandler *handlers.TodoHandler) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
		}

		// Todo routes (protected)
		todos := v1.Group("/todos")
		todos.Use(middleware.AuthMiddleware(authService))
		{
			todos.POST("", todoHandler.CreateTodo)
			todos.GET("", todoHandler.ListTodos)
			todos.GET("/:id", todoHandler.GetTodo)
			todos.PUT("/:id", todoHandler.UpdateTodo)
			todos.DELETE("/:id", todoHandler.DeleteTodo)
		}
	}
}
