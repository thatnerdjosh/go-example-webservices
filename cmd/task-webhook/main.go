package main

import (
	"net/http"
	"os"

	"github.com/thatnerdjosh/example-webservices/internal/controllers"
	"github.com/thatnerdjosh/example-webservices/internal/core"
)

const (
	defaultPort string = "8080"
)

// TODO: Add helm chart
func main() {
	port := defaultPort
	if p, has := os.LookupEnv("TASK_API_PORT"); has {
		port = p
	}

	// TODO: Consider adding structured logging
	mux := http.NewServeMux()
	mux.Handle("/", controllers.NewTaskController(nil, http.DefaultClient))

	//wrappedMux := middleware.WrapLogger(mux)

	httpServer := http.Server{
		Addr:    ":" + port,
		Handler: mux,
		//Handler: wrappedMux,
	}

	core.StartHTTPServer(&httpServer)
}
