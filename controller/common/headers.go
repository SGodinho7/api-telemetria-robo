package common

import "net/http"

// writeHeadersJSON writes all the necessary headers for a JSON response
func writeHeadersJSON(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
}
