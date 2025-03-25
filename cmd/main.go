package main

import (
	"fmt"
	"log"

	"github.com/andrew-tawfik/outreach-routing/internal/database"
	"github.com/andrew-tawfik/outreach-routing/internal/geoapi"
)

func main() {
	// Step 1. Open repository and initialize program

	//Retreive guests names and addresses who will require a service
	fmt.Println("Please provide Google SheetURL")

	var googleSheetURL string
	fmt.Scanln(&googleSheetURL)

	spreadsheetID, err := database.ExtractIDFromURL(googleSheetURL)

	if err != nil {
		log.Fatalf("error extracting ID: %v", err)
	}

	db, err := database.NewSheetClient(spreadsheetID)
	if err != nil {
		log.Fatalf("could not initialize sheet client: %v", err)
	}

	event, err := db.ProcessEvent()
	if err != nil {
		log.Fatalf("Could not process event: %v", err)
	}

	httpEvent := mapDatabaseEventToHttp(event)
	httpEvent.FilterGuestForService()

	// Print Array of Guests
	for i, g := range httpEvent.Guests {
		fmt.Println(i, g)
	}

	address := "171 George St"
	geoapi.GeocodeAddress(address)

	// Step 2. Fetch addresses exact coordinates that will be utilized

	// Step 3. Fetch distance matrix

	// Step 4. Determine the best route with RSP algorithm

	// Step 5. Output: display the routes
}
