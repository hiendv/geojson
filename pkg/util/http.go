// Package util provides utilities
package util

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HTTPResponse struct {
	HTTPError `json:",inline"`
	Data      interface{} `json:"data"`
}

func HTTPAbort(w http.ResponseWriter, message string, code int) {
	responseCode := code

	if code == http.StatusOK {
		responseCode = 0
	}

	resp, err := json.Marshal(HTTPError{responseCode, message})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	// nolint:errcheck
	w.Write(resp)
}

func HTTPRespondJSON(w http.ResponseWriter, message string, data interface{}) {
	resp, err := json.Marshal(HTTPResponse{HTTPError{0, message}, data})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	HTTPRespond(w, resp)
}

func HTTPRespond(w http.ResponseWriter, data []byte) {
	w.WriteHeader(200)

	// nolint:errcheck
	w.Write(data)
}
