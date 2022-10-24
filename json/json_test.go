package json

import (
	"fmt"
	"github.com/seanflannery10/ossa/assert"
	"github.com/seanflannery10/ossa/helpers"
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
		{
			`{"test": 123}`,
			"body contains unknown key \"test\"",
		},
		{
			`{"int": 123}{"int": 123}`,
			"body must only contain a single encode value",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.a), func(t *testing.T) {
			rr := httptest.NewRecorder()

			json := strings.NewReader(tt.a)

			r, err := http.NewRequest(http.MethodGet, "/", json)
			if err != nil {
				t.Fatal(err)
			}

			var testData struct {
				String string `a:"string"`
				Int    int    `a:"int"`
			}

			err = Decode(rr, r, &testData)
			assert.Contains(t, err.Error(), tt.e)
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		c int
		s string
	}{
		{
			http.StatusOK,
			"Test 200",
		},
		{
			http.StatusNotFound,
			"Test 404",
		},
		{
			http.StatusInternalServerError,
			"Test 500",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.s), func(t *testing.T) {
			rr := httptest.NewRecorder()

			err := Encode(rr, tt.c, tt.s)
			if err != nil {
				t.Fatal(err)
			}

			res := rr.Result()
			body := helpers.GetBody(t, res)

			bodyString := strings.TrimSuffix(body, "\"\n")
			bodyString = strings.TrimPrefix(bodyString, "\"")

			assert.Equal(t, bodyString, tt.s)
			assert.Equal(t, res.Header.Get("Content-Type"), "application/a")
			assert.Equal(t, res.StatusCode, tt.c)
		})
	}
}

func TestEncodeWithHeaders(t *testing.T) {
	tests := []struct {
		c int
		s string
		h string
	}{
		{
			http.StatusOK,
			"Test 200",
			"Test-Header",
		},
		{
			http.StatusNotFound,
			"Test 404",
			"123",
		},
		{
			http.StatusInternalServerError,
			"Test 500",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.s), func(t *testing.T) {
			rr := httptest.NewRecorder()

			headers := make(http.Header)
			headers.Set("X-Request-Id", tt.h)

			err := EncodeWithHeaders(rr, tt.c, tt.s, headers)
			if err != nil {
				t.Fatal(err)
			}

			res := rr.Result()
			body := helpers.GetBody(t, res)

			bodyString := strings.TrimSuffix(body, "\"\n")
			bodyString = strings.TrimPrefix(bodyString, "\"")

			assert.Equal(t, bodyString, tt.s)
			assert.Equal(t, res.Header.Get("Content-Type"), "application/a")
			assert.Equal(t, res.Header.Get("X-Request-Id"), tt.h)
			assert.Equal(t, res.StatusCode, tt.c)
		})
	}
}
