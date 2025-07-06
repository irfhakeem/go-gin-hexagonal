package database

import (
	"context"
	"go-gin-hexagonal/internal/domain/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) ports.BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Raw(ctx context.Context, query string) ([]T, error) {
	var entities []T
	if err := r.db.WithContext(ctx).Raw(query).Scan(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, limit, offset int, search string, query any, args ...any) ([]T, int64, error) {
	var entities []T
	var count int64

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	var entity T
	q := r.db.WithContext(ctx).Model(&entity)
	q = q.Where(query, args...)

	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Limit(limit).Offset(offset).Order("id asc").Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, count, nil
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) FindFirst(ctx context.Context, query any, args ...any) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).Where(query, args...).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Where(ctx context.Context, query any, args ...any) ([]*T, error) {
	var entities []*T
	if err := r.db.WithContext(ctx).Where(query, args...).Order("id asc").Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *BaseRepository[T]) WhereExisting(ctx context.Context, query any, args ...any) (bool, error) {
	var entity T
	err := r.db.WithContext(ctx).Where(query, args...).First(&entity).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) (*T, error) {
	if err := r.db.WithContext(ctx).Updates(&entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(new(T), "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
