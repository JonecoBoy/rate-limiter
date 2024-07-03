package limiter

import "time"

type Storage interface {
	Allow(key string, limit int, duration time.Duration) bool
}
