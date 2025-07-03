package test

import (
	"context"
	"go-gin-hexagonal/internal/application/service"
	"go-gin-hexagonal/internal/domain/dto"
	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (r *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, ports.ErrUserNotFound
}

func (r *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ports.ErrUserNotFound
}

func (r *MockUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, ports.ErrUserNotFound
}

func (r *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	if _, exists := r.users[user.ID]; !exists {
		return ports.ErrUserNotFound
	}
	user.UpdatedAt = time.Now()
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

func (r *MockUserRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, int64, error) {
	var users []*entity.User
	for _, user := range r.users {
		users = append(users, user)
	}

	total := int64(len(users))
	start := offset
	end := offset + limit
	if start > len(users) {
		start = len(users)
	}
	if end > len(users) {
		end = len(users)
	}

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

func TestUserRepositoryIntegration(t *testing.T) {
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

	ctx := context.Background()

	// Create user
	assert.NoError(t, mockRepo.Create(ctx, user))

	// Get user by ID
	retrievedUser, err := mockRepo.FindByID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, "johndoe@example.com", retrievedUser.Email)

	// List users
	_, total, err := mockRepo.FindAll(ctx, 10, 0)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))

	// Update user
	updatedUser := &entity.User{
		ID:       userID,
		Email:    "johndoeUPDATED@example.com",
		Username: "johndoeUPDATED",
		Name:     "John Doe Updated",
	}
	assert.NoError(t, mockRepo.Update(ctx, updatedUser))

	// Verify update
	retrievedUser, err = mockRepo.FindByID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, "johndoeUPDATED@example.com", retrievedUser.Email)

	// Delete user
	assert.NoError(t, mockRepo.Delete(ctx, userID))

	// Verify deletion
	_, err = mockRepo.FindByID(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrUserNotFound, err)
}

func TestUserRepository_EdgeCases(t *testing.T) {
	mockRepo := NewMockUserRepository()
	ctx := context.Background()

	t.Run("Create user with existing ID should fail", func(t *testing.T) {
		var existingID uuid.UUID
		for id := range mockRepo.users {
			existingID = id
			break
		}

		user := &entity.User{
			ID:       existingID,
			Email:    "duplicate@example.com",
			Username: "duplicate",
		}

		err := mockRepo.Create(ctx, user)
		assert.Equal(t, ports.ErrUserAlreadyExists, err)
	})

	t.Run("Get non-existent user should fail", func(t *testing.T) {
		nonExistentID := uuid.New()
		_, err := mockRepo.FindByID(ctx, nonExistentID)
		assert.Equal(t, ports.ErrUserNotFound, err)
	})

	t.Run("Update non-existent user should fail", func(t *testing.T) {
		user := &entity.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}
		err := mockRepo.Update(ctx, user)
		assert.Equal(t, ports.ErrUserNotFound, err)
	})

	t.Run("Delete non-existent user should fail", func(t *testing.T) {
		err := mockRepo.Delete(ctx, uuid.New())
		assert.Equal(t, ports.ErrUserNotFound, err)
	})

	t.Run("Find user by email", func(t *testing.T) {
		user, err := mockRepo.FindByEmail(ctx, "johndoe100@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "johndoe100@example.com", user.Email)
	})

	t.Run("Check if email exists", func(t *testing.T) {
		exists := mockRepo.ExistsByEmail(ctx, "johndoe100@example.com")
		assert.True(t, exists)

		exists = mockRepo.ExistsByEmail(ctx, "nonexistent@example.com")
		assert.False(t, exists)
	})
}

// Integration Tests
func TestUserServiceIntegration(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockHasher := NewMockSecurityService()

	// Setup mock hasher
	mockHasher.On("Hash", mock.AnythingOfType("string")).Return("hashedpassword", nil)
	mockHasher.On("Verify", "hashedpassword", "password123").Return(nil)

	userService := service.NewUserService(mockRepo, mockHasher)
	ctx := context.Background()

	t.Run("Complete user lifecycle", func(t *testing.T) {
		userID := uuid.New()

		user := &entity.User{
			ID:        userID,
			Email:     "lifecycle@example.com",
			Username:  "lifecycle",
			Password:  "hashedpassword",
			Name:      "Lifecycle User",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := mockRepo.Create(ctx, user)
		assert.NoError(t, err)

		// Read via service
		userInfo, err := userService.GetProfile(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, userInfo.Email)

		// Update via service
		newName := "Updated Lifecycle User"
		updateReq := &dto.UpdateUserRequest{
			Name: &newName,
		}

		updatedUserInfo, err := userService.UpdateProfile(ctx, userID, updateReq)
		assert.NoError(t, err)
		assert.Equal(t, newName, updatedUserInfo.Name)

		// Change password via service
		changePasswordReq := &dto.ChangePasswordRequest{
			CurrentPassword: "password123",
			NewPassword:     "newpassword123",
		}

		mockHasher.On("Hash", "newpassword123").Return("newhashedpassword", nil)

		err = userService.ChangePassword(ctx, userID, changePasswordReq)
		assert.NoError(t, err)

		// Delete via service
		err = userService.DeleteUser(ctx, userID)
		assert.NoError(t, err)

		// Verify deletion via service
		_, err = userService.GetProfile(ctx, userID)
		assert.Equal(t, ports.ErrUserNotFound, err)
	})
}
