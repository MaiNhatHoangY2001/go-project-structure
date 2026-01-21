package handlers

import (
	"net/http"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/models"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/services"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type TodoHandler struct {
	todoService services.TodoService
}

func NewTodoHandler(todoService services.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo item for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.TodoCreateRequest true "Todo creation request"
// @Success 201 {object} models.APIResponse{data=models.Todo}
// @Failure 400 {object} models.APIResponse{error=models.APIError}
// @Failure 401 {object} models.APIResponse{error=models.APIError}
// @Failure 500 {object} models.APIResponse{error=models.APIError}
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse(models.ErrCodeMissingUserID, "User ID not found"))
		return
	}

	var req models.TodoCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Invalid create todo request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.ErrCodeValidation, err.Error()))
		return
	}

	todo, err := h.todoService.Create(c.Request.Context(), userID.(primitive.ObjectID), &req)
	if err != nil {
		logger.Log.Error("Failed to create todo", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.ErrCodeInternalError, "Failed to create todo"))
		return
	}

	logger.Log.Info("Todo created successfully", zap.String("todoID", todo.ID.Hex()))
	c.JSON(http.StatusCreated, models.SuccessResponse(todo))
}

// GetTodo godoc
// @Summary Get a todo by ID
// @Description Get a specific todo item by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 200 {object} models.APIResponse{data=models.Todo}
// @Failure 400 {object} models.APIResponse{error=models.APIError}
// @Failure 401 {object} models.APIResponse{error=models.APIError}
// @Failure 404 {object} models.APIResponse{error=models.APIError}
// @Failure 500 {object} models.APIResponse{error=models.APIError}
// @Router /todos/{id} [get]
func (h *TodoHandler) GetTodo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse(models.ErrCodeMissingUserID, "User ID not found"))
		return
	}

	todoIDStr := c.Param("id")
	todoID, err := primitive.ObjectIDFromHex(todoIDStr)
	if err != nil {
		logger.Log.Error("Invalid todo ID", zap.String("todoID", todoIDStr))
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.ErrCodeBadRequest, "Invalid todo ID"))
		return
	}

	todo, err := h.todoService.GetByID(c.Request.Context(), userID.(primitive.ObjectID), todoID)
	if err != nil {
		if err.Error() == "todo not found" {
			logger.Log.Warn("Todo not found", zap.String("todoID", todoIDStr))
			c.JSON(http.StatusNotFound, models.ErrorResponse(models.ErrCodeNotFound, "Todo not found"))
			return
		}
		if err.Error() == "unauthorized access to todo" {
			logger.Log.Warn("Unauthorized access to todo", zap.String("todoID", todoIDStr))
			c.JSON(http.StatusForbidden, models.ErrorResponse(models.ErrCodeForbidden, "Access denied"))
			return
		}
		logger.Log.Error("Failed to get todo", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.ErrCodeInternalError, "Failed to get todo"))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(todo))
}

// ListTodos godoc
// @Summary List todos
// @Description Get a paginated list of todos for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10, max: 100)"
// @Param completed query bool false "Filter by completion status"
// @Success 200 {object} models.APIResponse{data=models.PaginatedData}
// @Failure 400 {object} models.APIResponse{error=models.APIError}
// @Failure 401 {object} models.APIResponse{error=models.APIError}
// @Failure 500 {object} models.APIResponse{error=models.APIError}
// @Router /todos [get]
func (h *TodoHandler) ListTodos(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse(models.ErrCodeMissingUserID, "User ID not found"))
		return
	}

	var query models.TodoListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Log.Error("Invalid list query", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.ErrCodeValidation, err.Error()))
		return
	}

	todos, total, err := h.todoService.List(c.Request.Context(), userID.(primitive.ObjectID), &query)
	if err != nil {
		logger.Log.Error("Failed to list todos", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.ErrCodeInternalError, "Failed to list todos"))
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse(query.Page, query.PageSize, total, todos))
}

// UpdateTodo godoc
// @Summary Update a todo
// @Description Update a specific todo item
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Param request body models.TodoUpdateRequest true "Todo update request"
// @Success 200 {object} models.APIResponse{data=models.Todo}
// @Failure 400 {object} models.APIResponse{error=models.APIError}
// @Failure 401 {object} models.APIResponse{error=models.APIError}
// @Failure 404 {object} models.APIResponse{error=models.APIError}
// @Failure 500 {object} models.APIResponse{error=models.APIError}
// @Router /todos/{id} [put]
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse(models.ErrCodeMissingUserID, "User ID not found"))
		return
	}

	todoIDStr := c.Param("id")
	todoID, err := primitive.ObjectIDFromHex(todoIDStr)
	if err != nil {
		logger.Log.Error("Invalid todo ID", zap.String("todoID", todoIDStr))
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.ErrCodeBadRequest, "Invalid todo ID"))
		return
	}

	var req models.TodoUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Invalid update todo request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.ErrCodeValidation, err.Error()))
		return
	}

	todo, err := h.todoService.Update(c.Request.Context(), userID.(primitive.ObjectID), todoID, &req)
	if err != nil {
		if err.Error() == "todo not found" {
			logger.Log.Warn("Todo not found", zap.String("todoID", todoIDStr))
			c.JSON(http.StatusNotFound, models.ErrorResponse(models.ErrCodeNotFound, "Todo not found"))
			return
		}
		if err.Error() == "unauthorized access to todo" {
			logger.Log.Warn("Unauthorized access to todo", zap.String("todoID", todoIDStr))
			c.JSON(http.StatusForbidden, models.ErrorResponse(models.ErrCodeForbidden, "Access denied"))
			return
		}
		logger.Log.Error("Failed to update todo", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.ErrCodeInternalError, "Failed to update todo"))
		return
	}

	logger.Log.Info("Todo updated successfully", zap.String("todoID", todo.ID.Hex()))
	c.JSON(http.StatusOK, models.SuccessResponse(todo))
}

// DeleteTodo godoc
// @Summary Delete a todo
// @Description Delete a specific todo item
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse{error=models.APIError}
// @Failure 401 {object} models.APIResponse{error=models.APIError}
// @Failure 404 {object} models.APIResponse{error=models.APIError}
// @Failure 500 {object} models.APIResponse{error=models.APIError}
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse(models.ErrCodeMissingUserID, "User ID not found"))
		return
	}

	todoIDStr := c.Param("id")
	todoID, err := primitive.ObjectIDFromHex(todoIDStr)
	if err != nil {
		logger.Log.Error("Invalid todo ID", zap.String("todoID", todoIDStr))
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.ErrCodeBadRequest, "Invalid todo ID"))
		return
	}

	err = h.todoService.Delete(c.Request.Context(), userID.(primitive.ObjectID), todoID)
	if err != nil {
		if err.Error() == "todo not found" {
			logger.Log.Warn("Todo not found", zap.String("todoID", todoIDStr))
			c.JSON(http.StatusNotFound, models.ErrorResponse(models.ErrCodeNotFound, "Todo not found"))
			return
		}
		if err.Error() == "unauthorized access to todo" {
			logger.Log.Warn("Unauthorized access to todo", zap.String("todoID", todoIDStr))
			c.JSON(http.StatusForbidden, models.ErrorResponse(models.ErrCodeForbidden, "Access denied"))
			return
		}
		logger.Log.Error("Failed to delete todo", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.ErrCodeInternalError, "Failed to delete todo"))
		return
	}

	logger.Log.Info("Todo deleted successfully", zap.String("todoID", todoIDStr))
	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{"message": "Todo deleted successfully"}))
}
