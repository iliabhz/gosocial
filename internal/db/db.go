package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(dsn string) (*pgxpool.Pool, error) {
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
