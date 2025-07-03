package test

import (
	"context"
	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockUserRepository struct {
	users map[uuid.UUID]*entity.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: map[uuid.UUID]*entity.User{
			uuid.New(): {
				Email:     "johndoe100@example.com",
				Username:  "johndoe100",
				Password:  "password123",
				Name:      "John Doe 100",
				IsActive:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: nil,
			},
		},
	}
}

func (r *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	if _, exists := r.users[user.ID]; exists {
		return ports.ErrUserAlreadyExists
	}
	r.users[user.ID] = user
	return nil
}

func (r *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, ports.ErrUserNotFound
}

func (r *MockUserRepository) ListUsers(ctx context.Context, page, pageSize int, search string) ([]*entity.User, int, error) {
	var users []*entity.User
	for _, user := range r.users {
		if search == "" || (user.Name != "" && user.Name == search) || (user.Email != "" && user.Email == search) {
			users = append(users, user)
		}
	}

	total := len(users)
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	return users[start:end], total, nil
}

func (r *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	if _, exists := r.users[user.ID]; !exists {
		return ports.ErrUserNotFound
	}
	r.users[user.ID] = user
	return nil
}

func (r *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, exists := r.users[id]; !exists {
		return ports.ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}

func TestCRUDUser(t *testing.T) {
	log.Println("Running User Repository Tests")

	mockRepo := NewMockUserRepository()
	userID := uuid.New()

	user := &entity.User{
		ID:        userID,
		Email:     "johndoe@example.com",
		Username:  "johndoe",
		Password:  "password123",
		Name:      "John Doe",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	// Create user
	assert.NoError(t, mockRepo.Create(context.Background(), user))

	// Get user by ID
	retrievedUser, err := mockRepo.GetByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, "johndoe@example.com", retrievedUser.Email)

	// List users
	_, total, err := mockRepo.ListUsers(context.Background(), 1, 10, "")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, total, 2)

	// Update user
	updatedUser := &entity.User{
		ID:       userID,
		Email:    "johndoeUPDATED@example.com",
		Username: "johndoeUPDATED",
		Name:     "John Doe Updated",
	}
	assert.NoError(t, mockRepo.Update(context.Background(), updatedUser))

	// Verify update
	retrievedUser, err = mockRepo.GetByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, "johndoeUPDATED@example.com", retrievedUser.Email)

	// Delete user
	assert.NoError(t, mockRepo.Delete(context.Background(), userID))

	// Verify deletion
	_, err = mockRepo.GetByID(context.Background(), userID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrUserNotFound, err)
}
