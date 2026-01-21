package services

import (
	"context"
	"testing"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestSignup_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, "test-secret", 24)

	req := &models.UserSignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, mongo.ErrNoDocuments)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	user, err := authService.Signup(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Name, user.Name)
	mockRepo.AssertExpectations(t)
}

func TestSignup_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, "test-secret", 24)

	req := &models.UserSignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	existingUser := &models.User{
		ID:    primitive.NewObjectID(),
		Email: req.Email,
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	user, err := authService.Signup(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, "test-secret", 24)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &models.User{
		ID:       primitive.NewObjectID(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Name:     "Test User",
	}

	req := &models.UserLoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	token, user, err := authService.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotNil(t, user)
	assert.Equal(t, existingUser.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, "test-secret", 24)

	req := &models.UserLoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, mongo.ErrNoDocuments)

	token, user, err := authService.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestValidateToken_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, "test-secret", 24)

	// Create a valid token
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &models.User{
		ID:       primitive.NewObjectID(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Name:     "Test User",
	}

	req := &models.UserLoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	token, _, _ := authService.Login(context.Background(), req)

	// Validate the token
	validatedToken, err := authService.ValidateToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, validatedToken)
	assert.True(t, validatedToken.Valid)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, "test-secret", 24)

	invalidToken := "invalid.token.here"

	validatedToken, err := authService.ValidateToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, validatedToken)
}
