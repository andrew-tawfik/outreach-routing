package main

import (
	"fmt"
	"log"

	"github.com/andrew-tawfik/outreach-routing/internal/app"
	"github.com/andrew-tawfik/outreach-routing/internal/database"
)

func main() {
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

	geoEvent.RetreiveDistanceMatrix()
	appEvent, lr := mapDatabaseGeoEventToApp(geoEvent)

	// appEvent, lr, err := app.LoadAppDataFromFile("data.json")
	// if err != nil {
	// 	log.Fatal("Error loading:", err)
	// } else {
	// 	fmt.Println("Successfully read")
	// }
	// Step 2. Fetch addresses exact coordinates that will be utilized

	err = app.SaveAppDataToFile("data.json", *appEvent, *lr)
	if err != nil {
		log.Fatalf("Could not save data: %v", err)
	}

	RouteManager := app.CreateRouteManager(lr, 8)
	RouteManager.DetermineSavingList(lr)
	RouteManager.StartRouteDispatch()

	RouteManager.DisplayResults()
	appEvent.Display()

	// RouteManager.TestRemoveAll()

	// Step 3. Fetch distance matrix

	// Step 4. Determine the best route with RSP algorithm

	// Step 5. Output: display the routes
}
