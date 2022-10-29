package context

import (
	"context"
	"fmt"
	"net/http"
)

const userContextKey = "authenticatedUser"

func SetAuthenticatedUser(r *http.Request, user string) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func GetAuthenticatedUser(r *http.Request) string {
	user := r.Context().Value(userContextKey)
	return fmt.Sprintf("%v", user)
}
