package postgres

import (
	"context"
	"database/sql"

	"github.com/iliabhz/gosocial/internal/store"
)

type PostgresPostRepository struct {
	db *sql.DB
}

func newPostgresPostRepository(db *sql.DB) *PostgresPostRepository {
	return &PostgresPostRepository{
		db: db,
	}
}

func (repo *PostgresPostRepository) Create(ctx context.Context, user *store.Post) error {
	return nil
}
