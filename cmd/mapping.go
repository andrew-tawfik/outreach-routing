package main

import (
	"github.com/andrew-tawfik/outreach-routing/internal/database"
	"github.com/andrew-tawfik/outreach-routing/internal/geoapi"
)

func mapDatabaseEventToHttp(dbEvent *database.Event) *geoapi.Event {
	numGuests := len((*dbEvent).Guests)
	httpGuests := make([]geoapi.Guest, 0, numGuests)

	for _, g := range (*dbEvent).Guests {
		convertedGuest := geoapi.Guest{
			Status:  geoapi.GuestStatus(g.Status),
			Name:    g.Name,
			Address: g.Address,
		}
		httpGuests = append(httpGuests, convertedGuest)
	}
	return &geoapi.Event{
		Guests:    httpGuests,
		EventType: dbEvent.EventType,
	}
}
