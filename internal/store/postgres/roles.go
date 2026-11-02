package postgres

import (
	"context"
	"errors"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRoleRepository struct {
	dbPool *pgxpool.Pool
}

func newPostgresRoleRepository(dbPool *pgxpool.Pool) *PostgresRoleRepository {
	return &PostgresRoleRepository{
		dbPool: dbPool,
	}
}

func (repo *PostgresRoleRepository) GetByName(ctx context.Context, name string) (*store.Role, error) {
	query := `
	SELECT id, level, description FROM roles WHERE name = $1
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	role := new(store.Role)
	err := repo.dbPool.QueryRow(
		ctx,
		query,
		name,
	).Scan(
		&role.ID,
		&role.Level,
		&role.Description,
	)
	role.Name = name

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrRoleNotFound
		}
		return nil, err
	}

	return role, nil
}
