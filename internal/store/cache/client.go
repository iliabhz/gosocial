package cache

import (
	"context"
	"time"
)

type Client interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value any, expiry time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
}
