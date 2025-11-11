package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api: ", log.LstdFlags)

	sm := http.NewServeMux()

	// Create a new server
	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// Start the server
	go func() {
		l.Println("Starting server on :8080")

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
	}()

	// Track interrupt and gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive a signal
	sig := <-c
	l.Printf("Got signal %v, initiating shutdown\n", sig)

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		l.Printf("Graceful shutdown failed: %v\n", err)
		os.Exit(1)
	}

	l.Println("Server stopped")
}
