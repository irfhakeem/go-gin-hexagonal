package routes

import (
	"go-gin-hexagonal/internal/adapters/http/handlers"
	"go-gin-hexagonal/internal/adapters/http/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler    *handlers.AuthHandler
	userHandler    *handlers.UserHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		authHandler:    authHandler,
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  true,
			"message": "Service is healthy",
		})
	})

	v1 := router.Group("/api")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/refresh", r.authHandler.RefreshToken)

			// Protected
			auth.Use(r.authMiddleware.Middleware())
			auth.POST("/logout", r.authHandler.Logout)
		}

		users := v1.Group("/users")
		users.Use(r.authMiddleware.Middleware())
		{
			users.GET("/profile", r.userHandler.GetProfile)
			users.PUT("/profile", r.userHandler.UpdateProfile)
			users.PUT("/change-password", r.userHandler.ChangePassword)
			users.GET("", r.userHandler.ListUsers)
			users.GET("/:id", r.userHandler.GetUserByID)
			users.DELETE("/:id", r.userHandler.DeleteUser)
		}
	}

	return router
}
