package limiter

import (
	"net/http"
	"strings"
	"time"
)

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		apiKey := r.Header.Get("API_KEY")

		limit := 10 // default limit per IP
		duration := time.Second

		if apiKey != "" {
			// Example of token-based rate limit, you can enhance it as per requirements
			limit = 100
		}

		if !rl.Allow(ip, limit, duration) {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
