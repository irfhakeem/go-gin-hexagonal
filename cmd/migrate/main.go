package main

import (
	"go-gin-hexagonal/pkg/config"
	"go-gin-hexagonal/pkg/database"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Only run this if you want to check database connection and run migrations.
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	db, err := database.NewPostgresConnection(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
}
