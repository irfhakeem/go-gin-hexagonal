package service

import (
	"context"
	"log"
	"math"

	"go-gin-hexagonal/internal/domain/dto"
	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports"
	"go-gin-hexagonal/pkg/utils"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo       ports.UserRepository
	passwordHasher ports.PasswordHasher
	emailService   ports.EmailService
}

func NewUserService(
	userRepo ports.UserRepository,
	passwordHasher ports.PasswordHasher,
	emailService ports.EmailService,
) ports.UserService {
	return &UserService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		emailService:   emailService,
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

func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ports.ErrUserNotFound
	}

	return FormatUserInfo(user), nil
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserInfo, error) {
	if req.Email == "" || req.Name == "" {
		return nil, ports.ErrInvalidInput
	}

	if s.userRepo.ExistsByEmail(ctx, req.Email) {
		return nil, ports.ErrUserAlreadyExists
	}

	var username string
	var usernameExists = true
	for usernameExists {
		username = utils.GenerateUsername(req.Name)

		existingUser, err := s.userRepo.FindByUsername(ctx, username)
		if err != nil && err != ports.ErrUserNotFound {
			return nil, err
		}
		if existingUser == nil {
			usernameExists = false
		}
	}

	password := utils.GeneratePassword(8, true)
	hashedPassword, err := s.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:    req.Email,
		Username: username,
		Password: hashedPassword,
		Name:     req.Name,
		IsActive: true,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	go func(email string, username string, password string) {
		newUserData := &dto.NewUserData{
			UserEmail: email,
			Username:  username,
			Password:  password,
		}

		err := s.emailService.SendNewUserEmail(email, newUserData)
		if err != nil {
			log.Printf("failed to send new user email: %v", err)
		}
	}(req.Email, username, password)

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

func (s *UserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return ports.ErrDeleteUser
	}
	return nil
}
