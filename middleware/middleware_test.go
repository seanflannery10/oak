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

	//t.Run("Bad JWKS", func(t *testing.T) {
	//	rr := httptest.NewRecorder()
	//	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	//
	//	r.Header.Add("Authorization", "Bearer 1234")
	//
	//	m.Authenticate(next).ServeHTTP(rr, r)
	//
	//	res := rr.Result()
	//	body := helpers.GetBody(t, res)
	//
	//	assert.Equal(t, res.Header.Get("Vary"), "Authorization")
	//	assert.Contains(t, body, "invalid or missing authentication token")
	//	assert.Equal(t, res.StatusCode, http.StatusUnauthorized)
	//})

	t.Run("Bad Token", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("Test"))
		}))
		defer srv.Close()

		m := New(srv.URL, nil)

		rr := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", nil)

		r.Header.Add("Authorization", "Bearer 1234")

		m.Authenticate(next).ServeHTTP(rr, r)

		res := rr.Result()
		body := helpers.GetBody(t, res)

		assert.Equal(t, res.Header.Get("Vary"), "Authorization")
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
