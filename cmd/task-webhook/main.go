package main

import (
	"context"

	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/thatnerdjosh/example-webservices/internal/controllers"
)

const (
	defaultPort string = "8080"
)

// TODO: Add helm chart
func main() {
	port := defaultPort
	if p, has := os.LookupEnv("PORT"); has {
		port = p
	}

	// TODO: Consider adding structured logging
	mux := http.NewServeMux()
	mux.Handle("/", controllers.NewTaskController())

	//wrappedMux := middleware.WrapLogger(mux)

	httpServer := http.Server{
		Addr:    ":" + port,
		Handler: mux,
		//Handler: wrappedMux,
	}

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	log.Printf("Starting HTTP Server on port %s", port)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	log.Println("Gracefully shutting down...")
	<-idleConnectionsClosed
	log.Println("Shut down.")
}
