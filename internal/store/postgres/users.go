package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/iliabhz/gosocial/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	SELECT id, username, email, password, is_active, created_at FROM users WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	var user store.User

	err := repo.dbPool.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.IsActive,
		&user.CreatedAt,
	)

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
		user.Password.Hash,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		e, ok := errors.AsType[*pgconn.PgError](err)
		if ok && e.Code == "23505" && e.Message == "duplicate key value violates unique constraint \"users_email_key\"" {
			return store.ErrDuplicateEmail
		}
		if ok && e.Code == "23505" && e.Message == "duplicate key value violates unique constraint \"users_username_key\"" {
			return store.ErrDuplicateUsername
		}
		return err
	}

	return nil
}

func (repo *PostgresUserRepository) CreateTx(ctx context.Context, tx pgx.Tx, user *store.User) error {
	query := `
	INSERT INTO users (username, password, email)
	 VALUES ($1, $2, $3) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRow(
		ctx,
		query,
		user.Username,
		user.Password.Hash,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		e, ok := errors.AsType[*pgconn.PgError](err)
		if ok && e.Code == "23505" && e.Message == "duplicate key value violates unique constraint \"users_email_key\"" {
			return store.ErrDuplicateEmail
		}
		if ok && e.Code == "23505" && e.Message == "duplicate key value violates unique constraint \"users_username_key\"" {
			return store.ErrDuplicateUsername
		}
		return err
	}

	return nil
}

func (repo *PostgresUserRepository) CreateAndInvite(ctx context.Context, invitationExp time.Duration, token string, user *store.User) error {
	return withTx(ctx, repo.dbPool, func(tx pgx.Tx) error {
		// create user
		if err := repo.CreateTx(ctx, tx, user); err != nil {
			return err
		}

		// invite user
		if err := repo.CreateUserInvitation(ctx, tx, invitationExp, token, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (repo *PostgresUserRepository) CreateUserInvitation(ctx context.Context, tx pgx.Tx, invitationExp time.Duration, token string, userID int64) error {
	query := `
	INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	_, err := tx.Exec(
		ctx,
		query,
		token,
		userID,
		time.Now().Add(invitationExp),
	)

	return err
}

func (repo *PostgresUserRepository) Activate(ctx context.Context, token string) error {
	return withTx(ctx, repo.dbPool, func(tx pgx.Tx) error {
		user, err := repo.GetUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		user.IsActive = true
		if err := repo.Update(ctx, tx, user); err != nil {
			return err
		}

		if err := repo.DeleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

// TODO: Make it a partial update with dynamic query building like posts
func (repo *PostgresUserRepository) Update(ctx context.Context, tx pgx.Tx, user *store.User) error {
	query := `
	UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	cmdTag, err := tx.Exec(
		ctx,
		query,
		user.Username,
		user.Email,
		user.IsActive,
		user.ID,
	)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return store.ErrUserNotFound
	}

	return nil
}

func (repo *PostgresUserRepository) GetUserFromInvitation(ctx context.Context, tx pgx.Tx, token string) (*store.User, error) {
	query := `
	SELECT u.id, u.username, u.email, u.created_at, u.is_active
	FROM users u
	JOIN user_invitations ui ON u.id = ui.user_id
	WHERE ui.token = $1 AND ui.expiry > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	user := store.User{}

	err := tx.QueryRow(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresUserRepository) DeleteUserInvitation(ctx context.Context, tx pgx.Tx, userID int64) error {
	query := `
	DELETE FROM user_invitations WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, store.QueryTimeoutDuration)
	defer cancel()

	cmdTag, err := tx.Exec(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return store.ErrUserNotFound
	}

	return nil
}
