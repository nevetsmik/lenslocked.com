package middleware

import (
	"fmt"
	"lenslocked.com/context"
	"net/http"

	"lenslocked.com/interfaces"
)

type RequireUser struct {
	interfaces.UserServiceInt
}

// ApplyFn will return an http.HandlerFunc that will
// check to see if a user is logged in and then either
// call next(w, r) if they are, or redirect them to the
// login page if they are not.
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	// We want to return a dynamically created
	// func(http.ResponseWriter, *http.Request)
	// but we also need to convert it into an
	// http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user, err := mw.UserServiceInt.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		fmt.Println("User found: ", user)
		// Grab the context
		ctx := r.Context()
		// Create a new context with the user from the existing context
		ctx = context.WithUser(ctx, user)
		// Create a new request with the new context from the exising request
		r = r.WithContext(ctx)
		next(w, r)
	})
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
