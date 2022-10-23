package read

import (
	"fmt"
	"github.com/seanflannery10/ossa/assert"
	"github.com/seanflannery10/ossa/validator"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestCSV(t *testing.T) {
	tests := []struct {
		key string
		csv string
	}{
		{
			"csv",
			"csv,foo,bar",
		},
		{
			"csv",
			"csv,foo,bar,test",
		},
		{
			"string",
			"string",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.csv), func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			split := strings.Split(tt.csv, ",")

			r.URL.RawQuery = url.Values{
				tt.key: {tt.csv},
			}.Encode()

			res := CSV(r, tt.key, nil)

			assert.Equal(t, len(res), len(split))
			assert.Equal(t, res[0], split[0])
			assert.Equal(t, res[len(res)-1], split[len(split)-1])
		})
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		key string
		int int
	}{
		{
			"test",
			42,
		},
		{
			"test",
			0,
		},
		{
			"test",
			-20,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.int), func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			r.URL.RawQuery = url.Values{
				tt.key: {strconv.Itoa(tt.int)},
			}.Encode()

			res := Int(r, tt.key, 1, &validator.Validator{})

			assert.Equal(t, res, tt.int)
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		key    string
		string string
	}{
		{
			"test",
			"42",
		},
		{
			"test",
			"test",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.string), func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			r.URL.RawQuery = url.Values{
				tt.key: {tt.string},
			}.Encode()

			res := String(r, tt.key, "")

			assert.Equal(t, res, tt.string)
		})
	}
}
