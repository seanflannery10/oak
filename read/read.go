package read

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/seanflannery10/ossa/validator"
)

var errInvalidIDparameter = errors.New("invalid id parameter")

//func IDParam(r *http.Request) (int64, error) {
//	params := httprouter.ParamsFromContext(r.Context())
//
//	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
//	if err != nil || id < 1 {
//		return 0, errInvalidIDparameter
//	}
//
//	return id, nil
//}

func IDParam(qs url.Values) (int64, error) {
	s := qs.Get("id")

	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil || id < 1 {
		return 0, errInvalidIDparameter
	}

	return id, nil
}

func String(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func CSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func Int(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

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
