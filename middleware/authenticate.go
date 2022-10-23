package middleware

import (
	"fmt"
	"github.com/pascaldekloe/jwt"
	"github.com/seanflannery10/ossa/errors"
	"github.com/seanflannery10/ossa/log"
	"net/http"
	"strings"
	"time"
)

type AuthenticateConfig struct {
	jwks string
}

var gAuthenticateConfig = &AuthenticateConfig{
	jwks: "",
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	if gAuthenticateConfig.jwks == "" {
		log.Fatal(fmt.Errorf("SetAuthenticateConfig not set"), nil)
	}

	return gAuthenticateConfig.Authenticate(next)
}

func SetAuthenticateConfig(jwks string) {
	gAuthenticateConfig.jwks = jwks
}

func (c *AuthenticateConfig) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			n, err := keys.LoadJWK([]byte(c.jwks))
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

		next(w, r)
	}
}
