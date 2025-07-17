package service

import (
	"encoding/json"
	"fmt"

	"go-gin-hexagonal/internal/domain/dto"
	"go-gin-hexagonal/internal/domain/ports"
)

type EmailService struct {
	application string
	mailer      ports.MailerManager
	mqManager   ports.MessageQueueManager
}

func NewEmailService(smtp ports.MailerManager, mq ports.MessageQueueManager) ports.EmailService {
	return &EmailService{
		application: "Go Gin Hexagonal Application",
		mailer:      smtp,
		mqManager:   mq,
	}
}

func (s *EmailService) SendNewUserEmail(to string, data *dto.NewUserData) error {
	subject := fmt.Sprintf("Here is your new account information for %s", s.application)

	emailMsg := dto.EmailMessage{
		To:       to,
		Template: "new_user",
		Subject:  subject,
		Data:     *data,
	}

	payload, err := json.Marshal(emailMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal email message: %w", err)
	}

	return s.mqManager.Publisher("", "email_queue", false, false, payload)
}

func (s *EmailService) SendVerifyEmail(to string, data *dto.VerifyEmailData) error {
	subject := fmt.Sprintf("Verify your email for %s", s.application)

	msg := dto.EmailMessage{
		To:       to,
		Subject:  subject,
		Template: "verify_email",
		Data:     *data,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal verify email message: %w", err)
	}

	return s.mqManager.Publisher("", "email_queue", false, false, payload)
}

func (s *EmailService) SendRequestResetPassword(to string, data *dto.ResetPasswordData) error {
	subject := fmt.Sprintf("Reset %s Account Password", s.application)

	msg := dto.EmailMessage{
		To:       to,
		Subject:  subject,
		Template: "reset_password",
		Data:     *data,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal reset password email message: %w", err)
	}

	return s.mqManager.Publisher("", "email_queue", false, false, payload)
}
