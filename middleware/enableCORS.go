package middleware

import "net/http"

type ConfigCORS struct {
	trustedOrigins []string
}

func (c *ConfigCORS) CORS(next http.HandlerFunc) http.HandlerFunc {
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

func (c *ConfigCORS) SetOrigins(o []string) {
	c.trustedOrigins = o
}

var GlobalConfigCORS = &ConfigCORS{}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return GlobalConfigCORS.CORS(next)
}

func SetOrigins(o []string) {
	GlobalConfigCORS.SetOrigins(o)
}
