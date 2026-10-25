package postgres

import (
	"context"
	"errors"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5"
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

func (repo *PostgresPostRepository) GetByID(ctx context.Context, postID int64) (*store.Post, error) {
	query := `
	SELECT id, content, title, user_id, tags, created_at, updated_at FROM posts WHERE id = $1;
	`

	var post store.Post

	err := repo.dbPool.QueryRow(
		ctx,
		query,
		postID,
	).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		&post.Tags,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, store.ErrPostNotFound
	default:
		return &post, err
	}
}

func (repo *PostgresPostRepository) Create(ctx context.Context, post *store.Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	 VALUES ($1, $2, $3, $4) RETURNING id, created_at::text, updated_at::text
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
