package postgres

import (
	"context"
	"database/sql"

	"github.com/iliabhz/tiago-backend/internal/store"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func newPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (repo *PostgresUserRepository) Create(ctx context.Context, user *store.User) error {
	return nil
}
