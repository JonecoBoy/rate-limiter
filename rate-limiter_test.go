package main

import (
	"fmt"
	"github.com/JonecoBoy/rate-limiter/limiter"
	"os"
	"testing"
	"time"
)

// Mock environment variables for testing
func setupEnvRedis() {
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("BLOCK_TIME", "60")
	os.Setenv("RATE_LIMIT_IP", "5")
	os.Setenv("RATE_LIMIT_TOKEN", "10")
	os.Setenv("RATE_LIMIT_STRATEGY", "redis")
	os.Setenv("RATE_LIMIT_DURATION", "1000")
}

func setupEnvInMemory() {
	os.Setenv("BLOCK_TIME", "60")
	os.Setenv("RATE_LIMIT_IP", "5")
	os.Setenv("RATE_LIMIT_TOKEN", "10")
	os.Setenv("RATE_LIMIT_DURATION", "1000")

}

func TestRateLimiterByIPRedis(t *testing.T) {
	setupEnvRedis()

	rs := limiter.NewStorage() // Assuming this creates a new Redis client

	ip := "192.168.1.1"
	ipLimit := 5
	duration := time.Second // 1 second window for rate limit

	rs.Cleanup(ip)

	// test ip
	t.Run("IP Rate Limit: Allowance", func(t *testing.T) {
		for i := 0; i < ipLimit; i++ {
			res := rs.Allow(ip, ipLimit, duration)
			if !res.Allowed {
				t.Errorf("IP rate limit failed at request %d: should not block", i+1)
			}
		}
	})

	// test ip block
	t.Run("IP Rate Limit: Blocking", func(t *testing.T) {
		res := rs.Allow(ip, ipLimit, duration)
		if res.Allowed {
			t.Errorf("IP rate limit failed: should block the sixth request")
		} else {
			fmt.Printf("Blocked until: %s\n", res.UnblockTime.Format(time.RFC1123))
			t.Logf("Successfully blocked as expected. Blocked until: %s", res.UnblockTime.Format(time.RFC1123))
		}
	})

	rs.Cleanup(ip)
}

func TestRateLimiterByApiTokenRedis(t *testing.T) {
	setupEnvRedis()

	rs := limiter.NewStorage() // Assuming this creates a new Redis client

	token := "abc123"
	tokenLimit := 10
	duration := time.Second // 1 second window for rate limit

	rs.Cleanup(token)
	// Test Token-based rate limiting for allowed requests
	t.Run("Token Rate Limit: Allowance", func(t *testing.T) {
		for i := 0; i < tokenLimit; i++ {
			res := rs.Allow(token, tokenLimit, duration)
			if !res.Allowed {
				t.Errorf("Token rate limit failed at request %d: should not block", i+1)
			}
		}
	})

	// Test Token-based rate limiting for blocking
	t.Run("Token Rate Limit: Blocking", func(t *testing.T) {
		res := rs.Allow(token, tokenLimit, duration)
		if res.Allowed {
			t.Errorf("Token rate limit failed: should block the eleventh request")
		} else {
			fmt.Printf("Blocked until: %s\n", res.UnblockTime.Format(time.RFC1123))
			t.Logf("Successfully blocked as expected. Blocked until: %s", res.UnblockTime.Format(time.RFC1123))
		}
	})
	rs.Cleanup(token)

}

func TestRateLimiterByIPInMemory(t *testing.T) {
	setupEnvInMemory()

	rs := limiter.NewStorage() // Assuming this creates a new Redis client

	ip := "192.168.1.1"
	ipLimit := 5
	duration := time.Second // 1 second window for rate limit

	rs.Cleanup(ip)

	// test ip
	t.Run("IP Rate Limit: Allowance", func(t *testing.T) {
		for i := 0; i < ipLimit; i++ {
			res := rs.Allow(ip, ipLimit, duration)
			if !res.Allowed {
				t.Errorf("IP rate limit failed at request %d: should not block", i+1)
			}
		}
	})

	// test ip block
	t.Run("IP Rate Limit: Blocking", func(t *testing.T) {
		res := rs.Allow(ip, ipLimit, duration)
		if res.Allowed {
			t.Errorf("IP rate limit failed: should block the sixth request")
		} else {
			fmt.Printf("Blocked until: %s\n", res.UnblockTime.Format(time.RFC1123))
			t.Logf("Successfully blocked as expected. Blocked until: %s", res.UnblockTime.Format(time.RFC1123))
		}
	})

	rs.Cleanup(ip)
}

func TestRateLimiterByApiTokenInMemory(t *testing.T) {
	setupEnvInMemory()

	rs := limiter.NewStorage() // Assuming this creates a new Redis client

	token := "abc123"
	tokenLimit := 10
	duration := time.Second // 1 second window for rate limit

	rs.Cleanup(token)
	// Test Token-based rate limiting for allowed requests
	t.Run("Token Rate Limit: Allowance", func(t *testing.T) {
		for i := 0; i < tokenLimit; i++ {
			res := rs.Allow(token, tokenLimit, duration)
			if !res.Allowed {
				t.Errorf("Token rate limit failed at request %d: should not block", i+1)
			}
		}
	})

	// Test Token-based rate limiting for blocking
	t.Run("Token Rate Limit: Blocking", func(t *testing.T) {
		res := rs.Allow(token, tokenLimit, duration)
		if res.Allowed {
			t.Errorf("Token rate limit failed: should block the eleventh request")
		} else {
			fmt.Printf("Blocked until: %s\n", res.UnblockTime.Format(time.RFC1123))
			t.Logf("Successfully blocked as expected. Blocked until: %s", res.UnblockTime.Format(time.RFC1123))
		}
	})
	rs.Cleanup(token)

}
