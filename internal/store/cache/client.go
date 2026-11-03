package cache

import (
	"context"
	"errors"
	"time"
)

var ValueNotFound = errors.New("value not found")
var KeyCollision = errors.New("key collision")

type Client interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value any, expiry time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
}
