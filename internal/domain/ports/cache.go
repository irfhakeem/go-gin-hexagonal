package ports

import (
	"context"
	"time"
)

type CacheManager interface {
	Set(ctx context.Context, key string, value any) error
	SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	SetJSON(ctx context.Context, key string, data any, ttl time.Duration) error
	GetJSON(ctx context.Context, key string, dest any) error
}
