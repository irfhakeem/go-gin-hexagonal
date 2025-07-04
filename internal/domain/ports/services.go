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
	VerifyEmail(ctx context.Context, otp string) error
	RequestVerifyEmail(ctx context.Context, email string) error
	RequestResetPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error
}

type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.UserInfo, error)
	GetAllUsers(ctx context.Context, req *dto.UserListRequest) (*dto.UserListResponse, error)
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserInfo, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserInfo, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type EmailService interface {
	SendNewUserEmail(to string, data *dto.NewUserData) error
	SendVerifyEmail(to string, data *dto.VerifyEmailData) error
	SendRequestResetPassword(to string, data *dto.RequestResetPasswordData) error
}
