package test

import (
	"context"
	"go-gin-hexagonal/internal/application/service"
	"go-gin-hexagonal/internal/domain/dto"
	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports"
	mockAdapter "go-gin-hexagonal/test/mock"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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

// README:
// CRUD = Create, Read, Update, Delete
// FO = Find Operations
// EO = Exists Operations
// EC = Error Cases
// CL = Complete Lifecycle

type MockUserRepository struct {
	users map[uuid.UUID]*entity.User
}

type UserTestSuite struct {
	suite.Suite
	mockRepo    *MockUserRepository
	mockHasher  *mockAdapter.MockSecurityService
	mockMailer  *mockAdapter.MockEmailService
	userService ports.UserService
	ctx         context.Context
}

func (suite *UserTestSuite) SetupTest() {
	suite.mockRepo = NewMockUserRepository()
	suite.mockHasher = mockAdapter.NewMockSecurityService()
	suite.mockMailer = mockAdapter.NewMockEmailService()
	suite.userService = service.NewUserService(suite.mockRepo, suite.mockHasher, suite.mockMailer)
	suite.ctx = context.Background()

	// Setup common mock expectations
	suite.mockHasher.On("Hash", mock.AnythingOfType("string")).Return("hashedpassword", nil)
	suite.mockHasher.On("Verify", "hashedpassword", testPassword).Return(nil)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func createTestUser() *entity.User {
	return &entity.User{
		Email:    testEmail,
		Username: testUsername,
		Password: testPassword,
		Name:     testName,
		IsActive: true,
		AuditInfo: entity.AuditInfo{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
			IsDeleted: false,
		},
	}
}

func createTestUserWithID(id uuid.UUID) *entity.User {
	user := createTestUser()
	user.ID = id
	return user
}

func NewMockUserRepository() *MockUserRepository {
	testUser := createTestUser()
	return &MockUserRepository{
		users: map[uuid.UUID]*entity.User{
			uuid.New(): testUser,
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

func (suite *UserTestSuite) assertUserEquals(expected, actual *entity.User) {
	suite.Equal(expected.Email, actual.Email)
	suite.Equal(expected.Username, actual.Username)
	suite.Equal(expected.Name, actual.Name)
}

func (suite *UserTestSuite) getExistingUserID() uuid.UUID {
	for id := range suite.mockRepo.users {
		return id
	}
	return uuid.New()
}

func (suite *UserTestSuite) TestUserRepoCRUD() {
	userID := uuid.New()
	user := createTestUserWithID(userID)

	// Create
	suite.NoError(suite.mockRepo.Create(suite.ctx, user))

	// Read
	retrievedUser, err := suite.mockRepo.FindByID(suite.ctx, userID)
	suite.NoError(err)
	suite.assertUserEquals(user, retrievedUser)

	// Update
	updatedUser := createTestUserWithID(userID)
	updatedUser.Name = updatedName
	suite.NoError(suite.mockRepo.Update(suite.ctx, updatedUser))

	// Verify update
	retrievedUser, err = suite.mockRepo.FindByID(suite.ctx, userID)
	suite.NoError(err)
	suite.Equal(updatedName, retrievedUser.Name)

	// Delete
	suite.NoError(suite.mockRepo.Delete(suite.ctx, userID))

	// Verify deletion
	_, err = suite.mockRepo.FindByID(suite.ctx, userID)
	suite.Equal(ports.ErrUserNotFound, err)
}

func (suite *UserTestSuite) TestUserRepoFO() {
	suite.Run("FindByEmail", func() {
		user, err := suite.mockRepo.FindByEmail(suite.ctx, testEmail)
		suite.NoError(err)
		suite.Equal(testEmail, user.Email)
	})

	suite.Run("FindByUsername", func() {
		user, err := suite.mockRepo.FindByUsername(suite.ctx, testUsername)
		suite.NoError(err)
		suite.Equal(testUsername, user.Username)
	})

	suite.Run("FindAll", func() {
		users, total, err := suite.mockRepo.FindAll(suite.ctx, 10, 0, "")
		suite.NoError(err)
		suite.GreaterOrEqual(total, int64(1))
		suite.LessOrEqual(len(users), 10)
	})
}

func (suite *UserTestSuite) TestUserRepoEO() {
	suite.Run("ExistsByEmail", func() {
		exists := suite.mockRepo.ExistsByEmail(suite.ctx, testEmail)
		suite.True(exists)

		exists = suite.mockRepo.ExistsByEmail(suite.ctx, "nonexistent@example.com")
		suite.False(exists)
	})

	suite.Run("ExistsByUsername", func() {
		exists := suite.mockRepo.ExistsByUsername(suite.ctx, testUsername)
		suite.True(exists)

		exists = suite.mockRepo.ExistsByUsername(suite.ctx, "nonexistentuser")
		suite.False(exists)
	})
}

func (suite *UserTestSuite) TestUserRepoEC() {
	suite.Run("CreateDuplicate", func() {
		existingID := suite.getExistingUserID()
		user := createTestUserWithID(existingID)
		err := suite.mockRepo.Create(suite.ctx, user)
		suite.Equal(ports.ErrUserAlreadyExists, err)
	})

	suite.Run("FindNonExistent", func() {
		_, err := suite.mockRepo.FindByID(suite.ctx, uuid.New())
		suite.Equal(ports.ErrUserNotFound, err)
	})

	suite.Run("UpdateNonExistent", func() {
		user := createTestUserWithID(uuid.New())
		err := suite.mockRepo.Update(suite.ctx, user)
		suite.Equal(ports.ErrUserNotFound, err)
	})

	suite.Run("DeleteNonExistent", func() {
		err := suite.mockRepo.Delete(suite.ctx, uuid.New())
		suite.Equal(ports.ErrUserNotFound, err)
	})
}

func (suite *UserTestSuite) TestUserServiceCL() {
	userID := uuid.New()
	user := createTestUserWithID(userID)
	user.Password = "hashedpassword"

	// Create
	suite.NoError(suite.mockRepo.Create(suite.ctx, user))

	// Read
	userInfo, err := suite.userService.GetUserByID(suite.ctx, userID)
	suite.NoError(err)
	suite.Equal(user.Email, userInfo.Email)

	// Update
	newName := "Updated Lifecycle User"
	updateReq := &dto.UpdateUserRequest{Name: &newName}

	updatedUserInfo, err := suite.userService.UpdateUser(suite.ctx, userID, updateReq)
	suite.NoError(err)
	suite.Equal(newName, updatedUserInfo.Name)

	// Change password
	suite.mockHasher.On("Hash", "newpassword123").Return("newhashedpassword", nil)
	changePasswordReq := &dto.ChangePasswordRequest{
		CurrentPassword: testPassword,
		NewPassword:     "newpassword123",
	}

	err = suite.userService.ChangePassword(suite.ctx, userID, changePasswordReq)
	suite.NoError(err)

	// Delete
	err = suite.userService.DeleteUser(suite.ctx, userID)
	suite.NoError(err)

	// Verify deletion
	_, err = suite.userService.GetUserByID(suite.ctx, userID)
	suite.Equal(ports.ErrUserNotFound, err)
}
