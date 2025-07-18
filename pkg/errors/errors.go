package errors

import "errors"

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
