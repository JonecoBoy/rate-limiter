package limiter

import (
	"time"
)

type RateLimiter struct {
	storage Storage
}

func NewRateLimiter(storageType string) *RateLimiter {
	var storage Storage
	if storageType == "redis" {
		storage = NewRedisStorage()
	} else {
		storage = NewMemoryStorage()
	}
	return &RateLimiter{storage: storage}
}

func (rl *RateLimiter) Allow(key string, limit int, duration time.Duration) bool {
	return rl.storage.Allow(key, limit, duration)
}
