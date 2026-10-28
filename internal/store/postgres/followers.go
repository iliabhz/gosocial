package postgres

import (
	"context"
	"errors"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresFollowerRepository struct {
	dbPool *pgxpool.Pool
}

func newPostgresFollwerRepository(dbPool *pgxpool.Pool) *PostgresFollowerRepository {
	return &PostgresFollowerRepository{
		dbPool: dbPool,
	}
}

func (repo *PostgresFollowerRepository) Follow(ctx context.Context, followerID int64, userID int64) error {
	query := `
	INSERT INTO followers (user_id, follower_id) VALUES ($1, $2) RETURNING followed_at;
	`
	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	_, err := repo.dbPool.Exec(
		ctx,
		query,
		userID,
		followerID,
	)

	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if ok && pgErr.Code == "23505" {
		return store.ErrAlreadyFollowing
	}

	return err
}

func (repo *PostgresFollowerRepository) Unfollow(ctx context.Context, followerID int64, userID int64) error {
	query := `
	DELETE FROM followers WHERE user_id = $1 AND follower_id = $2;
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	cmdTag, err := repo.dbPool.Exec(
		ctx,
		query,
		userID,
		followerID,
	)

	if cmdTag.RowsAffected() != 1 {
		return store.ErrNotFollowing
	}

	return err
}
