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

var UserCacheExpiry = time.Second * 10

type CacheStorage struct {
	store.Storage
	client *redis.Client
}

type CacheUserRepository struct {
	store.UserRepository
	client *redis.Client
	store  *CacheStorage
}

// The DSN can be empty if you want to opt out of using a cache
func NewRedisCacheStorage(dsn string, store store.Storage) *CacheStorage {
	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil
	}

	client := redis.NewClient(opts)
	cStore := CacheStorage{
		Storage: store,
		client:  client,
	}

	return &cStore
}

func (s *CacheStorage) Users() store.UserRepository {
	return CacheUserRepository{
		UserRepository: s.Storage.Users(),
		client:         s.client,
		store:          s,
	}
}

func (s CacheUserRepository) GetByID(ctx context.Context, userID int64) (*store.User, error) {
	key := fmt.Sprintf("user-%d", userID)
	data, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Cache Miss -> Caching
			user, err := s.store.Storage.Users().GetByID(ctx, userID) // Delegate to the store impl
			if err != nil {
				return nil, err
			}
			data, err := json.Marshal(user)
			if err != nil {
				return nil, err
			}
			if err := s.client.Set(ctx, key, data, UserCacheExpiry).Err(); err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}

	if data == "" {
		return nil, errors.New("The returned data on the cache for user was empty")
	}

	var user store.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *CacheStorage) Close() {
	s.client.Close()
}
