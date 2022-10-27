package middleware

import (
	"github.com/seanflannery10/ossa/assert"
	"github.com/seanflannery10/ossa/helpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	m    = New("", nil)
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
		const testJWKS = `{"keys":[{"kty":"RSA","use":"sig","kid":"AssJ7f0rCS4xNgo18rwdd9WdByDozxW4NxaB3C1vZDI","e":"AQAB","n":"1DyGj6W-HjwAG4bwLlkWB9i8WLY3ILkFWpp5sCSMXptnJZ4JHpMEcIw4ecv4y2mfxBHvx52QTuh7_mzGILE3Aso-igyrqMtye8JrIFJ3P44m7i10MQ3si8yoUfaEZg2iiCyHTW2kPLspi9d1f7VY-aUTRlJavzxYd1KqvJmuxQ3fyxhJX-XAaW6HG2jDzrW0xSYSFB-lO56hHk1Hdr3JAd9Emwp4BoCsCqDv92dEMurfJEjv1O2hJRaHZ-cC_IFeuBDAQSnmjMRTytkVFQ8yUPUK3ozUIveZPoqk6KvBWjBMKHQjA1xS77MhcIfcmKFR-0tef2546TFRwWxAaVSubrzBy5MsD9f_sK7ntFLtJjk1ljfn1F8-eBMX3H2rd75xSWiAy6sUbNzFRfJiKRnzc1PTVMuHnSzcg_9HTv1Ve1F121w_ykRU316Nz8y9MpSSMuQUVCsOadi9Z9LL5iJkH7VTzSR3GJbOb5LnRuJa_9FChFn_P6hB_e7EbJ_ItVd36kCpezv4Iw_NdfFliQ3YZhAWr2ZDtPIgfLH6350LsYCz93wsLUVSRQLOhCsUWcD8EfOHURMiTMEs6hz4UOw603HFzHH93qTRVuPg6pKYLtZybnPZolr_j7Bpb-0h4chgKDa2rHz7LcUhz87eVXHExaKaQz9lLngqtkigslg4HOk"}]}`
		const testHeader = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6ImF0K2p3dCIsImtpZCI6IkFzc0o3ZjByQ1M0eE5nbzE4cndkZDlXZEJ5RG96eFc0TnhhQjNDMXZaREkifQ.eyJyb2xlX25hbWVzIjpbXSwianRpIjoiNmlWVnNCMUlBbEZPTEx1YmVzdGFPIiwic3ViIjoiWFZ1QXpLTkdZWjVxMzd2SmRmc0N3IiwiaWF0IjoxNjY2ODMzMDEzLCJleHAiOjE2NjY4MzY2MTMsImNsaWVudF9pZCI6IlhWdUF6S05HWVo1cTM3dkpkZnNDdyIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6MzAwMS9vaWRjIiwiYXVkIjoiaHR0cHM6Ly90ZXN0LmFwaSJ9.AX8TraZX0qNDEfyq99WcaUHN3jy89AXNRkcqfxaFFZvfjuTKsptwmBObgqkG7MmMhDaJZjwma-skTLUJgyUn1Tt20lSbiO4ZRiZ4QSvSG2Sd4L8OzETIYM5t7lI-8ASQRDBCYEJwQrvcp_PjxXGDSvQ6QJTQ_nlNgW0m1bE74cSWvFlgQ3VHAFwHVSOMa8Kulbb_husDzCX4wBXIs6C5WtymPiHT_UvHliBiFVrdRDiNstzPqSwfMNv2W-H-7MbIJgNUKz10rZrgayFEGwj_gyx7QH6iJxxoQBLWXZLwunk2CC4ZAm9yUTyYpKLkYGrO7RRt98B67FkoNhYe5PCMnePXTOIJz06qp10YOc4vGLziFd_d5P-zrA3R6lQlmsJibk2tZKL5bvHlOOxHwTt9j9m4a36Pi6tNQbyIL11dBjlw68X4zR3ctJTqKi5mAbBrXXnrpO_qwcA4CQTh60dLZhmlyhwn-mgyjFOpT8VCunOLDcIgRFrYhF4cfP-S1uKBZGF29QC-NuA96zdB6K6B4d5kUMLjDJrRzhR--tRvrXnZz_b_PegsNbYT4vbEA0A_V5ZghbtgqgXAUgzKDEAgboATw0Wam-00fJXbdmJO9Q2UYusZy1S_8K1HerbXYMOtU9Mi0O2USnqYc6iSTzZb55gfIK2nO0qIGryx0Y_83Ww`

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte(testJWKS))
		}))
		defer srv.Close()

		m := New(srv.URL, nil)

		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		r.Header.Add("Authorization", testHeader)

		m.Authenticate(next).ServeHTTP(rr, r)

		res := rr.Result()
		body := helpers.GetBody(t, res)

		assert.Equal(t, res.Header.Get("Vary"), "Authorization")
		assert.Contains(t, body, "")
		assert.Equal(t, res.StatusCode, http.StatusOK)
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
