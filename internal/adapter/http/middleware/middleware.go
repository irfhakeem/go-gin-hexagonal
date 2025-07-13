package middleware

import (
	"strings"

	response "go-gin-hexagonal/internal/adapter/http"
	"go-gin-hexagonal/internal/adapter/http/message"
	"go-gin-hexagonal/internal/domain/ports"

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

func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, message.FAILED_GET_AUTHORIZATION_HEADER, ports.ErrAuthorizationHeaderNotFound.Error(), 401)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, message.FAILED_TOKEN_INVALID, ports.ErrTokenInvalid.Error(), 401)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			response.Error(c, message.FAILED_TOKEN_NOT_FOUND, ports.ErrTokenNotFound.Error(), 401)
			c.Abort()
			return
		}

		claims, err := m.tokenManager.ValidateAccessToken(token)
		if err != nil {
			response.Error(c, message.FAILED_TOKEN_INVALID, ports.ErrTokenInvalid.Error(), 401)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.Username)

		c.Next()
	}
}
