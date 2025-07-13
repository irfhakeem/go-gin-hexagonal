package entity

import (
	"time"
)

type AuditInfo struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	IsDeleted bool
}
