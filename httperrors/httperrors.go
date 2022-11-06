package httperrors

import (
	"fmt"
	"github.com/seanflannery10/ossa/json"
	"github.com/seanflannery10/ossa/logger"
	"github.com/seanflannery10/ossa/validator"
	"net/http"
)

func ErrorMessage(w http.ResponseWriter, r *http.Request, status int, message string) {
	ErrorMessageWithHeaders(w, r, status, message, nil)
}

func ErrorMessageWithHeaders(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	err := json.EncodeWithHeaders(w, status, map[string]string{"error": message}, headers)
	if err != nil {
		logger.Error(err, map[string]string{
			"request_method": r.Method,
			"request_url":    r.URL.String(),
		})
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	logger.Error(err, nil)

	message := "the server encountered a problem and could not process your json"
	ErrorMessage(w, r, http.StatusInternalServerError, message)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	ErrorMessage(w, r, http.StatusNotFound, message)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	ErrorMessage(w, r, http.StatusMethodNotAllowed, message)
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	ErrorMessage(w, r, http.StatusBadRequest, err.Error())
}

func FailedValidation(w http.ResponseWriter, r *http.Request, v *validator.Validator) {
	err := json.Encode(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		ServerError(w, r, err)
	}
}

func InvalidAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", "Bearer")

	ErrorMessageWithHeaders(w, r, http.StatusUnauthorized, "invalid or missing authentication token", headers)
}

func AuthenticationRequired(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	ErrorMessage(w, r, http.StatusUnauthorized, message)
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	ErrorMessage(w, r, http.StatusTooManyRequests, message)
}
