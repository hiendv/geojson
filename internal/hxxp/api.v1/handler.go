package v1

import (
	"net/http"
)

type Handler interface {
	Respond(w http.ResponseWriter, message string, data interface{})
	Abort(w http.ResponseWriter, message string, code int)
	Error(w http.ResponseWriter, err error, code int)
	Static(path string) string
}
