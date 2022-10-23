package middleware

import (
	"fmt"
	"github.com/seanflannery10/ossa/errors"
	"net/http"
)

func RecoverPanic(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errors.ServerError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next(w, r)
	}
}
