package handlers

import (
	"github.com/seanflannery10/ossa/httperrors"
	"github.com/seanflannery10/ossa/json"
	"github.com/seanflannery10/ossa/version"
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
		httperrors.ServerError(w, r, err)
	}
}
