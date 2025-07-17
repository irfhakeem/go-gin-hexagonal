package main

import (
	"log"

	"go-gin-hexagonal/internal/adapter/mailer"
	mqAdapter "go-gin-hexagonal/internal/adapter/message-queue"
	"go-gin-hexagonal/internal/adapter/message-queue/consumers"
	"go-gin-hexagonal/pkg/config"
	mq "go-gin-hexagonal/pkg/message-queue"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	rabbitMQChannel := mq.NewRabbitMQConnection(&cfg.RabbitMQ)
	defer rabbitMQChannel.Close()

	// Init Message Queue adapter
	mqManager := mqAdapter.NewRabbitMQ(rabbitMQChannel)

	// Init Mailer adapter
	mailerManager := mailer.NewSMTPMailer(&cfg.Mailer)

	// Initialize consumer service
	emailConsumer := consumers.NewEmailConsumer(mailerManager, mqManager)

	log.Println("Starting all consumer...")
	emailConsumer.StartEmailConsumer()
}
