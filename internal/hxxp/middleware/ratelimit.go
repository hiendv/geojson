package middleware

// Credit: https://www.alexedwards.net/blog/how-to-rate-limit-http-requests with some modification

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/hiendv/geojson/internal/hxxp/ctxx"
	"github.com/hiendv/geojson/pkg/util"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter   *rate.Limiter
	expiredAt time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.RWMutex
)

func init() {
	go cleanupLimit()
}

func limit(ctx context.Context, rps float64, burst int, ttl time.Duration) *rate.Limiter {
	ip, ok := ctxx.IP(ctx)
	if !ok {
		return nil
	}

	mu.RLock()
	v, ok := visitors[ip]
	mu.RUnlock()

	if !ok {
		limiter := rate.NewLimiter(rate.Limit(rps), burst)
		mu.Lock()
		visitors[ip] = &visitor{limiter, time.Now().Add(ttl)}
		mu.Unlock()
		return limiter
	}

	mu.Lock()
	v.expiredAt = time.Now().Add(ttl)
	mu.Unlock()
	return v.limiter
}

func cleanupLimit() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.expiredAt) > 0 {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func RateLimit(next http.Handler, rps float64, burst int, ttl time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := limit(r.Context(), rps, burst, ttl)
		if limiter == nil {
			util.HTTPAbort(w, "unable to identify client", http.StatusInternalServerError)
			return
		}

		if !limiter.Allow() {
			util.HTTPAbort(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
