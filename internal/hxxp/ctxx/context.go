package ctxx

import (
	"context"
)

type ctxKey string

const (
	ctxKeyOrigin ctxKey = "origin"
	ctxKeyIP     ctxKey = "ip"
)

// Origin retrieves "origin" value from request context.
func Origin(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyOrigin).(string)
	return v, ok
}

// IP retrieves "IP" value from context.
func IP(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyIP).(string)
	return v, ok
}

// SetOrigin sets "origin" value to request context.
func SetOrigin(ctx context.Context, origin string) context.Context {
	return context.WithValue(ctx, ctxKeyOrigin, origin)
}

// SetIP sets "IP" value to request context.
func SetIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, ctxKeyIP, ip)
}
