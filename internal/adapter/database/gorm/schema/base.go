package schema

import (
	"time"
)

type AuditInfo struct {
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	IsDeleted bool       `json:"is_deleted" gorm:"default:false"`
}
