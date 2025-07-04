package service

import (
	"fmt"

	"go-gin-hexagonal/internal/domain/dto"
	"go-gin-hexagonal/internal/domain/ports"
)

type EmailService struct {
	mailer ports.MailerManager
}

func NewEmailService(smtp ports.MailerManager) ports.EmailService {
	return &EmailService{
		mailer: smtp,
	}
}

func (s *EmailService) SendNewUserEmail(to string, data *dto.NewUserData) error {
	subject := "Here is your new account information for Go Gin Hexagonal Application"

	body, err := s.mailer.LoadEmailTemplate("new_user", data)
	if err != nil {
		return fmt.Errorf("failed to load welcome email template: %v", err)
	}

	return s.mailer.SendEmail(to, subject, body)
}
