package handler

import (
	"net/http"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/dto"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/usecase"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// Signup godoc
// @Summary User signup
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UserSignupRequest true "Signup request"
// @Success 201 {object} dto.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.APIError}
// @Failure 409 {object} dto.APIResponse{error=dto.APIError}
// @Failure 500 {object} dto.APIResponse{error=dto.APIError}
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.UserSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Invalid signup request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(dto.ErrCodeValidation, err.Error()))
		return
	}

	user, err := h.authUsecase.Signup(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "user already exists" {
			logger.Log.Warn("User already exists", zap.String("email", req.Email))
			c.JSON(http.StatusConflict, dto.ErrorResponse(dto.ErrCodeConflict, "User already exists"))
			return
		}
		logger.Log.Error("Failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(dto.ErrCodeInternalError, "Failed to create user"))
		return
	}

	logger.Log.Info("User created successfully", zap.String("userID", user.ID))
	c.JSON(http.StatusCreated, dto.SuccessResponse(user))
}

// Login godoc
// @Summary User login
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UserLoginRequest true "Login request"
// @Success 200 {object} dto.APIResponse{data=dto.UserLoginResponse}
// @Failure 400 {object} dto.APIResponse{error=dto.APIError}
// @Failure 401 {object} dto.APIResponse{error=dto.APIError}
// @Failure 500 {object} dto.APIResponse{error=dto.APIError}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(dto.ErrCodeValidation, err.Error()))
		return
	}

	response, err := h.authUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			logger.Log.Warn("Invalid credentials", zap.String("email", req.Email))
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse(dto.ErrCodeUnauthorized, "Invalid credentials"))
			return
		}
		logger.Log.Error("Failed to login", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(dto.ErrCodeInternalError, "Failed to login"))
		return
	}

	logger.Log.Info("User logged in successfully", zap.String("userID", response.User.ID))
	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}
