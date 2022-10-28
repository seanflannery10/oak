package context

import (
	"context"
	"net/http"
)

type User struct {
	ID string
}

type contextKey string

const authenticatedUserContextKey = contextKey("authenticatedUser")

func SetAuthenticatedUser(r *http.Request, user User) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func GetAuthenticatedUser(r *http.Request) *User {
	user, ok := r.Context().Value(authenticatedUserContextKey).(*User)
	if !ok {
		return nil
	}

	return user
}
