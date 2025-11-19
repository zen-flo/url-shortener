package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	"github.com/zen-flo/url-shortener/internal/handler"
	"github.com/zen-flo/url-shortener/internal/middleware"
)

// NewRouter Router creates and configures an HTTP router.
// Accepts a UrlService — this is important for tests.
func NewRouter(urlHandler *handler.URLHandler) http.Handler {
	r := chi.NewRouter()

	// Metrics
	r.Use(middleware.MetricsMiddleware)
	r.Handle("/metrics", middleware.MetricsHandler())

	// Test route to check if the server is running
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("URL Shortener API is running")); err != nil {
			fmt.Printf("Error writing response: %v\n", err)
		}
	})

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			fmt.Printf("Error writing response: %v\n", err)
		}
	})

	// Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Routes for URL Shortener
	urlHandler.RegisterRoutes(r)

	return r
}
