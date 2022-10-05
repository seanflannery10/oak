package middleware

import (
	"fmt"
	"github.com/pascaldekloe/jwt"
	"github.com/seanflannery10/oak/errors"
	"net/http"
	"strings"
	"time"
)

type Middleware struct {
	jwks string
}

func New(jwks string) *Middleware {
	return &Middleware{
		jwks: jwks,
	}
}

func (m *Middleware) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errors.ServerError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
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

			token := headerParts[1]

			var keys jwt.KeyRegister
			n, err := keys.LoadJWK([]byte(m.jwks))
			if n != 1 || err != nil {
				errors.InvalidAuthenticationToken(w, r)
				return
			}

			claims, err := keys.Check([]byte(token))
			if err != nil {
				errors.InvalidAuthenticationToken(w, r)
				return
			}

			if !claims.Valid(time.Now()) {
				errors.InvalidAuthenticationToken(w, r)
				return
			}

			//if claims.Issuer != app.config.baseURL {
			//	errors.InvalidAuthenticationToken(w, r)
			//	return
			//}
			//
			//if !claims.AcceptAudience(app.config.baseURL) {
			//	errors.InvalidAuthenticationToken(w, r)
			//	return
			//}

			//userID, err := strconv.Atoi(claims.Subject)
			//if err != nil {
			//	errors.ServerError(w, r, err)
			//	return
			//}

			//user, err := app.db.GetUser(userID)
			//if err != nil {
			//	errors.ServerError(w, r, err)
			//	return
			//}

			//if user != nil {
			//	r = contextSetAuthenticatedUser(r, user)
			//}
		}

		next.ServeHTTP(w, r)
	})
}
