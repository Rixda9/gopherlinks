package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5" 
	"github.com/Rixda9/url-shortener/internal/repository"
)

func NewRouter(postgresRepo *repository.PostgresRepo, redisRepo *repository.RedisRepo, baseURL string) http.Handler {
	r := chi.NewRouter()
	
	type Repos struct {
		DB repository.Repository
		Cache repository.CacheRepository
	}

	repo := Repos{
		DB: postgresRepo,
		Cache: redisRepo,
	}
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL Shortener API is running!"))
	})
  
	r.Get("/{shortCode}", RedirectHandler(repo.DB, repo.Cache))
	r.Post("/api/shorten", ShortenHandler(repo.DB, repo.Cache, baseURL))
		

	log.Println("Router endpoints initialized.")
	return r
}

