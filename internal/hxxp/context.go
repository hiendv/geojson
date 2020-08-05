package hxxp

import (
	"context"
	"time"

	"github.com/hiendv/geojson/internal/shared"
)

type ctxKey string

const (
	ctxKeyLog       ctxKey = "log"
	ctxKeyAddress   ctxKey = "address"
	ctxKeyOrigin    ctxKey = "origin"
	ctxKeyRate      ctxKey = "rate"
	ctxKeyRateBurst ctxKey = "rate-burst"
	ctxKeyRateTTL   ctxKey = "rate-ttl"
	ctxKeyOut       ctxKey = "out"
	ctxKeyPrefix    ctxKey = "prefix"
)

// NewContext is the utility to encapsulate pkg-scoped context values by preventing context key collision
func NewContext(ctx context.Context, log shared.Logger, address string, origin string, rate float64, burst int, ttl time.Duration, out string, prefix string) (context.Context, error) {
	ctxx := map[ctxKey]interface{}{
		ctxKeyLog:       log,
		ctxKeyAddress:   address,
		ctxKeyOrigin:    origin,
		ctxKeyRate:      rate,
		ctxKeyRateBurst: burst,
		ctxKeyRateTTL:   ttl,
		ctxKeyOut:       out,
		ctxKeyPrefix:    prefix,
	}

	if log != nil {
		log.Debugw("context", "values", ctxx)
	}

	for k, v := range ctxx {
		ctx = context.WithValue(ctx, k, v)
	}

	return ctx, nil
}

func ctxAddress(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyAddress).(string)
	return v, ok
}

func ctxLog(ctx context.Context) shared.Logger {
	v, ok := ctx.Value(ctxKeyLog).(shared.Logger)
	if !ok {
		return shared.LoggerNoop
	}

	return v
}

func ctxOrigin(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyOrigin).(string)
	return v, ok
}

func ctxRate(ctx context.Context) (float64, bool) {
	v, ok := ctx.Value(ctxKeyRate).(float64)
	return v, ok
}

func ctxRateBurst(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(ctxKeyRateBurst).(int)
	return v, ok
}

func ctxRateTTL(ctx context.Context) (time.Duration, bool) {
	v, ok := ctx.Value(ctxKeyRateTTL).(time.Duration)
	return v, ok
}

func ctxOutDir(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyOut).(string)
	return v, ok
}

func ctxPrefix(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyPrefix).(string)
	return v, ok
}
