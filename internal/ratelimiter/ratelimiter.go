package ratelimiter

import (
	"time"
)

type Limiter interface {
	Allow(ip string) (bool, time.Duration)
}

type RateLimitErr struct {
	RetryAfter time.Duration
}

func (e *RateLimitErr) Error() string {
	return "rate limit exceeded, retry after: " + e.RetryAfter.String()
}
