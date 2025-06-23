package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Global database connection pool
var appDB *sql.DB

func InitDatabase() error {
	connStr := os.Getenv("DBcon")
	if connStr == "" {
		log.Println("DBcon environment variable not set. Please set it or hardcode for testing.")
		return fmt.Errorf("DBcon environment variable not set")
	}

	var err error
	appDB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}

	err = appDB.Ping()
	if err != nil {
		appDB.Close()
		return fmt.Errorf("error connecting to the database: %w", err)
	}

	appDB.SetMaxOpenConns(25)
	appDB.SetMaxIdleConns(25)
	appDB.SetConnMaxLifetime(5 * time.Minute)
	appDB.SetConnMaxIdleTime(1 * time.Minute)

	log.Println("Successfully connected to PostgreSQL!")
	return nil
}

func main() {
	log.Println("Starting backend...")
	log.Println("Initializing database connection pool...")
	if err := InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Ensure the database connection is closed when main exits
	defer func() {
		if appDB != nil {
			err := appDB.Close()
			if err != nil {
				log.Printf("Error closing database connection: %v\n", err)
			}
			log.Println("Database connection pool closed.")
		}
	}()

	ticker := time.NewTicker(1 * time.Hour)
	var wg sync.WaitGroup

	log.Println("Starting Ingestion Routine every hour...")
	go func() {
		for range ticker.C {
			log.Println("Ingesting data at", time.Now().Format(time.RFC3339))
			wg.Add(1)
			go RunIngest()
		}
	}()

	wg.Add(1)
	log.Println("Starting API Routine...")
	go RunAPI()

	wg.Wait()
	log.Println("All routines completed.")
}
