package middleware

type Middleware struct {
	jwks           string
	trustedOrigins []string
	limiter        Limiter
}

type Limiter struct {
	enabled bool
	rps     float64
	burst   int
}

func New(jwks string, trustedOrigins []string, limiter Limiter) *Middleware {
	return &Middleware{
		jwks:           jwks,
		trustedOrigins: trustedOrigins,
		limiter:        limiter,
	}
}
