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
	RegisterAuthRoutes(v1, r.authHandler, r.authMiddleware)
	RegisterUserRoutes(v1, r.userHandler, r.authMiddleware)

	return router
}
