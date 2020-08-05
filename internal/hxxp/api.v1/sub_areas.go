package v1

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/hashicorp/golang-lru"
	"github.com/hiendv/geojson/internal/osm"
	"github.com/julienschmidt/httprouter"
)

type subAreasGroup struct {
	handler    Handler
	osmContext context.Context
	cache      Cache
	processing map[int64]bool
	mu         sync.RWMutex
}

type Cache interface {
	Add(key, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
}

func SubAreas(ctx context.Context, handler Handler) (*subAreasGroup, error) {
	if handler == nil {
		return nil, errors.New("invalid HTTP handler")
	}

	cache, err := lru.New2Q(5000)
	if err != nil {
		return nil, errors.New("invalid cache")
	}

	return &subAreasGroup{handler: handler, osmContext: ctx, cache: cache, processing: map[int64]bool{}}, nil
}

func (group *subAreasGroup) Query(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		group.handler.Error(w, errors.New("invalid ID"), http.StatusUnprocessableEntity)
		return
	}

	v, ok := group.cache.Get(id)
	if ok {
		path, ok := v.(string)
		if !ok {
			group.handler.Abort(w, "invalid path", http.StatusInternalServerError)
			return
		}

		group.handler.Respond(w, "", group.handler.Static(path))
		return
	}

	group.mu.RLock()
	working := group.processing[id]
	group.mu.RUnlock()

	if working {
		group.handler.Respond(w, "Please check back later", nil)
		return
	}

	path, err := osm.FindSubAreas(group.osmContext, id)
	if err == nil {
		group.cache.Add(id, path)
		group.handler.Respond(w, "", group.handler.Static(path))
		return
	}

	group.mu.Lock()
	group.processing[id] = true
	group.mu.Unlock()

	go func(group *subAreasGroup, id int64) {
		defer func() {
			group.mu.Lock()
			group.processing[id] = false
			group.mu.Unlock()
		}()

		// nolint:errcheck
		osm.SubAreas(group.osmContext, params.ByName("id"))
	}(group, id)

	group.handler.Respond(w, "Please check back later", nil)
}
