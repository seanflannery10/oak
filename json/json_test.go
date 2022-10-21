package json

import (
	"fmt"
	"github.com/seanflannery10/oak/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		a string
		e string
	}{
		{`["foo", "bar"]`, "body contains incorrect encode type (at character 1)"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.a), func(t *testing.T) {
			w := httptest.NewRecorder()

			json := strings.NewReader(tt.a)

			r, err := http.NewRequest(http.MethodGet, "/", json)
			if err != nil {
				t.Fatal(err)
			}

			var testData struct {
				String string `a:"string"`
				Int    int    `a:"int"`
			}

			err = Decode(w, r, &testData)
			assert.StringContains(t, err.Error(), tt.e)
		})
	}
}
