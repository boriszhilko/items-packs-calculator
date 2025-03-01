package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"items-packs-calculator/api"
)

const (
	shutdownTimeout = 5 * time.Second
	serverAddr      = ":8080"
)

func main() {
	mux := http.NewServeMux()

	handler, err := api.NewCalculateHandler("configs/packs.json")
	if err != nil {
		log.Fatalf("Could not create handler: %v", err)
	}
	mux.HandleFunc("/calculate", handler)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}

	// Listen for OS signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("HTTP server listening on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-done
	log.Println("Initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shut down: %v", err)
	}

	log.Println("Server exited properly")
}
