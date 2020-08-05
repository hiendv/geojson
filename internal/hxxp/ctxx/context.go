package ctxx

import (
	"context"
)

type ctxKey string

const (
	ctxKeyOrigin ctxKey = "origin"
	ctxKeyIP     ctxKey = "ip"
)

func Origin(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyOrigin).(string)
	return v, ok
}

func IP(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyIP).(string)
	return v, ok
}

func SetOrigin(ctx context.Context, origin string) context.Context {
	return context.WithValue(ctx, ctxKeyOrigin, origin)
}

func SetIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, ctxKeyIP, ip)
}
