package cache

import (
	"context"
	"encoding/json"
	"go-gin-hexagonal/internal/domain/ports"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacher struct {
	client *redis.Client
}

func NewRedisCacher(client *redis.Client) ports.CacheManager {
	return &RedisCacher{
		client: client,
	}
}

func (r *RedisCacher) Set(ctx context.Context, key string, value any) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisCacher) SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCacher) Get(ctx context.Context, key string) (any, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return val, nil
}

func (r *RedisCacher) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCacher) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *RedisCacher) SetJSON(ctx context.Context, key string, data any, ttl time.Duration) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, bytes, ttl).Err()
}

func (r *RedisCacher) GetJSON(ctx context.Context, key string, dest any) error {
	bytes, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	return json.Unmarshal(bytes, dest)
}
