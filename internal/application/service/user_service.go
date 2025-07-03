package service

import (
	"context"
	"math"

	"go-gin-hexagonal/internal/application/dto"
	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo       ports.UserRepository
	passwordHasher ports.PasswordHasher
}

func NewUserService(
	userRepo ports.UserRepository,
	passwordHasher ports.PasswordHasher,
) *UserService {
	return &UserService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func FormatUserInfo(user *entity.User) *dto.UserInfo {
	return &dto.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Name:      user.Name,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ports.ErrUserNotFound
	}

	return FormatUserInfo(user), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ports.ErrUserNotFound
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Username != nil {
		if s.userRepo.ExistsByUsername(ctx, *req.Username) {
			existingUser, _ := s.userRepo.FindByUsername(ctx, *req.Username)
			if existingUser != nil && existingUser.ID != userID {
				return nil, ports.ErrUserAlreadyExists
			}
		}
		user.Username = *req.Username
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return FormatUserInfo(user), nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ports.ErrUserNotFound
	}

	err = s.passwordHasher.Verify(user.Password, req.CurrentPassword)
	if err != nil {
		return ports.ErrInvalidCredentials
	}

	hashedPassword, err := s.passwordHasher.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) ListUsers(ctx context.Context, req *dto.UserListRequest) (*dto.UserListResponse, error) {
	offset := (req.Page - 1) * req.PageSize
	users, total, err := s.userRepo.FindAll(ctx, req.PageSize, offset)
	if err != nil {
		return nil, err
	}

	var userInfos []*dto.UserInfo
	for _, user := range users {
		userInfos = append(userInfos, FormatUserInfo(user))
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &dto.UserListResponse{
		Users:      userInfos,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}
