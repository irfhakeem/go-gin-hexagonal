package service

import (
	"context"
	"log"
	"math"

	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports"
	"go-gin-hexagonal/internal/domain/ports/repositories"
	"go-gin-hexagonal/internal/domain/ports/services"
	"go-gin-hexagonal/pkg/errors"
	"go-gin-hexagonal/pkg/utils"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo       repositories.UserRepository
	passwordHasher ports.PasswordHasher
	emailService   services.EmailService
}

func NewUserService(
	userRepo repositories.UserRepository,
	passwordHasher ports.PasswordHasher,
	emailService services.EmailService,
) services.UserService {
	return &UserService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		emailService:   emailService,
	}
}

func FormatUserInfo(user *entity.User) *services.UserInfo {
	return &services.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Name:      user.Name,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context, page, pageSize int, search string) (*services.UserPaginationResponse, error) {
	offset := (page - 1) * pageSize
	users, total, err := s.userRepo.FindAll(ctx, pageSize, offset, search)
	if err != nil {
		return nil, err
	}

	var userInfos []*services.UserInfo
	for _, user := range users {
		userInfos = append(userInfos, FormatUserInfo(user))
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &services.UserPaginationResponse{
		Datas:      userInfos,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*services.UserInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	return FormatUserInfo(user), nil
}

func (s *UserService) CreateUser(ctx context.Context, req *services.CreateUserRequest) (*services.UserInfo, error) {
	if req.Email == "" || req.Name == "" {
		return nil, errors.ErrInvalidInput
	}

	if s.userRepo.ExistsByEmail(ctx, req.Email) {
		return nil, errors.ErrUserAlreadyExists
	}

	var username string
	var usernameExists = true
	for usernameExists {
		username = utils.GenerateUsername(req.Name)

		existingUser := s.userRepo.ExistsByUsername(ctx, username)
		if !existingUser {
			usernameExists = false
			break
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

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	go func(email string, username string, password string) {
		newUserData := &services.NewUserEmailData{
			Email:    email,
			Password: password,
		}

		err := s.emailService.SendNewUserEmail(email, newUserData)
		if err != nil {
			log.Printf("failed to send new user email: %v", err)
		}
	}(req.Email, username, password)

	return FormatUserInfo(createdUser), nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID uuid.UUID, req *services.UpdateUserRequest) (*services.UserInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Username != nil {
		if s.userRepo.ExistsByUsername(ctx, *req.Username) {
			existingUser, _ := s.userRepo.FindByUsername(ctx, *req.Username)
			if existingUser != nil && existingUser.ID != userID {
				return nil, errors.ErrUserAlreadyExists
			}
		}
		user.Username = *req.Username
	}

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return FormatUserInfo(updatedUser), nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID uuid.UUID, req *services.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	err = s.passwordHasher.Verify(user.Password, req.CurrentPassword)
	if err != nil {
		return errors.ErrInvalidCredentials
	}

	hashedPassword, err := s.passwordHasher.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	_, err = s.userRepo.Update(ctx, user)

	return err
}

func (s *UserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return errors.ErrDeleteUser
	}
	return nil
}
