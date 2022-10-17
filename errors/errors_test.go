package errors

import (
	"github.com/seanflannery10/oak/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorMessage(t *testing.T) {
	tests := []struct {
		name    string
		sc      int
		message string
	}{
		{
			name:    "200",
			sc:      200,
			message: "testing 200 sc code",
		},
		{
			name:    "401",
			sc:      401,
			message: "testing 401 sc code",
		},
		{
			name:    "404",
			sc:      404,
			message: "testing 404 sc code",
		},
		{
			name:    "500",
			sc:      500,
			message: "testing 500 sc code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			ErrorMessage(w, r, tt.sc, tt.message)

			res := w.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.StringContains(t, string(body), tt.message)
			assert.Equal(t, res.StatusCode, tt.sc)
		})
	}

}
