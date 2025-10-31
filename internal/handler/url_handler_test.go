package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/zen-flo/url-shortener/internal/service"
	_ "modernc.org/sqlite" // SQLite driver
)

func setupRouter(t *testing.T) *chi.Mux {
	// Initializing the in-memory database
	db, err := sqlx.Connect("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to connect to in-memory DB: %v", err)
	}
	schema := `
	CREATE TABLE urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original TEXT NOT NULL,
		short TEXT NOT NULL UNIQUE,
		created_at DATETIME NOT NULL
	);`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	urlService := service.NewURLService(db)
	urlHandler := NewURLHandler(urlService)

	r := chi.NewRouter()
	urlHandler.RegisterRoutes(r)
	return r
}

func TestURLHandler(t *testing.T) {
	router := setupRouter(t)

	// Test CreateShortURL
	body := map[string]string{"original": "https://example.com"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/urls", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	t.Logf("Response body: %s", rec.Body.String())

	var created map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	short, ok := created["short"].(string)
	if !ok || short == "" {
		t.Fatalf("expected non-empty short code")
	}

	// Test GetOriginalURL
	req = httptest.NewRequest(http.MethodGet, "/urls/"+short, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode GET response: %v", err)
	}

	if got["original"] != "https://example.com" {
		t.Errorf("expected original URL https://example.com, got %v", got["original"])
	}

	// Test DeleteURL
	req = httptest.NewRequest(http.MethodDelete, "/urls/"+short, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	// Verify deletion
	req = httptest.NewRequest(http.MethodGet, "/urls/"+short, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 after deletion, got %d", rec.Code)
	}
}
