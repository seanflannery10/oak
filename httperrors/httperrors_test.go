package httperrors

import (
	"errors"
	"fmt"
	"github.com/seanflannery10/ossa/assert"
	"github.com/seanflannery10/ossa/validator"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorMessage(t *testing.T) {
	tests := []struct {
		sc   int
		body string
	}{
		{200, "testing status code 200"},
		{401, "testing status code 401"},
		{404, "testing status code 404"},
		{500, "testing status coOde 500"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.sc), func(t *testing.T) {
			rr := httptest.NewRecorder()

			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			ErrorMessage(rr, r, tt.sc, tt.body)

			assert.Contains(t, rr.Body.String(), tt.body)
			assert.Equal(t, rr.Result().StatusCode, tt.sc)
		})
	}
}

func TestFailedValidation(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	FailedValidation(rr, r, &validator.Validator{})

	assert.Equal(t, rr.Result().StatusCode, http.StatusUnprocessableEntity)
}

func TestStatusCodesWithError(t *testing.T) {
	tests := []struct {
		name string
		sc   int
		f    func(http.ResponseWriter, *http.Request, error)
	}{
		{
			name: "ServerError",
			sc:   http.StatusInternalServerError,
			f:    ServerError,
		},
		{
			name: "BadRequest",
			sc:   http.StatusBadRequest,
			f:    BadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			tt.f(rr, r, errors.New("test"))

			assert.Equal(t, rr.Result().StatusCode, tt.sc)
		})
	}
}

func TestStatusCodesWithoutError(t *testing.T) {
	tests := []struct {
		name string
		sc   int
		f    func(http.ResponseWriter, *http.Request)
	}{
		{
			name: "NotFound",
			sc:   http.StatusNotFound,
			f:    NotFound,
		},
		{
			name: "MethodNotAllowed",
			sc:   http.StatusMethodNotAllowed,
			f:    MethodNotAllowed,
		},
		{
			name: "InvalidAuthenticationToken",
			sc:   http.StatusUnauthorized,
			f:    InvalidAuthenticationToken,
		},
		{
			name: "AuthenticationRequired",
			sc:   http.StatusUnauthorized,
			f:    AuthenticationRequired,
		},
		{
			name: "RateLimitExceededResponse",
			sc:   http.StatusTooManyRequests,
			f:    RateLimitExceededResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			tt.f(rr, r)

			assert.Equal(t, rr.Result().StatusCode, tt.sc)
		})
	}
}
