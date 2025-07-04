package service

import (
	"fmt"

	"go-gin-hexagonal/internal/domain/dto"
	"go-gin-hexagonal/internal/domain/ports"
)

type EmailService struct {
	application string
	mailer      ports.MailerManager
}

func NewEmailService(smtp ports.MailerManager) ports.EmailService {
	return &EmailService{
		application: "Go Gin Hexagonal Application",
		mailer:      smtp,
	}
}

func (s *EmailService) SendNewUserEmail(to string, data *dto.NewUserData) error {
	subject := fmt.Sprintf("Here is your new account information for %s", s.application)

	body, err := s.mailer.LoadEmailTemplate("new_user", data)
	if err != nil {
		return fmt.Errorf("failed to load welcome email template: %v", err)
	}

	return s.mailer.SendEmail(to, subject, body)
}

func (s *EmailService) SendVerifyEmail(to string, data *dto.VerifyEmailData) error {
	subject := fmt.Sprintf("Verify your email for %s", s.application)

	body, err := s.mailer.LoadEmailTemplate("verify_email", data)
	if err != nil {
		return fmt.Errorf("failed to load verify email template: %v", err)
	}

	return s.mailer.SendEmail(to, subject, body)
}

func (s *EmailService) SendRequestResetPassword(to string, data *dto.RequestResetPasswordData) error {
	subject := fmt.Sprintf("Reset %s Account Password", s.application)

	body, err := s.mailer.LoadEmailTemplate("reset_password", data)
	if err != nil {
		return fmt.Errorf("failed to load password reset request email template: %v", err)
	}

	return s.mailer.SendEmail(to, subject, body)
}
