package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type item struct {
	value     []byte
	expiresAt time.Time
}

type InMemoryClient struct {
	kv map[string]item
	mu sync.Mutex
}

func NewInMemoryClient() Client {
	kv := make(map[string]item)
	return &InMemoryClient{
		kv: kv,
	}
}

func (c *InMemoryClient) Get(ctx context.Context, key string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

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
	if expiry < 0 {
		return c.Delete(ctx, key)
	}

	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Couldn't convert the data to []byte")
	}
	val := item{
		value: data,
	}

	if expiry == 0 {
		var time time.Time // Zero value of time so exp.IsZero() would be true
		val.expiresAt = time
	} else {
		val.expiresAt = time.Now().Add(expiry)
	}

	c.mu.Lock()
	c.kv[key] = val
	c.mu.Unlock()

	return nil
}

func (c *InMemoryClient) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	delete(c.kv, key)
	c.mu.Unlock()

	return nil
}

func (c *InMemoryClient) Close() error {
	return nil
}
