package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/thatnerdjosh/example-webservices/internal/core"
)

const (
	defaultHost string = "127.0.0.1"
	defaultPort string = "8080"
)

func main() {
	host := defaultHost
	if h, has := os.LookupEnv("AUTH_API_HOST"); has {
		host = h
	}

	port := defaultPort
	if p, has := os.LookupEnv("AUTH_API_PORT"); has {
		port = p
	}

	// TODO: Consider adding structured logging
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	httpServer := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: mux,
	}

	core.StartHTTPServer(&httpServer)
}
