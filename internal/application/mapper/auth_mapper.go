package mapper

import (
	"go-gin-hexagonal/internal/application/dto"
	"go-gin-hexagonal/internal/domain/ports/services"
)

func MapLoginRequestDTOToService(req *dto.LoginRequest) *services.LoginRequest {
	return &services.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

func MapLoginResponseServiceToDTO(res *services.LoginResponse) *dto.LoginResponse {
	return &dto.LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
}

func MapRegisterRequestDTOToService(req *dto.RegisterRequest) *services.RegisterRequest {
	return &services.RegisterRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Name:     req.Name,
	}
}

func MapRefreshTokenResponseServiceToDTO(res *services.RefreshTokenResponse) *dto.RefreshTokenResponse {
	return &dto.RefreshTokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
}

func MapResetPasswordRequestDTOToService(req *dto.ResetPasswordRequest) *services.ResetPasswordRequest {
	return &services.ResetPasswordRequest{
		Token:           req.Token,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}
}
