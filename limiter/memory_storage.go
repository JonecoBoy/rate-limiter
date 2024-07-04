package limiter

import (
	"os"
	"sync"
	"time"
)

type MemoryStorage struct {
	mu     sync.Mutex
	tokens map[string]*TokenBucket
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tokens: make(map[string]*TokenBucket),
	}
}

func (ms *MemoryStorage) Cleanup(keys ...string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, key := range keys {
		delete(ms.tokens, key)
	}
}

func (ms *MemoryStorage) IsBlocked(key string) bool {
	if bucket, exists := ms.tokens[key]; exists {
		return bucket.blocked && time.Now().Before(bucket.blockedUntil)
	}
	return false
}

func (ms *MemoryStorage) Block(key string, blockTime time.Duration) {
	if bucket, exists := ms.tokens[key]; exists {
		bucket.blocked = true
		bucket.blockedUntil = time.Now().Add(blockTime)
	}
}

func (ms *MemoryStorage) Allow(key string, limit int, duration time.Duration) AllowResponse {
	now := time.Now()
	response := AllowResponse{Allowed: true, IsBlocked: false}

	if ms.IsBlocked(key) {
		if bucket, exists := ms.tokens[key]; exists {
			response.Allowed = false
			response.IsBlocked = true
			response.UnblockTime = bucket.blockedUntil
		}
		return response
	}

	if bucket, exists := ms.tokens[key]; exists {
		if now.Sub(bucket.lastRefill) >= bucket.duration {
			bucket.tokens = limit
			bucket.lastRefill = now
		} else {
			ms.refill(bucket, now)
			if bucket.tokens > 0 {
				bucket.tokens--
			} else {
				response.Allowed = false
				blockTime, _ := time.ParseDuration(os.Getenv("BLOCK_TIME") + "s")
				ms.Block(key, blockTime)
				response.IsBlocked = true
				response.UnblockTime = now.Add(blockTime)
			}
		}
	} else {
		ms.tokens[key] = &TokenBucket{
			tokens:     limit - 1,
			lastRefill: now,
			refillRate: limit,
			bucketSize: limit,
			duration:   duration,
		}
	}
	return response
}

func (ms *MemoryStorage) refill(bucket *TokenBucket, now time.Time) {
	timeElapsed := now.Sub(bucket.lastRefill)
	tokensToAdd := int(timeElapsed.Seconds()) * bucket.refillRate / int(bucket.duration.Seconds())

	if tokensToAdd > 0 {
		bucket.tokens += tokensToAdd
		if bucket.tokens > bucket.bucketSize {
			bucket.tokens = bucket.bucketSize
		}
		bucket.lastRefill = now
	}
}
