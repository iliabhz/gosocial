package cache

import (
	"context"
	"time"

	"github.com/iliabhz/gosocial/internal/store"
)

var UserCacheExpiry = time.Second * 10

type CacheStorage interface {
	GetValue(ctx context.Context, key string) (string, error)
	SetValue(ctx context.Context, key string, value any, expiry time.Duration) error
	GetUser(ctx context.Context, userID int64) (*store.User, error)
	SetUser(ctx context.Context, user *store.User) error
	Close()
}
