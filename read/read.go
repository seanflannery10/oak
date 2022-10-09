package read

import (
	"errors"
	"github.com/seanflannery10/oak/validator"
	"net/http"
	"strconv"
	"strings"
)

func ID(r *http.Request) (int, error) {
	s := r.URL.Query().Get("id")

	i, err := strconv.Atoi(s)
	if err != nil || i < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return i, nil
}

func String(r *http.Request, key string, defaultValue string) string {
	s := r.URL.Query().Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func CSV(r *http.Request, key string, defaultValue []string) []string {
	csv := r.URL.Query().Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func Int(r *http.Request, key string, defaultValue int, v *validator.Validator) int {
	s := r.URL.Query().Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError("must be an integer value")
		return defaultValue
	}

	return i
}
