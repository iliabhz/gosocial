package store

import (
	"context"
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentRepository interface {
	GetByPostID(context.Context, int64) ([]Comment, error)
}
