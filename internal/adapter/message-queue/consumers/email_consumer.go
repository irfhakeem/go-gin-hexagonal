package consumers

import (
	"encoding/json"
	"fmt"
	"log"

	"go-gin-hexagonal/internal/domain/dto"
	"go-gin-hexagonal/internal/domain/ports"
)

type EmailConsumer struct {
	mailer    ports.MailerManager
	mqManager ports.MessageQueueManager
}

func NewEmailConsumer(mailer ports.MailerManager, mqManager ports.MessageQueueManager) *EmailConsumer {
	return &EmailConsumer{
		mailer:    mailer,
		mqManager: mqManager,
	}
}

func (s *EmailConsumer) StartEmailConsumer() {
	_, err := s.mqManager.QueueDeclare("email_queue", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare email queue: %v", err)
	}

	msgs, err := s.mqManager.Consumer("email_queue", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for delivery := range msgs {
			// Process the email message
			if err := s.ProcessEmailMessage(delivery.Body); err != nil {
				log.Printf("Failed to process email message: %v", err)
				delivery.Nack(false, true)
			}

			log.Printf("Email sent successfully")
			delivery.Ack(false)
		}
	}()

	log.Printf("Email consumer started. Waiting for messages...")
	<-forever
}

func (s *EmailConsumer) ProcessEmailMessage(messageBody []byte) error {
	var emailMsg dto.EmailMessage
	if err := json.Unmarshal(messageBody, &emailMsg); err != nil {
		return fmt.Errorf("failed to unmarshal email message: %w", err)
	}

	if emailMsg.To == "" {
		return fmt.Errorf("recipient email is required")
	}
	if emailMsg.Subject == "" {
		return fmt.Errorf("email subject is required")
	}
	if emailMsg.Template == "" {
		return fmt.Errorf("email template is required")
	}

	renderedBody, err := s.mailer.LoadEmailTemplate(emailMsg.Template, emailMsg.Data)
	if err != nil {
		return fmt.Errorf("failed to load email template '%s': %w", emailMsg.Template, err)
	}

	if err := s.mailer.SendEmail(emailMsg.To, emailMsg.Subject, renderedBody); err != nil {
		return fmt.Errorf("failed to send email to '%s': %w", emailMsg.To, err)
	}

	return nil
}
