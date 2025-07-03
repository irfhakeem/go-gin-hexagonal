package handlers

import (
	response "go-gin-hexagonal/internal/adapter/http"
	"go-gin-hexagonal/internal/adapter/http/message"
	"go-gin-hexagonal/internal/domain/dto"
	"go-gin-hexagonal/internal/domain/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, message.FAILED_INVALID_REQUEST_FORMAT, err.Error(), 400)
		return
	}

	result, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case ports.ErrInvalidCredentials:
			response.Error(c, message.FAILED_LOGIN_USER, err.Error(), 401)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	response.Success(c, message.SUCCESS_LOGIN, result, 200)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, message.FAILED_INVALID_REQUEST_FORMAT, err.Error(), 400)
		return
	}

	err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case ports.ErrUserAlreadyExists:
			response.Error(c, message.FAILED_REGISTER_USER, err.Error(), 409)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	response.Success(c, message.SUCCESS_REGISTER, nil, 201)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, message.FAILED_INVALID_REQUEST_FORMAT, err.Error(), 400)
		return
	}

	result, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case ports.ErrTokenInvalid:
			response.Error(c, message.FAILED_TOKEN_INVALID, err.Error(), 401)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	response.Success(c, message.SUCCESS_REFRESH_TOKEN, result, 200)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, message.FAILED_UNAUTHORIZED, "User not authenticated", 401)
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		response.Error(c, message.FAILED_INVALID_ID_FORMAT, "Invalid user ID format", 400)
		return
	}

	err := h.authService.Logout(c.Request.Context(), userUUID)
	if err != nil {
		response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		return
	}

	response.Success(c, message.SUCCESS_LOGOUT, nil, 200)
}
