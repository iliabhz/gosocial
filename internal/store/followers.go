package store

import (
	"context"
	"time"
)

// UserID --(Follows)--> FollowerID
type Follower struct {
	UserID     int64     `json:"user_id"`
	FollowerID int64     `json:"follower_id"`
	FollowedAt time.Time `json:"followed_at"`
}

type FollowerRepository interface {
	Follow(ctx context.Context, followerID int64, userID int64) error
	Unfollow(ctx context.Context, followerID int64, userID int64) error
}
