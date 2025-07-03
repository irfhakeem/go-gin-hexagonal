package database

import (
	"log"

	"go-gin-hexagonal/internal/adapter/database/model"
	"go-gin-hexagonal/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN()

	var logLevel logger.LogLevel
	if cfg.Host == "localhost" || cfg.Host == "127.0.0.1" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return nil, err
	}

	psqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	psqlDB.SetMaxOpenConns(25)
	psqlDB.SetMaxIdleConns(5)

	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
