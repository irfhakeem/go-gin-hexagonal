package database

import (
	"context"
	"time"

	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db       *gorm.DB
	baseRepo ports.BaseRepository[entity.User]
}

func NewUserRepository(db *gorm.DB, baseRepo ports.BaseRepository[entity.User]) ports.UserRepository {
	return &UserRepository{db: db, baseRepo: baseRepo}
}

func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return r.baseRepo.FindByID(ctx, id)
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	_, err := r.baseRepo.Create(ctx, user)
	return err
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.baseRepo.FindFirst(ctx, "email = ?", email)
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	return r.baseRepo.FindFirst(ctx, "username = ?", username)
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) bool {
	isExist, _ := r.baseRepo.WhereExisting(ctx, "email = ?", email)
	return isExist
}

func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) bool {
	isExist, _ := r.baseRepo.WhereExisting(ctx, "username = ?", username)
	return isExist
}
