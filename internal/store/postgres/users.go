package postgres

import (
	"context"
	"errors"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	dbPool *pgxpool.Pool
}

func newPostgresUserRepository(dbPool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		dbPool: dbPool,
	}
}

func (repo *PostgresUserRepository) GetByID(ctx context.Context, userID int64) (*store.User, error) {
	query := `
	SELECT id, username, email, password, created_at FROM users WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	var user store.User

	//Temp soltion
	var pass []byte

	err := repo.dbPool.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&pass,
		&user.CreatedAt,
	)

	//Temp soltion
	user.Password = string(pass)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, store.ErrUserNotFound
	default:
		return &user, err
	}
}

func (repo *PostgresUserRepository) Create(ctx context.Context, user *store.User) error {
	query := `
	INSERT INTO users (username, password, email)
	 VALUES ($1, $2, $3) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	err := repo.dbPool.QueryRow(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	return err
}
