package middleware

import (
	"strings"

	"go-gin-hexagonal/internal/domain/ports"
	"go-gin-hexagonal/pkg/message"
	"go-gin-hexagonal/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenManager ports.TokenManager
}

func NewAuthMiddleware(tokenManager ports.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, message.FAILED_GET_AUTHORIZATION_HEADER, "Authorization Header Not Found", 401)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, message.FAILED_TOKEN_INVALID, "Token Invalid", 401)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			response.Error(c, message.FAILED_TOKEN_NOT_FOUND, "Token Not Found", 401)
			c.Abort()
			return
		}

		claims, err := m.tokenManager.ValidateAccessToken(token)
		if err != nil {
			response.Error(c, message.FAILED_TOKEN_INVALID, "Token Invalid", 401)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.Username)

		c.Next()
	}
}
