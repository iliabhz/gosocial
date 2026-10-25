package postgres

import (
	"context"
	"time"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
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
	 VALUES ($1, $2, $3, $4) RETURNING id, created_at::text, updated_at::text
	`

	var createdAt pgtype.Timestamptz
	var updeatedAt pgtype.Timestamptz

	err := repo.dbPool.QueryRow(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		post.Tags,
	).Scan(
		&post.ID,
		&createdAt,
		&updeatedAt,
	)

	post.CreatedAt = createdAt.Time.Format(time.RFC3339)
	post.UpdatedAt = updeatedAt.Time.Format(time.RFC3339)

	return err
}
