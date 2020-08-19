package hxxp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/hiendv/geojson/internal/hxxp/api.v1"
	"github.com/hiendv/geojson/internal/hxxp/ctxx"
	"github.com/hiendv/geojson/internal/hxxp/middleware"
	"github.com/hiendv/geojson/internal/osm"
	"github.com/hiendv/geojson/pkg/util"
	"github.com/julienschmidt/httprouter"
)

// Handler is the application HTTP handler.
type Handler struct {
	ctx    context.Context
	router *httprouter.Router
}

// Listen opens and serves an HTTP handler.
func Listen(h *Handler) error {
	if h == nil {
		return errors.New("invalid handler")
	}

	log := ctxLog(h.ctx)
	address, ok := ctxAddress(h.ctx)
	if !ok {
		return errors.New("missing serving address")
	}

	srv := &http.Server{
		Addr:         address,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 40 * time.Second,
	}

	log.Infow("serving", "address", address)
	return srv.ListenAndServe()
}

// New constructs a new Handler.
func New(ctx context.Context) (handler *Handler, err error) {
	dir, ok := ctxOutDir(ctx)
	if !ok {
		return nil, errors.New("invalid output directory")
	}

	router := httprouter.New()
	handler = &Handler{ctx: ctx, router: router}
	osmContext, err := osm.NewContext(ctx, ctxLog(ctx), false, false, dir, false)
	if err != nil {
		return nil, err
	}

	v1SubAreas, err := v1.SubAreas(osmContext, handler)
	if err != nil {
		return nil, err
	}

	prefix, ok := ctxPrefix(ctx)
	if !ok {
		prefix = "/"
	}

	router.GET("/", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		util.HTTPRespond(w, []byte(`Hello`))
	})
	router.ServeFiles(fmt.Sprintf("%s/%s/*filepath", prefix, filepath.Base(dir)), http.Dir(dir))
	router.GET("/api/v1/subareas/:id", v1SubAreas.Query)
	return
}

// ServeHTTP serves HTTP requests.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// exclude the static serving from middleware
	_, params, _ := h.router.Lookup(r.Method, r.URL.Path)
	for _, param := range params {
		if param.Key == "filepath" {
			h.router.ServeHTTP(w, r)
			return
		}
	}

	origin, ok := ctxOrigin(h.ctx)
	if !ok {
		util.HTTPAbort(w, "missing origin", http.StatusInternalServerError)
		return
	}

	rate, ok := ctxRate(h.ctx)
	if !ok {
		util.HTTPAbort(w, "missing rate-limiting configuration", http.StatusInternalServerError)
		return
	}

	rateBurst, ok := ctxRateBurst(h.ctx)
	if !ok {
		util.HTTPAbort(w, "missing rate-limiting configuration (burst)", http.StatusInternalServerError)
		return
	}

	rateTTL, ok := ctxRateTTL(h.ctx)
	if !ok {
		util.HTTPAbort(w, "missing rate-limiting configuration (TTL)", http.StatusInternalServerError)
		return
	}

	middleware.IP(
		middleware.RateLimit(
			middleware.CORS(
				h.router,
			),
			rate,
			rateBurst,
			rateTTL,
		),
	).ServeHTTP(
		w,
		r.WithContext( // request-scoped context
			ctxx.SetOrigin(r.Context(), origin),
		),
	)
}

// JSON is a helper to encode data in JSON format.
func (h Handler) JSON(data io.ReadCloser, result interface{}) error {
	return json.NewDecoder(data).Decode(&result)
}

// Abort is a helper to respond HTTP requests with errors.
func (h Handler) Abort(w http.ResponseWriter, message string, code int) {
	util.HTTPAbort(w, message, code)
}

// Respond is a helper to respond HTTP requests.
func (h Handler) Respond(w http.ResponseWriter, message string, data interface{}) {
	util.HTTPRespondJSON(w, message, data)
}

// Error is a helper to abort HTTP requests with errors.
func (h Handler) Error(w http.ResponseWriter, err error, code int) {
	h.Abort(w, err.Error(), code)
}

// Static is a helper to serve static files.
func (h Handler) Static(path string) string {
	prefix, ok := ctxPrefix(h.ctx)
	if !ok {
		return path
	}

	return filepath.Join(prefix, path)
}
