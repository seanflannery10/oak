package middleware

import (
	"expvar"
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/felixge/httpsnoop"
	"github.com/golang-jwt/jwt/v4"
	"github.com/justinas/alice"
	"github.com/seanflannery10/ossa/errors"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Middleware struct {
	jwks           string
	trustedOrigins []string
	rlc            RateLimitConfig
}

type RateLimitConfig struct {
	enabled bool
	rps     float64
	burst   int
}

func New(jwks string, trustedOrigins []string) *Middleware {
	return &Middleware{
		jwks:           jwks,
		trustedOrigins: trustedOrigins,
		rlc: RateLimitConfig{
			enabled: true,
			rps:     2,
			burst:   4,
		},
	}
}

func (m *Middleware) SetRateLimitConfig(enabled bool, rps float64, burst int) {
	m.rlc.enabled = enabled
	m.rlc.rps = rps
	m.rlc.burst = burst
}

func (m *Middleware) Chain(constructors ...alice.Constructor) alice.Chain {
	return alice.New(constructors...)
}

func (m *Middleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader != "" {
			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				errors.InvalidAuthenticationToken(w, r)
				return
			}

			jwks, err := keyfunc.Get(m.jwks, keyfunc.Options{})
			if err != nil {
				errors.InvalidAuthenticationToken(w, r)
				return
			}

			tokenString := headerParts[1]

			_, err = jwt.Parse(tokenString, jwks.Keyfunc)
			if err != nil {
				errors.InvalidAuthenticationToken(w, r)
				return
			}
		}

		next(w, r)
	}
}

func (m *Middleware) CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range m.trustedOrigins {
				if origin == m.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						w.WriteHeader(http.StatusOK)
						return
					}

					break
				}
			}
		}

		next(w, r)
	}
}

func (m *Middleware) Metrics(next http.HandlerFunc) http.HandlerFunc {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")
	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

	return func(w http.ResponseWriter, r *http.Request) {
		metrics := httpsnoop.CaptureMetrics(next, w, r)

		totalRequestsReceived.Add(1)
		totalResponsesSent.Add(1)
		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())
		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	}
}

func (m *Middleware) RateLimit(next http.HandlerFunc) http.HandlerFunc {
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
		if m.rlc.enabled {
			ip := realip.FromRequest(r)

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(m.rlc.rps), m.rlc.burst),
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

func (m *Middleware) RecoverPanic(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errors.ServerError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next(w, r)
	}
}
