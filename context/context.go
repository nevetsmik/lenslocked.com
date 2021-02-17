package context

import (
	"context"

	"lenslocked.com/models"
)

// Unexported privateKey type ensures that exported WithUser cannot overwrite the context key/value pair with a key of "user"
type privateKey string

const (
	userKey privateKey = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// Retrieve a user from the context since *models.User cannot be retrieved from the context using a "user" key
func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
