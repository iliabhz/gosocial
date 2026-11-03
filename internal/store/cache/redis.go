package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(dsn string) (Client, error) {
	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse redis dsn: %w", err)
	}

	c := redis.NewClient(opts)

	return &RedisClient{
		client: c,
	}, nil
}

func (c *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil && errors.Is(err, redis.Nil) {
		return nil, ValueNotFound
	}
	return data, err
}

func (c *RedisClient) Set(ctx context.Context, key string, value any, expiry time.Duration) error {
	return c.client.Set(ctx, key, value, expiry).Err()
}

func (c *RedisClient) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}
