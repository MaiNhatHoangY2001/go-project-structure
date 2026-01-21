package http

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	pb "github.com/MaiNhatHoangY2001/go-project-structure/proto/todo"
	"github.com/MaiNhatHoangY2001/go-project-structure/services/todo-service/grpc"
)

type TodoHandler struct {
	todoClient pb.TodoServiceClient
	authClient *grpc.AuthClient
}

func NewTodoHandler(todoClient pb.TodoServiceClient, authClient *grpc.AuthClient) *TodoHandler {
	return &TodoHandler{
		todoClient: todoClient,
		authClient: authClient,
	}
}

func (h *TodoHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		userID, err := h.authClient.ValidateToken(context.Background(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.todoClient.CreateTodo(context.Background(), &pb.CreateTodoRequest{
		UserId:      userID,
		Title:       req.Title,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *TodoHandler) GetTodo(c *gin.Context) {
	userID := c.GetString("user_id")
	todoID := c.Param("id")

	resp, err := h.todoClient.GetTodo(context.Background(), &pb.GetTodoRequest{
		UserId: userID,
		TodoId: todoID,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TodoHandler) ListTodos(c *gin.Context) {
	userID := c.GetString("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var completed *bool
	if completedStr := c.Query("completed"); completedStr != "" {
		c := completedStr == "true"
		completed = &c
	}

	resp, err := h.todoClient.ListTodos(context.Background(), &pb.ListTodosRequest{
		UserId:    userID,
		Page:      int32(page),
		PageSize:  int32(pageSize),
		Completed: completed,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userID := c.GetString("user_id")
	todoID := c.Param("id")

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Completed   *bool   `json:"completed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.todoClient.UpdateTodo(context.Background(), &pb.UpdateTodoRequest{
		UserId:      userID,
		TodoId:      todoID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userID := c.GetString("user_id")
	todoID := c.Param("id")

	resp, err := h.todoClient.DeleteTodo(context.Background(), &pb.DeleteTodoRequest{
		UserId: userID,
		TodoId: todoID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}
