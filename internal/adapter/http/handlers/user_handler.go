package handlers

import (
	response "go-gin-hexagonal/internal/adapter/http"
	"go-gin-hexagonal/internal/adapter/http/message"
	"go-gin-hexagonal/internal/application/dto"
	"go-gin-hexagonal/internal/application/mapper"
	"go-gin-hexagonal/internal/domain/ports/services"
	"go-gin-hexagonal/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
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

	result, err := h.userService.GetUserByID(c.Request.Context(), userUUID)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			response.Error(c, message.FAILED_GET_USER_BY_ID, err.Error(), 404)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	mapResult := mapper.MapUserInfoToDTO(result)

	response.Success(c, message.SUCCESS_GET_USER_BY_ID, mapResult, 200)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(req); err != nil {
		response.Error(c, message.FAILED_INVALID_REQUEST_FORMAT, err.Error(), 400)
		return
	}

	mapReq := mapper.MapCreateUserRequestToService(&req)

	result, err := h.userService.CreateUser(c.Request.Context(), mapReq)
	if err != nil {
		switch err {
		case errors.ErrUserAlreadyExists:
			response.Error(c, message.FAILED_USER_ALREADY_EXISTS, err.Error(), 409)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}
	response.Success(c, message.SUCCESS_CREATE_USER, result, 201)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
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

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, message.FAILED_INVALID_REQUEST_FORMAT, err.Error(), 400)
		return
	}

	mapReq := mapper.MapUpdateUserRequestToService(&req)

	result, err := h.userService.UpdateUser(c.Request.Context(), userUUID, mapReq)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			response.Error(c, message.FAILED_GET_USER_BY_ID, err.Error(), 404)
		case errors.ErrUserAlreadyExists:
			response.Error(c, message.FAILED_USER_ALREADY_EXISTS, err.Error(), 409)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	response.Success(c, message.SUCCESS_UPDATE_USER, result, 200)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userUUID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, message.FAILED_INVALID_ID_FORMAT, errors.ErrInvalidIDFormat.Error(), 400)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, message.FAILED_INVALID_REQUEST_FORMAT, err.Error(), 400)
		return
	}

	mapReq := mapper.MapUpdateUserRequestToService(&req)

	result, err := h.userService.UpdateUser(c.Request.Context(), userUUID, mapReq)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			response.Error(c, message.FAILED_GET_USER_BY_ID, err.Error(), 404)
		case errors.ErrUserAlreadyExists:
			response.Error(c, message.FAILED_USER_ALREADY_EXISTS, err.Error(), 409)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	response.Success(c, message.SUCCESS_UPDATE_USER, result, 200)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
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

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, message.FAILED_INVALID_REQUEST_FORMAT, err.Error(), 400)
		return
	}

	mapReq := mapper.MapChangePasswordRequestToService(&req)

	err := h.userService.ChangePassword(c.Request.Context(), userUUID, mapReq)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			response.Error(c, message.FAILED_GET_USER_BY_ID, err.Error(), 404)
		case errors.ErrInvalidCredentials:
			response.Error(c, message.FAILED_PASSWORD_INCORRECT, "", 400)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	response.Success(c, message.SUCCESS_CHANGE_PASSWORD, nil, 200)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var req dto.PaginationRequest

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	req.Search = c.Query("search")

	result, err := h.userService.GetAllUsers(c.Request.Context(), page, pageSize, req.Search)
	if err != nil {
		response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		return
	}

	meta := &response.Meta{
		Page:       result.Page,
		PageSize:   result.PageSize,
		Total:      result.Total,
		TotalPages: result.TotalPages,
	}

	response.SuccessWithMeta(c, message.SUCCESS_GET_ALL_USERS, result.Datas, meta)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, message.FAILED_INVALID_ID_FORMAT, errors.ErrInvalidIDFormat.Error(), 400)
		return
	}

	result, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			response.Error(c, message.FAILED_GET_USER_BY_ID, err.Error(), 404)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}

	response.Success(c, message.SUCCESS_GET_USER_BY_ID, result, 200)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDstr := c.Param("id")
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		response.Error(c, message.FAILED_INVALID_ID_FORMAT, errors.ErrInvalidIDFormat.Error(), 400)
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			response.Error(c, message.FAILED_GET_USER_BY_ID, err.Error(), 404)
		case errors.ErrDeleteUser:
			response.Error(c, message.FAILED_DELETE_USER, err.Error(), 500)
		default:
			response.Error(c, message.FAILED_INTERNAL_SERVER_ERROR, err.Error(), 500)
		}
		return
	}
	response.Success(c, message.SUCCESS_DELETE_USER, nil, 204)
}
