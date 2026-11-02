package db

import (
	"fmt"
	"strings"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/iliabhz/gosocial/internal/store/postgres"
)

var (
	POSTGRESQL = "postgres"
	MYSQL      = "mysql"
	REDIS      = "redis"
)

func NewDB(dsn string) (store.Storage, error) {
	dbTypeStr := strings.Split(dsn, "://")[0]

	var store store.Storage
	var err error

	switch dbTypeStr {
	case POSTGRESQL:
		store, err = postgres.NewPostgresStorage(dsn)
		if err != nil {
			return nil, fmt.Errorf("Failed to initialize database: %w", err)
		}
	case MYSQL:
		return nil, fmt.Errorf("MySQL driver not yet implemented")
	case REDIS:
		return nil, fmt.Errorf("Redis driver not yet implemented")
	default:
		return nil, fmt.Errorf("Invalid driver: %s", dbTypeStr)
	}

	return store, nil
}
