package server

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
  _ "github.com/Rixda9/url-shortener/docs"

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

	// doc route	
	r.Get("/swagger/*", httpSwagger.Handler(
        httpSwagger.URL(baseURL + "/swagger/doc.json"), // Points to the generated JSON
    ))
	// Public routes 
	
	// html for route path
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})
 	// 
	r.Get("/{shortCode}", RedirectHandler(repo.DB, repo.Cache))
	//
	r.Post("/api/shorten", ShortenHandler(repo.DB, repo.Cache, baseURL))
		

	log.Println("Router endpoints initialized.")
	return r
}

