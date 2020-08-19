package osm

import (
	"context"
	"errors"

	"github.com/hiendv/geojson/internal/shared"
	"github.com/paulmach/osm"
)

type ctxKey string

const (
	ctxKeyRaw       ctxKey = "raw"
	ctxKeySeparated ctxKey = "separated"
	ctxKeyOut       ctxKey = "out"
	ctxKeyRewind    ctxKey = "rewind"
	ctxKeyRoot      ctxKey = "root"
	ctxKeyLog       ctxKey = "log"
)

// NewContext is the utility to encapsulate pkg-scoped context values by preventing context key collision
func NewContext(ctx context.Context, log shared.Logger, raw bool, separated bool, out string, rewind bool) (context.Context, error) {
	ctxx := map[ctxKey]interface{}{
		ctxKeyLog:       log,
		ctxKeyRaw:       raw,
		ctxKeySeparated: separated,
		ctxKeyOut:       out,
		ctxKeyRewind:    rewind,
	}

	if log != nil {
		log.Debugw("context", "values", ctxx)
	}

	err := validateOut(out)
	if err != nil {
		return ctx, err
	}

	for k, v := range ctxx {
		ctx = context.WithValue(ctx, k, v)
	}

	return ctx, nil
}

func ctxShouldNormalize(ctx context.Context) bool {
	raw, ok := ctx.Value(ctxKeyRaw).(bool)
	return !(ok && raw)
}

func ctxShouldPrint(ctx context.Context) bool {
	out, ok := ctx.Value(ctxKeyOut).(string)
	return ok && out == ""
}

func ctxShouldCombine(ctx context.Context) bool {
	separated, ok := ctx.Value(ctxKeySeparated).(bool)
	return !(ok && separated)
}

func ctxOutDir(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyOut).(string)
	return v, ok
}

func ctxShouldRewind(ctx context.Context) bool {
	rewind, ok := ctx.Value(ctxKeyRewind).(bool)
	return ok && rewind
}

func ctxRoot(ctx context.Context) (*osm.Relation, bool) {
	v, ok := ctx.Value(ctxKeyRoot).(*osm.Relation)
	return v, ok
}

func ctxLog(ctx context.Context) shared.Logger {
	v, ok := ctx.Value(ctxKeyLog).(shared.Logger)
	if !ok {
		return shared.LoggerNoop
	}

	return v
}

func CtxSetRewind(ctx context.Context, rewind bool) context.Context {
	return context.WithValue(ctx, ctxKeyRewind, rewind)
}

func CtxSetRoot(ctx context.Context, root *osm.Relation) context.Context {
	return context.WithValue(ctx, ctxKeyRoot, root)
}

func CtxBareClone(ctx context.Context) (context.Context, error) {
	log, ok := ctx.Value(ctxKeyLog).(shared.Logger)
	if !ok {
		return ctx, errors.New("invalid context: logger")
	}

	raw, ok := ctx.Value(ctxKeyRaw).(bool)
	if !ok {
		return ctx, errors.New("invalid context: raw")
	}

	separated, ok := ctx.Value(ctxKeySeparated).(bool)
	if !ok {
		return ctx, errors.New("invalid context: separated")
	}

	out, ok := ctx.Value(ctxKeyOut).(string)
	if !ok {
		return ctx, errors.New("invalid context: out")
	}

	rewind, ok := ctx.Value(ctxKeyRewind).(bool)
	if !ok {
		return ctx, errors.New("invalid context: rewind")
	}

	return NewContext(context.Background(), log, raw, separated, out, rewind)
}
