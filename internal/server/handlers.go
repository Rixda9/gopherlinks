package server 

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-chi/chi/v5"
	"github.com/Rixda9/url-shortener/internal/repository" 
)
var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ShortenRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}
// Shorten the url
func ShortenHandler(dbRepo repository.Repository, cacheRepo repository.CacheRepository, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ShortenRequest
		
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(req); err != nil {
			http.Error(w, fmt.Sprintf("Validation failed: %s", err.Error()), http.StatusBadRequest)
			return
		}

		// save to db
		slug, err := dbRepo.SaveURL(req.URL)
		if err != nil {
			http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
			return
		}
		// save to cache for 24h
		cacheRepo.Set(slug, req.URL, 24 * time.Hour)
		shortURL := baseURL + "/" + slug

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		
		resp := ShortenResponse{ShortURL: shortURL}
		json.NewEncoder(w).Encode(resp)
	}
}

// Redirect short codes to actual address 
func RedirectHandler(dbRepo repository.Repository, cacheRepo repository.CacheRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "shortCode")
		if slug == "" {
			http.Error(w, "Short code not provided", http.StatusBadRequest)
			return
		}

		originalURL, err := cacheRepo.Get(slug)
		if err != nil {
			log.Printf("Cache read error for slug: %s: %v", slug, err)
		}
		// cache hit
		if originalURL != "" {
			log.Printf("Cache HIT for slug: %s", slug)
			http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
			return
		}
		
		// cache miss 
		log.Printf("Cache MISS for slug: %s. Hitting database.", slug)
		originalURL, err = dbRepo.RetrieveURL(slug)

		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "Short code not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Server error retrieving URL", http.StatusInternalServerError)
			return
		}
		// cache for 24 hours
		cacheRepo.Set(slug, originalURL, 24 * time.Hour)
 		
		http.Redirect(w, r, originalURL, http.StatusMovedPermanently) 
	}
}

