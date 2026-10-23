package postgres

import (
	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	postRepo *PostgresPostRepository
	userRepo *PostgresUserRepository
}

func NewPostgresStorage(dbPool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{
		postRepo: newPostgresPostRepository(dbPool),
		userRepo: newPostgresUserRepository(dbPool),
	}
}

func (storage *PostgresStorage) Posts() store.PostRepository {
	return storage.postRepo
}

func (storage *PostgresStorage) Users() store.UserRepository {
	return storage.userRepo
}
