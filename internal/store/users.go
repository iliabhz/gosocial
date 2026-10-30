package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"`
	Password  Password  `json:"-"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at,omitzero"`
}

type Password struct {
	Text *string
	Hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Text = &text
	p.Hash = hash

	return nil
}

type UserRepository interface {
	GetByID(context.Context, int64) (*User, error)
	Create(context.Context, *User) error
	CreateTx(context.Context, pgx.Tx, *User) error
	CreateAndInvite(ctx context.Context, invitationExp time.Duration, token string, user *User) error
	CreateUserInvitation(ctx context.Context, tx pgx.Tx, invitationExp time.Duration, token string, userID int64) error
	GetUserFromInvitation(ctx context.Context, tx pgx.Tx, token string) (*User, error)
	DeleteUserInvitation(ctx context.Context, tx pgx.Tx, userID int64) error
	Update(context.Context, pgx.Tx, *User) error
	Activate(ctx context.Context, token string) error
}
