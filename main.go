package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	serverAddress := fmt.Sprintf("%s:%s", "127.0.0.1", "5000")
	server := http.Server{
		Addr:    serverAddress,
		Handler: mux,
	}

	server.ListenAndServe()
}
