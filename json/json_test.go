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
		{
			``,
			"body must not be empty",
		},
		{
			`<?xml version="1.0">`,
			"body contains badly-formed encode (at character 1)",
		},
		{
			`{"string": "test", }`,
			"body contains badly-formed encode (at character 20)",
		},
		{
			`["foo", "bar"]`,
			"body contains incorrect encode type (at character 1)",
		},
		{
			`{"string": 123}`,
			"body contains incorrect encode type for field \"String\"",
		},
		{
			`{"int": "123"}`,
			"body contains incorrect encode type for field \"Int\"",
		},
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
