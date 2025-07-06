package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Email    string
	Username string
	Password string
	Name     string
	IsActive bool

	AuditInfo
}
