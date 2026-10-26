package store

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")
var ErrPostNotFound = fmt.Errorf("post: %w", ErrNotFound)

type Storage interface {
	Posts() PostRepository
	Comments() CommentRepository
	Users() UserRepository
	Close()
}
