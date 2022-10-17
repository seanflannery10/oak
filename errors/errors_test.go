package errors

import (
	"github.com/seanflannery10/oak/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorMessage(t *testing.T) {
	w := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ErrorMessage(w, r, 200, "Test")

	res, err := io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.StringContains(t, string(res), "Test")
}
