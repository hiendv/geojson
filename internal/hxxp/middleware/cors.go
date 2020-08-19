package middleware

import (
	"net/http"

	"github.com/hiendv/geojson/internal/hxxp/ctxx"
	"github.com/hiendv/geojson/pkg/util"
)

// CORS is an HTTP middleware which specifies related headers.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin, ok := ctxx.Origin(r.Context())
		if !ok {
			util.HTTPAbort(w, "missing origin", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Access-Control-Allow-Origin", origin)
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}
