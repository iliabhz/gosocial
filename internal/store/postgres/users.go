package postgres

import (
	"context"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	dbPool *pgxpool.Pool
}

func newPostgresUserRepository(dbPool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		dbPool: dbPool,
	}
}

func (repo *PostgresUserRepository) Create(ctx context.Context, user *store.User) error {
	query := `
	INSERT INTO users (username, password, email)
	 VALUES ($1, $2, $3) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	err := repo.dbPool.QueryRow(
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
