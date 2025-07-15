package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Mailer   MailerConfig
	AES      AESConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	AppUrl       string
	Host         string
	Port         int
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type MailerConfig struct {
	Host     string
	Port     int
	Sender   string
	Auth     string
	Password string
}

type AESConfig struct {
	Key string
	IV  string
}

type RedisConfig struct {
	Address  string
	Username string
	Password string
	Database int
	Protocol int
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "localhost"),
			Port:         getEnvAsInt("SERVER_PORT", 3000),
			Environment:  getEnv("APP_ENV", "development"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "go_gin_hexagonal"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			AccessTokenSecret:  getEnv("JWT_ACCESS_SECRET", "your-access-secret-key"),
			RefreshTokenSecret: getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
			AccessTokenExpiry:  getEnvAsDuration("JWT_ACCESS_EXPIRY", 1*time.Hour),
			RefreshTokenExpiry: getEnvAsDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
		},
		Mailer: MailerConfig{
			Host:     getEnv("MAILER_HOST", "smtp.example.com"),
			Port:     getEnvAsInt("MAILER_PORT", 587),
			Sender:   getEnv("MAILER_SENDER", "Go.Gin.Hexagonal <no-reply@testing.com>"),
			Auth:     getEnv("MAILER_AUTH", "your-authentication-string"),
			Password: getEnv("MAILER_PASSWORD", "your-email-password"),
		},
		AES: AESConfig{
			Key: getEnv("AES_KEY", "your-aes-encryption-key"),
			IV:  getEnv("AES_IV", "your-aes-initialization-vector"),
		},
		Redis: RedisConfig{
			Address:  getEnv("REDIS_ADDRESS", "localhost:6379"),
			Username: getEnv("REDIS_USERNAME", "default"),
			Password: getEnv("REDIS_PASSWORD", ""),
			Database: getEnvAsInt("REDIS_DATABASE", 0),
			Protocol: getEnvAsInt("REDIS_PROTOCOL", 0),
		},
	}

	return config, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode)
}

func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func GetAppURL() string {
	return getEnv("APP_FE_URL", "http://localhost:5000")
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
