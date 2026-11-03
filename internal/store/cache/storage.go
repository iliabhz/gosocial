package cache

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"time"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var userCacheExp time.Duration
var roleCacheExp time.Duration

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

type CacheRoleRepository struct {
	store.RoleRepository

	client *redis.Client
	logger *zap.SugaredLogger
}

// A empty DSN will opt out of caching and the original store would be returned
func NewRedisCacheStorage(dsn string, userCacheExpiry time.Duration, roleCacheExpiry time.Duration, dbStore store.Storage, logger *zap.SugaredLogger) (store.Storage, error) {
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

	userCacheExp = userCacheExpiry
	roleCacheExp = roleCacheExpiry
	return &cStore, nil
}

func (s *CacheStorage) Users() store.UserRepository {
	return &CacheUserRepository{
		UserRepository: s.Storage.Users(),
		client:         s.client,
		logger:         s.logger,
	}
}

func (r *CacheUserRepository) GetByID(ctx context.Context, userID int64) (*store.User, error) {
	key := fmt.Sprintf("user-%d", userID)
	data, err := r.client.Get(ctx, key).Bytes()

	if err == nil {
		var user store.User

		if err := json.Unmarshal(data, &user); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached user: %w", err)
		}
		return &user, nil
	} else if !errors.Is(err, redis.Nil) {
		r.logger.Warnw("failed to get user from cache", "user_id", userID, "error", err.Error())
	}

	user, err := r.UserRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	data, err = json.Marshal(user)
	if err != nil {
		r.logger.Warnw("failed to marshal user for cache", "error", err)
		return user, nil
	}

	if err := r.client.Set(ctx, key, data, userCacheExp).Err(); err != nil {
		r.logger.Warnw("failed to cache user", "user_id", userID, "error", err.Error())
	}

	return user, nil
}

func (r *CacheUserRepository) Update(ctx context.Context, tx pgx.Tx, user *store.User) error {
	err := r.UserRepository.Update(ctx, tx, user)
	if err != nil {
		key := fmt.Sprintf("user-%d", user.ID)
		if err := r.client.Del(ctx, key).Err(); err != nil {
			r.logger.Warn("failed to delete a user key for invalidation", "user_id", user.ID, "error", err.Error())
		}
		return nil
	}
	return err
}

func (r *CacheUserRepository) Delete(ctx context.Context, userID int64) error {
	err := r.UserRepository.Delete(ctx, userID)
	if err != nil {
		key := fmt.Sprintf("user-%d", userID)
		if err := r.client.Del(ctx, key).Err(); err != nil {
			r.logger.Warn("failed to delete a user key for invalidation", "user_id", userID, "error", err.Error())
		}
		return nil
	}
	return err
}

func (s *CacheStorage) Roles() store.RoleRepository {
	return &CacheRoleRepository{
		RoleRepository: s.Storage.Roles(),
		client:         s.client,
		logger:         s.logger,
	}
}

func (r *CacheRoleRepository) GetByName(ctx context.Context, name string) (*store.Role, error) {
	key := fmt.Sprintf("role-%s", name)
	data, err := r.client.Get(ctx, key).Bytes()

	if err == nil {
		var role store.Role

		if err := json.Unmarshal(data, &role); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached role: %w", err)
		}
		return &role, nil
	} else if !errors.Is(err, redis.Nil) {
		r.logger.Warnw("failed to get role from cache", "role_name", name, "error", err.Error())
	}

	role, err := r.RoleRepository.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	data, err = json.Marshal(role)
	if err != nil {
		r.logger.Warnw("failed to marshal role for cache", "error", err)
		return role, nil
	}

	if err := r.client.Set(ctx, key, data, roleCacheExp).Err(); err != nil {
		r.logger.Warnw("failed to cache role", "role", name, "error", err.Error())
	}

	return role, nil
}

func (s *CacheStorage) Close() {
	if s.client != nil {
		if err := s.client.Close(); err != nil {
			s.logger.Errorw("Failed to close the redis client", "error", err.Error())
		}
	}
	s.Storage.Close()
}
