package postgres

import (
	"context"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCommentRepository struct {
	dbPool *pgxpool.Pool
}

func newPostgresCommentRepository(dbPool *pgxpool.Pool) *PostgresCommentRepository {
	return &PostgresCommentRepository{
		dbPool: dbPool,
	}
}

func (repo *PostgresCommentRepository) GetByPostID(ctx context.Context, postID int64) ([]store.Comment, error) {
	query := `
	SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id FROM comments c
	JOIN users on users.id = c.user_id
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	rows, err := repo.dbPool.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []store.Comment{}
	for rows.Next() {
		var c store.Comment
		c.User = store.User{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username, &c.User.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}
