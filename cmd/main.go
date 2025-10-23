package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/zen-flo/url-shortener/internal/db"
	"net/http"
)

/*
Main function starts the HTTP server for the URL Shortener API.
*/
func main() {
	// Create a new router
	r := chi.NewRouter()

	// Initialize database
	database := db.InitDB("urls.db")
	defer func() {
		if err := database.Close(); err != nil {
			fmt.Printf("Error closing database: %v\n", err)
		}
	}()

	// Test route to check if the server is running
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("URL Shortener API is running")); err != nil {
			fmt.Printf("Error writing response: %v\n", err)
		}
	})

	// Server port
	port := 8080
	fmt.Printf("Starting server on port %d...\n", port)

	// Start HTTP server
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
