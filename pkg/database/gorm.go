package database

import (
	"log"

	"go-gin-hexagonal/internal/adapter/database/model"
	"go-gin-hexagonal/pkg/config"
	seeders "go-gin-hexagonal/pkg/database/seeder"

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

func RunSeeders(db *gorm.DB) {
	log.Println("Running database seeders...")

	err := seeders.UserSeeder(db)
	if err != nil {
		log.Printf("Error seeding user data: %v", err)
	}

	log.Println("Database seeding completed successfully")
}

func RunFreshMigrations(db *gorm.DB) {
	log.Println("Running fresh migrations...")

	err := db.Migrator().DropTable(
		&model.User{},
		&model.RefreshToken{},
	)
	if err != nil {
		log.Printf("Error dropping tables: %v", err)
		return
	}

	err = RunMigrations(db)
	if err != nil {
		log.Printf("Error running migrations: %v", err)
		return
	}

	RunSeeders(db)
}
