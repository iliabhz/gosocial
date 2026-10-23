package postgres

import (
	"context"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPostRepository struct {
	dbPool *pgxpool.Pool
}

func newPostgresPostRepository(dbPool *pgxpool.Pool) *PostgresPostRepository {
	return &PostgresPostRepository{
		dbPool: dbPool,
	}
}

func (repo *PostgresPostRepository) Create(ctx context.Context, post *store.Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	 VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := repo.dbPool.QueryRow(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		post.Tags,
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	return err
}
