package mapper

import (
	"go-gin-hexagonal/internal/application/dto"
	"go-gin-hexagonal/internal/domain/ports/services"
)

func MapUserInfoToDTO(user *services.UserInfo) *dto.UserInfo {
	return &dto.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func MapUserPaginationResponseToDTO(response *services.UserPaginationResponse) *dto.PaginationResponse[services.UserInfo] {
	return &dto.PaginationResponse[services.UserInfo]{
		Datas:      response.Datas,
		Total:      response.Total,
		Page:       response.Page,
		PageSize:   response.PageSize,
		TotalPages: response.TotalPages,
	}
}

func MapCreateUserRequestToService(req *dto.CreateUserRequest) *services.CreateUserRequest {
	return &services.CreateUserRequest{
		Email: req.Email,
		Name:  req.Name,
	}
}

func MapUpdateUserRequestToService(req *dto.UpdateUserRequest) *services.UpdateUserRequest {
	return &services.UpdateUserRequest{
		Name:     req.Name,
		Username: req.Username,
	}
}

func MapChangePasswordRequestToService(req *dto.ChangePasswordRequest) *services.ChangePasswordRequest {
	return &services.ChangePasswordRequest{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}
}

func MapPaginationRequestToService(req *dto.PaginationRequest) (int, int, string) {
	return req.Page, req.PageSize, req.Search
}
