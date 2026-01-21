package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/dto"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthUsecase is a mock implementation of AuthUsecase
type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Signup(ctx context.Context, req *dto.UserSignupRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockAuthUsecase) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserLoginResponse), args.Error(1)
}

func (m *MockAuthUsecase) ValidateToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	// Initialize logger for tests
	_ = logger.InitLogger("error", "console")
	return gin.New()
}

func TestSignup_Success(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := NewAuthHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/signup", handler.Signup)

	expectedResponse := &dto.UserResponse{
		ID:        "123",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}

	mockUsecase.On("Signup", mock.Anything, mock.AnythingOfType("*dto.UserSignupRequest")).Return(expectedResponse, nil)

	reqBody := dto.UserSignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)

	mockUsecase.AssertExpectations(t)
}

func TestSignup_InvalidJSON(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := NewAuthHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/signup", handler.Signup)

	invalidJSON := []byte(`{"email": "test@example.com", "password": }`)

	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, dto.ErrCodeValidation, response.Error.Code)
}

func TestSignup_ValidationError(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := NewAuthHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/signup", handler.Signup)

	reqBody := dto.UserSignupRequest{
		Email:    "invalid-email",
		Password: "123", // too short
		Name:     "",    // empty
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, dto.ErrCodeValidation, response.Error.Code)
}

func TestSignup_UserAlreadyExists(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := NewAuthHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/signup", handler.Signup)

	mockUsecase.On("Signup", mock.Anything, mock.AnythingOfType("*dto.UserSignupRequest")).Return(nil, errors.New("user already exists"))

	reqBody := dto.UserSignupRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, dto.ErrCodeConflict, response.Error.Code)
	assert.Equal(t, "User already exists", response.Error.Message)

	mockUsecase.AssertExpectations(t)
}

func TestSignup_InternalError(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := NewAuthHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/signup", handler.Signup)

	mockUsecase.On("Signup", mock.Anything, mock.AnythingOfType("*dto.UserSignupRequest")).Return(nil, errors.New("database error"))

	reqBody := dto.UserSignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, dto.ErrCodeInternalError, response.Error.Code)
	assert.Equal(t, "Failed to create user", response.Error.Message)

	mockUsecase.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := NewAuthHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/login", handler.Login)

	expectedResponse := &dto.UserLoginResponse{
		Token: "fake-jwt-token",
		User: dto.UserResponse{
			ID:        "123",
			Email:     "test@example.com",
			Name:      "Test User",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
		},
	}

	mockUsecase.On("Login", mock.Anything, mock.AnythingOfType("*dto.UserLoginRequest")).Return(expectedResponse, nil)

	reqBody := dto.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)

	mockUsecase.AssertExpectations(t)
}

func TestLogin_InvalidJSON(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := NewAuthHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/login", handler.Login)

	invalidJSON := []byte(`{"email": "test@example.com", "password": }`)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, dto.ErrCodeValidation, response.Error.Code)
}

func TestLogin_ValidationError(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := NewAuthHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/login", handler.Login)

	reqBody := dto.UserLoginRequest{
		Email:    "invalid-email",
		Password: "", // empty
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, dto.ErrCodeValidation, response.Error.Code)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := NewAuthHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/login", handler.Login)

	mockUsecase.On("Login", mock.Anything, mock.AnythingOfType("*dto.UserLoginRequest")).Return(nil, errors.New("invalid credentials"))

	reqBody := dto.UserLoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, dto.ErrCodeUnauthorized, response.Error.Code)
	assert.Equal(t, "Invalid credentials", response.Error.Message)

	mockUsecase.AssertExpectations(t)
}

func TestLogin_InternalError(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := NewAuthHandler(mockUsecase)
	router := setupTestRouter()
	router.POST("/login", handler.Login)

	mockUsecase.On("Login", mock.Anything, mock.AnythingOfType("*dto.UserLoginRequest")).Return(nil, errors.New("database error"))

	reqBody := dto.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, dto.ErrCodeInternalError, response.Error.Code)
	assert.Equal(t, "Failed to login", response.Error.Message)

	mockUsecase.AssertExpectations(t)
}
