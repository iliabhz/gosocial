package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	clients map[string]int
	limit   int
	window  time.Duration

	mu sync.RWMutex
}

func NewFixedWindwoRateLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.mu.RLock()
	count, exists := rl.clients[ip]
	rl.mu.RUnlock()

	if !exists || count < rl.limit {
		rl.mu.Lock()
		if !exists {
			go rl.resetCount(ip)
		}

		rl.clients[ip]++
		rl.mu.Unlock()
		return true, 0
	}

	return false, rl.window
}

func (rl *FixedWindowRateLimiter) resetCount(ip string) {
	time.Sleep(rl.window)
	rl.mu.Lock()
	delete(rl.clients, ip)
	rl.mu.Unlock()
}
