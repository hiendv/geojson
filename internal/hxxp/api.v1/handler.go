package v1

import (
	"net/http"
)

type Handler interface {
	Respond(w http.ResponseWriter, message string, data interface{})
	Error(w http.ResponseWriter, err error, code int)
}
