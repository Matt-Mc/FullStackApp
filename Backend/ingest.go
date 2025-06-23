package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
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
	Name          string `json:"name"`
	URL           string `json:"url"`
	CurrentParty  Party  `json:"current_party"`
	CurrentRiding Riding `json:"current_riding"`
	Image         string `json:"image"`
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

func getBills() error {
	resp, err := http.Get("https://api.openparliament.ca/bills/?format=json")
	if err != nil {
		return fmt.Errorf("failed to make HTTP request for Bills: %w", err)
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
		_, err = appDB.Exec(`
		INSERT INTO bills (Session, Legisinfo_id, Introduced, Name, Number, Url) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		ON CONFLICT (Legisinfo_id) DO UPDATE SET  
		session = EXCLUDED.session, 
		introduced = EXCLUDED.introduced, 
		name = EXCLUDED.name, 
		number = EXCLUDED.number, 
		url = EXCLUDED.url `,
			bill.Session, bill.Legisinfo_id, bill.Introduced, bill.Name.EN, bill.Number, bill.Url)
		if err != nil {
			fmt.Printf("Error inserting bill into database: %v\n", err)
			fmt.Printf("Bill: %v\n", bill)
		}
	}
	return nil
}

func getMPs() error {
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
		_, err = appDB.Exec(`
		INSERT INTO mps (Name, CurrentParty, CurrentRiding, Url, Image, province) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		ON CONFLICT (Name) DO UPDATE SET  
		CurrentParty = EXCLUDED.CurrentParty, 
		CurrentRiding = EXCLUDED.CurrentRiding, 
		Url = EXCLUDED.Url, 
		image = EXCLUDED.Image,
		province = EXCLUDED.province`,
			mp.Name, mp.CurrentParty.ShortName.EN, mp.CurrentRiding.Name.EN, mp.URL, mp.Image, mp.CurrentRiding.Province)
		if err != nil {
			fmt.Printf("Error inserting MP into database: %v\n", err)
			fmt.Printf("MP: %v\n", mp)
		}
	}
	return nil
}

func getData(wg *sync.WaitGroup) {
	defer wg.Done()

	var internalWg sync.WaitGroup

	internalWg.Add(1)
	go func() {
		defer internalWg.Done()
		if err := getBills(); err != nil {
			log.Printf("Error getting bills: %v\n", err)
		}
	}()

	internalWg.Add(1)
	go func() {
		defer internalWg.Done()
		if err := getMPs(); err != nil {
			log.Printf("Error getting MPs: %v\n", err)
		}
	}()

	internalWg.Wait()

}

func RunIngest() {
	log.Println("Starting data pulling routines...")

	var wg sync.WaitGroup

	wg.Add(1)
	go getData(&wg)

	wg.Wait()
	log.Println("All data pulling routines completed.")
}
