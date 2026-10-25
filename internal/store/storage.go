package store

import "errors"

var ErrPostNotFound = errors.New("post not found")

type Storage interface {
	Posts() PostRepository
	Users() UserRepository
	Close()
}
