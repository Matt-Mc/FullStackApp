package main // Declares the package this file belongs to. 'main' is special for executables.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
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

func getBills() error {
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
	for i, bill := range billResponse.Objects {
		fmt.Printf("Bill %d: %s\n", i, bill.Number)
		fmt.Printf("Bill %d: %s\n", i, bill.Name.EN)
		fmt.Printf("Bill %d: %s\n", i, bill.Name.FR)
		fmt.Printf("Bill %d: %s\n", i, bill.Introduced)
		fmt.Printf("Bill %d: %s\n", i, bill.Url)
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
	for i, mp := range mps.Objects {
		fmt.Printf("MP %d: %s\n", i, mp.Name)
		fmt.Printf("MP %d: %s\n", i, mp.CurrentParty.ShortName.EN)
		fmt.Printf("MP %d: %s\n", i, mp.CurrentRiding.Name.EN)
		fmt.Printf("MP %d: %s\n", i, mp.CurrentRiding.Province)
		fmt.Printf("MP %d: %s\n", i, mp.Image)
	}
	return nil
}

func getData(wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := getBills(); err != nil {
			fmt.Printf("Error getting bills: %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := getMPs(); err != nil {
			fmt.Printf("Error getting MPs: %v\n", err)
		}
	}()

}

func main() { // This is the main function, the entry point of your program.
	fmt.Println("Starting Go Routines")

	var wg sync.WaitGroup

	wg.Add(1)
	go getData(&wg)

	wg.Wait()
	fmt.Println("All Go Routines completed")
}
