package middleware

import (
	"strings"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/dto"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/usecase"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware(authUsecase usecase.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, dto.ErrorResponse(dto.ErrCodeMissingAuthHeader, "Authorization header is required"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, dto.ErrorResponse(dto.ErrCodeInvalidAuthHeader, "Invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := authUsecase.ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, dto.ErrorResponse(dto.ErrCodeInvalidToken, "Invalid or expired token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, dto.ErrorResponse(dto.ErrCodeInvalidTokenClaims, "Invalid token claims"))
			c.Abort()
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(401, dto.ErrorResponse(dto.ErrCodeMissingUserID, "User ID not found in token"))
			c.Abort()
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(401, dto.ErrorResponse(dto.ErrCodeInvalidTokenClaims, "Invalid user ID format"))
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
