package core

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func StartHTTPServer(httpServer *http.Server) {
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

	log.Printf("Starting HTTP Server on port %s", strings.Split(httpServer.Addr, ":")[1])
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	log.Println("Gracefully shutting down...")
	<-idleConnectionsClosed
	log.Println("Shut down.")
}
