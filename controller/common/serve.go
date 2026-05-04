package common

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// ServeJSON serves any struct thought the API as JSON
func ServeJSON(w *http.ResponseWriter, v any) {
	writeHeadersJSON(w)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		serveError(w, err.Error())
	}

	(*w).WriteHeader(200)
	(*w).Write(buf.Bytes())
}

// ServeError serves an errorMessage thought the API
func ServeError(w *http.ResponseWriter, errorMessage string) {
	writeHeadersJSON(w)

	err := struct {
		Error string
	}{
		Error: errorMessage,
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(err)

	(*w).WriteHeader(500)
	(*w).Write(buf.Bytes())
}

func serveError(w *http.ResponseWriter, errorMessage string) {
	err := struct {
		Error string
	}{
		Error: errorMessage,
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(err)

	(*w).WriteHeader(500)
	(*w).Write(buf.Bytes())
}
