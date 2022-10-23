package middleware

import (
	"github.com/seanflannery10/ossa/errors"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

type RateLimitConfig struct {
	enabled bool
	rps     float64
	burst   int
}

var gRateLimitConfig = &RateLimitConfig{
	enabled: true,
	rps:     4,
	burst:   2,
}

func RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return gRateLimitConfig.RateLimit(next)
}

func SetRateLimitConfig(enabled bool, rps float64, burst int) {
	gRateLimitConfig.enabled = enabled
	gRateLimitConfig.rps = rps
	gRateLimitConfig.burst = burst
}

func (c *RateLimitConfig) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		if c.enabled {
			ip := realip.FromRequest(r)

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(c.rps), c.burst),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				errors.RateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next(w, r)
	}
}
