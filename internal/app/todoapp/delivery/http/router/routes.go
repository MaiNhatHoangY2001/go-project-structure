package router

import (
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/delivery/http/handler"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/delivery/middleware"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/usecase"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, authUsecase usecase.AuthUsecase, authHandler *handler.AuthHandler, todoHandler *handler.TodoHandler) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
		}

		todos := v1.Group("/todos")
		todos.Use(middleware.AuthMiddleware(authUsecase))
		{
			todos.POST("", todoHandler.CreateTodo)
			todos.GET("", todoHandler.ListTodos)
			todos.GET("/:id", todoHandler.GetTodo)
			todos.PUT("/:id", todoHandler.UpdateTodo)
			todos.DELETE("/:id", todoHandler.DeleteTodo)
		}
	}
}
