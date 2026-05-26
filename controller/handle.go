package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func serveSuccess(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(`{"status": "success"}`); err != nil {
		serveError(w, err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// ServeJSON serves any struct thought the API as JSON
func serveJSON(w http.ResponseWriter, v any) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		serveError(w, err.Error())
	}

	w.WriteHeader(200)
	w.Write(buf.Bytes())
}

// ServeError serves an errorMessage thought the API
func serveError(w http.ResponseWriter, errorMessage string) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	err := struct {
		Error string
	}{
		Error: errorMessage,
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(err)

	w.WriteHeader(500)
	w.Write(buf.Bytes())
}
