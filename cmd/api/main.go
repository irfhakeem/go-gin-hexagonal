package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	dbAdapter "go-gin-hexagonal/internal/adapter/database"
	"go-gin-hexagonal/internal/adapter/http/handlers"
	"go-gin-hexagonal/internal/adapter/http/middleware"
	"go-gin-hexagonal/internal/adapter/http/routes"
	"go-gin-hexagonal/internal/adapter/mailer"
	"go-gin-hexagonal/internal/adapter/security"
	"go-gin-hexagonal/internal/application/service"
	"go-gin-hexagonal/internal/domain/entity"

	"go-gin-hexagonal/pkg/config"
	"go-gin-hexagonal/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.NewPostgresConnection(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Init adapters
	// Database adapters
	userRepo := dbAdapter.NewUserRepository(db, dbAdapter.NewBaseRepository[entity.User](db))
	refreshTokenRepo := dbAdapter.NewRefreshTokenRepository(db, dbAdapter.NewBaseRepository[entity.RefreshToken](db))

	// Security adapters
	passwordHasher := security.NewBcryptHasher()
	tokenManager := security.NewJWTTokenManager(cfg.JWT)
	encryptor := security.NewAESEncryptor(cfg.AES)

	// Mailer adapter
	mailerManager := mailer.NewSMTPMailer(&cfg.Mailer)

	// Init services
	emailService := service.NewEmailService(mailerManager)
	authService := service.NewAuthService(userRepo, refreshTokenRepo, tokenManager, passwordHasher, emailService, encryptor)
	userService := service.NewUserService(userRepo, passwordHasher, emailService)

	// Init Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Init middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenManager)

	// Init router
	appRouter := routes.NewRouter(authHandler, userHandler, authMiddleware)
	ginRouter := appRouter.SetupRoutes()

	srv := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      ginRouter,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("Server starting on %s", cfg.Server.Address())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
