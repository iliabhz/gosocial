package postgres

import (
	"context"
	"fmt"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	dbPool   *pgxpool.Pool
	postRepo *PostgresPostRepository
	userRepo *PostgresUserRepository
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	dbPool, err := newDBPool(dsn)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		dbPool:   dbPool,
		postRepo: newPostgresPostRepository(dbPool),
		userRepo: newPostgresUserRepository(dbPool),
	}, nil
}

func newDBPool(dsn string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to create an connection pool: %w", err)
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Failed to establish an connection to the database: %w", err)
	}

	return dbPool, err
}

func (storage *PostgresStorage) Posts() store.PostRepository {
	return storage.postRepo
}

func (storage *PostgresStorage) Users() store.UserRepository {
	return storage.userRepo
}

func (storage *PostgresStorage) Close() {
	storage.dbPool.Close()
}
