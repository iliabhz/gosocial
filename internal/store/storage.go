package store

type Storage interface {
	Posts() PostRepository
	Users() UserRepository
	Close()
}
