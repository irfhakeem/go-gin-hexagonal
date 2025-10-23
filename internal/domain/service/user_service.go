package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	domerr "go-gin-clean/internal/domain/error"
	"go-gin-clean/internal/domain/model"
	"go-gin-clean/internal/ports/primary"
	"go-gin-clean/internal/ports/secondary"
	"go-gin-clean/pkg/config"
)

type UserService struct {
	userRepo         secondary.UserRepository
	emailService     primary.EmailUseCase
	refreshTokenRepo secondary.RefreshTokenRepository
	jwtService       secondary.JWTService
	bcryptService    secondary.BcryptService
	aesService       secondary.EncryptionService
	mediaService     secondary.MediaService
}

func NewUserService(
	userRepo secondary.UserRepository,
	emailService primary.EmailUseCase,
	refreshTokenRepo secondary.RefreshTokenRepository,
	jwtService secondary.JWTService,
	bcryptService secondary.BcryptService,
	aesService secondary.EncryptionService,
	mediaService secondary.MediaService,
) primary.UserUseCase {
	return &UserService{
		userRepo:         userRepo,
		emailService:     emailService,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		bcryptService:    bcryptService,
		aesService:       aesService,
		mediaService:     mediaService,
	}
}

func FormatUserInfo(user *model.User) *primary.UserInfo {
	return &primary.UserInfo{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Gender:   user.Gender,
		Avatar:   user.Avatar,
		IsActive: user.IsActive,
	}
}

func (s *UserService) Login(ctx context.Context, req *primary.LoginRequest) (*primary.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, domerr.ErrUserNotFound
	}

	if !user.IsActive {
		return nil, domerr.ErrUserNotFound
	}

	if err := s.bcryptService.ValidatePassword(req.Password, user.Password); err != nil {
		return nil, domerr.ErrPasswordNotMatch
	}

	accessToken, _, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, expiryAt, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	hashedRefreshToken, err := s.aesService.EncryptInternal(refreshToken)
	if err != nil {
		return nil, err
	}

	token := model.NewRefreshToken(user.ID, hashedRefreshToken, expiryAt, false, *user)

	if err := s.refreshTokenRepo.Save(ctx, token); err != nil {
		return nil, err
	}

	return &primary.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *FormatUserInfo(user),
	}, nil
}

func (s *UserService) Register(ctx context.Context, req *primary.RegisterRequest) error {
	if s.userRepo.ExistsByEmail(ctx, req.Email) {
		return domerr.ErrEmailAlreadyExists
	}

	hashedPassword, err := s.bcryptService.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user, err := model.NewUser(req.Name, req.Email, hashedPassword, "", model.Unknown)
	if err != nil {
		return err
	}

	savedUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	plainText := fmt.Sprintf("%d_%s", savedUser.ID, time.Now().Add(24*time.Hour).Format(time.RFC3339))

	token, err := s.aesService.EncryptURLSafe(plainText)
	if err != nil {
		return err
	}

	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", config.GetAppURL(), token)

	go func() {
		if err := s.emailService.SendVerifyEmail(user.Email, user.Name, verificationURL); err != nil {
			log.Printf("Failed to send verification email to %s: %v", user.Email, err)
		}
	}()

	return nil
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (*primary.RefreshTokenResponse, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, domerr.ErrTokenInvalid
	}

	if !s.refreshTokenRepo.IsTokenValid(ctx, refreshToken) {
		return nil, domerr.ErrTokenInvalid
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, domerr.ErrUserNotFound
	}

	newAccessToken, _, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, expiryAt, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.RevokeByToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	hashedRefreshToken, err := s.aesService.EncryptInternal(newRefreshToken)
	if err != nil {
		return nil, err
	}

	token := model.NewRefreshToken(user.ID, hashedRefreshToken, expiryAt, false, *user)

	if err := s.refreshTokenRepo.Save(ctx, token); err != nil {
		return nil, err
	}

	return &primary.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *UserService) Logout(ctx context.Context, userID int64) error {
	return s.refreshTokenRepo.RevokeAllByUserID(ctx, userID)
}

func (s *UserService) VerifyEmail(ctx context.Context, token string) error {
	token, err := s.aesService.DecryptURLSafe(token)
	if err != nil {
		return domerr.ErrTokenInvalid
	}

	payloads := strings.Split(token, "_")
	if len(payloads) != 2 {
		return domerr.ErrTokenInvalid
	}

	expiryAt, err := time.Parse(time.RFC3339, payloads[1])
	if err != nil {
		return domerr.ErrTokenInvalid
	}

	if time.Now().After(expiryAt) {
		return domerr.ErrTokenExpired
	}

	userID, err := strconv.ParseInt(payloads[0], 10, 64)
	if err != nil {
		return domerr.ErrInvalidIDFormat
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return domerr.ErrUserNotFound
	}

	user.Activate()

	_, err = s.userRepo.Update(ctx, user)
	return err
}

func (s *UserService) SendVerifyEmail(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return domerr.ErrUserNotFound
	}

	plainText := fmt.Sprintf("%d_%s", user.ID, time.Now().Add(24*time.Hour).Format(time.RFC3339))

	token, err := s.aesService.EncryptURLSafe(plainText)
	if err != nil {
		return err
	}

	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", config.GetAppURL(), token)

	go func() {
		if err := s.emailService.SendVerifyEmail(user.Email, user.Name, verificationURL); err != nil {
			log.Printf("Failed to send verification email to %s: %v", user.Email, err)
		}
	}()

	return nil
}

func (s *UserService) SendResetPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return domerr.ErrUserNotFound
	}

	plainText := fmt.Sprintf("%s_%s", user.Email, time.Now().Add(1*time.Hour).Format(time.RFC3339))

	token, err := s.aesService.EncryptURLSafe(plainText)
	if err != nil {
		return err
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", config.GetAppURL(), token)

	go func() {
		if err := s.emailService.SendResetPasswordEmail(user.Email, user.Name, resetURL); err != nil {
			log.Printf("Failed to send password reset email to %s: %v", user.Email, err)
		}
	}()

	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, req *primary.ResetPasswordRequest) error {
	payload, err := s.aesService.DecryptURLSafe(req.Token)
	if err != nil {
		return domerr.ErrTokenInvalid
	}

	parts := strings.Split(payload, "_")
	if len(parts) != 2 {
		return domerr.ErrTokenInvalid
	}

	email := parts[0]
	expiryAt, err := time.Parse(time.RFC3339, parts[1])
	if err != nil {
		return domerr.ErrTokenInvalid
	}

	if time.Now().After(expiryAt) {
		return domerr.ErrTokenExpired
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return domerr.ErrUserNotFound
	}

	hashedPassword, err := s.bcryptService.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	if err := user.ChangePassword(hashedPassword); err != nil {
		return err
	}

	_, err = s.userRepo.Update(ctx, user)
	return err
}

func (s *UserService) GetAllUsers(ctx context.Context, page, pageSize int, search string) (*primary.PaginationResponse[primary.UserInfo], error) {
	offset := primary.Offset(page, pageSize)
	users, total, err := s.userRepo.FindAll(ctx, pageSize, offset, search)
	if err != nil {
		return nil, err
	}

	userInfos := make([]primary.UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = *FormatUserInfo(user)
	}

	return primary.NewPaginationResponse(userInfos, page, pageSize, int(total)), nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*primary.UserInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, domerr.ErrUserNotFound
	}

	return FormatUserInfo(user), nil
}

func (s *UserService) CreateUser(ctx context.Context, req *primary.CreateUserRequest) (*primary.UserInfo, error) {
	if s.userRepo.ExistsByEmail(ctx, req.Email) {
		return nil, domerr.ErrEmailAlreadyExists
	}

	hashedPassword, err := s.bcryptService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := model.NewUser(req.Name, req.Email, hashedPassword, "", model.Unknown)
	if err != nil {
		return nil, err
	}

	user.IsActive = true

	savedUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return FormatUserInfo(savedUser), nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID int64, req *primary.UpdateUserRequest) (*primary.UserInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, domerr.ErrUserNotFound
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Avatar != nil {
		path := fmt.Sprintf("avatars/user_%d/", user.ID)

		filePath, err := s.mediaService.UploadFile(req.Avatar.Filename, req.Avatar.Size, req.Avatar.Content, path)
		if err != nil {
			return nil, err
		}

		user.Avatar = *filePath
	}

	if req.Gender != nil {
		user.Gender = *req.Gender
	}

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return FormatUserInfo(updatedUser), nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID int64, req *primary.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return domerr.ErrUserNotFound
	}

	if err := s.bcryptService.ValidatePassword(req.OldPassword, user.Password); err != nil {
		return domerr.ErrPasswordNotMatch
	}

	hashedPassword, err := s.bcryptService.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	if err := user.ChangePassword(hashedPassword); err != nil {
		return err
	}

	_, err = s.userRepo.Update(ctx, user)
	return err
}

func (s *UserService) DeleteUser(ctx context.Context, userID int64) error {
	return s.userRepo.Delete(ctx, userID)
}
