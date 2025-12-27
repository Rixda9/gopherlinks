package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"  
	
)

func RunMigrations(dbURL string) {
	sourceURL := "file://migrations"

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		log.Fatalf("FATAL: Couldnt not create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("FATAL: Failed to run migration: %v", err)
	}
	
	if err == migrate.ErrNoChange {
		fmt.Println("Database schema is up to date")
	} else {
		fmt.Println("Database migrations applied successfully")
	}
}

func Connect() *sql.DB {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("FATAL: DATABSE_URL not set in .env file!")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("FATAL: Could not connect to the database: %v", err)
	}

	db.SetMaxOpenConns(25) // max total connections

	db.SetMaxIdleConns(25) 
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Fatalf("FATAL: Database connection failed during ping: %v", err)
	}

	fmt.Println("Database connected successfully")
	return db
}

func Close(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("ERROR: Failed to close database: %v", err)
		}
	}
}
