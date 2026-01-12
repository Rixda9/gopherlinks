package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
	"github.com/Rixda9/url-shortener/internal/database"
	"github.com/Rixda9/url-shortener/internal/repository"
	"github.com/Rixda9/url-shortener/internal/server"
)

// @title           URL Shortener API
// @version         1.0
// @description     A high-performance URL shortener using Go, PostgreSQL, and Redis.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

func main() {
	// Load Config from .env file
	if err := godotenv.Load(); err != nil {
		// Log the error but don't stop if .env is missing,
		
		log.Println("No .env file found, using system environment variables.")
	}
	
	dbURL := os.Getenv("DATABASE_URL")
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = os.Getenv("REDIS_URL")
	}
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%s", port)
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

  db := database.Connect()
  defer database.Close(db) 

	database.RunMigrations(dbURL)
  
	postgresRepo := repository.NewPostgresRepo(db)
	redisRepo := repository.NewRedisRepo(redisAddr)
	

	router := server.NewRouter(postgresRepo, redisRepo, baseURL)
	listenAddr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("Server starting on %s\n", listenAddr)
		 
	log.Fatal(http.ListenAndServe(listenAddr,router))

}

