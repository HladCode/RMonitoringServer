package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type authLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	interval time.Duration
}

func newAuthLimiter(limit int, interval time.Duration) *authLimiter {
	return &authLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		interval: interval,
	}
}

// TODO: анулювати користувачів, які успішно залогінились чи оновили токен
func (al *authLimiter) allow(ip string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-al.interval)

	requestTimes := al.requests[ip]

	var recent []time.Time
	for _, t := range requestTimes {
		if t.After(windowStart) {
			recent = append(recent, t)
		}
	}

	if len(recent) >= al.limit {
		return false
	}

	recent = append(recent, now)
	al.requests[ip] = recent
	return true
}

func AuthenticationRateLimiter(limit int, interval time.Duration) func(http.Handler) http.Handler {
	limiter := newAuthLimiter(limit, interval)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Invalid IP", http.StatusInternalServerError)
				return
			}

			if !limiter.allow(ip) {
				http.Error(w, "Too many authentication attempts. Try again later.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
