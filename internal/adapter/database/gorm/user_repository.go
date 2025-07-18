package gorm

import (
	"context"

	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db       *gorm.DB
	baseRepo repositories.BaseRepository[entity.User]
}

func NewUserRepository(db *gorm.DB, baseRepo repositories.BaseRepository[entity.User]) repositories.UserRepository {
	return &UserRepository{db: db, baseRepo: baseRepo}
}

func (r *UserRepository) FindAll(ctx context.Context, limit, offset int, search string) ([]*entity.User, int64, error) {
	return r.baseRepo.FindAll(ctx, limit, offset, "username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return r.baseRepo.FindByID(ctx, id)
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	return r.baseRepo.Create(ctx, user)
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	return r.baseRepo.Update(ctx, user)
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.baseRepo.Delete(ctx, id)
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
