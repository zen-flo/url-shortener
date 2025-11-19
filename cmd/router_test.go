package main

import (
	"github.com/zen-flo/url-shortener/internal/handler"
	"github.com/zen-flo/url-shortener/internal/model"
	"net/http/httptest"
	"testing"
)

// mockService заглушка для URLService
type mockService struct{}

func (m *mockService) CreateShortURL(original string) (*model.URL, error) {
	return &model.URL{ID: 1, Original: original, Short: "abc123"}, nil
}

func (m *mockService) GetOriginalURL(short string) (*model.URL, error) {
	return &model.URL{ID: 1, Original: "https://example.com", Short: short}, nil
}

func (m *mockService) DeleteURL(_ string) error {
	return nil
}

func (m *mockService) UpdateURLCount() {}

func TestRouterRoutes(t *testing.T) {
	// Создаём мок-сервис и handler
	svc := &mockService{}
	h := handler.NewURLHandler(svc)

	// Router
	r := NewRouter(h)

	tests := []struct {
		method     string
		url        string
		wantStatus int
	}{
		{"GET", "/", 200},
		{"GET", "/metrics", 200},
		{"GET", "/health", 200},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.url, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != tt.wantStatus {
			t.Errorf("expected status %d for %s, got %d", tt.wantStatus, tt.url, w.Code)
		}
	}
}
