package primary

import (
	"go-gin-clean/internal/domain/model"
	"io"
	"time"
)

// DTOs and shared types for primary ports
type (
	UserInfo struct {
		ID       int64
		Name     string
		Email    string
		Avatar   string
		Gender   model.Gender
		IsActive bool
	}

	LoginRequest struct {
		Email    string
		Password string
	}

	LoginResponse struct {
		AccessToken  string
		RefreshToken string
		User         UserInfo
	}

	RegisterRequest struct {
		Name     string
		Email    string
		Password string
	}

	RefreshTokenResponse struct {
		AccessToken  string
		RefreshToken string
	}

	ResetPasswordRequest struct {
		Token       string
		NewPassword string
	}

	ChangePasswordRequest struct {
		OldPassword string
		NewPassword string
	}

	CreateUserRequest struct {
		Name     string
		Email    string
		Password string
		Gender   model.Gender
	}

	UpdateUserRequest struct {
		Name   *string
		Gender *model.Gender
		Avatar *FileUpload
	}

	FileUpload struct {
		Filename string
		Size     int64
		Content  io.Reader
	}

	PaginationRequest struct {
		Page    int
		PerPage int
		Search  string
	}

	PaginationResponse[T any] struct {
		Data       []T
		Page       int
		PerPage    int
		Total      int
		TotalPages int
	}

	AccessTokenClaims struct {
		UserID    int64
		Email     string
		TokenType string
		ExpiresAt time.Time
		IssuedAt  time.Time
		NotBefore time.Time
		Issuer    string
		Subject   string
	}

	RefreshTokenClaims struct {
		UserID    int64
		TokenType string
		ExpiresAt time.Time
		IssuedAt  time.Time
		NotBefore time.Time
		Issuer    string
		Subject   string
	}
)

// Helper functions
func NewPaginationResponse[T any](data []T, page, perPage, total int) *PaginationResponse[T] {
	totalPages := (total + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}

	return &PaginationResponse[T]{
		Data:       data,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

func Offset(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * pageSize
}
