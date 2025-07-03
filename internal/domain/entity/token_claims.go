package entity

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AccessTokenClaims struct {
	UserID    uuid.UUID
	Email     string
	Username  string
	TokenType string
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID    uuid.UUID
	TokenType string
	jwt.RegisteredClaims
}
