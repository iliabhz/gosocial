package postgres

import (
	"database/sql"

	"github.com/iliabhz/tiago-backend/internal/store"
)

type PostgresStorage struct {
	postRepo *PostgresPostRepository
	userRepo *PostgresUserRepository
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		postRepo: newPostgresPostRepository(db),
		userRepo: newPostgresUserRepository(db),
	}
}

func (storage *PostgresStorage) Posts() store.PostRepository {
	return storage.postRepo
}

func (storage *PostgresStorage) Users() store.UserRepository {
	return storage.userRepo
}
