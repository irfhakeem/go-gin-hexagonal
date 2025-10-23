package secondary

import (
	"go-gin-clean/internal/domain/model"
	"go-gin-clean/internal/ports/primary"
	"time"
)

// JWTService defines the secondary port for JWT operations
type JWTService interface {
	GenerateAccessToken(user *model.User) (string, time.Time, error)
	GenerateRefreshToken(userID int64) (string, time.Time, error)
	ValidateAccessToken(token string) (*primary.AccessTokenClaims, error)
	ValidateRefreshToken(token string) (*primary.RefreshTokenClaims, error)
}

// BcryptService defines the secondary port for password hashing
type BcryptService interface {
	HashPassword(password string) (string, error)
	ValidatePassword(password, hashedPassword string) error
}

// EncryptionService defines the secondary port for encryption operations
type EncryptionService interface {
	EncryptInternal(plaintext string) (string, error)
	DecryptInternal(ciphertext string) (string, error)
	EncryptURLSafe(plaintext string) (string, error)
	DecryptURLSafe(ciphertext string) (string, error)
}

// MailerService defines the secondary port for email operations
type MailerService interface {
	SendEmail(to string, subject string, body string) error
	LoadTemplate(templateName string, data any) (string, error)
}

// MediaService defines the secondary port for media/file operations
type MediaService interface {
	UploadFile(filename string, size int64, content any, filePath string) (*string, error)
	DeleteFile(filePath string) error
}
