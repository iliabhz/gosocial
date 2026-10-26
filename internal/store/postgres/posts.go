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
	var updateStrBulilder strings.Builder
	updateStrBulilder.WriteString("SET ")

	var args []any
	i := 1

	if title != nil {
		updateStrBulilder.WriteString(fmt.Sprintf("title = $%v", i))
		args = append(args, *title)
		i++
	}
	if content != nil {
		if i == 1 {
			updateStrBulilder.WriteString(fmt.Sprintf("content = $%v", i))
		} else {
			updateStrBulilder.WriteString(fmt.Sprintf(", content = $%v", i))
		}

		args = append(args, *content)
		i++
	}
	if len(tags) != 0 {
		if i == 1 {
			updateStrBulilder.WriteString(fmt.Sprintf("tags = $%v", i))
		} else {
			updateStrBulilder.WriteString(fmt.Sprintf(", tags = $%v", i))
		}

		args = append(args, tags)
	}

	query := "UPDATE posts " + updateStrBulilder.String() + fmt.Sprintf(" WHERE id = %v", postID)

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

	cmdTag, err := repo.dbPool.Exec(ctx, query, postID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() != 1 {
		return store.ErrPostNotFound
	}

	return nil
}
