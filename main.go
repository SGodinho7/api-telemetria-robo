package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	mux := mux.NewRouter()

	serverAddress := fmt.Sprintf("%s:%s", "127.0.0.1", "5000")
	server := http.Server{
		Addr:    serverAddress,
		Handler: mux,
	}

	server.ListenAndServe()
}
