package main

import (
	"fmt"
	"log"

	"github.com/andrew-tawfik/outreach-routing/internal/app"
	"github.com/andrew-tawfik/outreach-routing/internal/database"
)

func tui() {

	// Prompt user for the Google Sheet URL that contains event and guest information
	fmt.Println("Please provide Google SheetURL")

	var googleSheetURL string
	fmt.Scanln(&googleSheetURL)

	// Extract the spreadsheet ID from the provided URL
	spreadsheetID, err := database.ExtractIDFromURL(googleSheetURL)

	if err != nil {
		log.Fatalf("error extracting ID: %v", err)
	}

	// Initialize Google Sheets client
	db, err := database.NewSheetClient(spreadsheetID)
	if err != nil {
		log.Fatalf("could not initialize sheet client: %v", err)
	}

	// Parse event and guest data from the sheet
	event, err := db.ProcessEvent()
	if err != nil {
		log.Fatalf("Could not process event: %v", err)
	}

	// Map database event to geo event structure for coordinate lookup
	geoEvent := mapDatabaseEventToHttp(event)

	// Filter guests who require transportation service (ie Confirmed or GroceryOnly)
	geoEvent.FilterGuestForService()

	// Request GPS coordinates for all guest addresses
	geoEvent.RequestGuestCoordiantes()

	// Request distance matrix for all coordinates
	geoEvent.RetreiveDistanceMatrix()

	// Map geo-level data to app-level event and location registry
	appEvent, lr := mapDatabaseGeoEventToApp(geoEvent)

	// Set vehicle count
	vehicleCount := 8

	// Initialize the route manager
	RouteManager := app.CreateRouteManager(lr, vehicleCount)

	// Determine savings between guest pairings using Clarke-Wright Algorithm
	RouteManager.DetermineSavingList(lr)

	// Start Route Dispatch Algorithm
	RouteManager.StartRouteDispatch()

	// Display Results
	RouteManager.Display(appEvent, lr)
}
