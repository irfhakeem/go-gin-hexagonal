package main

import (
	"go-gin-hexagonal/pkg/config"
	"go-gin-hexagonal/pkg/database"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	args := os.Args

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

	switch args[1] {
	case "--migrate":
		database.RunMigrations(db)
	case "--seed":
		database.RunSeeders(db)
	case "--fresh":
		database.RunFreshMigrations(db)
	default:
		log.Println("Unknown command. Use --migrate, --seed, or --fresh.")
		return
	}
}
