package repositories

import (
	"context"
	"go-gin-hexagonal/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, limit, offset int, search string) ([]*entity.User, int64, error)
	ExistsByEmail(ctx context.Context, email string) bool
	ExistsByUsername(ctx context.Context, username string) bool
}
