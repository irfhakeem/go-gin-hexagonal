package handlers

import (
	response "go-gin-hexagonal/internal/adapter/http"
	"go-gin-hexagonal/internal/adapter/http/message"
	"go-gin-hexagonal/internal/application/dto"
	"go-gin-hexagonal/internal/application/mapper"
	"go-gin-hexagonal/internal/domain/ports/services"
	"go-gin-hexagonal/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
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

	result, err := h.authService.Login(c.Request.Context(), (*services.LoginRequest)(&req))
	if err != nil {
		switch err {
		case errors.ErrPasswordMismatch:
			response.Error(c, message.FAILED_PASSWORD_INCORRECT, err.Error(), 400)
		case errors.ErrUserNotFound:
			response.Error(c, message.FAILED_USER_NOT_FOUND, err.Error(), 404)
		case errors.ErrUserNotVerified:
			response.Error(c, message.FAILED_FORBIDDEN, err.Error(), 403)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	mapResult := mapper.MapLoginResponseServiceToDTO(result)

	response.Success(c, message.SUCCESS_LOGIN, mapResult, 200)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, message.FAILED_INVALID_REQUEST_FORMAT, err.Error(), 400)
		return
	}

	mapReq := mapper.MapRegisterRequestDTOToService(&req)

	err := h.authService.Register(c.Request.Context(), mapReq)
	if err != nil {
		switch err {
		case errors.ErrUserAlreadyExists:
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

	result, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case errors.ErrTokenInvalid:
			response.Error(c, message.FAILED_TOKEN_INVALID, err.Error(), 401)
		case errors.ErrUserNotFound:
			response.Error(c, message.FAILED_USER_NOT_FOUND, err.Error(), 404)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	mapResult := mapper.MapRefreshTokenResponseServiceToDTO(result)

	response.Success(c, message.SUCCESS_REFRESH_TOKEN, mapResult, 200)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, message.FAILED_UNAUTHORIZED, errors.ErrInvalidCredentials.Error(), 401)
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		response.Error(c, message.FAILED_INVALID_ID_FORMAT, errors.ErrInvalidIDFormat.Error(), 400)
		return
	}

	err := h.authService.Logout(c.Request.Context(), userUUID)
	if err != nil {
		response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		return
	}

	response.Success(c, message.SUCCESS_LOGOUT, nil, 200)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, message.FAILED_INVALID_REQUEST_FORMAT, err.Error(), 400)
		return
	}

	err := h.authService.VerifyEmail(c.Request.Context(), req.Token)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			response.Error(c, message.FAILED_USER_NOT_FOUND, err.Error(), 404)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	response.Success(c, message.SUCCESS_VERIFY_USER, nil, 200)
}

func (h *AuthHandler) SendVerifyEmail(c *gin.Context) {
	var req dto.SendVerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, message.FAILED_INVALID_REQUEST_FORMAT, err.Error(), 400)
		return
	}

	err := h.authService.SendVerifyEmail(c.Request.Context(), req.Email)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			response.Error(c, message.FAILED_USER_NOT_FOUND, err.Error(), 404)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	response.Success(c, message.SUCCESS_SENT_VERIFY_EMAIL, nil, 200)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, message.FAILED_INVALID_REQUEST_FORMAT, err.Error(), 400)
		return
	}

	mapReq := mapper.MapResetPasswordRequestDTOToService(&req)

	err := h.authService.ResetPassword(c.Request.Context(), mapReq)
	if err != nil {
		switch err {
		case errors.ErrTokenInvalid:
			response.Error(c, message.FAILED_TOKEN_INVALID, err.Error(), 401)
		case errors.ErrUserNotFound:
			response.Error(c, message.FAILED_USER_NOT_FOUND, err.Error(), 404)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	response.Success(c, message.SUCCESS_RESET_PASSWORD, nil, 200)
}

func (h *AuthHandler) SendResetPassword(c *gin.Context) {
	var req dto.SendResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, message.FAILED_INVALID_REQUEST_FORMAT, err.Error(), 400)
		return
	}

	err := h.authService.SendResetPassword(c.Request.Context(), req.Email)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			response.Error(c, message.FAILED_USER_NOT_FOUND, err.Error(), 404)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	response.Success(c, message.SUCCESS_SENT_RESET_PASSWORD, nil, 200)
}
