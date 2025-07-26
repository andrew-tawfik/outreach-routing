package ui

import (
	

	"fmt"
	"log"

	"github.com/andrew-tawfik/outreach-routing/internal/app"
	"github.com/andrew-tawfik/outreach-routing/internal/converter"
	"github.com/andrew-tawfik/outreach-routing/internal/database"
)

type RoutingProcess struct {
	rm *app.RouteManager
	ae *app.Event
	lr *app.LocationRegistry
}

func ProcessEvent(googleSheetURL string) (*RoutingProcess, error) {

	
	spreadsheetID, err := database.ExtractIDFromURL(googleSheetURL)

	if err != nil {
		return nil, fmt.Errorf("error extracting ID: %v", err)
	}

	
	db, err := database.NewSheetClient(spreadsheetID)
	if err != nil {
		return nil, fmt.Errorf("could not initialize sheet client: %v", err)
	}

	
	event, err := db.ProcessEvent()
	if err != nil {
		return nil, fmt.Errorf("could not process event: %v", err)
	}

	
	geoEvent := converter.MapDatabaseEventToHttp(event)

	
	err = geoEvent.RequestGuestCoordiantes()
	if err != nil {
		return nil, fmt.Errorf("could not geocode addresses: %w", err)
	}

	
	err = geoEvent.RetreiveDistanceMatrix()
	if err != nil {
		return nil, fmt.Errorf("could not retreive distance matrix: %w", err)
	}

	
	appEvent, lr := converter.MapDatabaseGeoEventToApp(geoEvent)

	err = app.SaveAppDataToFile("data.json", *appEvent, *lr)
	if err != nil {
		log.Printf("Warning: Could not save data to JSON: %v", err)
	} else {
		log.Println("Data saved to data.json for testing")
	}

	RouteManager := app.OrchestateDispatch(lr, appEvent)

	return &RoutingProcess{
		rm: RouteManager,
		ae: appEvent,
		lr: lr,
	}, nil
}

func (rp *RoutingProcess) String() string {
	return rp.rm.Display(rp.ae, rp.lr)
}

func ProcessJsonEvent(eventType int) (*RoutingProcess, error) {

	var appEvent app.Event
	var lr app.LocationRegistry
	var err error
	if eventType == 0 {
		appEvent, lr, err = app.LoadAppDataFromFile("data_dinner.json")
		if err != nil {
			return nil, fmt.Errorf("could not load json event information. %w", err)
		}
	} else {
		appEvent, lr, err = app.LoadAppDataFromFile("data_grocery.json")
		if err != nil {
			return nil, fmt.Errorf("could not load json event information. %w", err)
		}
	}

	RouteManager := app.OrchestateDispatch(&lr, &appEvent)

	RouteManager.Display(&appEvent, &lr)

	return &RoutingProcess{
		rm: RouteManager,
		ae: &appEvent,
		lr: &lr,
	}, nil

}
