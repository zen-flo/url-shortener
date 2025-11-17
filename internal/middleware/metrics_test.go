package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	dto "github.com/prometheus/client_model/go"
)

// Проверяем, что middleware учитывает запросы
func TestMetricsMiddleware(t *testing.T) {
	handler := MetricsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Result().StatusCode)
	}

	// Проверяем, что httpRequestsTotal увеличился
	m := &dto.Metric{}
	err := httpRequestsTotal.WithLabelValues("/test", "GET", strconv.Itoa(http.StatusOK)).Write(m)
	if err != nil {
		t.Fatalf("failed to get metric: %v", err)
	}

	if *m.Counter.Value != 1 {
		t.Errorf("expected http_requests_total=1, got %v", *m.Counter.Value)
	}
}

// Проверяем MetricsHandler
func TestMetricsHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/metrics", nil)
	rec := httptest.NewRecorder()

	MetricsHandler().ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Result().StatusCode)
	}

	if rec.Body.Len() == 0 {
		t.Errorf("expected non-empty metrics body")
	}
}
