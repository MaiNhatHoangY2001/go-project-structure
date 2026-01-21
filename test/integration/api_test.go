package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/handlers"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/models"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/routes"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/services"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository for integration tests
type MockUserRepo struct {
	users map[string]*models.User
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		users: make(map[string]*models.User),
	}
}

func (r *MockUserRepo) Create(ctx context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	r.users[user.Email] = user
	return nil
}

func (r *MockUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, exists := r.users[email]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *MockUserRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

// MockTodoRepository for integration tests
type MockTodoRepo struct {
	todos map[string]*models.Todo
}

func NewMockTodoRepo() *MockTodoRepo {
	return &MockTodoRepo{
		todos: make(map[string]*models.Todo),
	}
}

func (r *MockTodoRepo) Create(ctx context.Context, todo *models.Todo) error {
	todo.ID = primitive.NewObjectID()
	r.todos[todo.ID.Hex()] = todo
	return nil
}

func (r *MockTodoRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Todo, error) {
	todo, exists := r.todos[id.Hex()]
	if !exists {
		return nil, nil
	}
	return todo, nil
}

func (r *MockTodoRepo) FindByUserID(ctx context.Context, userID primitive.ObjectID, page, pageSize int, completed *bool) ([]*models.Todo, int64, error) {
	var result []*models.Todo
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

	authService := services.NewAuthService(userRepo, "test-secret", 24)
	todoService := services.NewTodoService(todoRepo)

	authHandler := handlers.NewAuthHandler(authService)
	todoHandler := handlers.NewTodoHandler(todoService)

	router := gin.Default()
	routes.SetupRoutes(router, authService, authHandler, todoHandler)

	return router
}

func TestAuthSignupAndLogin_Integration(t *testing.T) {
	router := setupRouter()

	// Test Signup
	signupReq := models.UserSignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	signupBody, _ := json.Marshal(signupReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(signupBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var signupResp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &signupResp)
	assert.True(t, signupResp.Success)

	// Test Login
	loginReq := models.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	loginBody, _ := json.Marshal(loginReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var loginResp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.True(t, loginResp.Success)

	// Extract token
	loginData := loginResp.Data.(map[string]interface{})
	assert.NotEmpty(t, loginData["token"])
}

func TestTodoCRUD_Integration(t *testing.T) {
	router := setupRouter()

	// Signup first
	signupReq := models.UserSignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	signupBody, _ := json.Marshal(signupReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(signupBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Login to get token
	loginReq := models.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	loginBody, _ := json.Marshal(loginReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var loginResp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	loginData := loginResp.Data.(map[string]interface{})
	token := loginData["token"].(string)

	// Create Todo
	createReq := models.TodoCreateRequest{
		Title:       "Test Todo",
		Description: "Test Description",
	}
	createBody, _ := json.Marshal(createReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.True(t, createResp.Success)

	// List Todos
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/todos?page=1&pageSize=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var listResp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &listResp)
	assert.True(t, listResp.Success)
}

func TestTodoUnauthorizedAccess_Integration(t *testing.T) {
	router := setupRouter()

	// Try to create todo without token
	createReq := models.TodoCreateRequest{
		Title:       "Test Todo",
		Description: "Test Description",
	}
	createBody, _ := json.Marshal(createReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp models.APIResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, models.ErrCodeMissingAuthHeader, resp.Error.Code)
}
