package errors

import (
	"bytes"
	"github.com/seanflannery10/oak/log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var output bytes.Buffer

func TestErrorMessage(t *testing.T) {
	log.SetOutput(&output)
	w := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ErrorMessage(w, r, 200, "test")
	//assert.StringContains(t, output.String(), "200")
}
