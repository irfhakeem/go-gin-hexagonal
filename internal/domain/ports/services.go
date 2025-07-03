package ports

import (
	"go-gin-hexagonal/internal/domain/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AuthService interface {
	GenerateAccessToken(user *entity.User) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	ValidateAccessToken(tokenString string) (*AccessTokenClaims, error)
	ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error)
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}

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

type AccessTokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}
