package limiter

import (
	"sync"
	"time"
)

type MemoryStorage struct {
	mu     sync.Mutex
	tokens map[string]*TokenBucket
}

type TokenBucket struct {
	tokens     int
	lastRefill time.Time
	refillRate int
	bucketSize int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tokens: make(map[string]*TokenBucket),
	}
}

func (ms *MemoryStorage) Allow(key string, limit int, duration time.Duration) bool {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	now := time.Now()
	if bucket, exists := ms.tokens[key]; exists {
		ms.refill(bucket, now)
		if bucket.tokens > 0 {
			bucket.tokens--
			return true
		}
		return false
	}

	ms.tokens[key] = &TokenBucket{
		tokens:     limit - 1,
		lastRefill: now,
		refillRate: limit,
		bucketSize: limit,
	}
	return true
}

func (ms *MemoryStorage) refill(bucket *TokenBucket, now time.Time) {
	timeElapsed := now.Sub(bucket.lastRefill)
	tokensToAdd := int(timeElapsed.Seconds()) * bucket.refillRate

	if tokensToAdd > 0 {
		bucket.tokens += tokensToAdd
		if bucket.tokens > bucket.bucketSize {
			bucket.tokens = bucket.bucketSize
		}
		bucket.lastRefill = now
	}
}
