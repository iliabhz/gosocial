package cache

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"time"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/redis/go-redis/v9"
)

type RedisCacheStorage struct {
	client *redis.Client
}

// The DSN can be empty if you want to opt out of using a cache
func NewRedisCacheStorage(dsn string) *RedisCacheStorage {
	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil
	}

	client := redis.NewClient(opts)
	store := RedisCacheStorage{
		client: client,
	}

	return &store
}

func (s *RedisCacheStorage) GetValue(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *RedisCacheStorage) SetValue(ctx context.Context, key string, value any, expiry time.Duration) error {
	return s.client.Set(ctx, key, value, expiry).Err()
}

func (s *RedisCacheStorage) GetUser(ctx context.Context, userID int64) (*store.User, error) {
	key := fmt.Sprintf("user-%d", userID)
	data, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var user store.User
	if data == "" {
		return nil, errors.New("The returned data on the cache for user was empty")
	}

	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *RedisCacheStorage) SetUser(ctx context.Context, user *store.User) error {
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("user-%d", user.ID)
	return s.client.Set(ctx, key, json, UserCacheExpiry).Err()
}

func (s *RedisCacheStorage) Close() {
	s.client.Close()
}
