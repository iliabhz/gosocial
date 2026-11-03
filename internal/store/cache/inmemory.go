package cache

import (
	"context"
	"fmt"
	"time"
)

type item struct {
	value     []byte
	expiresAt time.Time
}

type InMemoryClient struct {
	kv map[string]item
}

func NewInMemoryClient() Client {
	kv := make(map[string]item)
	return &InMemoryClient{
		kv: kv,
	}
}

func (c *InMemoryClient) Get(ctx context.Context, key string) ([]byte, error) {
	item, ok := c.kv[key]
	if !ok {
		return nil, ValueNotFound
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(c.kv, key)
		return nil, ValueNotFound
	}

	return item.value, nil
}

func (c *InMemoryClient) Set(ctx context.Context, key string, value any, expiry time.Duration) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("couldn't convert the data to []byte")
	}
	val := item{}
	val.value = data

	if expiry == 0 {
		var time time.Time // Zero value of time so exp.IsZero() would be true
		val.expiresAt = time
	}
	val.expiresAt = time.Now().Add(expiry)

	c.kv[key] = val
	return nil
}

func (c *InMemoryClient) Delete(ctx context.Context, key string) error {
	delete(c.kv, key)
	return nil
}

func (c *InMemoryClient) Close() error {
	return nil
}
