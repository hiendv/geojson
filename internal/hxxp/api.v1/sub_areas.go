package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/hashicorp/golang-lru"
	"github.com/hiendv/geojson/internal/osm"
	"github.com/hiendv/geojson/internal/shared"
	"github.com/julienschmidt/httprouter"
)

type subAreasGroup struct {
	handler    Handler
	logger     shared.Logger
	osmContext context.Context
	cache      Cache
	processing map[int64]bool
	mu         sync.RWMutex
}

// Cache is the contract of cache.
type Cache interface {
	Add(key, value interface{})
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
}

// SubAreas constructs the routing group itself.
func SubAreas(logger shared.Logger, osmContext context.Context, handler Handler) (*subAreasGroup, error) {
	if handler == nil {
		return nil, errors.New("invalid HTTP handler")
	}

	cache, err := lru.New2Q(5000)
	if err != nil {
		return nil, errors.New("invalid cache")
	}

	return &subAreasGroup{handler: handler, logger: logger, osmContext: osmContext, cache: cache, processing: map[int64]bool{}}, nil
}

func (group *subAreasGroup) Query(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	osmContext, err := osm.CtxBareClone(group.osmContext)
	if err != nil {
		group.handler.Abort(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		group.handler.Error(w, errors.New("invalid ID"), http.StatusUnprocessableEntity)
		return
	}

	_, rewind := r.URL.Query()["rewind"]
	if rewind {
		osmContext = osm.CtxSetRewind(osmContext, true)
	}

	cacheKey := fmt.Sprintf("%d-%v", id, rewind)
	v, ok := group.cache.Get(cacheKey)
	if ok {
		path, ok := v.(string)
		if !ok {
			group.handler.Abort(w, "invalid path", http.StatusInternalServerError)
			return
		}

		err := osm.VerifyOutput(osmContext, path)
		if err != nil {
			group.cache.Remove(cacheKey)
			group.handler.Abort(w, "missing outputs. try again.", http.StatusInternalServerError)
			return
		}

		group.handler.Respond(w, "", group.handler.Static(path))
		return
	}

	group.mu.RLock()
	working := group.processing[id]
	group.mu.RUnlock()

	if working {
		group.handler.Respond(w, "check back later", nil)
		return
	}

	path, err := osm.FindSubAreas(osmContext, id)
	if err == nil {
		group.cache.Add(cacheKey, path)
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

		if rewind {
			err := osm.SubAreas(osmContext, params.ByName("id"))
			if err != nil {
				group.logger.Error(err)
			}
			return
		}

		err := osm.SubAreas(osmContext, params.ByName("id"))
		if err != nil {
			group.logger.Error(err)
		}
	}(group, id)

	group.handler.Respond(w, "check back later", nil)
}
