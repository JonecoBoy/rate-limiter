package limiter

import (
	"time"
)

type RateLimiter struct {
	storage Storage
}

type AllowResponse struct {
	Allowed     bool
	IsBlocked   bool
	UnblockTime time.Time
}

type TokenBucket struct {
	tokens       int
	lastRefill   time.Time
	refillRate   int
	bucketSize   int
	blocked      bool
	blockedUntil time.Time
	duration     time.Duration
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

func (rl *RateLimiter) Allow(key string, limit int, duration time.Duration) AllowResponse {
	return rl.storage.Allow(key, limit, duration)
}
