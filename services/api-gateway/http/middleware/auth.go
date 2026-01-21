package middleware

import (
	"net/http"
	"strings"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/dto"
	"github.com/MaiNhatHoangY2001/go-project-structure/services/api-gateway/grpc/clients"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authClient *clients.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse(dto.ErrCodeMissingAuthHeader, "Authorization header is required"))
			c.Abort()
			return
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse(dto.ErrCodeInvalidAuthHeader, "Invalid authorization header format"))
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token via Auth Service
		userID, err := authClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse(dto.ErrCodeInvalidToken, "Invalid or expired token"))
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("user_id", userID)
		c.Next()
	}
}
