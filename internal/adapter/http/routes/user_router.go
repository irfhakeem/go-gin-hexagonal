package routes

import (
	"go-gin-hexagonal/internal/adapter/http/handlers"
	"go-gin-hexagonal/internal/adapter/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, userHandler *handlers.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	users := rg.Group("/users")
	users.Use(authMiddleware.Middleware())
	{
		users.GET("", userHandler.GetAllUsers)
		users.GET("/profile", userHandler.GetProfile)
		users.GET("/:id", userHandler.GetUserByID)
		users.POST("", userHandler.CreateUser)
		users.PUT("/profile", userHandler.UpdateProfile)
		users.PUT("/change-password", userHandler.ChangePassword)
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}
