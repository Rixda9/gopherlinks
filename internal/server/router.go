package server

import (
	"log"
	"net/http"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth_chi"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/Rixda9/url-shortener/docs"

	"github.com/go-chi/chi/v5"	
	"github.com/Rixda9/url-shortener/internal/repository"
)


func NewRouter(postgresRepo *repository.PostgresRepo, redisRepo *repository.RedisRepo, baseURL string) http.Handler {
	r := chi.NewRouter()

	// rate limiter
	  limiter := tollbooth.NewLimiter(10, nil)
    limiter.SetMessage("Rate limit exceeded. Try again in a minute.")
    
    r.Use(tollbooth_chi.LimitHandler(limiter))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(baseURL + "/swagger/doc.json"),
	))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	r.Post("/api/shorten", ShortenHandler(postgresRepo, redisRepo, baseURL))

	r.Get("/{shortCode}", RedirectHandler(postgresRepo, redisRepo))


	log.Println("Router endpoints initialized.")
	return r
}
