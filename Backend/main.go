package main // Declares the package this file belongs to. 'main' is special for executables.

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
	"net/http"
	"sync"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	_ "github.com/joho/godotenv/autoload"
)

type BillName struct {
	EN string `json:"en"`
	FR string `json:"fr"`
}

type Bill struct {
	Session      string   `json:"session"`
	Legisinfo_id int      `json:"legisinfo_id"`
	Introduced   string   `json:"introduced"`
	Name         BillName `json:"name"`
	Number       string   `json:"number"`
	Url          string   `json:"url"`
}

type BillResponse struct {
	Objects []Bill `json:"objects"`
}

type MPResponse struct {
	Objects []MP `json:"objects"`
}

type MP struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	CurrentParty Party `json:"current_party"`
	CurrentRiding Riding `json:"current_riding"`
	Image        string   `json:"image"`
}

type Riding struct {
	Province string `json:"province"`
	Name     struct {
		EN string `json:"en"`
	} `json:"name"`
}

type Party struct {
	ShortName struct {
		EN string `json:"en"`
	} `json:"short_name"`
}

func getBills(database *sql.DB) error {	
	resp, err := http.Get("https://api.openparliament.ca/bills/?format=json")
	if err != nil {
		return fmt.Errorf("failed to make HTTP request for Bills: %w", err)
		print(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var billResponse BillResponse

	err = json.Unmarshal(body, &billResponse)
	if err != nil {
		return fmt.Errorf("failed to read response body for bills: %w", err)
	}
	if len(billResponse.Objects) == 0 {
		fmt.Println("No bills found in response.")
		return nil
	}
	
	fmt.Printf("successfully got %d bills\n", len(billResponse.Objects))
	for _, bill := range billResponse.Objects {
		_, err = database.Exec("INSERT INTO bills (Session, Legisinfo_id, Introduced, Name, Number, Url) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (Legisinfo_id) DO UPDATE SET legisinfo_id = EXCLUDED.legisinfo_id, session = EXCLUDED.session, introduced = EXCLUDED.introduced, name = EXCLUDED.name, number = EXCLUDED.number, url = EXCLUDED.url ", bill.Session, bill.Legisinfo_id, bill.Introduced, bill.Name.EN, bill.Number, bill.Url)
		if err != nil {
			fmt.Printf("Error inserting bill into database: %v\n", err)
			fmt.Printf("Bill: %v\n", bill)
		}
	}
	return nil
}

func getMPs(database *sql.DB) error {
	resp, err := http.Get("https://api.openparliament.ca/politicians/?format=json")
	if err != nil {
		return fmt.Errorf("failed to make HTTP request for MPs: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var mps MPResponse

	err = json.Unmarshal(body, &mps)
	if err != nil {
		return fmt.Errorf("failed to read response body for MPs: %w", err)

	}
	if len(mps.Objects) == 0 {
		fmt.Println("No MPs found in response.")
		return nil
	}


	fmt.Printf("successfully got %d MPs\n", len(mps.Objects))
	for _, mp := range mps.Objects {
		_, err = database.Exec("INSERT INTO mps (Name, CurrentParty, CurrentRiding, Url, Image) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (Name) DO UPDATE SET Name = EXCLUDED.Name, CurrentParty = EXCLUDED.CurrentParty, CurrentRiding = EXCLUDED.CurrentRiding, Url = EXCLUDED.Url, image = EXCLUDED.Image", mp.Name, mp.CurrentParty.ShortName.EN, mp.CurrentRiding.Name.EN, mp.URL, mp.Image)
		if err != nil {
			fmt.Printf("Error inserting MP into database: %v\n", err)
			fmt.Printf("MP: %v\n", mp)
		}
	}
	return nil
}



func getData(wg *sync.WaitGroup, database *sql.DB) {
	defer wg.Done()

	var internalWg sync.WaitGroup

	internalWg.Add(1)
	go func() {
		defer internalWg.Done()
		if err := getBills(database); err != nil { // Pass 'database' here
			log.Printf("Error getting bills: %v\n", err)
		}
	}()

	internalWg.Add(1)
	go func() {
		defer internalWg.Done()
		if err := getMPs(database); err != nil { // Pass 'database' here
			log.Printf("Error getting MPs: %v\n", err)
		}
	}()

	internalWg.Wait()

}


func initDatabase() error {
	connStr := os.Getenv("DBcon")
	if connStr == "" {
		log.Println("DBcon environment variable not set. Please set it or hardcode for testing.")
		return fmt.Errorf("DBcon environment variable not set")
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return fmt.Errorf("error connecting to the database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	log.Println("Successfully connected to PostgreSQL!")
	return nil
}



func main() { 
	// Initialize database connection pool once
	if err := initDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Ensure the database connection is closed when main exits
	defer func() {
		if db != nil {
			err := db.Close()
			if err != nil {
				log.Printf("Error closing database connection: %v\n", err)
			}
			log.Println("Database connection pool closed.")
		}
	}()

	log.Println("Starting data pulling routines...")

	var wg sync.WaitGroup

	wg.Add(1)
	go getData(&wg, db) // Pass the global 'db' instance here

	wg.Wait()
	log.Println("All data pulling routines completed.")
}