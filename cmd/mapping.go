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
			Coordinates: g.Coordinates,
		}
		appGuests = append(appGuests, convertedGuest)
	}

	appCoordMap := app.CoordinateMapping{
		DestinationOccupancy: geoEvent.GuestLocations.CoordianteMap.DestinationOccupancy,
		CoordinateToAddress:  geoEvent.GuestLocations.CoordianteMap.CoordinateToAddress,
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
