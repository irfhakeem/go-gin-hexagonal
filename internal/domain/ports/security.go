package ports

import (
	"go-gin-hexagonal/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) error
}

type TokenManager interface {
	GenerateAccessToken(user *entity.User) (string, time.Time, error)
	GenerateRefreshToken(userID uuid.UUID) (string, time.Time, error)
	ValidateAccessToken(token string) (*AccessTokenClaims, error)
	ValidateRefreshToken(token string) (*RefreshTokenClaims, error)
}

type Encryptor interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type ExtractInfo struct {
	UserID   uuid.UUID
	TenantID uuid.UUID
	RoleID   int64
}

type AccessTokenClaims struct {
	UserID    uuid.UUID
	Email     string
	Username  string
	TokenType string
	ExpiresAt time.Time
	IssuedAt  time.Time
	NotBefore time.Time
	Issuer    string
	Subject   string
}

type RefreshTokenClaims struct {
	UserID    uuid.UUID
	TokenType string
	ExpiresAt time.Time
	IssuedAt  time.Time
	NotBefore time.Time
	Issuer    string
	Subject   string
}
