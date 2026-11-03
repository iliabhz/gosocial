package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type MockStorage struct{}

func NewMockStore() Storage {
	return &MockStorage{}
}

func (s *MockStorage) Users() UserRepository {
	return &MockUserRepository{}
}

func (s *MockStorage) Posts() PostRepository {
	return nil
}

func (s *MockStorage) Comments() CommentRepository {
	return nil
}

func (s *MockStorage) Roles() RoleRepository {
	return nil
}

func (s *MockStorage) Followers() FollowerRepository {
	return nil
}

func (s *MockStorage) Close() {

}

type MockUserRepository struct{}

func (r *MockUserRepository) GetByID(context.Context, int64) (*User, error) {
	return &User{}, nil
}

func (r *MockUserRepository) GetByEmail(context.Context, string) (*User, error) {
	return &User{}, nil
}

func (r *MockUserRepository) Create(context.Context, *User) error {
	return nil
}

func (r *MockUserRepository) CreateTx(context.Context, pgx.Tx, *User) error {
	return nil
}

func (r *MockUserRepository) CreateAndInvite(ctx context.Context, invitationExp time.Duration, token string, user *User) error {
	return nil
}

func (r *MockUserRepository) CreateUserInvitation(ctx context.Context, tx pgx.Tx, invitationExp time.Duration, token string, userID int64) error {
	return nil
}

func (r *MockUserRepository) GetUserFromInvitation(ctx context.Context, tx pgx.Tx, token string) (*User, error) {
	return &User{}, nil
}

func (r *MockUserRepository) DeleteUserInvitation(ctx context.Context, tx pgx.Tx, userID int64) error {
	return nil
}

func (r *MockUserRepository) Update(context.Context, pgx.Tx, *User) error {
	return nil
}

func (r *MockUserRepository) Activate(ctx context.Context, token string) error {
	return nil
}

func (r *MockUserRepository) Delete(ctx context.Context, userID int64) error {
	return nil
}
