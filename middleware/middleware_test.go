package middleware

import (
	"github.com/seanflannery10/ossa/assert"
	"github.com/seanflannery10/ossa/helpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	m    = New("test", nil)
	r, _ = http.NewRequest(http.MethodGet, "/", nil)

	next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
)

func TestMiddleware_Authenticate(t *testing.T) {
	rr := httptest.NewRecorder()

	m.Authenticate(next).ServeHTTP(rr, r)

	body := helpers.GetBody(t, rr.Result())

	assert.Equal(t, body, "OK")
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
