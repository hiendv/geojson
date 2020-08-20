package osm

import (
	"github.com/paulmach/osm/osmapi"
)

// ErrIsClient determines if an error thrown by osmapi is client error
func ErrIsClient(err interface{}) bool {
	switch err.(type) {
	case *osmapi.NotFoundError, *osmapi.ForbiddenError:
		return true
	default:
		return false
	}
}
