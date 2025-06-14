package main // Declares the package this file belongs to. 'main' is special for executables.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

type MP struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Constituency string `json:"constituency"`
	Party        string `json:"party"`
}

func getBills() {
	resp, err := http.Get("https://api.openparliament.ca/bills/?format=json")
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var billResponse BillResponse

	err = json.Unmarshal(body, &billResponse)
	if err != nil {
		print(err)
	}
	if len(billResponse.Objects) == 0 {
		print("No bills found in response")
		return
	}

	fmt.Printf("successfully got %d bills\n", len(billResponse.Objects))
	for i, bill := range billResponse.Objects {
		fmt.Printf("Bill %d: %s\n", i, bill.Number)
		fmt.Printf("Bill %d: %s\n", i, bill.Name.EN)
		fmt.Printf("Bill %d: %s\n", i, bill.Name.FR)
		fmt.Printf("Bill %d: %s\n", i, bill.Introduced)
		fmt.Printf("Bill %d: %s\n", i, bill.Url)
	}
}

func getMPs() {
	resp, err := http.Get("https://api.openparliament.ca/politicians/?format=json")
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var mps []MP
	err = json.Unmarshal(body, &mps)
	if err != nil {
		print(err)
	}
}

func getData() {
	go getBills()
	//go getMPs()

}

func main() { // This is the main function, the entry point of your program.
	print("Starting Go Routines")
	go getData()
	time.Sleep(time.Second * 10)
}
