package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	errBadlyFormed     = errors.New("body contains badly-formed encode")
	errIncorrectEncode = errors.New("body contains incorrect encode type")
	errEmptyBody       = errors.New("body must not be empty")
	errUnknownKey      = errors.New("body contains unknown key")
	errBodyToLarge     = errors.New("body must not be larger than")
	errToManyValues    = errors.New("body must only contain a single encode value")
)

func Read(w http.ResponseWriter, r *http.Request, dst interface{}) error {
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

func Write(w http.ResponseWriter, status int, data any) error {
	return WriteWithHeaders(w, status, data, nil)
}

func WriteWithHeaders(w http.ResponseWriter, status int, data any, headers http.Header) error {
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
