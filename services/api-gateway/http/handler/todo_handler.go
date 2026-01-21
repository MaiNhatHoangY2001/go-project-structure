package handler

import (
	"net/http"
	"strconv"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/dto"
	"github.com/MaiNhatHoangY2001/go-project-structure/services/api-gateway/grpc/clients"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	todoClient *clients.TodoClient
}

func NewTodoHandler(todoClient *clients.TodoClient) *TodoHandler {
	return &TodoHandler{
		todoClient: todoClient,
	}
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo item for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TodoCreateRequest true "Create todo request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/todos [post]
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse(dto.ErrCodeMissingUserID, "User ID not found in context"))
		return
	}

	var req dto.TodoCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(dto.ErrCodeBadRequest, err.Error()))
		return
	}

	todo, err := h.todoClient.CreateTodo(c.Request.Context(), userID.(string), req.Title, req.Description)
	if err != nil {
		statusCode, errCode := mapGRPCError(err)
		c.JSON(statusCode, dto.ErrorResponse(errCode, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(map[string]interface{}{
		"id":          todo.Id,
		"user_id":     todo.UserId,
		"title":       todo.Title,
		"description": todo.Description,
		"completed":   todo.Completed,
		"created_at":  todo.CreatedAt,
		"updated_at":  todo.UpdatedAt,
	}))
}

// GetTodo godoc
// @Summary Get a todo by ID
// @Description Get details of a specific todo item
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/todos/{id} [get]
func (h *TodoHandler) GetTodo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse(dto.ErrCodeMissingUserID, "User ID not found in context"))
		return
	}

	todoID := c.Param("id")

	todo, err := h.todoClient.GetTodo(c.Request.Context(), userID.(string), todoID)
	if err != nil {
		statusCode, errCode := mapGRPCError(err)
		c.JSON(statusCode, dto.ErrorResponse(errCode, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(map[string]interface{}{
		"id":          todo.Id,
		"user_id":     todo.UserId,
		"title":       todo.Title,
		"description": todo.Description,
		"completed":   todo.Completed,
		"created_at":  todo.CreatedAt,
		"updated_at":  todo.UpdatedAt,
	}))
}

// ListTodos godoc
// @Summary List all todos
// @Description Get all todo items for the authenticated user with pagination
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/todos [get]
func (h *TodoHandler) ListTodos(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse(dto.ErrCodeMissingUserID, "User ID not found in context"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	todos, total, err := h.todoClient.ListTodos(c.Request.Context(), userID.(string), int32(page), int32(pageSize))
	if err != nil {
		statusCode, errCode := mapGRPCError(err)
		c.JSON(statusCode, dto.ErrorResponse(errCode, err.Error()))
		return
	}

	todosData := make([]map[string]interface{}, len(todos))
	for i, todo := range todos {
		todosData[i] = map[string]interface{}{
			"id":          todo.Id,
			"user_id":     todo.UserId,
			"title":       todo.Title,
			"description": todo.Description,
			"completed":   todo.Completed,
			"created_at":  todo.CreatedAt,
			"updated_at":  todo.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse(page, pageSize, total, todosData))
}

// UpdateTodo godoc
// @Summary Update a todo
// @Description Update an existing todo item
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Param request body dto.TodoUpdateRequest true "Update todo request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/todos/{id} [put]
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse(dto.ErrCodeMissingUserID, "User ID not found in context"))
		return
	}

	todoID := c.Param("id")

	var req dto.TodoUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(dto.ErrCodeBadRequest, err.Error()))
		return
	}

	// Get default values if not provided
	title := ""
	if req.Title != nil {
		title = *req.Title
	}
	description := ""
	if req.Description != nil {
		description = *req.Description
	}
	completed := false
	if req.Completed != nil {
		completed = *req.Completed
	}

	todo, err := h.todoClient.UpdateTodo(c.Request.Context(), userID.(string), todoID, title, description, completed)
	if err != nil {
		statusCode, errCode := mapGRPCError(err)
		c.JSON(statusCode, dto.ErrorResponse(errCode, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(map[string]interface{}{
		"id":          todo.Id,
		"user_id":     todo.UserId,
		"title":       todo.Title,
		"description": todo.Description,
		"completed":   todo.Completed,
		"created_at":  todo.CreatedAt,
		"updated_at":  todo.UpdatedAt,
	}))
}

// DeleteTodo godoc
// @Summary Delete a todo
// @Description Delete a specific todo item
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 403 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse(dto.ErrCodeMissingUserID, "User ID not found in context"))
		return
	}

	todoID := c.Param("id")

	err := h.todoClient.DeleteTodo(c.Request.Context(), userID.(string), todoID)
	if err != nil {
		statusCode, errCode := mapGRPCError(err)
		c.JSON(statusCode, dto.ErrorResponse(errCode, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(map[string]interface{}{"message": "Todo deleted successfully"}))
}
