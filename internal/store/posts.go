package store

import (
	"context"
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	Comments  []Comment `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostRepository interface {
	GetByID(context.Context, int64) (*Post, error)
	UpdateByID(ctx context.Context, title *string, content *string, tags []string, postID int64) error
	DeleteByID(context.Context, int64) error
	Create(context.Context, *Post) error
}
