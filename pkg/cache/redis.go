package cache

import (
	"fmt"
	"go-gin-hexagonal/pkg/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(cfg *config.RedisConfig) *redis.Client {
	client := redis.NewClient(
		&redis.Options{
			Addr:     cfg.Address,
			Username: cfg.Username,
			Password: cfg.Password,
			DB:       cfg.Database,
		},
	)

	return client
}

func CloseRedisConnection(client *redis.Client) error {
	if client != nil {
		if err := client.Close(); err != nil {
			return fmt.Errorf("failed to close Redis connection: %w", err)
		}
	}
	return nil
}
