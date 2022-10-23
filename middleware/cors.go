package middleware

import (
	"fmt"
	"github.com/seanflannery10/ossa/log"
	"net/http"
)

type CORSConfig struct {
	trustedOrigins []string
}

var gCORSConfig = &CORSConfig{
	trustedOrigins: nil,
}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	if gCORSConfig.trustedOrigins == nil {
		log.Fatal(fmt.Errorf("SetCORSConfig not set"), nil)
	}

	return gCORSConfig.CORS(next)
}

func SetCORSConfig(trustedOrigins []string) {
	gCORSConfig.trustedOrigins = trustedOrigins
}

func (c *CORSConfig) CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range c.trustedOrigins {
				if origin == c.trustedOrigins[i] {
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
