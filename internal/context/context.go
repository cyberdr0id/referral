package context

import "context"

type key int

var id key

// Set sets the value in application context.
func Set(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, id, userID)
}

// Get gets the value by key from application context.
func GetUserID(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(id).(string)
	return val, ok
}
