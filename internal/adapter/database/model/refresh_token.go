package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Token     string         `json:"token" gorm:"unique;not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	IsRevoked bool           `json:"is_revoked" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	return nil
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
