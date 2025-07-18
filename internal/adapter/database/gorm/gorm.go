package gorm

import (
	"log"
	"strings"

	"go-gin-hexagonal/internal/adapter/database/gorm/schema"
	"go-gin-hexagonal/internal/adapter/database/gorm/seeder"
	"go-gin-hexagonal/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// Enums
	enums = map[string][]string{
		// Example:
		// "gender": {
		// 	"Male",
		// 	"Female",
		// 	"Prefer not to say",
		// },
	}

	// Models
	models = []any{
		&schema.User{},
		&schema.RefreshToken{},
	}
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

	for name, values := range enums {
		quotedValues := make([]string, len(values))
		for i, value := range values {
			quotedValues[i] = "'" + value + "'"
		}
		if err := db.Exec("CREATE TYPE " + name + " AS ENUM (" + strings.Join(quotedValues, ", ") + ")").Error; err != nil {
			log.Print(err)
			return err
		}
	}

	err := db.AutoMigrate(
		models...,
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func RunSeeders(db *gorm.DB) {
	log.Println("Running database seeders...")

	err := seeder.UserSeeder(db)
	if err != nil {
		log.Printf("Error seeding user data: %v", err)
	}

	log.Println("Database seeding completed successfully")
}

func RunFreshMigrations(db *gorm.DB) {
	log.Println("Running fresh migrations...")

	err := db.Migrator().DropTable(
		models...,
	)

	if err != nil {
		log.Printf("Error dropping tables: %v", err)
		return
	}

	for name := range enums {
		if err := db.Exec("DROP TYPE IF EXISTS " + name).Error; err != nil {
			log.Print(err)
			return
		}
	}

	err = RunMigrations(db)
	if err != nil {
		log.Printf("Error running migrations: %v", err)
		return
	}

	RunSeeders(db)
}
