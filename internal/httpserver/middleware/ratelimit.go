package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimit returns a middleware that limits requests per IP to `max` within `window`.
// This is an in-memory implementation suitable for single-instance deployments.
// For production or multiple instances, use a distributed store like Redis.
func RateLimit(max int, window time.Duration) func(http.Handler) http.Handler {
	type entry struct {
		mu    sync.Mutex
		times []time.Time
	}

	var store sync.Map // map[string]*entry

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
			}
			v, _ := store.LoadOrStore(ip, &entry{})
			e := v.(*entry)

			e.mu.Lock()
			now := time.Now()
			cutoff := now.Add(-window)
			// drop old
			i := 0
			for ; i < len(e.times); i++ {
				if e.times[i].After(cutoff) {
					break
				}
			}
			if i > 0 {
				e.times = e.times[i:]
			}
			if len(e.times) >= max {
				// calculate retry-after
				retryAfter := int(window.Seconds())
				if len(e.times) > 0 {
					earliest := e.times[0]
					retrySecs := int(earliest.Add(window).Sub(now).Seconds())
					if retrySecs > 0 {
						retryAfter = retrySecs
					} else {
						retryAfter = 1
					}
				}
				e.mu.Unlock()
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			e.times = append(e.times, now)
			e.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
