package handlers

import (
	"net/http"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/models"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/services"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Signup godoc
// @Summary User signup
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserSignupRequest true "Signup request"
// @Success 201 {object} models.APIResponse{data=models.User}
// @Failure 400 {object} models.APIResponse{error=models.APIError}
// @Failure 409 {object} models.APIResponse{error=models.APIError}
// @Failure 500 {object} models.APIResponse{error=models.APIError}
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req models.UserSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Invalid signup request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.ErrCodeValidation, err.Error()))
		return
	}

	user, err := h.authService.Signup(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "user already exists" {
			logger.Log.Warn("User already exists", zap.String("email", req.Email))
			c.JSON(http.StatusConflict, models.ErrorResponse(models.ErrCodeConflict, "User already exists"))
			return
		}
		logger.Log.Error("Failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.ErrCodeInternalError, "Failed to create user"))
		return
	}

	logger.Log.Info("User created successfully", zap.String("userID", user.ID.Hex()))
	c.JSON(http.StatusCreated, models.SuccessResponse(user))
}

// Login godoc
// @Summary User login
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserLoginRequest true "Login request"
// @Success 200 {object} models.APIResponse{data=models.UserLoginResponse}
// @Failure 400 {object} models.APIResponse{error=models.APIError}
// @Failure 401 {object} models.APIResponse{error=models.APIError}
// @Failure 500 {object} models.APIResponse{error=models.APIError}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse(models.ErrCodeValidation, err.Error()))
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			logger.Log.Warn("Invalid credentials", zap.String("email", req.Email))
			c.JSON(http.StatusUnauthorized, models.ErrorResponse(models.ErrCodeUnauthorized, "Invalid credentials"))
			return
		}
		logger.Log.Error("Failed to login", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(models.ErrCodeInternalError, "Failed to login"))
		return
	}

	response := models.UserLoginResponse{
		Token: token,
		User:  *user,
	}

	logger.Log.Info("User logged in successfully", zap.String("userID", user.ID.Hex()))
	c.JSON(http.StatusOK, models.SuccessResponse(response))
}
