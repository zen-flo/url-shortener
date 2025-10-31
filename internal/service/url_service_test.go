package service

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // SQLite driver
)

func setupTestDB(t *testing.T) *sqlx.DB {
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
	return db
}

func TestCreateGetDeleteURL(t *testing.T) {
	db := setupTestDB(t)
	service := NewURLService(db)

	original := "https://example.com"

	// Test CreateShortURL
	url, err := service.CreateShortURL(original)
	if err != nil {
		t.Fatalf("CreateShortURL failed: %v", err)
	}

	if url.Original != original {
		t.Errorf("expected original %s, got %s", original, url.Original)
	}
	if url.Short == "" {
		t.Errorf("expected non-empty short code")
	}
	if time.Since(url.CreatedAt) > time.Minute {
		t.Errorf("CreatedAt seems incorrect")
	}

	// Test GetOriginalURL
	got, err := service.GetOriginalURL(url.Short)
	if err != nil {
		t.Fatalf("GetOriginalURL failed: %v", err)
	}
	if got.Original != original {
		t.Errorf("expected original %s, got %s", original, got.Original)
	}

	// Test DeleteURL
	err = service.DeleteURL(url.Short)
	if err != nil {
		t.Fatalf("DeleteURL failed: %v", err)
	}

	// Verify deletion
	_, err = service.GetOriginalURL(url.Short)
	if err == nil {
		t.Errorf("expected error after deletion, got nil")
	}
}
