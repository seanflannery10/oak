package context

import (
	"context"
	"fmt"
	"net/http"
)

type key string

var userKey key

func SetAuthenticatedUser(r *http.Request, user string) *http.Request {
	ctx := context.WithValue(r.Context(), userKey, user)
	return r.WithContext(ctx)
}

func GetAuthenticatedUser(r *http.Request) string {
	user := r.Context().Value(userKey)
	if user == nil {
		return ""
	}

	return fmt.Sprintf("%v", user)
}
