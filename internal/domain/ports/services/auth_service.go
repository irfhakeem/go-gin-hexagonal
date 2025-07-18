package services

import (
	"context"

	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterRequest) error
	RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	VerifyEmail(ctx context.Context, token string) error
	SendVerifyEmail(ctx context.Context, email string) error
	SendResetPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
}

type RegisterRequest struct {
	Email    string
	Username string
	Password string
	Name     string
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenResponse struct {
	AccessToken  string
	RefreshToken string
}

type ResetPasswordRequest struct {
	Token           string
	NewPassword     string
	ConfirmPassword string
}
