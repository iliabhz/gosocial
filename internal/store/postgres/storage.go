package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	dbPool       *pgxpool.Pool
	postRepo     *PostgresPostRepository
	commentRepo  *PostgresCommentRepository
	userRepo     *PostgresUserRepository
	followerRepo *PostgresFollowerRepository
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	dbPool, err := newDBPool(dsn)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		dbPool:       dbPool,
		postRepo:     newPostgresPostRepository(dbPool),
		commentRepo:  newPostgresCommentRepository(dbPool),
		userRepo:     newPostgresUserRepository(dbPool),
		followerRepo: newPostgresFollwerRepository(dbPool),
	}, nil
}

func withTx(ctx context.Context, dbPool *pgxpool.Pool, fn func(pgx.Tx) error) error {
	tx, err := dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func newDBPool(dsn string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to create an connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to establish an connection to the database: %w", err)
	}

	return dbPool, err
}

func (storage *PostgresStorage) Posts() store.PostRepository {
	return storage.postRepo
}

func (storage *PostgresStorage) Comments() store.CommentRepository {
	return storage.commentRepo
}

func (storage *PostgresStorage) Users() store.UserRepository {
	return storage.userRepo
}

func (storage *PostgresStorage) Followers() store.FollowerRepository {
	return storage.followerRepo
}

func (storage *PostgresStorage) Close() {
	storage.dbPool.Close()
}
