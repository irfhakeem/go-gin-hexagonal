package security

import (
	"errors"
	"time"

	"go-gin-hexagonal/internal/domain/dto"
	"go-gin-hexagonal/internal/domain/entity"
	"go-gin-hexagonal/internal/domain/ports"
	"go-gin-hexagonal/pkg/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTToken struct {
	accessTokenSecret  string
	refreshTokenSecret string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewJWTToken(config config.JWTConfig) ports.TokenManager {
	return &JWTToken{
		accessTokenSecret:  config.AccessTokenSecret,
		refreshTokenSecret: config.RefreshTokenSecret,
		accessTokenExpiry:  config.AccessTokenExpiry,
		refreshTokenExpiry: config.RefreshTokenExpiry,
	}
}

func (tm *JWTToken) GenerateAccessToken(user *entity.User) (string, time.Time, error) {
	expiryDate := time.Now().Add(tm.accessTokenExpiry)
	domainClaims := &dto.AccessTokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		TokenType: "access",
		ExpiresAt: expiryDate,
		IssuedAt:  time.Now(),
		NotBefore: time.Now(),
		Issuer:    "go-gin-hexagonal",
		Subject:   user.ID.String(),
	}

	claims := jwt.MapClaims{
		"user_id":    domainClaims.UserID.String(),
		"email":      domainClaims.Email,
		"username":   domainClaims.Username,
		"token_type": domainClaims.TokenType,
		"expires_at": domainClaims.ExpiresAt.Unix(),
		"issued_at":  domainClaims.IssuedAt.Unix(),
		"not_before": domainClaims.NotBefore.Unix(),
		"issuer":     domainClaims.Issuer,
		"subject":    domainClaims.Subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tm.accessTokenSecret))
	return tokenString, expiryDate, err
}

func (tm *JWTToken) GenerateRefreshToken(userID uuid.UUID) (string, time.Time, error) {
	expiryDate := time.Now().Add(tm.refreshTokenExpiry)
	domainClaims := &dto.RefreshTokenClaims{
		UserID:    userID,
		TokenType: "refresh",
		ExpiresAt: expiryDate,
		IssuedAt:  time.Now(),
		NotBefore: time.Now(),
		Issuer:    "go-gin-hexagonal",
		Subject:   userID.String(),
	}

	claims := jwt.MapClaims{
		"user_id":    domainClaims.UserID.String(),
		"token_type": domainClaims.TokenType,
		"expires_at": domainClaims.ExpiresAt.Unix(),
		"issued_at":  domainClaims.IssuedAt.Unix(),
		"not_before": domainClaims.NotBefore.Unix(),
		"issuer":     domainClaims.Issuer,
		"subject":    domainClaims.Subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tm.refreshTokenSecret))
	return tokenString, expiryDate, err
}

func (tm *JWTToken) ValidateAccessToken(tokenString string) (*dto.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ports.ErrUnexpectedSinginMethod
		}
		return []byte(tm.accessTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["token_type"].(string)
		if !ok || tokenType != "access" {
			return nil, ports.ErrTokenInvalid
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return nil, ports.ErrInvalidClaims
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, ports.ErrInvalidIDFormat
		}

		email, _ := claims["email"].(string)
		username, _ := claims["username"].(string)
		issuer, _ := claims["issuer"].(string)
		subject, _ := claims["subject"].(string)

		expiresAt := time.Unix(int64(claims["expires_at"].(float64)), 0)
		issuedAt := time.Unix(int64(claims["issued_at"].(float64)), 0)
		notBefore := time.Unix(int64(claims["not_before"].(float64)), 0)

		return &dto.AccessTokenClaims{
			UserID:    userID,
			Email:     email,
			Username:  username,
			TokenType: tokenType,
			ExpiresAt: expiresAt,
			IssuedAt:  issuedAt,
			NotBefore: notBefore,
			Issuer:    issuer,
			Subject:   subject,
		}, nil
	}

	return nil, ports.ErrTokenInvalid
}

func (tm *JWTToken) ValidateRefreshToken(tokenString string) (*dto.RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tm.refreshTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["token_type"].(string)
		if !ok || tokenType != "refresh" {
			return nil, ports.ErrTokenInvalid
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return nil, ports.ErrInvalidClaims
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, ports.ErrInvalidIDFormat
		}

		issuer, _ := claims["issuer"].(string)
		subject, _ := claims["subject"].(string)

		expiresAt := time.Unix(int64(claims["expires_at"].(float64)), 0)
		issuedAt := time.Unix(int64(claims["issued_at"].(float64)), 0)
		notBefore := time.Unix(int64(claims["not_before"].(float64)), 0)

		return &dto.RefreshTokenClaims{
			UserID:    userID,
			TokenType: tokenType,
			ExpiresAt: expiresAt,
			IssuedAt:  issuedAt,
			NotBefore: notBefore,
			Issuer:    issuer,
			Subject:   subject,
		}, nil
	}

	return nil, ports.ErrTokenInvalid
}
