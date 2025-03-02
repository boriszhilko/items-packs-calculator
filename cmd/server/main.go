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

const shutdownTimeout = 5 * time.Second

// setupServer configures the HTTP server and returns it.
func setupServer(configPath string) (*http.Server, error) {
	mux := http.NewServeMux()

	handler, err := api.NewCalculateHandler(configPath)
	if err != nil {
		return nil, err
	}
	mux.HandleFunc("/calculate", handler)

	// Serve the "frontend" folder so users can access index.html at "/"
	fileServer := http.FileServer(http.Dir("frontend"))
	mux.Handle("/", fileServer)

	// Read assigned port from environment (fallback to 8080 locally)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	serverAddr := ":" + port

	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}

	return server, nil
}

func main() {
	log.Println("Starting server...")

	server, err := setupServer("configs/packs.json")
	if err != nil {
		log.Fatalf("Could not create server: %v", err)
	}

	// Listen for OS signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("HTTP server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
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
