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

func (repo *PostgresPostRepository) Create(ctx context.Context, post *store.Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	 VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := repo.db.QueryRowContext(
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
