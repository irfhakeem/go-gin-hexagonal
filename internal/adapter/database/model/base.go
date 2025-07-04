package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditInfo struct {
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	IsDeleted bool       `json:"is_deleted" gorm:"default:false"`
	CreatedBy uuid.UUID  `json:"created_by" gorm:"type:uuid"`
	UpdatedBy uuid.UUID  `json:"updated_by" gorm:"type:uuid"`
	DeletedBy *uuid.UUID `json:"deleted_by,omitempty" gorm:"type:uuid"`
}
