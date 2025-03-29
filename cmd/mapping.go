package main

import (
	"github.com/andrew-tawfik/outreach-routing/internal/app"
	"github.com/andrew-tawfik/outreach-routing/internal/database"
	"github.com/andrew-tawfik/outreach-routing/internal/geoapi"
)

func mapDatabaseEventToHttp(dbEvent *database.Event) *geoapi.Event {
	numGuests := len((*dbEvent).Guests)
	httpGuests := make([]geoapi.Guest, 0, numGuests)

	for _, g := range (*dbEvent).Guests {
		convertedGuest := geoapi.Guest{
			Status:    geoapi.GuestStatus(g.Status),
			Name:      g.Name,
			GroupSize: g.GroupSize,
			Address:   g.Address,
		}
		httpGuests = append(httpGuests, convertedGuest)
	}
	return &geoapi.Event{
		Guests:    httpGuests,
		EventType: dbEvent.EventType,
	}
}

func mapDatabaseGeoEventToApp(geoEvent *geoapi.Event) (*app.Event, *app.LocationRegistry) {
	numGuests := len(geoEvent.Guests)
	appGuests := make([]app.Guest, 0, numGuests)

	for _, g := range geoEvent.Guests {
		convertedGuest := app.Guest{
			Name:        g.Name,
			GroupSize:   g.GroupSize,
			Coordinates: app.GuestCoordinates(g.Coordinates),
		}
		appGuests = append(appGuests, convertedGuest)
	}

	// Convert DestinationOccupancy
	appDestOccupancy := make(map[app.GuestCoordinates]int, len(geoEvent.GuestLocations.CoordianteMap.DestinationOccupancy))
	for coord, count := range geoEvent.GuestLocations.CoordianteMap.DestinationOccupancy {
		appCoord := app.GuestCoordinates(coord)
		appDestOccupancy[appCoord] = count
	}

	// Convert CoordinateToAddress
	appCoordToAddr := make(map[app.GuestCoordinates]string, len(geoEvent.GuestLocations.CoordianteMap.CoordinateToAddress))
	for coord, address := range geoEvent.GuestLocations.CoordianteMap.CoordinateToAddress {
		appCoord := app.GuestCoordinates(coord)
		appCoordToAddr[appCoord] = address
	}

	appCoordMap := app.CoordinateMapping{
		DestinationOccupancy: appDestOccupancy,
		CoordinateToAddress:  appCoordToAddr,
		AddressOrder:         geoEvent.GuestLocations.CoordianteMap.AddressOrder,
	}

	return &app.Event{
			Guests:    appGuests,
			EventType: geoEvent.EventType,
		}, &app.LocationRegistry{
			DistanceMatrix: geoEvent.GuestLocations.DistanceMatrix,
			CoordianteMap:  appCoordMap,
		}
}
