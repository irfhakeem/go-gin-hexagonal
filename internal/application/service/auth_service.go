package service

import (
	"context"

	"go-gin-hexagonal/internal/application/dto"
	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo         ports.UserRepository
	refreshTokenRepo ports.RefreshTokenRepository
	tokenManager     ports.TokenManager
	passwordHasher   ports.PasswordHasher
}

func NewAuthService(
	userRepo ports.UserRepository,
	refreshTokenRepo ports.RefreshTokenRepository,
	tokenManager ports.TokenManager,
	passwordHasher ports.PasswordHasher,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		tokenManager:     tokenManager,
		passwordHasher:   passwordHasher,
	}
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ports.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ports.ErrInvalidCredentials
	}

	if err := s.passwordHasher.Verify(user.Password, req.Password); err != nil {
		return nil, ports.ErrInvalidCredentials
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshTokenEntity := &entity.RefreshToken{
		UserID: user.ID,
		Token:  refreshToken,
	}
	if err := s.refreshTokenRepo.Save(ctx, refreshTokenEntity); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) error {
	if s.userRepo.ExistsByEmail(ctx, req.Email) {
		return ports.ErrUserAlreadyExists
	}

	if s.userRepo.ExistsByUsername(ctx, req.Username) {
		return ports.ErrUserAlreadyExists
	}

	hashedPassword, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		return err
	}

	user := &entity.User{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Name:     req.Name,
		IsActive: true,
	}

	return s.userRepo.Create(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	claims, err := s.tokenManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, ports.ErrTokenInvalid
	}

	if !s.refreshTokenRepo.IsTokenValid(ctx, req.RefreshToken) {
		return nil, ports.ErrTokenInvalid
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, ports.ErrUserNotFound
	}

	newAccessToken, err := s.tokenManager.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	s.refreshTokenRepo.RevokeByToken(ctx, req.RefreshToken)

	newRefreshTokenEntity := &entity.RefreshToken{
		UserID: user.ID,
		Token:  newRefreshToken,
	}
	if err := s.refreshTokenRepo.Save(ctx, newRefreshTokenEntity); err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.refreshTokenRepo.RevokeAllByUserID(ctx, userID)
}
