package go_http_utils

import (
	"encoding/json"
	"net/http"
)

const (
	contentTypeJSON = "application/json"
)

func FormatError(w http.ResponseWriter, status int, err error) {
	body, _ := json.Marshal(map[string]string{
		"status": http.StatusText(status),
		"reason": err.Error(),
	})
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	if _, werr := w.Write(body); werr != nil {
		http.Error(w, werr.Error(), http.StatusInternalServerError)
	}
}
