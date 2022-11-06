package middleware

import (
	"github.com/seanflannery10/ossa/assert"
	"github.com/seanflannery10/ossa/auth"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

var next = http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}))

func TestMiddleware_Chain(t *testing.T) {
	rr := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	m := New()
	m.Chain(m.RecoverPanic).Then(next).ServeHTTP(rr, r)

	assert.Contains(t, rr.Body.String(), "OK")
	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestMiddleware_Authenticate(t *testing.T) {
	tests := []struct {
		name   string
		header string
		jwks   string
		apiURL string
		body   string
		code   int
	}{
		{
			name:   "Missing Bearer",
			header: `Test Header`,
			jwks:   ``,
			apiURL: "",
			body:   "invalid or missing authentication token",
			code:   http.StatusUnauthorized,
		},
		{
			name:   "Bad Token",
			header: `Bearer Test`,
			jwks:   ``,
			apiURL: "",
			body:   "invalid or missing authentication token",
			code:   http.StatusUnauthorized,
		},
		{
			name:   "Bad Audience",
			header: `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6ImF0K2p3dCIsImtpZCI6IkFzc0o3ZjByQ1M0eE5nbzE4cndkZDlXZEJ5RG96eFc0TnhhQjNDMXZaREkifQ.eyJyb2xlX25hbWVzIjpbXSwianRpIjoidkVEQ0luM3RlOXdscG9MN2ZhUW1PIiwic3ViIjoiWFZ1QXpLTkdZWjVxMzd2SmRmc0N3IiwiaWF0IjoxNjY2ODM3NzIwLCJleHAiOjEwMDAwMDAwMDE2NjY4Mzc4MDAsImNsaWVudF9pZCI6IlhWdUF6S05HWVo1cTM3dkpkZnNDdyIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6MzAwMS9vaWRjIiwiYXVkIjoiaHR0cHM6Ly90ZXN0LmFwaSJ9.ffTskNjO-t9AFx6jpkmwqMfoY5DxCGJiVY3gjTRvasN6LahYWlZtFx7zaD9_MjJftpNzFAwJsqYmYLnB6GNvIM_oY_WUk_OP9_qGcTKvMtWYTRkI018q1zQLP_IpjQ3KVzGw_xIszbBXQ7wQrOKl28yN0UTQ-iraiPUVSVFQEfvoUZARc7zDtjlAHX-__fN5JoNMxQJfVdED_QDU8v9XbKI0ngwl99JwSeFxP9w8ByHmL3vgLXq8aQxODrfTZO-ev292Ziyy6iQWbLM6c-OnixYbPnugbpKF-d3rkIVsk6xx_cqyoAUnJc-Sz2un3ouqiLNMmK8IXqFnBW-fRMf6ApGDb3hTvH0wVxS4etJazH9Q0Np4CbB0O3B2dpi9gzHTqvNiTY6fPKrKYEvPKCVfewpp34r7-6hScfAEfluuz8y438OKWyANoOSwQ0ws58mh9F_3h59-sRYK23Mb_2yngmlJTvtY72RqzIEC617_lo0-8ABDqgr3ojJFrYbIhnYGxr2Wye715Ovu02Al3p580qGvqV3fCYtm7DUavfH7dyXQeS7Yb7YvejeIVrz-m_SRHO2FHYUDg268J9LGDpTVAPGPZVxD8jPMWNUeYCBtfPhlC0gTJI3By3qnq8hmke5gJTHonrzWTCHwkBcQUIdNzh1cAZe3uREzl7EbX2Hty5M`,
			jwks:   `{"keys":[{"kty":"RSA","use":"sig","kid":"AssJ7f0rCS4xNgo18rwdd9WdByDozxW4NxaB3Ccu1vZDI","e":"AQAB","n":"1DyGj6W-HjwAG4bwLlkWB9i8WLY3ILkFWpp5sCSMXptnJZ4JHpMEcIw4ecv4y2mfxBHvx52QTuh7_mzGILE3Aso-igyrqMtye8JrIFJ3P44m7i10MQ3si8yoUfaEZg2iiCyHTW2kPLspi9d1f7VY-aUTRlJavzxYd1KqvJmuxQ3fyxhJX-XAaW6HG2jDzrW0xSYSFB-lO56hHk1Hdr3JAd9Emwp4BoCsCqDv92dEMurfJEjv1O2hJRaHZ-cC_IFeuBDAQSnmjMRTytkVFQ8yUPUK3ozUIveZPoqk6KvBWjBMKHQjA1xS77MhcIfcmKFR-0tef2546TFRwWxAaVSubrzBy5MsD9f_sK7ntFLtJjk1ljfn1F8-eBMX3H2rd75xSWiAy6sUbNzFRfJiKRnzc1PTVMuHnSzcg_9HTv1Ve1F121w_ykRU316Nz8y9MpSSMuQUVCsOadi9Z9LL5iJkH7VTzSR3GJbOb5LnRuJa_9FChFn_P6hB_e7EbJ_ItVd36kCpezv4Iw_NdfFliQ3YZhAWr2ZDtPIgfLH6350LsYCz93wsLUVSRQLOhCsUWcD8EfOHURMiTMEs6hz4UOw603HFzHH93qTRVuPg6pKYLtZybnPZolr_j7Bpb-0h4chgKDa2rHz7LcUhz87eVXHExaKaQz9lLngqtkigslg4HOk"}]}`,
			apiURL: "Test URL",
			body:   "invalid or missing authentication token",
			code:   http.StatusUnauthorized,
		},
		{
			name:   "Bad Issuer",
			header: `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6ImF0K2p3dCIsImtpZCI6IkFzc0o3ZjByQ1M0eE5nbzE4cndkZDlXZEJ5RG96eFc0TnhhQjNDMXZaREkifQ.eyJyb2xlX25hbWVzIjpbXSwianRpIjoidkVEQ0luM3RlOXdscG9MN2ZhUW1PIiwic3ViIjoiWFZ1QXpLTkdZWjVxMzd2SmRmc0N3IiwiaWF0IjoxNjY2ODM3NzIwLCJleHAiOjEwMDAwMDAwMDE2NjY4Mzc4MDAsImNsaWVudF9pZCI6IlhWdUF6S05HWVo1cTM3dkpkZnNDdyIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6MzAwMS9vaWRjIiwiYXVkIjoiaHR0cHM6Ly90ZXN0LmFwaSJ9.ffTskNjO-t9AFx6jpkmwqMfoY5DxCGJiVY3gjTRvasN6LahYWlZtFx7zaD9_MjJftpNzFAwJsqYmYLnB6GNvIM_oY_WUk_OP9_qGcTKvMtWYTRkI018q1zQLP_IpjQ3KVzGw_xIszbBXQ7wQrOKl28yN0UTQ-iraiPUVSVFQEfvoUZARc7zDtjlAHX-__fN5JoNMxQJfVdED_QDU8v9XbKI0ngwl99JwSeFxP9w8ByHmL3vgLXq8aQxODrfTZO-ev292Ziyy6iQWbLM6c-OnixYbPnugbpKF-d3rkIVsk6xx_cqyoAUnJc-Sz2un3ouqiLNMmK8IXqFnBW-fRMf6ApGDb3hTvH0wVxS4etJazH9Q0Np4CbB0O3B2dpi9gzHTqvNiTY6fPKrKYEvPKCVfewpp34r7-6hScfAEfluuz8y438OKWyANoOSwQ0ws58mh9F_3h59-sRYK23Mb_2yngmlJTvtY72RqzIEC617_lo0-8ABDqgr3ojJFrYbIhnYGxr2Wye715Ovu02Al3p580qGvqV3fCYtm7DUavfH7dyXQeS7Yb7YvejeIVrz-m_SRHO2FHYUDg268J9LGDpTVAPGPZVxD8jPMWNUeYCBtfPhlC0gTJI3By3qnq8hmke5gJTHonrzWTCHwkBcQUIdNzh1cAZe3uREzl7EbX2Hty5M`,
			jwks:   `{"keys":[{"kty":"RSA","use":"sig","kid":"AssJ7f0rCS4xNgo18rwdd9WdByDozxW4NxaB3C1vZDI","e":"AQAB","n":"1DyGj6W-HjwAG4bwLlkWB9i8WLY3ILkFWpp5sCSMXptnJZ4JHpMEcIw4ecv4y2mfxBHvx52QTuh7_mzGILE3Aso-igyrqMtye8JrIFJ3P44m7i10MQ3si8yoUfaEZg2iiCyHTW2kPLspi9d1f7VY-aUTRlJavzxYd1KqvJmuxQ3fyxhJX-XAaW6HG2jDzrW0xSYSFB-lO56hHk1Hdr3JAd9Emwp4BoCsCqDv92dEMurfJEjv1O2hJRaHZ-cC_IFeuBDAQSnmjMRTytkVFQ8yUPUK3ozUIveZPoqk6KvBWjBMKHQjA1xS77MhcIfcmKFR-0tef2546TFRwWxAaVSubrzBy5MsD9f_sK7ntFLtJjk1ljfn1F8-eBMX3H2rd75xSWiAy6sUbNzFRfJiKRnzc1PTVMuHnSzcg_9HTv1Ve1F121w_ykRU316Nz8y9MpSSMuQUVCsOadi9Z9LL5iJkH7VTzSR3GJbOb5LnRuJa_9FChFn_P6hB_e7EbJ_ItVd36kCpezv4Iw_NdfFliQ3YZhAWr2ZDtPIgfLH6350LsYCz93wsLUVSRQLOhCsUWcD8EfOHURMiTMEs6hz4UOw603HFzHH93qTRVuPg6pKYLtZybnPZolr_j7Bpb-0h4chgKDa2rHz7LcUhz87eVXHExaKaQz9lLngqtkigslg4HOk"}]}`,
			apiURL: "https://test.api",
			body:   "invalid or missing authentication token",
			code:   http.StatusUnauthorized,
		},
		{
			name:   "Good Token",
			header: `Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6ImM3OTg5ZDc4ZjFlY2VkYTU2MWYzZWVhZmIyYmRhNDYxIn0.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjUwMDAwL29pZGMiLCJhdWQiOiJodHRwczovL3Rlc3QuYXBpIiwic3ViIjoiWFZ1QXpLTkdZWjVxMzd2SmRmc0N3IiwiY2xpZW50X2lkIjoiWFZ1QXpLTkdZWjVxMzd2SmRmc0N3IiwiZXhwIjpudWxsLCJpYXQiOjE2NjcwNzY4MDEsImp0aSI6IjAzOGFkMmViOTg4NWU5MGRmOGNkNGI0NDlmMzQzY2IyIn0.LXOvfiJ912c5C2qmiKPByPvhx7t2fJLEyy0rPVkIiP-JBEOmkFaQGWgLW6Kjn2Cg7FB8o1kdo1rmk6X-zR4eDk4gXHgdlmK3uahYyvm1KXvPDFu5A6vDhVPZER43RSr_sd_HrdkXMP1uHA_i57TEOfWlr6hRkV9m8CYqrBhQIenWjk2JYlTmE-q630fEWAZxJqpLkoZfprY49e6WQnTggRvU6_zogd9N0PdQc1sZVMSEdoCvH7TyJEJploo9NgAHAh5UAt7lr9HEtdKSLEH9pAT29X5Oatz5MiwV7iXD3v3N1X1eLI-V62Ag0z58NQQvzHSDkaxkknPQCYZ7abkLmQ`,
			jwks:   `{"keys":[{"kty": "RSA", "use": "sig", "kid": "c7989d78f1eceda561f3eeafb2bda461", "e": "AQAB", "n": "41d7rkhUP9lZT0f8vMNrGQolLUWRCEYJzgg6fU-XaOrUK1elk94LjW4TXDLSmRRUDQT_gfkYgVfzoj2iLgKP1owD6ND25rNLizD8G9vQxGhHYWKTAqb8Vf11XBtS86LNmrhfmjavKPo9IOeE7jzcEHULm4IjgQPak9YyQjYmNjoXtReLjyEkSybVu3r_xEKlwWPlI7COpX-u6xzDBZMSCnI-EyS2Wm-Rz3xGchvRhTmH9mD890AuvC44bULhhez-Xb9H2-9oCEjRDcYrD3jDSAPMA636RRR0Npg1PeGwhfvJesFTQ5HS17oldh90_VXlIoUJrZB9SVdJOWcIy8P5wQ"}]}`,
			apiURL: "https://test.api",
			body:   "OK",
			code:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, _ := net.Listen("tcp", "127.0.0.1:50000")

			srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.jwks))
			}))

			srv.Listener.Close()
			srv.Listener = l

			srv.Start()
			defer srv.Close()

			m := New()
			m.SetAuthenticateConfig(srv.URL+"/oidc/jwks", tt.apiURL)

			rr := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/oidc/jwks", nil)

			r.Header.Add("Authorization", tt.header)

			m.Authenticate(next).ServeHTTP(rr, r)

			assert.Equal(t, rr.Header().Get("Vary"), "Authorization")
			assert.Contains(t, rr.Body.String(), tt.body)
			assert.Equal(t, rr.Code, tt.code)
		})
	}
}

func TestMiddleware_RequireAuthenticatedUser(t *testing.T) {
	t.Run("Bad Auth", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		m := New()

		m.RequireAuthenticatedUser(next).ServeHTTP(rr, r)

		assert.Contains(t, rr.Body.String(), "you must be authenticated to access this resource")
		assert.Equal(t, rr.Code, http.StatusUnauthorized)
	})

	t.Run("Good Auth", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		r = auth.SetUser(r, "Test")

		m := New()

		m.RequireAuthenticatedUser(next).ServeHTTP(rr, r)

		assert.Equal(t, rr.Body.String(), "OK")
		assert.Equal(t, rr.Code, http.StatusOK)
	})
}

func TestMiddleware_CORS(t *testing.T) {
	t.Run("MethodGet", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		r.Header.Set("Origin", "127.0.0.1")

		m := New()
		m.SetCorsConfig([]string{"127.0.0.1"})

		m.CORS(next).ServeHTTP(rr, r)

		assert.Equal(t, rr.Body.String(), "OK")
		assert.Equal(t, rr.Code, http.StatusOK)
	})

	t.Run("MethodOptions", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodOptions, "/", nil)

		r.Header.Set("Origin", "127.0.0.1")
		r.Header.Set("Access-Control-Request-Method", "Test")

		m := New()
		m.SetCorsConfig([]string{"127.0.0.1"})

		m.CORS(next).ServeHTTP(rr, r)

		assert.Equal(t, rr.Body.String(), "")
		assert.Equal(t, rr.Code, http.StatusOK)
	})
}

func TestMiddleware_Metrics(t *testing.T) {
	rr := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	m := New()

	m.Metrics(next).ServeHTTP(rr, r)

	assert.Equal(t, rr.Body.String(), "OK")
}

func TestMiddleware_RateLimit(t *testing.T) {
	rr := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	m := New()
	m.SetRateLimitConfig(true, 2, 4)

	m.RateLimit(next).ServeHTTP(rr, r)

	assert.Equal(t, rr.Body.String(), "OK")
	//TODO: Test connection being limited
}

func TestMiddleware_RecoverPanic(t *testing.T) {
	t.Run("No Panic", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		m := New()

		m.RecoverPanic(next).ServeHTTP(rr, r)

		assert.Equal(t, rr.Body.String(), "OK")
	})

	t.Run("Panic", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		homeHandler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("test error") })

		m := New()

		m.RecoverPanic(homeHandler).ServeHTTP(rr, r)

		assert.Contains(t, rr.Body.String(), "the server encountered a problem and could not process your json")
	})
}
