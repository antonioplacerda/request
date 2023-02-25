package context

import (
	"context"
)

type key struct{}

type userContext struct {
	UserID string
}

// SetUserID returns a new context with the specified user id value set.
func SetUserID(ctx context.Context, id string) context.Context {
	u, _ := ctx.Value(key{}).(userContext)

	u.UserID = id
	return context.WithValue(ctx, key{}, u)
}

// UserID returns the authenticated user id.
func UserID(ctx context.Context) string {
	u, _ := ctx.Value(key{}).(userContext)
	return u.UserID
}
