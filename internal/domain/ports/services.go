package ports

import (
	"context"

	"go-gin-hexagonal/internal/domain/dto"

	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) error
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
}

type UserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserInfo, error)
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserInfo, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserInfo, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) error
	ListUsers(ctx context.Context, req *dto.UserListRequest) (*dto.UserListResponse, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type EmailService interface {
	SendNewUserEmail(to string, data *dto.NewUserData) error
}
