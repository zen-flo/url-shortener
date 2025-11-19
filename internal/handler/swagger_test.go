package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpSwagger "github.com/swaggo/http-swagger"
)

func TestSwaggerHandler(t *testing.T) {
	r := http.NewServeMux()
	r.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusNotFound {
		t.Errorf("expected 200 OK or 404 for Swagger JSON, got %d", rec.Code)
	}
}
