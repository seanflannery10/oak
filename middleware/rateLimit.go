package middleware

import (
	"github.com/seanflannery10/oak/errors"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

func RateLimit(next http.HandlerFunc) http.HandlerFunc {
	type Limiter struct {
		enabled bool
		rps     float64
		burst   int
	}

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
		if m.limiter.enabled {
			ip := realip.FromRequest(r)

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(m.limiter.rps), m.limiter.burst),
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

		next.ServeHTTP(w, r)
	}
}
