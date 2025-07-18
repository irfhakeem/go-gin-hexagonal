package repositories

import (
	"context"
	"go-gin-hexagonal/internal/domain/entity"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *entity.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.RefreshToken, error)
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
	RevokeByToken(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
	IsTokenValid(ctx context.Context, token string) bool
}
