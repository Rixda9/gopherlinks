package server 

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/Rixda9/url-shortener/internal/repository" 
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}
// Shorten the url
func ShortenHandler(repo repository.Repository, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ShortenRequest
		
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.URL == "" {
			http.Error(w, "Invalid request body or missing URL", http.StatusBadRequest)
			return
		}

		slug, err := repo.SaveURL(req.URL)
		if err != nil {
			http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
			return
		}

		shortURL := baseURL + "/" + slug

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		
		resp := ShortenResponse{ShortURL: shortURL}
		json.NewEncoder(w).Encode(resp)
	}
}

// Redirect short codes to actual address 
func RedirectHandler(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "shortCode")
		if slug == "" {
			http.Error(w, "Short code not provided", http.StatusBadRequest)
			return
		}

		originalURL, err := repo.RetrieveURL(slug)
		
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "Short code not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Server error retrieving URL", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusMovedPermanently) 
	}
}

