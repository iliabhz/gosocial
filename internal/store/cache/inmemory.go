package cache

import (
	"context"
	"fmt"
	"time"
)

type InMemoryClient struct {
	kv  map[string][]byte
	exp map[string]time.Time
}

func NewInMemoryClient() Client {
	kv := make(map[string][]byte)
	exp := make(map[string]time.Time)
	return &InMemoryClient{
		kv:  kv,
		exp: exp,
	}
}

func (c *InMemoryClient) Get(ctx context.Context, key string) ([]byte, error) {
	data, ok := c.kv[key]
	if !ok {
		return nil, ValueNotFound
	}

	if exp, ok := c.exp[key]; ok {
		if time.Now().After(exp) {
			delete(c.kv, key)
			delete(c.exp, key)
			return nil, ValueNotFound
		}
	}
	return data, nil
}

func (c *InMemoryClient) Set(ctx context.Context, key string, value any, expiry time.Duration) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("couldn't convert the data to []byte")
	}
	c.kv[key] = data
	c.exp[key] = time.Now().Add(expiry)
	return nil
}

func (c *InMemoryClient) Delete(ctx context.Context, key string) error {
	delete(c.kv, key)
	delete(c.exp, key)
	return nil
}

func (c *InMemoryClient) Close() error {
	return nil
}
