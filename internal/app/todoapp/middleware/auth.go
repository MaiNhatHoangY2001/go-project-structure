package middleware

import (
	"strings"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/models"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, models.ErrorResponse(models.ErrCodeMissingAuthHeader, "Authorization header is required"))
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, models.ErrorResponse(models.ErrCodeInvalidAuthHeader, "Invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, models.ErrorResponse(models.ErrCodeInvalidToken, "Invalid or expired token"))
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, models.ErrorResponse(models.ErrCodeInvalidTokenClaims, "Invalid token claims"))
			c.Abort()
			return
		}

		// Get user ID from claims
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(401, models.ErrorResponse(models.ErrCodeMissingUserID, "User ID not found in token"))
			c.Abort()
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(401, models.ErrorResponse(models.ErrCodeInvalidTokenClaims, "Invalid user ID format"))
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("userID", userID)
		c.Next()
	}
}
