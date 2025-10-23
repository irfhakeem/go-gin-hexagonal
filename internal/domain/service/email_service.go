package service

import (
	"fmt"
	"go-gin-clean/internal/ports/secondary"
)

type EmailService struct {
	application string
	smtp        secondary.MailerService
}

func NewEmailService(smtp secondary.MailerService) *EmailService {
	return &EmailService{
		application: "Go Gin Clean App",
		smtp:        smtp,
	}
}

func (e *EmailService) SendVerifyEmail(to, name, url string) error {
	subject := fmt.Sprintf("Verify User %s Email", e.application)

	data := map[string]any{
		"Name":            name,
		"VerificationURL": url,
	}

	body, err := e.smtp.LoadTemplate("verify_email", data)
	if err != nil {
		return fmt.Errorf("failed to load email template: %v", err)
	}

	return e.smtp.SendEmail(to, subject, body)
}

func (e *EmailService) SendResetPasswordEmail(to, name, url string) error {
	subject := fmt.Sprintf("Reset %s Account Password", e.application)

	data := map[string]any{
		"Name":     name,
		"ResetURL": url,
	}

	body, err := e.smtp.LoadTemplate("reset_password", data)
	if err != nil {
		return fmt.Errorf("failed to load password reset request email template: %v", err)
	}

	return e.smtp.SendEmail(to, subject, body)
}
