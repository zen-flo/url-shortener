package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zen-flo/url-shortener/internal/service"
)

/*
URLHandler provides HTTP endpoints for managing shortened URLs.
*/
type URLHandler struct {
	Service *service.URLService
}

/*
NewURLHandler creates a new instance of URLHandler.
*/
func NewURLHandler(s *service.URLService) *URLHandler {
	return &URLHandler{Service: s}
}

/*
RegisterRoutes registers all URL-related routes to the given router.
*/
func (h *URLHandler) RegisterRoutes(r chi.Router) {
	r.Post("/urls", h.CreateShortURL)
	r.Get("/urls/{short}", h.GetOriginalURL)
	r.Delete("/urls/{short}", h.DeleteURL)
}

/*
CreateShortURL handles POST /urls requests and creates a new shortened URL.
Expected JSON body: {"original": "https://example.com"}
*/
func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Original string `json:"original"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	url, err := h.Service.CreateShortURL(req.Original)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(url); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
GetOriginalURL handles GET /urls/{short} requests and returns the original URL.
*/
func (h *URLHandler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	url, err := h.Service.GetOriginalURL(short)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(url); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
DeleteURL handles DELETE /urls/{short} requests and deletes a shortened URL.
*/
func (h *URLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	if err := h.Service.DeleteURL(short); err != nil {
		http.Error(w, "failed to delete URL", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
