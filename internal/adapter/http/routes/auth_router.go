package routes

import (
	"go-gin-hexagonal/internal/adapter/http/handlers"
	"go-gin-hexagonal/internal/adapter/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/verify-email", authHandler.VerifyEmail)
		auth.POST("/send-verify-email", authHandler.SendVerifyEmail)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/send-reset-password", authHandler.SendResetPassword)

		authProtected := auth.Group("")
		authProtected.Use(authMiddleware.Middleware())
		{
			authProtected.POST("/refresh", authHandler.RefreshToken)
			authProtected.POST("/logout", authHandler.Logout)
		}
	}
}
