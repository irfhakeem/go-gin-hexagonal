package mock_repository

import (
	"context"
	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/pkg/errors"
	"time"

	"github.com/google/uuid"
)

const (
	testEmail       = "johndoe100@example.com"
	testUsername    = "johndoe100"
	testPassword    = "password123"
	testName        = "John Doe 100"
	updatedEmail    = "johndoeUPDATED@example.com"
	updatedUsername = "johndoeUPDATED"
	updatedName     = "John Doe Updated"
)

type MockUserRepository struct {
	users map[uuid.UUID]*entity.User
}

func NewMockUserRepository() *MockUserRepository {
	id := uuid.New()
	return &MockUserRepository{
		users: map[uuid.UUID]*entity.User{
			id: {
				ID:       id,
				Email:    testEmail,
				Username: testUsername,
				Password: "hashedpassword",
				Name:     testName,
				IsActive: true,
				AuditInfo: entity.AuditInfo{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					DeletedAt: nil,
					IsDeleted: false,
				},
			},
		},
	}
}

func (r *MockUserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	_, exists := r.users[user.ID]
	if exists {
		return nil, errors.ErrUserAlreadyExists
	}

	r.users[user.ID] = user
	return user, nil
}

func (r *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, errors.ErrUserNotFound
}

func (r *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.ErrUserNotFound
}

func (r *MockUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.ErrUserNotFound
}

func (r *MockUserRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	_, exists := r.users[user.ID]
	if !exists {
		return nil, errors.ErrUserNotFound
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return user, nil
}

func (r *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, exists := r.users[id]; !exists {
		return errors.ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}

func (r *MockUserRepository) FindAll(ctx context.Context, limit, offset int, search string) ([]*entity.User, int64, error) {
	var users []*entity.User
	for _, user := range r.users {
		users = append(users, user)
	}

	total := int64(len(users))
	start := min(offset, len(users))
	end := min(offset+limit, len(users))

	return users[start:end], total, nil
}

func (r *MockUserRepository) ExistsByEmail(ctx context.Context, email string) bool {
	for _, user := range r.users {
		if user.Email == email {
			return true
		}
	}
	return false
}

func (r *MockUserRepository) ExistsByUsername(ctx context.Context, username string) bool {
	for _, user := range r.users {
		if user.Username == username {
			return true
		}
	}
	return false
}
