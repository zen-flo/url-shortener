package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/zen-flo/url-shortener/internal/handler"
	"github.com/zen-flo/url-shortener/internal/service"
)

func setupTestRouter(t *testing.T) http.Handler {
	db, err := sqlx.Connect("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to initialize the database: %v", err)
	}

	schema := `
	CREATE TABLE urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original TEXT NOT NULL,
		short TEXT NOT NULL UNIQUE,
		created_at DATETIME NOT NULL
	);`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("couldn't create a table: %v", err)
	}

	svc := service.NewURLService(db)
	urlHandler := handler.NewURLHandler(svc)

	return NewRouter(urlHandler)
}

func TestServerRoutes(t *testing.T) {
	r := setupTestRouter(t)

	tests := []struct {
		path       string
		wantStatus int
		wantBody   string
	}{
		{"/", http.StatusOK, "URL Shortener API is running"},
		{"/health", http.StatusOK, "OK"},
		{"/metrics", http.StatusOK, ""},                              // проверим только код ответа
		{"/swagger/doc.json", http.StatusOK, "\"swagger\": \"2.0\""}, // Swagger HTML
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != tt.wantStatus {
			t.Errorf("expected %d for %s, got %d", tt.wantStatus, tt.path, rec.Code)
		}

		if tt.wantBody != "" && !strings.Contains(rec.Body.String(), tt.wantBody) {
			t.Errorf("expected response body with %q for %s, got %q", tt.wantBody, tt.path, rec.Body.String())
		}
	}
}

func TestMetricsContent(t *testing.T) {
	r := setupTestRouter(t)

	// Сделаем несколько запросов, чтобы метрики появились
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	req = httptest.NewRequest("GET", "/health", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Проверим /metrics
	req = httptest.NewRequest("GET", "/metrics", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "http_requests_total") || !strings.Contains(body, "http_request_duration_seconds") {
		t.Errorf("Prometheus metrics not found in /metrics")
	}
}
