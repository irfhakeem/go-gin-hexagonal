package services

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers(ctx context.Context, page, pageSize int, search string) (*UserPaginationResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*UserInfo, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*UserInfo, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req *UpdateUserRequest) (*UserInfo, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *ChangePasswordRequest) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type UserInfo struct {
	ID        uuid.UUID
	Email     string
	Username  string
	Name      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUserRequest struct {
	Email string
	Name  string
}

type UpdateUserRequest struct {
	Name     *string
	Username *string
}

type ChangePasswordRequest struct {
	CurrentPassword string
	NewPassword     string
}

type UserPaginationResponse struct {
	Datas      []*UserInfo
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
