package middleware

import (
	"net/http"

	"lenslocked.com/context"
	"lenslocked.com/interfaces"
)

type RequireUser struct {}

// Redirect to the login page if no user is on the context
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

type User struct {
	interfaces.UserServiceInt
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn will return an http.HandlerFunc that will check to see if a user is logged in via remember token. If a
// remember token exists then add the user to the context and next(w, r), o.w., go to the next function
// call next(w, r)
// By writing the user to the context, we can add the request to the views.Render method, add a User to views.Data struct,
// show a Galleries link in the navbar when a user is already logged in
func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We want to return a dynamically created
		// func(http.ResponseWriter, *http.Request)
		// but we also need to convert it into an
		// http.HandlerFunc
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}
		user, err := mw.UserServiceInt.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}

