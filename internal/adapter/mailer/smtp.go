package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"go-gin-hexagonal/internal/domain/ports"
	"go-gin-hexagonal/pkg/config"

	"gopkg.in/gomail.v2"
)

type SMTPMailer struct {
	cfg *config.MailerConfig
}

func NewSMTPMailer(cfg *config.MailerConfig) ports.MailerManager {
	return &SMTPMailer{
		cfg: cfg,
	}
}

func (m *SMTPMailer) LoadEmailTemplate(templateName string, data any) (string, error) {
	templatePath := filepath.Join("internal", "adapter", "mailer", "template", templateName+".html")

	// Check if template file exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template file not found: %s", templatePath)
	}

	// Read and parse template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	var renderedBody string
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}
	renderedBody = buf.String()

	return renderedBody, nil
}

func (m *SMTPMailer) SendEmail(to string, subject string, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", m.cfg.Sender)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		m.cfg.Host,
		m.cfg.Port,
		m.cfg.Auth,
		m.cfg.Password,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
