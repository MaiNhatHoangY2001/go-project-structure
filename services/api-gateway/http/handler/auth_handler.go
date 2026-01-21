package handler

import (
	"net/http"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/dto"
	"github.com/MaiNhatHoangY2001/go-project-structure/services/api-gateway/grpc/clients"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authClient *clients.AuthClient
}

func NewAuthHandler(authClient *clients.AuthClient) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
	}
}

// Signup godoc
// @Summary User signup
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UserSignupRequest true "Signup request"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 409 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.UserSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(dto.ErrCodeBadRequest, err.Error()))
		return
	}

	userID, token, err := h.authClient.Signup(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		statusCode, errCode := mapGRPCError(err)
		c.JSON(statusCode, dto.ErrorResponse(errCode, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(map[string]interface{}{
		"user_id": userID,
		"token":   token,
	}))
}

// Login godoc
// @Summary User login
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UserLoginRequest true "Login request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(dto.ErrCodeBadRequest, err.Error()))
		return
	}

	userID, token, err := h.authClient.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		statusCode, errCode := mapGRPCError(err)
		c.JSON(statusCode, dto.ErrorResponse(errCode, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(map[string]interface{}{
		"user_id": userID,
		"token":   token,
	}))
}

// mapGRPCError maps gRPC status codes to HTTP status codes and custom error codes
func mapGRPCError(err error) (int, int) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, dto.ErrCodeInternalError
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return http.StatusBadRequest, dto.ErrCodeValidation
	case codes.AlreadyExists:
		return http.StatusConflict, dto.ErrCodeConflict
	case codes.NotFound:
		return http.StatusNotFound, dto.ErrCodeNotFound
	case codes.Unauthenticated:
		return http.StatusUnauthorized, dto.ErrCodeUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden, dto.ErrCodeForbidden
	default:
		return http.StatusInternalServerError, dto.ErrCodeInternalError
	}
}
