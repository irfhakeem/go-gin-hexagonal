package routes

import (
	"go-gin-hexagonal/internal/adapter/http/handlers"
	"go-gin-hexagonal/internal/adapter/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
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
	router.Use(middleware.CSRFMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  true,
			"message": "Service is healthy",
		})
	})

	router.GET("/csrf-token", r.authMiddleware.Middleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  true,
			"message": "CSRF token retrieved successfully",
			"data":    csrf.Token(c.Request),
		})
	})

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/verify-email", r.authHandler.VerifyEmail)
			auth.POST("/send-verify-email", r.authHandler.SendVerifyEmail)
			auth.POST("/reset-password", r.authHandler.ResetPassword)
			auth.POST("/send-reset-password", r.authHandler.SendResetPassword)

			authProtected := auth.Group("")
			authProtected.Use(r.authMiddleware.Middleware())
			{
				authProtected.POST("/refresh", r.authHandler.RefreshToken)
				authProtected.POST("/logout", r.authHandler.Logout)
			}
		}

		users := v1.Group("/users")
		users.Use(r.authMiddleware.Middleware())
		{
			users.GET("", r.userHandler.GetAllUsers)
			users.GET("/profile", r.userHandler.GetProfile)
			users.GET("/:id", r.userHandler.GetUserByID)
			users.POST("", r.userHandler.CreateUser)
			users.PUT("/profile", r.userHandler.UpdateProfile)
			users.PUT("/change-password", r.userHandler.ChangePassword)
			users.DELETE("/:id", r.userHandler.DeleteUser)
		}
	}

	return router
}
