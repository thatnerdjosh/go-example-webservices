package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/thatnerdjosh/example-webservices/internal/controllers"
	"github.com/thatnerdjosh/example-webservices/internal/core"
)

const (
	defaultHost string = "127.0.0.1"
	defaultPort string = "8080"
)

func main() {
	host := defaultHost
	if h, has := os.LookupEnv("TASK_API_HOST"); has {
		host = h
	}

	port := defaultPort
	if p, has := os.LookupEnv("TASK_API_PORT"); has {
		port = p
	}

	// TODO: Consider adding structured logging
	mux := http.NewServeMux()
	mux.Handle("/", controllers.NewTaskController(nil, http.DefaultClient))

	httpServer := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: mux,
	}

	core.StartHTTPServer(&httpServer)
}
