package repository

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
	"errors"
)

type Repository interface {
	SaveURL(originalURL string) (string, error)
	RetrieveURL(slug string) (string, error)
}

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{DB: db}
}

func generateSlug() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (r *PostgresRepo) SaveURL(originalURL string) (string, error) {
	slug := generateSlug()
	
	query := `INSERT INTO urls (short_code, original_url, created_at) VALUES ($1, $2, NOW()) RETURNING short_code`
	
	err := r.DB.QueryRow(query, slug, originalURL).Scan(&slug)
	
	if err != nil {
		log.Printf("Error saving URL %s with slug %s: %v", originalURL, slug, err)
		return "", fmt.Errorf("could not save URL: %w", err)
	}
	
	return slug, nil
}

func (r *PostgresRepo) RetrieveURL(slug string) (string, error) {
	var originalURL string 
	query := `SELECT original_url FROM urls WHERE short_code = $1;`

	err := r.DB.QueryRow(query, slug).Scan(&originalURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("Short code not found")
		}
		log.Printf("Error retrieving URL for slug: %s: %v", slug, err)
		return "", fmt.Errorf("Could not retrieve URL: %w", err)
	}
	
	return originalURL, nil 
}

