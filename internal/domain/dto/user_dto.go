package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required,min=3,max=100"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,min=3,max=100"`
	Username *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type UserListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Search   string `form:"search,omitempty"`
}

type UserListResponse struct {
	Users      []*UserInfo `json:"users"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}
