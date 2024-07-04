package limiter

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		apiKey := r.Header.Get("API_KEY")

		limit, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_IP")) // Default limit per IP
		ms, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_DURATION"))
		duration := time.Millisecond * time.Duration(ms)

		if apiKey != "" {
			limit, _ = strconv.Atoi(os.Getenv("RATE_LIMIT_TOKEN"))
		}

		result := rl.Allow(ip, limit, duration)
		if !result.Allowed {
			w.WriteHeader(http.StatusTooManyRequests)
			if result.IsBlocked {
				//w.Write([]byte(fmt.Sprintf("You have reached the maximum number of requests or actions allowed. You are blocked until %s.", result.UnblockTime.Format(time.RFC1123))))
				w.Write([]byte(fmt.Sprintf("You have reached the maximum number of requests or actions allowed.")))
			} else {
				w.Write([]byte("You have reached the maximum number of requests or actions allowed."))
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
