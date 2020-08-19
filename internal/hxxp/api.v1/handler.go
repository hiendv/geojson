package v1

import (
	"net/http"
)

// Handler is the contract of an HTTP handler.
type Handler interface {
	Respond(w http.ResponseWriter, message string, data interface{})
	Abort(w http.ResponseWriter, message string, code int)
	Error(w http.ResponseWriter, err error, code int)
	Static(path string) string
}
