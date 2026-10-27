package store

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at,omitzero"`
}

type UserRepository interface {
	GetByID(context.Context, int64) (*User, error)
	Create(context.Context, *User) error
}
