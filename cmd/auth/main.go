package main

import (
	"net/http"
	"os"

	"github.com/thatnerdjosh/example-webservices/internal/core"
)

const (
	defaultPort string = "8080"
)

func main() {
	port := defaultPort
	if p, has := os.LookupEnv("AUTH_API_PORT"); has {
		port = p
	}

	// TODO: Consider adding structured logging
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	//wrappedMux := middleware.WrapLogger(mux)

	httpServer := http.Server{
		Addr:    ":" + port,
		Handler: mux,
		//Handler: wrappedMux,
	}

	core.StartHTTPServer(&httpServer)
}
