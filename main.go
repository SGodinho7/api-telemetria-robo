package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	serverAddress := fmt.Sprintf("%s:%s", "0.0.0.0", "8080")
	server := http.Server{
		Addr:    serverAddress,
		Handler: mux,
	}

	server.ListenAndServe()
}
