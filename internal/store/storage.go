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
var ErrRoleNotFound = fmt.Errorf("role not found") // Not wrapping ErrNotFound so we get internal server error
//  instead of Resource not found which can cause confusion

var ErrAlreadyFollowing = errors.New("Already following the user")
var ErrNotFollowing = errors.New("Not following the user")

var ErrDuplicateEmail = errors.New("duplicate email")
var ErrDuplicateUsername = errors.New("duplicate username")

type Storage interface {
	Posts() PostRepository
	Comments() CommentRepository
	Users() UserRepository
	Followers() FollowerRepository
	Roles() RoleRepository
	Close()
}
