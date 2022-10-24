package handlers

import (
	"github.com/seanflannery10/ossa/assert"
	"github.com/seanflannery10/ossa/helpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	Healthcheck(rr, r)

	rs := rr.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	body := helpers.GetBody(t, rr.Result())

	assert.Contains(t, body, "available")
}
