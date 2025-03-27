package main

import (
	"fmt"
	"log"

	"github.com/andrew-tawfik/outreach-routing/internal/database"
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

	geoEvent := mapDatabaseEventToHttp(event)
	geoEvent.FilterGuestForService()
	geoEvent.RequestGuestCoordiantes()

	// Print Array of Guests
	for i, g := range geoEvent.Guests {
		fmt.Println(i, g)
	}

	fmt.Println("Coordinate list:", geoEvent.GuestLocations.CoordListString)
	fmt.Println("String size: ", len(geoEvent.GuestLocations.CoordListString))
	fmt.Println(geoEvent.GuestLocations.GuestCountByCoord)

	// Step 2. Fetch addresses exact coordinates that will be utilized

	// Step 3. Fetch distance matrix

	// Step 4. Determine the best route with RSP algorithm

	// Step 5. Output: display the routes
}
