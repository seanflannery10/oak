package helpers

import (
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/seanflannery10/ossa/validator"
)

var (
	errBadlyFormed        = errors.New("body contains badly-formed encode")
	errIncorrectEncode    = errors.New("body contains incorrect encode type")
	errEmptyBody          = errors.New("body must not be empty")
	errUnknownKey         = errors.New("body contains unknown key")
	errBodyToLarge        = errors.New("body must not be larger than")
	errToManyValues       = errors.New("body must only contain a single encode value")
	errInvalidIDParameter = errors.New("invalid id parameter")
)

func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var (
			syntaxError           *json.SyntaxError
			unmarshalTypeError    *json.UnmarshalTypeError
			invalidUnmarshalError *json.InvalidUnmarshalError
		)

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("%w (at character %d)", errBadlyFormed, syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errBadlyFormed

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("%w for field %q", errIncorrectEncode, unmarshalTypeError.Field)
			}

			return fmt.Errorf("%w (at character %d)", errIncorrectEncode, unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errEmptyBody

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("%w %s", errUnknownKey, fieldName)

		case err.Error() == "http: json body too large":
			return fmt.Errorf("%w %d bytes", errBodyToLarge, maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errToManyValues
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	return WriteJSONWithHeaders(w, status, data, nil)
}

func WriteJSONWithHeaders(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/a")
	w.WriteHeader(status)

	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func ReadIDParam(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errInvalidIDParameter
	}

	return id, nil
}

func ReadStringParam(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func ReadCSVParam(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func ReadIntParam(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
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

func GetVersion() string {
	var (
		revision string
		modified bool
	)

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	if revision == "" {
		return "unavailable"
	}

	if modified {
		return fmt.Sprintf("%s-dirty", revision)
	}

	return revision
}

func PublishCommonMetrics() {
	expvar.NewString("version").Set(GetVersion())

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))
}
