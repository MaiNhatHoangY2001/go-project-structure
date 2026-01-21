package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/delivery/http/handler"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/delivery/http/router"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/dto"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/entity"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/usecase"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository for integration tests
type MockUserRepo struct {
	users map[string]*entity.User
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		users: make(map[string]*entity.User),
	}
}

func (r *MockUserRepo) Create(ctx context.Context, user *entity.User) error {
	user.ID = primitive.NewObjectID()
	r.users[user.Email] = user
	return nil
}

func (r *MockUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, exists := r.users[email]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *MockUserRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

// MockTodoRepository for integration tests
type MockTodoRepo struct {
	todos map[string]*entity.Todo
}

func NewMockTodoRepo() *MockTodoRepo {
	return &MockTodoRepo{
		todos: make(map[string]*entity.Todo),
	}
}

func (r *MockTodoRepo) Create(ctx context.Context, todo *entity.Todo) error {
	todo.ID = primitive.NewObjectID()
	r.todos[todo.ID.Hex()] = todo
	return nil
}

func (r *MockTodoRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Todo, error) {
	todo, exists := r.todos[id.Hex()]
	if !exists {
		return nil, nil
	}
	return todo, nil
}

func (r *MockTodoRepo) FindByUserID(ctx context.Context, userID primitive.ObjectID, page, pageSize int, completed *bool) ([]*entity.Todo, int64, error) {
	var result []*entity.Todo
	for _, todo := range r.todos {
		if todo.UserID == userID {
			if completed == nil || todo.Completed == *completed {
				result = append(result, todo)
			}
		}
	}
	return result, int64(len(result)), nil
}

func (r *MockTodoRepo) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	todo, exists := r.todos[id.Hex()]
	if !exists {
		return nil
	}
	// Simple update logic for testing
	r.todos[id.Hex()] = todo
	return nil
}

func (r *MockTodoRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	delete(r.todos, id.Hex())
	return nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Initialize logger for tests
	if err := logger.InitLogger("error", "console"); err != nil {
		panic(err)
	}

	userRepo := NewMockUserRepo()
	todoRepo := NewMockTodoRepo()

	authUsecase := usecase.NewAuthUsecase(userRepo, "test-secret", 24)
	todoUsecase := usecase.NewTodoUsecase(todoRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	todoHandler := handler.NewTodoHandler(todoUsecase)

	ginRouter := gin.Default()
	router.SetupRoutes(ginRouter, authUsecase, authHandler, todoHandler)

	return ginRouter
}

func TestAuthSignupAndLogin_Integration(t *testing.T) {
	ginRouter := setupRouter()

	// Test Signup
	signupReq := dto.UserSignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	signupBody, _ := json.Marshal(signupReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(signupBody))
	req.Header.Set("Content-Type", "application/json")
	ginRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var signupResp dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &signupResp)
	assert.NoError(t, err)
	assert.True(t, signupResp.Success)

	// Test Login
	loginReq := dto.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	loginBody, _ := json.Marshal(loginReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	ginRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var loginResp dto.APIResponse
	err = json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NoError(t, err)
	assert.True(t, loginResp.Success)

	// Extract token
	loginData := loginResp.Data.(map[string]interface{})
	assert.NotEmpty(t, loginData["token"])
}

func TestTodoCRUD_Integration(t *testing.T) {
	ginRouter := setupRouter()

	// Signup first
	signupReq := dto.UserSignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	signupBody, _ := json.Marshal(signupReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(signupBody))
	req.Header.Set("Content-Type", "application/json")
	ginRouter.ServeHTTP(w, req)

	// Login to get token
	loginReq := dto.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	loginBody, _ := json.Marshal(loginReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	ginRouter.ServeHTTP(w, req)

	var loginResp dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NoError(t, err)
	loginData := loginResp.Data.(map[string]interface{})
	token := loginData["token"].(string)

	// Create Todo
	createReq := dto.TodoCreateRequest{
		Title:       "Test Todo",
		Description: "Test Description",
	}
	createBody, _ := json.Marshal(createReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	ginRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp dto.APIResponse
	err = json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.NoError(t, err)
	assert.True(t, createResp.Success)

	// List Todos
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/todos?page=1&pageSize=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	ginRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var listResp dto.APIResponse
	err = json.Unmarshal(w.Body.Bytes(), &listResp)
	assert.NoError(t, err)
	assert.True(t, listResp.Success)
}

func TestTodoUnauthorizedAccess_Integration(t *testing.T) {
	ginRouter := setupRouter()

	// Try to create todo without token
	createReq := dto.TodoCreateRequest{
		Title:       "Test Todo",
		Description: "Test Description",
	}
	createBody, _ := json.Marshal(createReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	ginRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, dto.ErrCodeMissingAuthHeader, resp.Error.Code)
}
