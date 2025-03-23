package main

import (
	"github.com/andrew-tawfik/outreach-routing/internal/http"
	"github.com/andrew-tawfik/outreach-routing/internal/repository"
)

func mapRepoEventToHttp(repoEvent *repository.Event) *http.Event {
	numGuests := len((*repoEvent).Guests)
	httpGuests := make([]http.Guest, 0, numGuests)

	for _, g := range (*repoEvent).Guests {
		convertedGuest := http.Guest{
			Status:  http.GuestStatus(g.Status),
			Name:    g.Name,
			Address: g.Address,
		}
		httpGuests = append(httpGuests, convertedGuest)
	}
	return &http.Event{
		Guests:    httpGuests,
		EventType: repoEvent.EventType,
	}
}
