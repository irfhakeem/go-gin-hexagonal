package primary

import "context"

// UserUseCase defines the primary port for user operations
type UserUseCase interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterRequest) error
	RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error)
	Logout(ctx context.Context, userID int64) error
	VerifyEmail(ctx context.Context, token string) error
	SendVerifyEmail(ctx context.Context, email string) error
	SendResetPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
	GetAllUsers(ctx context.Context, page, pageSize int, search string) (*PaginationResponse[UserInfo], error)
	GetUserByID(ctx context.Context, userID int64) (*UserInfo, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*UserInfo, error)
	UpdateUser(ctx context.Context, userID int64, req *UpdateUserRequest) (*UserInfo, error)
	ChangePassword(ctx context.Context, userID int64, req *ChangePasswordRequest) error
	DeleteUser(ctx context.Context, userID int64) error
}

// EmailUseCase defines the primary port for email operations
type EmailUseCase interface {
	SendVerifyEmail(toEmail, toName, verifyToken string) error
	SendResetPasswordEmail(toEmail, toName, resetToken string) error
}
