package middleware

import (
	"net"
	"net/http"

	"github.com/hiendv/geojson/internal/hxxp/ctxx"
)

// IP is an HTTP middleware which keeps track of client IP addresses.
func IP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement IP resolution behind trusted proxies. E.g. X-Real-IP or X-Forwarded-For
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctxx.SetIP(r.Context(), ip)))
	})
}
