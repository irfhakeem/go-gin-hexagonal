package mock_service

import (
	"context"
	"go-gin-hexagonal/internal/domain/ports/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func NewMockUserService() *MockUserService {
	return &MockUserService{}
}

func (m *MockUserService) GetAllUsers(ctx context.Context, page, pageSize int, search string) (*services.UserPaginationResponse, error) {
	args := m.Called(ctx, page, pageSize, search)
	return args.Get(0).(*services.UserPaginationResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*services.UserInfo, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*services.UserInfo), args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, req *services.CreateUserRequest) (*services.UserInfo, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*services.UserInfo), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, userID uuid.UUID, req *services.UpdateUserRequest) (*services.UserInfo, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).(*services.UserInfo), args.Error(1)
}

func (m *MockUserService) ChangePassword(ctx context.Context, userID uuid.UUID, req *services.ChangePasswordRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
