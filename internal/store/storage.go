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

var ErrAlreadyFollowing = errors.New("Already following the user")
var ErrNotFollowing = errors.New("Not following the user")

type Storage interface {
	Posts() PostRepository
	Comments() CommentRepository
	Users() UserRepository
	Followers() FollowerRepository
	Close()
}
