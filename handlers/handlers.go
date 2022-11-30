package handlers

import (
	"net/http"

	"github.com/seanflannery10/ossa/httperrors"
	"github.com/seanflannery10/ossa/jsonutil"
	"github.com/seanflannery10/ossa/version"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	env := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"version": version.Get(),
		},
	}

	err := jsonutil.Write(w, http.StatusOK, env)
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}
