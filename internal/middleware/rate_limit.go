package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	buckets   = make(map[string]*bucket)
	bucketsMu sync.Mutex
)

type bucket struct {
	count   int
	resetAt int64
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		windowMs := int64(60000)
		max := 30
		key := "global"
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
		// Remove port from ip if present
		if idx := strings.LastIndex(ip, ":"); idx != -1 {
			ip = ip[:idx]
		}
		bucketKey := key + ":" + ip
		now := time.Now().UnixMilli()

		bucketsMu.Lock()
		b, ok := buckets[bucketKey]
		if !ok || now > b.resetAt {
			b = &bucket{count: 0, resetAt: now + windowMs}
			buckets[bucketKey] = b
		}
		if b.count >= max {
			retryAfter := (b.resetAt - now) / 1000
			w.Header().Set("Retry-After", string(retryAfter))
			bucketsMu.Unlock()
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "rate_limited"})
			return
		}
		b.count++
		bucketsMu.Unlock()
		next.ServeHTTP(w, r)
	})
}
