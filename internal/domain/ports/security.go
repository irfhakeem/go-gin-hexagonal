package ports

import (
	"go-gin-hexagonal/internal/domain/dto"
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
	ValidateAccessToken(token string) (*dto.AccessTokenClaims, error)
	ValidateRefreshToken(token string) (*dto.RefreshTokenClaims, error)
}

type Encryptor interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}
