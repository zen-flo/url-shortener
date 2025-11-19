package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/zen-flo/url-shortener/internal/model"
	"github.com/zen-flo/url-shortener/internal/service"
)

/*
URLHandler provides HTTP endpoints for managing shortened URLs.
*/
type URLHandler struct {
	Service service.URLServiceInterface
}

/*
NewURLHandler creates a new instance of URLHandler.
*/
func NewURLHandler(s service.URLServiceInterface) *URLHandler {
	return &URLHandler{Service: s}
}

/*
RegisterRoutes registers all URL-related routes to the given router.
*/
func (h *URLHandler) RegisterRoutes(r chi.Router) {
	// URL routes
	r.Post("/urls", h.CreateShortURL)
	r.Get("/urls/{short}", h.GetOriginalURL)
	r.Delete("/urls/{short}", h.DeleteURL)
}

/*
CreateShortURL handles POST /urls requests and creates a new shortened URL.
Expected JSON body: {"original": "https://example.com"}
*/
// CreateShortURL handles POST /urls requests and creates a new shortened URL.
// @Summary Create a shortened URL
// @Description Generate a short link from the original URL
// @Tags URLs
// @Accept json
// @Produce json
// @Param url body map[string]string true "Original URL" example({"original": "https://example.com"})
// @Success 201 {object} model.URL "Successfully created"
// @Failure 400 {string} string "invalid JSON"
// @Failure 500 {string} string "internal server error"
// @Router /urls [post]
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
// GetOriginalURL handles GET /urls/{short} requests.
// @Summary Get original URL
// @Description Retrieve the original URL by short code
// @Tags URLs
// @Produce json
// @Param short path string true "Short code" example("abc123")
// @Success 200 {object} model.URL "Original URL retrieved successfully"
// @Failure 404 {string} string "URL not found"
// @Router /urls/{short} [get]
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
// DeleteURL handles DELETE /urls/{short} requests.
// @Summary Delete a shortened URL
// @Description Remove a short URL by its code
// @Tags URLs
// @Param short path string true "Short code" example("abc123")
// @Success 204 {string} string "No Content"
// @Failure 500 {string} string "failed to delete URL"
// @Router /urls/{short} [delete]
func (h *URLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	if err := h.Service.DeleteURL(short); err != nil {
		http.Error(w, "failed to delete URL", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
