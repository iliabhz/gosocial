package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

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

func (repo *PostgresPostRepository) UpdateByID(ctx context.Context, title *string, content *string, tags []string, postID int64) error {
	var args []any
	var sets []string

	if title != nil {
		args = append(args, *title)
		sets = append(sets, fmt.Sprintf("title = $%d", len(args)))
	}
	if content != nil {
		args = append(args, *content)
		sets = append(sets, fmt.Sprintf("content = $%d", len(args)))
	}
	if len(tags) != 0 {
		args = append(args, tags)
		sets = append(sets, fmt.Sprintf("tags = $%d", len(args)))
	}

	sets = append(sets, "updated_at = NOW()")

	query := fmt.Sprintf(`
        UPDATE posts
        SET %s
        WHERE id = %d
    `,
		strings.Join(sets, ", "),
		postID,
	)

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	cmdTag, err := repo.dbPool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return store.ErrPostNotFound
	}

	return nil
}

func (repo *PostgresPostRepository) Create(ctx context.Context, post *store.Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	 VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

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

func (repo *PostgresPostRepository) DeleteByID(ctx context.Context, postID int64) error {
	query := `
	DELETE FROM posts WHERE id=$1;
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	cmdTag, err := repo.dbPool.Exec(ctx, query, postID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() != 1 {
		return store.ErrPostNotFound
	}

	return nil
}
