package hxxp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/hiendv/geojson/internal/hxxp/api.v1"
	"github.com/hiendv/geojson/internal/hxxp/ctxx"
	"github.com/hiendv/geojson/internal/hxxp/middleware"
	"github.com/hiendv/geojson/internal/osm"
	"github.com/hiendv/geojson/pkg/util"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	ctx    context.Context
	router *httprouter.Router
}

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

func New(ctx context.Context) (handler *Handler, err error) {
	dir, ok := ctxOutDir(ctx)
	if !ok {
		return nil, errors.New("invalid output directory")
	}

	router := httprouter.New()
	handler = &Handler{ctx: ctx, router: router}
	v1SubAreas, err := v1.SubAreas(osm.NewContext(ctx, ctxLog(ctx), false, false, dir), handler)
	if err != nil {
		return nil, err
	}

	router.GET("/", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		util.HTTPRespond(w, []byte(`Hello`))
	})

	router.GET("/api/v1/subareas/:id", v1SubAreas.Query)
	return
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h Handler) JSON(data io.ReadCloser, result interface{}) error {
	return json.NewDecoder(data).Decode(&result)
}

func (h Handler) Abort(w http.ResponseWriter, message string, code int) {
	util.HTTPAbort(w, message, code)
}

func (h Handler) Respond(w http.ResponseWriter, message string, data interface{}) {
	util.HTTPRespondJSON(w, message, data)
}

func (h Handler) Error(w http.ResponseWriter, err error, code int) {
	h.Abort(w, err.Error(), code)
}
