package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5" 
	"github.com/Rixda9/url-shortener/internal/repository"
)

func NewRouter(db *sql.DB, baseURL string) http.Handler {
	r := chi.NewRouter()

	repo := repository.NewPostgresRepo(db)    
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL Shortener API is running!"))
	})
  
	r.Get("/{shortCode}", RedirectHandler(repo))
	r.Post("/api/shorten", ShortenHandler(repo, baseURL))
		

	log.Println("Router endpoints initialized.")
	return r
}

