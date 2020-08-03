package osm

import (
	"context"

	"github.com/hiendv/geojson/internal/shared"
	"github.com/paulmach/osm"
)

type ctxKey string

const DEFAULT_OUTDIR = "./geojson"

const ctxKeyRaw ctxKey = "raw"
const ctxKeySeparated ctxKey = "separated"
const ctxKeyOut ctxKey = "out"
const ctxKeyRoot ctxKey = "root"
const ctxKeyLog ctxKey = "log"

// NewContext is the utility to encapsulate pkg-scoped context values by preventing context key collision
func NewContext(ctx context.Context, log Logger, raw bool, separated bool, out string) context.Context {
	if log != nil {
		log.Debugw("context", "raw", raw, "separated", separated, "out", out)
	}

	return context.WithValue(
		context.WithValue(
			context.WithValue(
				context.WithValue(
					ctx,
					ctxKeyOut,
					out,
				),
				ctxKeySeparated,
				separated,
			),
			ctxKeyRaw,
			raw,
		),
		ctxKeyLog,
		log,
	)
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

func ctxOutDir(ctx context.Context) string {
	out, ok := ctx.Value(ctxKeyOut).(string)
	if !ok {
		return DEFAULT_OUTDIR
	}

	return out
}

func ctxRoot(ctx context.Context) (*osm.Relation, bool) {
	v, ok := ctx.Value(ctxKeyRoot).(*osm.Relation)
	return v, ok
}

func ctxLog(ctx context.Context) Logger {
	v, ok := ctx.Value(ctxKeyLog).(Logger)
	if !ok {
		return shared.LoggerNoop
	}

	return v
}

func ctxSetRoot(ctx context.Context, root *osm.Relation) context.Context {
	return context.WithValue(ctx, ctxKeyRoot, root)
}
