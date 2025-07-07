package ports

import (
	"context"
	"errors"

	"go-gin-hexagonal/internal/domain/entity"

	"github.com/google/uuid"
)

// Error
var (
	// General
	ErrGenerateToken               = errors.New("failed to create token")
	ErrCreateRefreshToken          = errors.New("failed to create refresh token")
	ErrPasswordMismatch            = errors.New("password mismatch")
	ErrPasswordTooShort            = errors.New("password too short, minimum 8 characters")
	ErrPasswordTooLong             = errors.New("password too long, maximum 20 characters")
	ErrEmailAlreadyExists          = errors.New("email already exists")
	ErrEmailNotFound               = errors.New("email not found")
	ErrPasswordWeak                = errors.New("password must contain at least one uppercase letter and one number")
	ErrTokenExpired                = errors.New("token expired")
	ErrTokenInvalid                = errors.New("token invalid")
	ErrTokenNotFound               = errors.New("token not found")
	ErrInvalidCredentials          = errors.New("invalid credentials")
	ErrAuthorizationHeaderNotFound = errors.New("authorization header not found")
	ErrInvalidIDFormat             = errors.New("invalid ID format")
	ErrUnexpectedSinginMethod      = errors.New("unexpected signin method")
	ErrInvalidClaims               = errors.New("invalid claims in token")
	ErrInvalidInput                = errors.New("invalid input provided")

	// User
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUpdateUser        = errors.New("failed to update user")
	ErrDeleteUser        = errors.New("failed to delete user")
	ErrCreateUser        = errors.New("failed to create user")
	ErrUserNotVerified   = errors.New("user not verified")
)

type BaseRepository[T any] interface {
	Raw(ctx context.Context, query string) ([]*T, error)
	FindAll(ctx context.Context, limit, offset int, query any, args ...any) ([]*T, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*T, error)
	FindFirst(ctx context.Context, query any, args ...any) (*T, error)
	Where(ctx context.Context, query any, args ...any) ([]*T, error)
	WhereExisting(ctx context.Context, query any, args ...any) (bool, error)
	Create(ctx context.Context, entity *T) (*T, error)
	Update(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, limit, offset int, search string) ([]*entity.User, int64, error)
	ExistsByEmail(ctx context.Context, email string) bool
	ExistsByUsername(ctx context.Context, username string) bool
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *entity.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.RefreshToken, error)
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
	RevokeByToken(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
	IsTokenValid(ctx context.Context, token string) bool
}
