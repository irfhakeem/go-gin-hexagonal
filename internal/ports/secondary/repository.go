package secondary

import (
	"context"
	"go-gin-clean/internal/domain/model"
)

// UserRepository defines the secondary port for user persistence
type UserRepository interface {
	FindAll(ctx context.Context, limit, offset int, search string) ([]*model.User, int64, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	ExistsByEmail(ctx context.Context, email string) bool
}

// RefreshTokenRepository defines the secondary port for refresh token persistence
type RefreshTokenRepository interface {
	Save(ctx context.Context, token *model.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	FindByUserID(ctx context.Context, userID int64) ([]*model.RefreshToken, error)
	RevokeAllByUserID(ctx context.Context, userID int64) error
	RevokeByToken(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
	IsTokenValid(ctx context.Context, token string) bool
}

// BaseRepository is a generic repository interface
type BaseRepository[T any] interface {
	FindAll(ctx context.Context, limit, offset int, query any, args ...any) ([]*T, int64, error)
	FindByID(ctx context.Context, id int64) (*T, error)
	FindFirst(ctx context.Context, query any, args ...any) (*T, error)
	Where(ctx context.Context, query any, args ...any) ([]*T, error)
	WhereExisting(ctx context.Context, query any, args ...any) (bool, error)
	Create(ctx context.Context, entity *T) (*T, error)
	Update(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id int64) error
}
