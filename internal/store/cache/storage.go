package cache

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"time"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var userCacheExpiry = time.Second * 10

type CacheStorage struct {
	store.Storage

	client *redis.Client
	logger *zap.SugaredLogger
}

type CacheUserRepository struct {
	store.UserRepository

	client *redis.Client
	logger *zap.SugaredLogger
}

// A empty DSN will opt out of caching and the original store would be returned
func NewRedisCacheStorage(dsn string, dbStore store.Storage, logger *zap.SugaredLogger) (store.Storage, error) {
	if dsn == "" {
		return dbStore, nil
	}

	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse redis dsn: %w", err)
	}

	client := redis.NewClient(opts)

	cStore := CacheStorage{
		Storage: dbStore,
		client:  client,
		logger:  logger,
	}

	return &cStore, nil
}

func (s *CacheStorage) Users() store.UserRepository {
	return CacheUserRepository{
		UserRepository: s.Storage.Users(),
		client:         s.client,
		logger:         s.logger,
	}
}

func (r CacheUserRepository) GetByID(ctx context.Context, userID int64) (*store.User, error) {
	key := fmt.Sprintf("user-%d", userID)
	data, err := r.client.Get(ctx, key).Bytes()

	if err == nil {
		var user store.User

		if err := json.Unmarshal(data, &user); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached user: %w", err)
		}
		r.logger.Info("CACHE HIT")
		return &user, nil
	} else if !errors.Is(err, redis.Nil) {
		r.logger.Warnw("failed to get user from cache", "user_id", userID, "error", err.Error())
	}
	r.logger.Info("CACHE MISS")

	user, err := r.UserRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	data, err = json.Marshal(user)
	if err != nil {
		r.logger.Warnw("failed to marshal user for cache", "error", err)
		return user, nil
	}

	if err := r.client.Set(ctx, key, data, userCacheExpiry).Err(); err != nil {
		r.logger.Warnw("failed to cache user", "user_id", userID, "error", err.Error())
	}

	return user, nil
}

func (s *CacheStorage) Close() {
	if s.client != nil {
		if err := s.client.Close(); err != nil {
			s.logger.Errorw("Failed to close the redis client", "error", err.Error())
		}
	}
	s.Storage.Close()
}
