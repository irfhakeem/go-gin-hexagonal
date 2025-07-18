package services

type EmailService interface {
	SendNewUserEmail(to string, data *NewUserEmailData) error
	SendVerifyEmail(to string, data *VerifyEmailData) error
	SendRequestResetPassword(to string, data *ResetPasswordData) error
}

type NewUserEmailData struct {
	Email    string
	Password string
}

type VerifyEmailData struct {
	VerificationURL string
}

type ResetPasswordData struct {
	ResetLink string
}
