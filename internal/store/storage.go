package store

import (
	"errors"
	"fmt"
	"time"
)

var QueryTimeoutDuration = time.Second * 15

var ErrNotFound = errors.New("not found")
var ErrPostNotFound = fmt.Errorf("post: %w", ErrNotFound)
var ErrUserNotFound = fmt.Errorf("user: %w", ErrNotFound)

type Storage interface {
	Posts() PostRepository
	Comments() CommentRepository
	Users() UserRepository
	Close()
}
