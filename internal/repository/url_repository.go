package repository

import (
	"database/sql"
	"fmt"
	"log"
	//"time"
	"errors"
	"crypto/rand"
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


func (r *PostgresRepo) SaveURL(originalURL string) (string, error) {
	existingSlug, err := r.GetByOriginalURL(originalURL)
	if err != nil {
		log.Printf("Error checking for existing URL: %v", err)
	}
	if existingSlug != "" {
		log.Printf("URL already exists, returning existing slug: %s", existingSlug)
		return existingSlug, nil
	}

	slug := generateSlug()

	query := `INSERT INTO urls (short_code, original_url, created_at) VALUES ($1, $2, NOW()) RETURNING short_code`

	err = r.DB.QueryRow(query, slug, originalURL).Scan(&slug)

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

func (r *PostgresRepo) GetByOriginalURL(originalURL string) (string, error) {
	var shortCode string
	query := `SELECT short_code FROM urls WHERE original_url = $1 LIMIT 1`

	err := r.DB.QueryRow(query, originalURL).Scan(&shortCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return shortCode, nil
}

func generateSlug() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error generating random bytes: %v", err)
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}


