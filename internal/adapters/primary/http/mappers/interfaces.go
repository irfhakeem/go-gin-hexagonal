package mappers

import (
	"go-gin-clean/internal/adapters/primary/http/dto"
	"go-gin-clean/internal/ports/primary"
)

type UserMapper interface {
	// DTO to Contract mappings
	LoginRequestToContract(req *dto.LoginRequest) *primary.LoginRequest
	RegisterRequestToContract(req *dto.RegisterRequest) *primary.RegisterRequest
	ResetPasswordRequestToContract(req *dto.ResetPasswordRequest) *primary.ResetPasswordRequest
	ChangePasswordRequestToContract(req *dto.ChangePasswordRequest) *primary.ChangePasswordRequest
	CreateUserRequestToContract(req *dto.CreateUserRequest) *primary.CreateUserRequest
	UpdateUserRequestToContract(req *dto.UpdateUserRequest) *primary.UpdateUserRequest
	PaginationRequestToContract(req *dto.PaginationRequest) *primary.PaginationRequest

	// Contract to DTO mappings
	LoginResponseToDTO(resp *primary.LoginResponse) *dto.LoginResponse
	RefreshTokenResponseToDTO(resp *primary.RefreshTokenResponse) *dto.RefreshTokenResponse
	UserInfoToDTO(user *primary.UserInfo) *dto.UserInfo
	PaginationResponseToDTO(resp *primary.PaginationResponse[primary.UserInfo]) *dto.PaginationResponse[dto.UserInfo]
}

type PaginationMapper interface {
	RequestToContract(req *dto.PaginationRequest) *primary.PaginationRequest
	UserInfoResponseToDTO(resp *primary.PaginationResponse[primary.UserInfo]) *dto.PaginationResponse[dto.UserInfo]
}
