package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuditInfo struct {
	CreatedAt time.Time
	CreatedBy uuid.UUID
	UpdatedAt time.Time
	UpdatedBy uuid.UUID
	DeletedAt *time.Time
	DeletedBy uuid.UUID
	IsDeleted bool
}
