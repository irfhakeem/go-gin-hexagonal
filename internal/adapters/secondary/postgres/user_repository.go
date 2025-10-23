package postgres

import (
	"context"
	"go-gin-clean/internal/domain/model"
	"go-gin-clean/internal/ports/secondary"

	"gorm.io/gorm"
)

type UserRepository struct {
	db       *gorm.DB
	baseRepo secondary.BaseRepository[model.User]
}

func NewUserRepository(db *gorm.DB) secondary.UserRepository {
	baseRepo := NewBaseRepository[model.User](db)
	return &UserRepository{
		db:       db,
		baseRepo: baseRepo,
	}
}

func (r *UserRepository) FindAll(ctx context.Context, limit, offset int, search string) ([]*model.User, int64, error) {
	return r.baseRepo.FindAll(ctx, limit, offset, "name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return r.baseRepo.FindByID(ctx, id)
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	return r.baseRepo.Create(ctx, user)
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	return r.baseRepo.Update(ctx, user)
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	return r.baseRepo.Delete(ctx, id)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.baseRepo.FindFirst(ctx, "email = ?", email)
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) bool {
	isExist, _ := r.baseRepo.WhereExisting(ctx, "email = ?", email)
	return isExist
}
