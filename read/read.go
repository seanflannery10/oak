package read

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/seanflannery10/ossa/validator"
)

var errInvalidIDparameter = errors.New("invalid id parameter")

func IDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errInvalidIDparameter
	}

	return id, nil
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
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}
