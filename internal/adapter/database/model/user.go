package model

import (
	"go-gin-hexagonal/internal/adapter/security"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email      string    `json:"email" gorm:"unique;not null;type:varchar(100)"`
	Username   string    `json:"username" gorm:"unique;not null;type:varchar(50)"`
	Password   string    `json:"password" gorm:"not null;type:varchar(255)"`
	Name       string    `json:"name" gorm:"not null;type:varchar(100)"`
	IsActive   bool      `json:"is_active" gorm:"default:true"`
	TrialCount int       `json:"trial_count" gorm:"default:0"`

	AuditInfo
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	if u.Password != "" {
		hashedPassword, err := security.NewBcryptHasher().Hash(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return nil
}

func (User) TableName() string {
	return "users"
}
