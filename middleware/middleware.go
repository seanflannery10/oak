package middleware

import (
	"expvar"
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/felixge/httpsnoop"
	"github.com/golang-jwt/jwt/v4"
	"github.com/justinas/alice"
	"github.com/seanflannery10/ossa/context"
	"github.com/seanflannery10/ossa/errors"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	AuthenticateConfig struct {
		JWKSURL string
		APIURL  string
	}

	CorsConfig struct {
		TrustedOrigins []string
	}

	RateLimitConfig struct {
		Enabled bool
		RPS     float64
		Burst   int
	}

	Middleware struct {
		authenticate AuthenticateConfig
		cors         CorsConfig
		rateLimit    RateLimitConfig
	}
)

func New() *Middleware {
	return &Middleware{}
}

func (m *Middleware) SetAuthenticateConfig(jwksURL, apiURL string) {
	m.authenticate.JWKSURL = jwksURL
	m.authenticate.APIURL = apiURL
}

func (m *Middleware) SetCorsConfig(trustedOrigins []string) {
	m.cors.TrustedOrigins = trustedOrigins
}

func (m *Middleware) SetRateLimitConfig(enabled bool, rps float64, burst int) {
	m.rateLimit.Enabled = enabled
	m.rateLimit.RPS = rps
	m.rateLimit.Burst = burst
}

func (m *Middleware) Chain(constructors ...alice.Constructor) alice.Chain {
	return alice.New(constructors...)
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader != "" {
			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				errors.InvalidAuthenticationToken(w, r)
				return
			}

			jwks, err := keyfunc.Get(m.authenticate.JWKSURL, keyfunc.Options{})
			if err != nil {
				errors.InvalidAuthenticationToken(w, r)
				return
			}

			tokenString := headerParts[1]

			token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, jwks.Keyfunc)
			if err != nil {
				errors.InvalidAuthenticationToken(w, r)
				return
			}

			claims, ok := token.Claims.(*jwt.RegisteredClaims)

			if !ok && !token.Valid {
				errors.InvalidAuthenticationToken(w, r)
				return
			}

			if !claims.VerifyAudience(m.authenticate.APIURL, false) {
				errors.InvalidAuthenticationToken(w, r)
				return
			}

			issuer := strings.TrimRight(m.authenticate.JWKSURL, "/jwks")

			if !claims.VerifyIssuer(issuer, false) {
				errors.InvalidAuthenticationToken(w, r)
				return
			}

			r = context.SetAuthenticatedUser(r, claims.Subject)
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticatedUser := context.GetAuthenticatedUser(r)

		if authenticatedUser == "" {
			errors.AuthenticationRequired(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range m.cors.TrustedOrigins {
				if origin == m.cors.TrustedOrigins[i] {
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

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Metrics(next http.Handler) http.Handler {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")
	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics := httpsnoop.CaptureMetrics(next, w, r)

		totalRequestsReceived.Add(1)
		totalResponsesSent.Add(1)
		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())
		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	})
}

func (m *Middleware) RateLimit(next http.Handler) http.Handler {
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

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.rateLimit.Enabled {
			ip := realip.FromRequest(r)

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(m.rateLimit.RPS), m.rateLimit.Burst),
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
	})
}

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errors.ServerError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
