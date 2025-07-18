package test

import (
	"context"
	"go-gin-hexagonal/internal/application/service"
	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports/services"
	"go-gin-hexagonal/pkg/errors"
	mock_external "go-gin-hexagonal/test/mock/external"
	mock_repository "go-gin-hexagonal/test/mock/repository"
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

type UserTestSuite struct {
	suite.Suite
	mockRepo    *mock_repository.MockUserRepository
	mockHasher  *mock_external.MockSecurityService
	mockMailer  *mock_external.MockEmailService
	userService services.UserService
	ctx         context.Context
}

func (suite *UserTestSuite) SetupTest() {
	suite.mockRepo = mock_repository.NewMockUserRepository()
	suite.mockHasher = mock_external.NewMockSecurityService()
	suite.mockMailer = mock_external.NewMockEmailService()
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

func (suite *UserTestSuite) assertUserEquals(expected, actual *entity.User) {
	suite.Equal(expected.Email, actual.Email)
	suite.Equal(expected.Username, actual.Username)
	suite.Equal(expected.Name, actual.Name)
}

func (suite *UserTestSuite) TestUserRepoCRUD() {
	userID := uuid.New()
	user := createTestUserWithID(userID)

	// Create
	_, err := suite.mockRepo.Create(suite.ctx, user)
	suite.NoError(err)

	// Read
	retrievedUser, err := suite.mockRepo.FindByID(suite.ctx, userID)
	suite.NoError(err)
	suite.assertUserEquals(user, retrievedUser)

	// Update
	updatedUser := createTestUserWithID(userID)
	updatedUser.Name = updatedName
	_, err = suite.mockRepo.Update(suite.ctx, updatedUser)
	suite.NoError(err)

	// Verify update
	retrievedUser, err = suite.mockRepo.FindByID(suite.ctx, userID)
	suite.NoError(err)
	suite.Equal(updatedName, retrievedUser.Name)

	// Delete
	suite.NoError(suite.mockRepo.Delete(suite.ctx, userID))

	// Verify deletion
	_, err = suite.mockRepo.FindByID(suite.ctx, userID)
	suite.Equal(errors.ErrUserNotFound, err)
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
	id := uuid.New()
	user := createTestUserWithID(id)
	_, err := suite.mockRepo.Create(suite.ctx, user)
	suite.NoError(err)

	suite.Run("CreateDuplicate", func() {
		user = createTestUserWithID(id)
		_, err = suite.mockRepo.Create(suite.ctx, user)
		suite.Equal(errors.ErrUserAlreadyExists, err)
	})

	suite.Run("FindNonExistent", func() {
		_, err := suite.mockRepo.FindByID(suite.ctx, uuid.New())
		suite.Equal(errors.ErrUserNotFound, err)
	})

	suite.Run("UpdateNonExistent", func() {
		user := createTestUserWithID(uuid.New())
		_, err := suite.mockRepo.Update(suite.ctx, user)
		suite.Equal(errors.ErrUserNotFound, err)
	})

	suite.Run("DeleteNonExistent", func() {
		err := suite.mockRepo.Delete(suite.ctx, uuid.New())
		suite.Equal(errors.ErrUserNotFound, err)
	})
}

func (suite *UserTestSuite) TestUserServiceCL() {
	userID := uuid.New()
	user := createTestUserWithID(userID)
	user.Password = "hashedpassword"

	// Create
	_, err := suite.mockRepo.Create(suite.ctx, user)
	suite.NoError(err)

	// Read
	userInfo, err := suite.userService.GetUserByID(suite.ctx, userID)
	suite.NoError(err)
	suite.Equal(user.Email, userInfo.Email)

	// Update
	newName := "Updated Lifecycle User"
	updateReq := &services.UpdateUserRequest{Name: &newName}

	updatedUserInfo, err := suite.userService.UpdateUser(suite.ctx, userID, updateReq)
	suite.NoError(err)
	suite.Equal(newName, updatedUserInfo.Name)

	// Change password
	suite.mockHasher.On("Hash", "newpassword123").Return("newhashedpassword", nil)
	changePasswordReq := &services.ChangePasswordRequest{
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
	suite.Equal(errors.ErrUserNotFound, err)
}
