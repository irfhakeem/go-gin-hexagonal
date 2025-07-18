package repositories

import (
	"context"

	"github.com/google/uuid"
)

type BaseRepository[T any] interface {
	Raw(ctx context.Context, query string) ([]*T, error)
	FindAll(ctx context.Context, limit, offset int, query any, args ...any) ([]*T, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*T, error)
	FindFirst(ctx context.Context, query any, args ...any) (*T, error)
	Where(ctx context.Context, query any, args ...any) ([]*T, error)
	WhereExisting(ctx context.Context, query any, args ...any) (bool, error)
	Create(ctx context.Context, entity *T) (*T, error)
	Update(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
