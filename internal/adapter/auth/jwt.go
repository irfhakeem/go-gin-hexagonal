package auth

import (
	"errors"
	"time"

	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTTokenManager struct {
	accessTokenSecret  string
	refreshTokenSecret string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func NewJWTTokenManager(config JWTConfig) ports.TokenManager {
	return &JWTTokenManager{
		accessTokenSecret:  config.AccessTokenSecret,
		refreshTokenSecret: config.RefreshTokenSecret,
		accessTokenExpiry:  config.AccessTokenExpiry,
		refreshTokenExpiry: config.RefreshTokenExpiry,
	}
}

func (tm *JWTTokenManager) GenerateAccessToken(user *entity.User) (string, time.Time, error) {
	expiryDate := time.Now().Add(tm.accessTokenExpiry)
	claims := &ports.AccessTokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryDate),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-gin-hexagonal",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tm.accessTokenSecret))
	return tokenString, expiryDate, err
}

func (tm *JWTTokenManager) GenerateRefreshToken(userID uuid.UUID) (string, time.Time, error) {
	expiryDate := time.Now().Add(tm.refreshTokenExpiry)
	claims := &ports.RefreshTokenClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryDate),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-gin-hexagonal",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tm.refreshTokenSecret))
	return tokenString, expiryDate, err
}

func (tm *JWTTokenManager) ValidateAccessToken(tokenString string) (*ports.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ports.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tm.accessTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ports.AccessTokenClaims); ok && token.Valid {
		if claims.TokenType != "access" {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (tm *JWTTokenManager) ValidateRefreshToken(tokenString string) (*ports.RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ports.RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tm.refreshTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ports.RefreshTokenClaims); ok && token.Valid {
		if claims.TokenType != "refresh" {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
