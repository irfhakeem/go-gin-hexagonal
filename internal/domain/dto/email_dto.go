package dto

type NewUserData struct {
	UserEmail string
	Username  string
	Password  string
}

type VerifyEmailData struct {
	Token string
}

type RequestResetPasswordData struct {
	Email string
	Token string
}
