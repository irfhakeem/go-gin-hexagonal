package infrastructure

import (
	"go-gin-clean/internal/adapters/secondary/crypto"
	"go-gin-clean/internal/adapters/secondary/localstorage"
	"go-gin-clean/internal/adapters/secondary/postgres"
	"go-gin-clean/internal/adapters/secondary/smtp"
	"go-gin-clean/internal/domain/service"
	"go-gin-clean/internal/ports/primary"
	"go-gin-clean/internal/ports/secondary"
	"go-gin-clean/pkg/config"

	"gorm.io/gorm"
)

type Container struct {
	UserUseCase   primary.UserUseCase
	EmailUseCase  primary.EmailUseCase
	JWTService    secondary.JWTService
	MailerService secondary.MailerService
}

func NewContainer(db *gorm.DB, cfg *config.Config) *Container {
	// Init repositories
	userRepo := postgres.NewUserRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)

	// Init services (secondary adapters)
	jwtService := crypto.NewJWTService(&cfg.JWT)
	bcryptService := crypto.NewBcryptService()
	aesService := crypto.NewAESService(&cfg.AES)
	smtpService := smtp.NewSMTPService(&cfg.Mailer)
	localStorageService := localstorage.NewLocalStorageService()

	// Init use cases (domain services)
	emailUseCase := service.NewEmailService(smtpService)
	userUseCase := service.NewUserService(
		userRepo,
		emailUseCase,
		refreshTokenRepo,
		jwtService,
		bcryptService,
		aesService,
		localStorageService,
	)

	return &Container{
		UserUseCase:   userUseCase,
		EmailUseCase:  emailUseCase,
		JWTService:    jwtService,
		MailerService: smtpService,
	}
}
