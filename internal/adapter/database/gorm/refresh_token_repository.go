package gorm

import (
	"context"
	"time"

	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db       *gorm.DB
	baseRepo repositories.BaseRepository[entity.RefreshToken]
}

func NewRefreshTokenRepository(db *gorm.DB, baseRepo repositories.BaseRepository[entity.RefreshToken]) repositories.RefreshTokenRepository {
	return &RefreshTokenRepository{db: db, baseRepo: baseRepo}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, token *entity.RefreshToken) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	_, err := r.baseRepo.Create(ctx, token)
	return err
}

func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	return r.baseRepo.FindFirst(ctx, "token = ? AND is_revoked = false AND expires_at > ?", token, time.Now())
}

func (r *RefreshTokenRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.RefreshToken, error) {
	return r.baseRepo.Where(ctx, "user_id = ? AND is_revoked = false AND expires_at > ?", userID, time.Now())
}

func (r *RefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("user_id = ? AND is_revoked = false", userID).
		Update("is_revoked", true).Error
}

func (r *RefreshTokenRepository) RevokeByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error
}

func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entity.RefreshToken{}).Error
}

func (r *RefreshTokenRepository) IsTokenValid(ctx context.Context, token string) bool {
	var count int64
	r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("token = ? AND is_revoked = false AND expires_at > ?", token, time.Now()).
		Count(&count)
	return count > 0
}
