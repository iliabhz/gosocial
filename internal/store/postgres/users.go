package postgres

import (
	"context"
	"database/sql"

	"github.com/iliabhz/gosocial/internal/store"
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
	query := `
	INSERT INTO users (username, password, email)
	 VALUES ($1, $2, $3) RETURNING id, created_at
	`

	err := repo.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	return err
}
