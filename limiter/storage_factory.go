package limiter

import (
	"os"
)

// NewStorage returns a Storage implementation based on environment configuration.
func NewStorage() Storage {
	switch os.Getenv("STORAGE_TYPE") {
	case "redis":
		return NewRedisStorage()
	default:
		return NewMemoryStorage()
	}
}
