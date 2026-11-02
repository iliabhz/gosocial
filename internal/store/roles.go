package store

import "context"

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int64  `json:"level"`
	Description string `json:"description"`
}

type RoleRepository interface {
	GetByName(context.Context, string) (*Role, error)
}
