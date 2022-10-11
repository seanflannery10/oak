package handler

import (
	"github.com/seanflannery10/oak/errors"
	"github.com/seanflannery10/oak/json"
	"github.com/seanflannery10/oak/version"
	"net/http"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	env := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"version": version.Get(),
		},
	}

	err := json.Encode(w, http.StatusOK, env)
	if err != nil {
		errors.ServerError(w, r, err)
	}
}
