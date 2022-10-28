package middleware

import (
	"github.com/seanflannery10/ossa/assert"
	"github.com/seanflannery10/ossa/helpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	m = New()

	r, _ = http.NewRequest(http.MethodGet, "/", nil)

	next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
)

func TestMiddleware_Authenticate(t *testing.T) {
	t.Run("Missing Bearer", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		r.Header.Add("Authorization", "Test Header")

		m.Authenticate(next).ServeHTTP(rr, r)

		res := rr.Result()
		body := helpers.GetBody(t, res)

		assert.Equal(t, res.Header.Get("Vary"), "Authorization")
		assert.Contains(t, body, "invalid or missing authentication token")
		assert.Equal(t, res.StatusCode, http.StatusUnauthorized)
	})

	t.Run("Bad Token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		r.Header.Add("Authorization", "Bearer 1234")

		m.Authenticate(next).ServeHTTP(rr, r)

		res := rr.Result()
		body := helpers.GetBody(t, res)

		assert.Contains(t, body, "invalid or missing authentication token")
		assert.Equal(t, res.StatusCode, http.StatusUnauthorized)
	})

	t.Run("Good Token, Bad Audience", func(t *testing.T) {
		const testJWKS = `{"keys":[{"kty":"RSA","use":"sig","kid":"AssJ7f0rCS4xNgo18rwdd9WdByDozxW4NxaB3C1vZDI","e":"AQAB","n":"1DyGj6W-HjwAG4bwLlkWB9i8WLY3ILkFWpp5sCSMXptnJZ4JHpMEcIw4ecv4y2mfxBHvx52QTuh7_mzGILE3Aso-igyrqMtye8JrIFJ3P44m7i10MQ3si8yoUfaEZg2iiCyHTW2kPLspi9d1f7VY-aUTRlJavzxYd1KqvJmuxQ3fyxhJX-XAaW6HG2jDzrW0xSYSFB-lO56hHk1Hdr3JAd9Emwp4BoCsCqDv92dEMurfJEjv1O2hJRaHZ-cC_IFeuBDAQSnmjMRTytkVFQ8yUPUK3ozUIveZPoqk6KvBWjBMKHQjA1xS77MhcIfcmKFR-0tef2546TFRwWxAaVSubrzBy5MsD9f_sK7ntFLtJjk1ljfn1F8-eBMX3H2rd75xSWiAy6sUbNzFRfJiKRnzc1PTVMuHnSzcg_9HTv1Ve1F121w_ykRU316Nz8y9MpSSMuQUVCsOadi9Z9LL5iJkH7VTzSR3GJbOb5LnRuJa_9FChFn_P6hB_e7EbJ_ItVd36kCpezv4Iw_NdfFliQ3YZhAWr2ZDtPIgfLH6350LsYCz93wsLUVSRQLOhCsUWcD8EfOHURMiTMEs6hz4UOw603HFzHH93qTRVuPg6pKYLtZybnPZolr_j7Bpb-0h4chgKDa2rHz7LcUhz87eVXHExaKaQz9lLngqtkigslg4HOk"}]}`
		const testHeader = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6ImF0K2p3dCIsImtpZCI6IkFzc0o3ZjByQ1M0eE5nbzE4cndkZDlXZEJ5RG96eFc0TnhhQjNDMXZaREkifQ.eyJyb2xlX25hbWVzIjpbXSwianRpIjoidkVEQ0luM3RlOXdscG9MN2ZhUW1PIiwic3ViIjoiWFZ1QXpLTkdZWjVxMzd2SmRmc0N3IiwiaWF0IjoxNjY2ODM3NzIwLCJleHAiOjEwMDAwMDAwMDE2NjY4Mzc4MDAsImNsaWVudF9pZCI6IlhWdUF6S05HWVo1cTM3dkpkZnNDdyIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6MzAwMS9vaWRjIiwiYXVkIjoiaHR0cHM6Ly90ZXN0LmFwaSJ9.ffTskNjO-t9AFx6jpkmwqMfoY5DxCGJiVY3gjTRvasN6LahYWlZtFx7zaD9_MjJftpNzFAwJsqYmYLnB6GNvIM_oY_WUk_OP9_qGcTKvMtWYTRkI018q1zQLP_IpjQ3KVzGw_xIszbBXQ7wQrOKl28yN0UTQ-iraiPUVSVFQEfvoUZARc7zDtjlAHX-__fN5JoNMxQJfVdED_QDU8v9XbKI0ngwl99JwSeFxP9w8ByHmL3vgLXq8aQxODrfTZO-ev292Ziyy6iQWbLM6c-OnixYbPnugbpKF-d3rkIVsk6xx_cqyoAUnJc-Sz2un3ouqiLNMmK8IXqFnBW-fRMf6ApGDb3hTvH0wVxS4etJazH9Q0Np4CbB0O3B2dpi9gzHTqvNiTY6fPKrKYEvPKCVfewpp34r7-6hScfAEfluuz8y438OKWyANoOSwQ0ws58mh9F_3h59-sRYK23Mb_2yngmlJTvtY72RqzIEC617_lo0-8ABDqgr3ojJFrYbIhnYGxr2Wye715Ovu02Al3p580qGvqV3fCYtm7DUavfH7dyXQeS7Yb7YvejeIVrz-m_SRHO2FHYUDg268J9LGDpTVAPGPZVxD8jPMWNUeYCBtfPhlC0gTJI3By3qnq8hmke5gJTHonrzWTCHwkBcQUIdNzh1cAZe3uREzl7EbX2Hty5M`

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(testJWKS))
		}))
		defer srv.Close()

		m.SetAuthenticateConfig(srv.URL+"/oidc/jwksURL", "test")

		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/oidc/jwks", nil)

		r.Header.Add("Authorization", testHeader)

		m.Authenticate(next).ServeHTTP(rr, r)

		res := rr.Result()
		body := helpers.GetBody(t, res)

		assert.Contains(t, body, "invalid or missing authentication token")
		assert.Equal(t, res.StatusCode, http.StatusUnauthorized)
	})

	t.Run("Good Token, Bad Issuer", func(t *testing.T) {
		const testJWKS = `{"keys":[{"kty":"RSA","use":"sig","kid":"AssJ7f0rCS4xNgo18rwdd9WdByDozxW4NxaB3C1vZDI","e":"AQAB","n":"1DyGj6W-HjwAG4bwLlkWB9i8WLY3ILkFWpp5sCSMXptnJZ4JHpMEcIw4ecv4y2mfxBHvx52QTuh7_mzGILE3Aso-igyrqMtye8JrIFJ3P44m7i10MQ3si8yoUfaEZg2iiCyHTW2kPLspi9d1f7VY-aUTRlJavzxYd1KqvJmuxQ3fyxhJX-XAaW6HG2jDzrW0xSYSFB-lO56hHk1Hdr3JAd9Emwp4BoCsCqDv92dEMurfJEjv1O2hJRaHZ-cC_IFeuBDAQSnmjMRTytkVFQ8yUPUK3ozUIveZPoqk6KvBWjBMKHQjA1xS77MhcIfcmKFR-0tef2546TFRwWxAaVSubrzBy5MsD9f_sK7ntFLtJjk1ljfn1F8-eBMX3H2rd75xSWiAy6sUbNzFRfJiKRnzc1PTVMuHnSzcg_9HTv1Ve1F121w_ykRU316Nz8y9MpSSMuQUVCsOadi9Z9LL5iJkH7VTzSR3GJbOb5LnRuJa_9FChFn_P6hB_e7EbJ_ItVd36kCpezv4Iw_NdfFliQ3YZhAWr2ZDtPIgfLH6350LsYCz93wsLUVSRQLOhCsUWcD8EfOHURMiTMEs6hz4UOw603HFzHH93qTRVuPg6pKYLtZybnPZolr_j7Bpb-0h4chgKDa2rHz7LcUhz87eVXHExaKaQz9lLngqtkigslg4HOk"}]}`
		const testHeader = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6ImF0K2p3dCIsImtpZCI6IkFzc0o3ZjByQ1M0eE5nbzE4cndkZDlXZEJ5RG96eFc0TnhhQjNDMXZaREkifQ.eyJyb2xlX25hbWVzIjpbXSwianRpIjoidkVEQ0luM3RlOXdscG9MN2ZhUW1PIiwic3ViIjoiWFZ1QXpLTkdZWjVxMzd2SmRmc0N3IiwiaWF0IjoxNjY2ODM3NzIwLCJleHAiOjEwMDAwMDAwMDE2NjY4Mzc4MDAsImNsaWVudF9pZCI6IlhWdUF6S05HWVo1cTM3dkpkZnNDdyIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6MzAwMS9vaWRjIiwiYXVkIjoiaHR0cHM6Ly90ZXN0LmFwaSJ9.ffTskNjO-t9AFx6jpkmwqMfoY5DxCGJiVY3gjTRvasN6LahYWlZtFx7zaD9_MjJftpNzFAwJsqYmYLnB6GNvIM_oY_WUk_OP9_qGcTKvMtWYTRkI018q1zQLP_IpjQ3KVzGw_xIszbBXQ7wQrOKl28yN0UTQ-iraiPUVSVFQEfvoUZARc7zDtjlAHX-__fN5JoNMxQJfVdED_QDU8v9XbKI0ngwl99JwSeFxP9w8ByHmL3vgLXq8aQxODrfTZO-ev292Ziyy6iQWbLM6c-OnixYbPnugbpKF-d3rkIVsk6xx_cqyoAUnJc-Sz2un3ouqiLNMmK8IXqFnBW-fRMf6ApGDb3hTvH0wVxS4etJazH9Q0Np4CbB0O3B2dpi9gzHTqvNiTY6fPKrKYEvPKCVfewpp34r7-6hScfAEfluuz8y438OKWyANoOSwQ0ws58mh9F_3h59-sRYK23Mb_2yngmlJTvtY72RqzIEC617_lo0-8ABDqgr3ojJFrYbIhnYGxr2Wye715Ovu02Al3p580qGvqV3fCYtm7DUavfH7dyXQeS7Yb7YvejeIVrz-m_SRHO2FHYUDg268J9LGDpTVAPGPZVxD8jPMWNUeYCBtfPhlC0gTJI3By3qnq8hmke5gJTHonrzWTCHwkBcQUIdNzh1cAZe3uREzl7EbX2Hty5M`

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(testJWKS))
		}))
		defer srv.Close()

		m.SetAuthenticateConfig(srv.URL+"/oidc/jwksURL", "https://test.api")

		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/oidc/jwks", nil)

		r.Header.Add("Authorization", testHeader)

		m.Authenticate(next).ServeHTTP(rr, r)

		res := rr.Result()
		body := helpers.GetBody(t, res)

		assert.Contains(t, body, "invalid or missing authentication token")
		assert.Equal(t, res.StatusCode, http.StatusUnauthorized)
	})
}

func TestMiddleware_CORS(t *testing.T) {
	rr := httptest.NewRecorder()

	m.CORS(next).ServeHTTP(rr, r)

	body := helpers.GetBody(t, rr.Result())

	assert.Equal(t, body, "OK")
}

func TestMiddleware_Metrics(t *testing.T) {
	rr := httptest.NewRecorder()

	m.Metrics(next).ServeHTTP(rr, r)

	body := helpers.GetBody(t, rr.Result())

	assert.Equal(t, body, "OK")
}

func TestMiddleware_RateLimit(t *testing.T) {
	rr := httptest.NewRecorder()

	m.RateLimit(next).ServeHTTP(rr, r)

	body := helpers.GetBody(t, rr.Result())

	assert.Equal(t, body, "OK")
}

func TestMiddleware_RecoverPanic(t *testing.T) {
	rr := httptest.NewRecorder()

	m.RecoverPanic(next).ServeHTTP(rr, r)

	body := helpers.GetBody(t, rr.Result())

	assert.Equal(t, body, "OK")
}
