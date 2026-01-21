package http

import (
	pb "github.com/MaiNhatHoangY2001/go-project-structure/proto/todo"
	"github.com/MaiNhatHoangY2001/go-project-structure/services/todo-service/grpc"
	"github.com/gin-gonic/gin"
)

func SetupRouter(todoClient pb.TodoServiceClient, authClient *grpc.AuthClient) *gin.Engine {
	r := gin.Default()

	handler := NewTodoHandler(todoClient, authClient)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Todo routes
	api := r.Group("/api/v1")
	{
		todos := api.Group("/todos")
		todos.Use(handler.AuthMiddleware())
		{
			todos.POST("", handler.CreateTodo)
			todos.GET("", handler.ListTodos)
			todos.GET("/:id", handler.GetTodo)
			todos.PUT("/:id", handler.UpdateTodo)
			todos.DELETE("/:id", handler.DeleteTodo)
		}
	}

	return r
}
