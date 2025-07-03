package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	IsRevoked bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// Relations
	User User
}
